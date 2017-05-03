FROM golang:1.8

COPY . /go/src/github.com/alvinfeng/mosaic
RUN go install -v github.com/alvinfeng/mosaic
EXPOSE 8080
ENTRYPOINT ["mosaic"]
