FROM golang:1.23

WORKDIR /app

COPY go.mod ./
RUN go mod download && go mod verify

COPY /internal ./internal
COPY /cmd ./cmd

COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /protohackers

EXPOSE 8080

CMD ["/protohackers"]
