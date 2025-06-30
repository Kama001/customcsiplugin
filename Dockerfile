FROM ubuntu

RUN apt-get update && apt-get install -y ca-certificates

RUN mkdir -p /home/csi

COPY go.mod /home/csi/
COPY go.sum /home/csi/
COPY main.go /home/csi
COPY pkg /home/csi

COPY bsos /usr/local/bin/bsos

ENTRYPOINT [ "/usr/local/bin/bsos" ]