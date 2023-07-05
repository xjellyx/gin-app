FROM ubuntu:22.04

ENV TZ="Asia/Shanghai" \
    LANG="en_US.utf8" \
    LC_ALL="en_US.utf8"
RUN echo 'LANG="en_US.utf8"' > /etc/locale.conf
WORKDIR /app
COPY bin/gin-app /app/server