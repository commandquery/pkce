# scratch is an empty image
# FROM scratch
# If you need /bin/sh and a few utilities, uncomment
# the following line. It increases the image by 5.5 MB
FROM alpine:latest

RUN apk add gcompat

COPY login /bin/login
# copy other files if needed

LABEL org.opencontainers.image.source=https://github.com/bookworkhq/login
LABEL org.opencontainers.image.description="Bookwork Login"
# LABEL org.opencontainers.image.licenses=MIT

EXPOSE 8080

ENTRYPOINT ["/bin/login"]
