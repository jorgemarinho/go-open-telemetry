FROM golang:1.22.4 as build
WORKDIR /app
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o service-a cmd/servicea/main.go

FROM scratch
WORKDIR /app
COPY --from=build /app/service-a .
ENTRYPOINT ["./service-a"]