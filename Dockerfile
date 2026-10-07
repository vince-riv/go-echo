FROM golang:1.26 AS build
WORKDIR /src
COPY go.mod main.go ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /go-echo .

FROM scratch
COPY --from=build /go-echo /go-echo
EXPOSE 8080
USER 65534:65534
ENTRYPOINT ["/go-echo"]
