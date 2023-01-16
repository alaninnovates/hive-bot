FROM golang:1.19-alpine

WORKDIR /hive-builder

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /hive-bot

CMD [ "/hive-bot" ]