FROM API_MAIN/base:latest
MAINTAINER Kwawdwo Amo-Addai <kwadwoamoad@gmail.com>


# Add application source
COPY src /application
WORKDIR /application

# Install app dependencies
RUN npm install

# OUTPUT: Build artefacts (Node Artifact) are output here
VOLUME /nodeartifact

# OUTPUT: Build cache
VOLUME /build

# OUTPUT: Test reports are output here
VOLUME /reports

# Add test entrypoint script
COPY scripts/test.sh /usr/local/bin/test.sh
RUN chmod +x /usr/local/bin/test.sh

EXPOSE 8080

# Set defaults for entrypoint and command string
ENTRYPOINT ["test.sh"]
CMD ["npm", "test" ]


