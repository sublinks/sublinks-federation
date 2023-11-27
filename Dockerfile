FROM golang:1.21.4

COPY federation /app/

EXPOSE 8080

ENTRYPOINT [ "/app/federation" ]