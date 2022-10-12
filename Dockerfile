FROM golang:latest
WORKDIR /app

COPY . .
RUN go mod download

ENV PORT 8080
CMD [ "go","run","cmd/main.go" ]