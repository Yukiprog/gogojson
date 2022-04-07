FROM golang:alpine3.15
WORKDIR /api
COPY ./main.go .
RUN go mod init todo && go mod tidy
CMD ["go", "run", "./main.go"]