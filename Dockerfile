FROM golang

WORKDIR /app

COPY gogabo/* ./
    
RUN go mod download

RUN go get gogabo

RUN go build

RUN ls

ENV DG_TOKEN=""

ENTRYPOINT [ "sh", "-c", "go run /app/main.go"]
