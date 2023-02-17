FROM golang:1.19

ENV config=docker

WORKDIR /app

RUN mkdir "/build"

COPY ./ /app

RUN go mod download

RUN go get github.com/githubnemo/CompileDaemon

EXPOSE 5555

ENTRYPOINT CompileDaemon --build="go build -o /build/app" --command=./main