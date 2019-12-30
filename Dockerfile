FROM golang:1.11-stretch AS build

WORKDIR /go/src/github.com/sigmonsays/paste/pasted
RUN apt-get install \
    git gcc g++ binutils
COPY . /go/src/github.com/sigmonsays/paste/
RUN go get -d .
ENV GOPATH=/go
RUN go install -ldflags '-w -extldflags "-static"' github.com/sigmonsays/paste/...

# ---

FROM golang:1.11-stretch
COPY --from=build /go/bin/pasted /pasted
EXPOSE 7001
CMD [ "/pasted", "-bindaddr", ":7001"  ]
