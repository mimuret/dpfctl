FROM golang:1.16-alpine3.13 as build
WORKDIR /var/src
RUN apk add git
ADD go.mod .
ADD go.sum .
RUN go mod download
ADD . .
RUN go build

FROM alpine:3.13
COPY --from=build /var/src/dpfctl /bin
ENTRYPOINT [ "/bin/dpfctl" ]