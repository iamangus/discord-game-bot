FROM golang

WORKDIR /app

COPY gogabo/* ./
    
RUN go mod download

RUN go build

ENV DG_TOKEN=""

ENTRYPOINT [ "sh", "-c", "go run /app/main.go"]
