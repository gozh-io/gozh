FROM golang:1.10-alpine3.7 AS builder

#param
ARG PROJECT_URL=github.com/gozh-io/gozh
ARG PROJECT_NAME=gozh

# Install tools required to build the project
# We need to run `docker build --no-cache .` to update those dependencies
#RUN apk add --no-cache git  && \
#    go get -u github.com/tools/godep

ENV APP_DIR          $GOPATH/src/$PROJECT_URL
ENV APP_CONFIG_DIR   $APP_DIR

RUN mkdir -p $APP_DIR
COPY . $APP_DIR

WORKDIR $APP_DIR
RUN go build -o $PROJECT_NAME .  



FROM alpine:latest  
RUN apk --no-cache add ca-certificates
RUN mkdir /myapp
WORKDIR   /myapp


COPY --from=builder /go/src/github.com/gozh-io/gozh/gozh .
COPY --from=builder /go/src/github.com/gozh-io/gozh/conf conf

VOLUME /myapp 
EXPOSE 80

CMD  ./gozh conf/cf.json >> stdout.log 2>&1