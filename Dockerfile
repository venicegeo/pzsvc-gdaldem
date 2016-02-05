FROM klokantech/gdal

ENTRYPOINT ["/pzsvc-gdaldem"]

COPY pzsvc-gdaldem /pzsvc-gdaldem
RUN chmod a+x /pzsvc-gdaldem
