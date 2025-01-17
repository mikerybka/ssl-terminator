FROM alpine:latest

RUN apk update
RUN apk add go

COPY . /app
WORKDIR /app

RUN go build -o /bin/app main.go

ENTRYPOINT ["/bin/app"]
