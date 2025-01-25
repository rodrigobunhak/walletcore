FROM golang:1.23.0

WORKDIR /app/

RUN apt-get update && apt-get install -y librdkafka-dev

CMD ["tail", "-f", "/dev/null"]

