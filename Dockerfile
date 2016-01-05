FROM venicegeo/gdaldem
EXPOSE 8080
WORKDIR /app
# copy binary into image
COPY gdaldem /app/
ENTRYPOINT ["./gdaldem"]
