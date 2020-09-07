FROM alpine:latest as certs
RUN apk add -U --no-cache ca-certificates

FROM scratch AS amd64
COPY ri-linux-amd64 /ri
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["/ri"]
