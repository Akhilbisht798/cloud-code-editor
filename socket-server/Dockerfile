FROM golang

WORKDIR /app

COPY . . 

RUN go mod download

RUN go build -o socket-server

CMD ["./socket-server"]