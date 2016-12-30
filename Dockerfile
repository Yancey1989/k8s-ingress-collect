FROM alpine
RUN mkdir -p /go/bin
COPY k8s-ingress-collect /go/bin
ENTRYPOINT ["/go/bin/k8s-ingress-collect"]
