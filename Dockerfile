FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o note-tracker cmd/notes/main.go
RUN ls -la

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/note-tracker .
CMD [ "./note-tracker" ]
EXPOSE 8080