FROM golang:1.7.3-alpine

ARG GOOS

COPY . /go/src/github.com/git-jiby-me/swarm
WORKDIR /go/src/github.com/git-jiby-me/swarm

RUN set -ex \
	&& apk add --no-cache --virtual .build-deps \
	git \
	&& GOOS=$GOOS CGO_ENABLED=0 go install -v -a -tags netgo -installsuffix netgo -ldflags '-w -X "github.com/git-jiby-me/swarm/version.GITCOMMIT=`git rev-parse --short HEAD`" -X "github.com/git-jiby-me/swarm/version.BUILDTIME=`date -u`"'  \
	&& apk del .build-deps

ENV SWARM_HOST :2375
EXPOSE 2375

VOLUME $HOME/.swarm

ENTRYPOINT ["swarm"]
CMD ["--help"]
