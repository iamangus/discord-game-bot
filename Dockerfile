FROM golang

RUN go get github.com/bwmarrin/discordgo && \
    find / -name *discordgo*

WORKDIR /app

COPY gogabo/main.go ./gogabo/src/

ENV DG_TOKEN=""
ENV GO111MODULE=auto

ENTRYPOINT [ "sh", "-c", "echo $GOPATH && export $GOPATH=$GOPATH:/go/pkg/mod/ && go run /app/gogabo/src/main.go"]
