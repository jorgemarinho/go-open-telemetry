FROM golang:1.23.0-alpine as build
WORKDIR /app
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o service-a cmd/servicea/main.go

FROM scratch
WORKDIR /app
COPY --from=build /app/service-a .
ENTRYPOINT ["./service-a"]