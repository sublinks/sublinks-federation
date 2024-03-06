FROM golang:1.22.1

COPY federation /app/

EXPOSE 8080

ENTRYPOINT [ "/app/federation" ]
