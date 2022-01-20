FROM golang

WORKDIR /app

RUN go get github.com/bwmarrin/discordgo

COPY gogabo/main.go ./

ENV DG_TOKEN=""
ENV GO111MODULE=auto

ENTRYPOINT [ "sh", "-c", "go run /app/main.go"]
