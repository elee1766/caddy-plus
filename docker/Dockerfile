# syntax=docker/dockerfile:1

# CA certificates are required for ACME (Let's Encrypt etc.) outbound TLS.
FROM alpine:3.24@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b AS certs
RUN apk add --no-cache ca-certificates

FROM scratch
ARG TARGETOS
ARG TARGETARCH

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=certs /etc/passwd /etc/group /etc/
COPY dist/${TARGETOS}-${TARGETARCH}/caddy /usr/bin/caddy

# Run as nobody. Binding to 80/443 works in containers because Docker sets
# net.ipv4.ip_unprivileged_port_start=0 by default.
USER 65534:65534

EXPOSE 80 443 443/udp 2019

ENTRYPOINT ["/usr/bin/caddy"]
CMD ["run", "--config", "/etc/caddy/Caddyfile", "--adapter", "caddyfile"]
