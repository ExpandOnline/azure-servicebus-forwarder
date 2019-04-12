FROM golang:1.11-alpine as builder

RUN apk add --no-cache git openssh-client curl


RUN mkdir /root/.ssh && ssh-keyscan github.com >> /root/.ssh/known_hosts
ARG COMMIT_HASH
RUN git clone https://github.com/ExpandOnline/azure-servicebus-forwarder.git $GOPATH/src/github.com/expandonline/azure-servicebus-forwarder \
  && cd $GOPATH/src/github.com/expandonline/azure-servicebus-forwarder \
  && git checkout $COMMIT_HASH

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN cd $GOPATH/src/github.com/expandonline/azure-servicebus-forwarder \
  && dep ensure \
  && go install -v -ldflags "-X main.version=${COMMIT_HASH}" github.com/expandonline/azure-servicebus-forwarder


FROM alpine
WORKDIR /app
RUN apk add --no-cache tzdata ca-certificates supervisor \
  && cp /usr/share/zoneinfo/Europe/Amsterdam /etc/localtime \
  && echo "Europe/Amsterdam" > /etc/timezone \
  && apk del tzdata
RUN mkdir -p /aws && \
	apk -Uuv add groff less python py-pip && \
	pip install awscli && \
	apk --purge -v del py-pip && \
	rm /var/cache/apk/*
COPY --from=builder /go/bin/azure-servicebus-forwarder /app/redirect
COPY ./entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

ENTRYPOINT ["/entrypoint.sh"]
