FROM golang:1.12.4 as build
ENV GO111MODULE on

WORKDIR /go/cache

ADD go.mod .
ADD go.sum .
RUN go mod download

WORKDIR /go/release

ADD . .

RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -installsuffix cgo -o scheduler main.go

FROM centos:7 as prod
WORKDIR /
COPY --from=build /go/release/scheduler .
CMD ["/scheduler daemon"]
