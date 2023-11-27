FROM golang:1.21.4-bullseye

COPY federation /app/

EXPOSE 8080

ENTRYPOINT [ "/app/federation" ]