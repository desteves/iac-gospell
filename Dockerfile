FROM --platform=linux/amd64 golang
WORKDIR /usr/src/app
COPY gospell .
ENTRYPOINT ["gospell"]