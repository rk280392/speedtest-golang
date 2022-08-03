FROM golang:1.19.0-alpine3.16
WORKDIR /usr/src/app
COPY go.mod go.sum ./
# RUN go mod download && go mod verify
# RUN go build -v -o /usr/local/bin/app ./...
COPY speedtest-app.go speedtest-go crontab.txt entry.sh ./
RUN chmod 755 ./entry.sh
RUN /usr/bin/crontab ./crontab.txt
CMD ["./entry.sh"]
