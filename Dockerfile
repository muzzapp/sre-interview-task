# Pin version of Go
FROM golang

WORKDIR /go/src/app

COPY . .

RUN go mod tidy && go build -o ./app ./cmd/api/

# Copy image insto scratch
CMD ["./app"]
