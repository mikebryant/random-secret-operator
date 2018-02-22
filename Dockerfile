FROM scratch

ADD random-secret-operator /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/random-secret-operator"]
