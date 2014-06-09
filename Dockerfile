# BUILD-USING:    docker build -t firefight/firestarter .
# TEST-USING:     docker run --rm -i -t -v /var/run/docker.sock:/docker.sock --name=firestarter-dev --entrypoint=/bin/bash firefight/firestarter -s
# RUN-USING:      docker run -v /var/run/docker.sock:/docker.sock --name=firestarter firefight/firestarter

FROM google/golang

WORKDIR /gopath/src/firestarter
ADD . /gopath/src/firestarter/
RUN go get
RUN go install firestarter

ENTRYPOINT ["/gopath/bin/firestarter"]
