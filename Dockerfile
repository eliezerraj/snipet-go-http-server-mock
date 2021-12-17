# docker build -t go_http_server_mock .
# docker run -dit --name go_http_server_mock -p 8900:8900 go_http_server_mock /bin/bash

FROM golang:1.17 As builder
WORKDIR /app

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-w' -o go_http_server_mock

CMD ["/app/go_http_server_mock"]