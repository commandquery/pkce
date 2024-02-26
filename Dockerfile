# scratch is an empty image
# FROM scratch
# If you need /bin/sh and a few utilities, uncomment
# the following line. It increases the image by 5.5 MB
FROM alpine:latest

RUN apk add gcompat

COPY pkce /bin/pkce
# copy other files if needed

LABEL org.opencontainers.image.source=https://github.com/commandquery/pkce
LABEL org.opencontainers.image.description="PKCE code exchange implementation"
LABEL org.opencontainers.image.licenses=MIT

EXPOSE 8080

ENTRYPOINT ["/bin/pkce"]