FROM golang:1.22.0

COPY federation /app/

EXPOSE 8080

ENTRYPOINT [ "/app/federation" ]
