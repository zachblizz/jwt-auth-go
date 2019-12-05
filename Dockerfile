FROM golang:1.8

WORKDIR /go/src/bitbucket/zblizz/jwt-go
COPY . .

RUN go get
RUN go build

ENTRYPOINT [ "./jwt-go" ]

# TODO: need to add a mongo service to the mix
