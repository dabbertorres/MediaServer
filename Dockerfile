FROM centos:7

COPY radio /srv
WORKDIR /srv

ENTRYPOINT ["./radio"]
