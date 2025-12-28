FROM debian:13-slim
WORKDIR /app
COPY smarthome-p1-receiver /app/smarthome-p1-receiver
ENTRYPOINT ["/app/smarthome-p1-receiver"]