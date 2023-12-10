install:
	cd src/cli && go build -o /usr/local/bin/dacrane main.go
install-plugin:
	docker build -t docker:latest -f src/plugin/docker/Dockerfile .
	docker build -t local:latest -f src/plugin/local/Dockerfile .
	docker build -t terraform:latest -f src/plugin/terraform/Dockerfile .
uninstall:
	rm /usr/local/bin/dacrane
