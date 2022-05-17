FROM golang:latest
WORKDIR /test_project
COPY ./ ./
RUN go clean --modcache
RUN go mod download
RUN go get github.com/lib/pq
RUN go build -o /build cmd/main.go
RUN ls
EXPOSE 7000

RUN git clone https://github.com/vishnubob/wait-for-it.git

CMD ["/build"]