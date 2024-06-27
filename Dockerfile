FROM golang:1.22

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -buildmode=plugin -o ./compiled_handlers/postgresql.so postgresql/postgresql.go
RUN go build -v -o /usr/local/bin/app ./relay

CMD ["app"]