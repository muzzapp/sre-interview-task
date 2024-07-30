FROM golang

WORKDIR /go/src/app

COPY . .

RUN go mod tidy && go build -o ./app ./cmd/api/

CMD ["./app"]
