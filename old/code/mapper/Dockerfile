FROM golang:1.16-alpine

WORKDIR /code

COPY ./points/ ./points/
COPY ./mapper/ ./mapper/
COPY ./cluster/ ./cluster/
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum
RUN go mod download gioui.org

CMD [ "sh", "-c", "go run ./mapper/*.go " ]
