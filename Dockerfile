FROM golang:latest
WORKDIR /build
ENV GOPROXY https://goproxy.cn
COPY go.mod .
COPY go.sum .
RUN go mod download
ADD . .
WORKDIR ./cmd
RUN go build .
EXPOSE 8080
CMD ["-c","./config.json"]
ENTRYPOINT ["./cmd"]