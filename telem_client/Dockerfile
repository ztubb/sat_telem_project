FROM golang:1.23.2

ADD client.go client.go
ADD api.go api.go
ADD go.mod go.mod
ADD go.sum go.sum

RUN go mod download

CMD [ "go", "run", "client.go", "api.go" ]