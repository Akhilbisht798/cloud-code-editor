FROM golang
# use ubuntu that will have many command
# add a user
# add packages like npm,node,go,python for easy head start. 

WORKDIR /app

COPY . . 

RUN go mod download

RUN go build -o socket-server

CMD ["./socket-server"]