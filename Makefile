TAG=v1.3
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
	docker build -t docker.baifendian.com/k8s-ingress-collector:${TAG} .
clean:
	rm ./k8s-ingress-collect
