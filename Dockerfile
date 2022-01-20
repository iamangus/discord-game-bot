FROM golang

WORKDIR /app

COPY gogabo/* ./

RUN go mod download github.com/bwmarrin/discordgo

ENV DG_TOKEN=""

ENTRYPOINT [ "sh", "-c", "go run /app/main.go"]