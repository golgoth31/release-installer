FROM alpine:latest as certs
RUN apk add -U --no-cache ca-certificates

FROM scratch AS arm64
COPY ri-linux-arm64 /ri
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["/homedynip"]
