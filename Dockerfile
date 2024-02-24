FROM golang:1.18.10-alpine3.17
WORKDIR /app

COPY go.mod go.sum ./
# RUN GOPROXY=https://goproxy.cn go mod download

COPY . .
RUN go build -o main .

CMD ["./main"]