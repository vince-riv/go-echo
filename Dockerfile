FROM --platform=$BUILDPLATFORM golang:1.26@sha256:eb36c1664dd974cde625f736e02c204383deebe03977365caaec5bf49f794348 AS build
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
