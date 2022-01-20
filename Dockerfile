FROM golang

RUN go get github.com/bwmarrin/discordgo && \
    go install github.com/bwmarrin/discordgo@latest

WORKDIR /app

COPY gogabo/main.go ./gogabo/src/

ENV DG_TOKEN=""
ENV GO111MODULE=auto

ENTRYPOINT [ "sh", "-c", "go run /app/gogabo/src/main.go"]