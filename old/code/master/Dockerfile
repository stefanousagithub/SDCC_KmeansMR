FROM golang:1.16-alpine

ARG NUMMAP
ARG THRESHOLD
ARG MAXITER
ENV NUMMAP ${NUMMAP}
ENV THRESHOLD ${THRESHOLD}
ENV MAXITER ${MAXITER}

WORKDIR /code

COPY ./points/ ./points/
COPY ./master/ ./master/
COPY ./cluster/ ./cluster/
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum
RUN go mod download gioui.org
CMD [ "sh", "-c", "go run ./master/*.go ${NUMMAP} ${MAXITER} ${THRESHOLD} " ]
