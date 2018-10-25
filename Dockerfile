# This file is only for build purposes

FROM node:slim as build-node
ARG NODE_ENV=production
ADD web /srv/web
RUN cd /srv/web && npm install && npm run build && rm -rf ./node_modules

FROM golang:1.11-stretch as build-go
WORKDIR $GOPATH/src/github.com/misha-plus/power-start
ADD . .
COPY --from=build-node /srv/web $GOPATH/src/github.com/misha-plus/power-start/web/build
RUN go get github.com/gobuffalo/packr/packr
RUN make server agent
RUN ls artifacts
