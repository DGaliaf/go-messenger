# Build stage
FROM golang:1.17-alpine as build

WORKDIR /app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

RUN go build -o myapp ./app/cmd/app/main.go

# Final stage
FROM alpine:3.14

WORKDIR /app

COPY --from=build /app/myapp .

EXPOSE 8080

CMD ["./myapp"]
