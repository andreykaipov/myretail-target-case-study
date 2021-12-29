FROM golang:alpine AS build
WORKDIR /work
COPY . .
RUN go build -o myretail-api

FROM alpine:latest
COPY --from=build /work/myretail-api /
CMD ["/myretail-api"]
