# we need the root certificates from Alpine linux.
FROM alpine:latest
RUN apk add gcompat

COPY pkce /bin/pkce

LABEL org.opencontainers.image.source=https://github.com/commandquery/pkce
LABEL org.opencontainers.image.description="PKCE code exchange implementation"
LABEL org.opencontainers.image.licenses=MIT

EXPOSE 8080

ENTRYPOINT ["/bin/pkce"]