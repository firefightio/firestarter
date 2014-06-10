# BUILD-USING:    docker build -t firefightio/firestarter .
# TEST-USING:     docker run --rm -i -t -v /var/run/docker.sock:/docker.sock --name=firestarter-dev --link=rabbitmq:rabbitmq --entrypoint=/bin/bash firefightio/firestarter -s
# RUN-USING:      docker run --rm -v /var/run/docker.sock:/docker.sock --link=rabbitmq:rabbitmq --name=firestarter firefightio/firestarter

FROM google/golang

WORKDIR /gopath/src/firestarter
ADD . /gopath/src/firestarter/
RUN go get
RUN go install firestarter

# Overwrite with runtime arguements
ENV RABBITMQ_USER firestarter
ENV RABBITMQ_PASS test

ENTRYPOINT ["/gopath/bin/firestarter"]
