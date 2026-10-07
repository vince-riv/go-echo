FROM --platform=$BUILDPLATFORM golang:1.27@sha256:162be5298a40ed317005c8339c6de4d10d3eef336d66dc8e9259b03ab9d3a6d2 AS build
ARG TARGETOS
ARG TARGETARCH
WORKDIR /src
COPY go.mod main.go ./
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w" -o /go-echo .

FROM scratch
COPY --from=build /go-echo /go-echo
EXPOSE 8080
USER 65534:65534
ENTRYPOINT ["/go-echo"]
