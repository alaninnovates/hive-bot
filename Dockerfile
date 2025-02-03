FROM golang:1.22-alpine as build

WORKDIR /hive-builder

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /hive-builder-bot .


FROM alpine

COPY --from=build /hive-builder-bot /hive-builder-bot
COPY --from=build /hive-builder/assets /assets

RUN mkdir -p /data && echo "[]" > /data/hives.json

VOLUME /data

CMD ["/hive-builder-bot"]