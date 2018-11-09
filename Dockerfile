FROM scratch
COPY jsonfmt /usr/local/bin/jsonfmt
ENTRYPOINT ["/usr/local/bin/jsonfmt"]
