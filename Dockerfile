FROM golang:1.9

WORKDIR /go/src/github.com/gormDemo
COPY . .

RUN go get -u github.com/golang/dep/cmd/dep

RUN dep ensure -vendor-only
RUN go install 

CMD gormDemo