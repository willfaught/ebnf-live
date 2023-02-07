FROM golang:1.20
WORKDIR /app
COPY go.mod go.sum .
RUN go mod download && go mod verify
COPY *.go .
RUN go build -o /usr/local/bin/app
EXPOSE 80
CMD ["app"]
