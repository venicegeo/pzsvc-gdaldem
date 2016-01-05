package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/venicegeo/pzsvc-gdaldem/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws"
	"github.com/venicegeo/pzsvc-gdaldem/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/venicegeo/pzsvc-gdaldem/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws/session"
	"github.com/venicegeo/pzsvc-gdaldem/Godeps/_workspace/src/github.com/aws/aws-sdk-go/service/s3"
	"github.com/venicegeo/pzsvc-gdaldem/Godeps/_workspace/src/github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/venicegeo/pzsvc-gdaldem/Godeps/_workspace/src/github.com/julienschmidt/httprouter"
	"github.com/venicegeo/pzsvc-gdaldem/Godeps/_workspace/src/github.com/venicegeo/pzsvc-pdal/objects"
)

// S3Bucket defines the expected JSON structure for S3 buckets.
// An S3 bucket can be used for source (input) and destination (output) files.
type S3Bucket struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
}

// JobInput defines the expected input JSON structure.
// We currently support S3 input (bucket/key), though provider-specific (e.g.,
// GRiD) may be legitimate.
type JobInput struct {
	Source      S3Bucket         `json:"source"`
	Function    *string          `json:"function"`
	Options     *json.RawMessage `json:"options"`
	Destination S3Bucket         `json:"destination"`
}

// JobOutput defines the expected output JSON structure.
type JobOutput struct {
	Input      JobInput                    `json:"input"`
	StartedAt  time.Time                   `json:"started_at"`
	FinishedAt time.Time                   `json:"finished_at"`
	Code       int                         `json:"code"`
	Message    string                      `json:"message"`
	Response   map[string]*json.RawMessage `json:"response"`
}

// GdalDemHandler handles PDAL jobs.
func GdalDemHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Create the job output message. No matter what happens, we should always be
	// able to populate the StartedAt field.
	var res objects.JobOutput
	res.StartedAt = time.Now()

	// There should always be a body, else how are we to know what to do? Throw
	// 400 if missing.
	if r.Body == nil {
		http.Error(w, "no body", http.StatusBadRequest)
		return
	}

	// Throw 500 if we cannot read the body.
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error with body", http.StatusInternalServerError)
		return
	}

	// Throw 400 if we cannot unmarshal the body as a valid JobInput.
	var msg objects.JobInput
	if err := json.Unmarshal(b, &msg); err != nil {
		http.Error(w, "error with json", http.StatusBadRequest)
		return
	}

	// Throw 400 if the JobInput does not specify a function.
	if msg.Function == nil {
		http.Error(w, "must provide a function", http.StatusBadRequest)
		return
	}

	fmt.Println(msg.Source)
	fmt.Println(*msg.Function)
	fmt.Println(msg.Destination)

	// This block stolen from pzsvc-pdal. Maybe we should break it out into it's own package?
	var fileIn, fileOut *os.File

	// Split the source S3 key string, interpreting the last element as the
	// input filename. Create the input file, throwing 500 on error.
	fileIn, err = os.Create("elevation.tif")
	if err != nil {
		log.Fatal("Cannot create elevation.tif")
	}
	defer fileIn.Close()

	// If provided, split the destination S3 key string, interpreting the last
	// element as the output filename. Create the output file, throwing 500 on
	// error.
	fileOut, err = os.Create("hillshade.tif")
	if err != nil {
		log.Fatal("Cannot create hillshade.tif")
	}
	defer fileOut.Close()

	log.Println(msg.Source.Bucket)
	log.Println(msg.Source.Key)
	log.Println(fileIn.Name())
	downloader := s3manager.NewDownloader(session.New(&aws.Config{Region: aws.String("us-east-1")}))
	numBytes, err := downloader.Download(fileIn,
		&s3.GetObjectInput{
			Bucket: aws.String(msg.Source.Bucket),
			Key:    aws.String(msg.Source.Key),
		})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			log.Println("Error:", awsErr.Code(), awsErr.Message())
		} else {
			fmt.Println(err.Error())
		}
		return
	}
	log.Println("Downloaded", numBytes, "bytes")
	// end stolen from pzsvc-pdal

	var args []string
	args = append(args, *msg.Function)
	args = append(args, fileIn.Name())
	args = append(args, fileOut.Name())
	out, _ := exec.Command("gdaldem", args...).CombinedOutput()
	fmt.Println(string(out))

	// more from pzsvc-pdal
	// If an output has been created, upload the destination data to S3,
	// throwing 500 on error.
	uploader := s3manager.NewUploader(session.New(&aws.Config{Region: aws.String("us-east-1")}))
	result, err := uploader.Upload(&s3manager.UploadInput{
		Body:   fileOut,
		Bucket: aws.String(msg.Destination.Bucket),
		Key:    aws.String(msg.Destination.Key),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			log.Println("Error:", awsErr.Code(), awsErr.Message())
		} else {
			fmt.Println(err.Error())
		}
		return
	}
	log.Println("Successfully uploaded to", result.Location)
	// end from pzsvc-pdal
}
