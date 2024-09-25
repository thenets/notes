FROM docker.io/golang:latest as builder

ENV GOPATH=/go

WORKDIR /go/pkg/mod/github.com/thenets/notes/

# COPY go.mod go.sum main.go static/ kvstore/ .
COPY ./* /go/pkg/mod/github.com/thenets/notes/

RUN set -ex \
    && go mod download \
    && CGO_ENABLED=0 GOOS=linux go build -o /notes


# Final image
FROM gcr.io/distroless/base-debian11 AS build-release-stage
WORKDIR /
COPY --from=builder /notes /notes
USER nonroot:nonroot
EXPOSE 8080
ENV PORT=8080
ENTRYPOINT ["/notes"]
