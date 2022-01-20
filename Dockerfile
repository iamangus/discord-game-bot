FROM golang

RUN go get github.com/bwmarrin/discordgo

WORKDIR /app

COPY gogabo/main.go ./

ENV DG_TOKEN=""
ENV GO111MODULE=auto

ENTRYPOINT [ "sh", "-c", "go run /app/main.go"]
