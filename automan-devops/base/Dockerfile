FROM node:boron
MAINTAINER Kwawdwo Amo-Addai <kwadwoamoad@gmail.com>

# Prevent dpkg errors
ENV TERM=xterm-256color

# Set mirrors to NZ
# RUN sed -i "s/http:\/\/archive./http:\/\/nz.archive./g" /etc/apt/sources.list

# Install node.js
RUN apt-get update && \
    apt-get install curl -y && \
    curl -sL https://deb.nodesource.com/setup_4.x | sudo -E bash - && \
    apt-get install -y nodejs

# Install Python runtime
RUN apt-get update && \
    apt-get install -qy \
    -o APT::Install-Recommend=false -o APT::Install-Suggests=false \
    python python-virtualenv libpython2.7 python-mysqldb

# Create virtual environment
# Upgrade PIP in virtual environment to latest version
RUN virtualenv /appenv && \
    . /appenv/bin/activate && \
    pip install pip --upgrade


LABEL application=API_MAIN