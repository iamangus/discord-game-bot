FROM golang

RUN go get github.com/bwmarrin/discordgo

copy in gogabo/main.go

ENV DG_TOKEN=""

ENTRYPOINT [ "sh", "-c", "go run main.go"]
