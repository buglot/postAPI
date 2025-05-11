FROM golang:1.24
WORKDIR /api
COPY . .
RUN go mod download
EXPOSE 8080
CMD [ "go","run","." ]
