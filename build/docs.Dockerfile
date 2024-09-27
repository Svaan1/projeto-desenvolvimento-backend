FROM golang:1.23-alpine

RUN go install golang.org/x/tools/cmd/godoc@latest

ENV PATH="$PATH:/home/app/go/bin"
WORKDIR /docs

COPY . .

CMD ["godoc", "-http", "0.0.0.0:6060"]