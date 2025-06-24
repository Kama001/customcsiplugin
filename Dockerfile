FROM golang

COPY bsos /usr/local/bin/bsos

ENTRYPOINT [ "/usr/local/bin/bsos" ]