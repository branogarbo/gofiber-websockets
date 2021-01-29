FROM golang
WORKDIR /go/src/websockets
COPY . .
RUN go install
EXPOSE 3004
CMD ["/go/bin/websockets"]