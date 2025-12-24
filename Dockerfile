FROM scratch
ARG TARGETPLATFORM
ENTRYPOINT ["/usr/local/bin/jsonfmt"]
COPY $TARGETPLATFORM/jsonfmt /usr/local/bin/jsonfmt
