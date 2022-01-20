FROM golang

RUN go get github.com/bwmarrin/discordgo

WORKDIR /app

COPY gogabo/* ./gogabo/src/

ENV DG_TOKEN=""

ENTRYPOINT [ "sh", "-c", "go run /app/gogabo/src/main.go"]