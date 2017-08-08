FROM golang
ADD . /go/src/github.com/tobyjsullivan/ues-projection-svc
RUN  go install github.com/tobyjsullivan/ues-projection-svc
CMD /go/bin/ues-projection-svc
