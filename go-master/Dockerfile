FROM golang:1.21.4

WORKDIR /app

COPY . /app/


RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /main cmd/main/main.go

EXPOSE 9001

ENTRYPOINT [ "/main" ]