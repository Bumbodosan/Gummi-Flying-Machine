FROM golang:1.15.2

WORKDIR /gfm

COPY . .

RUN go build -o /build/gfm .

CMD ["/build/gfm"]
