FROM golang

RUN go mod download github.com/bwmarrin/discordgo

WORKDIR /app

COPY gogabo/* ./

ENV DG_TOKEN=""

ENTRYPOINT [ "sh", "-c", "go run /app/main.go"]