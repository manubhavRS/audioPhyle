FROM golang:1.17-alpine3.13
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
COPY . .
RUN go mod download
RUN go build -o /go/bin/audioPhile
EXPOSE 8080
ENTRYPOINT [ "/go/bin/audioPhile" ]