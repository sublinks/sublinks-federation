FROM golang:1.21.4-bullseye

COPY install-turso.sh /tmp/install-turso.sh

RUN chmod +x /tmp/install-turso.sh && \
    /tmp/install-turso.sh && \
    rm /tmp/install-turso.sh && \
    ln -s /root/.turso/turso /bin/turso
