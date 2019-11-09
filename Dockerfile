FROM golang:1.8

WORKDIR /go/src/app
COPY . .

RUN go build

CMD ["app"]


# docker build -t . auth-app .
# docker run -it --rm --name auth-running-app auth-app