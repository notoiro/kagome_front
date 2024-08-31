FROM golang:1.23

COPY . /app
WORKDIR /app
RUN go build main.go
CMD /app/main
