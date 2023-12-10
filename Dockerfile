FROM golang:1.21.5

COPY federation /app/

EXPOSE 8080

ENTRYPOINT [ "/app/federation" ]