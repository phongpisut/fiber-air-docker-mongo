FROM golang:1.16-alpine
RUN mkdir /app
ADD ./api /app/
WORKDIR /app
RUN go get -v github.com/cosmtrek/air
ENTRYPOINT ["air"]