unit:
	go test -v --race $(shell go list ./... | grep -v /vendor/)

build_linux:
	CGO_ENABLED=0 GOOS=linux go build -o auth cmd/auth.go

build_docker:build_linux
	docker build -t dikaeinstein/gomicroservice-auth:latest .

RSA_PRIVATE_KEY := $(RSA_PRIVATE_KEY)
DATADOG_API_KEY := $(DATADOG_API_KEY)
run_docker:
	docker run -p 8081:8081 \
	-e "DOGSTATSD=localhost:8125" \
	-e "DD_SITE=datadoghq.eu" \
	-e "DD_API_KEY=$$DATADOG_API_KEY" \
	-e "RSA_PRIVATE_KEY=$$RSA_PRIVATE_KEY" \
	dikaeinstein/gomicroservice-auth:latest

run:
	DOGSTATSD=localhost:8125 go run cmd/auth.go

staticcheck:
	staticcheck $(shell go list ./... | grep -v /vendor/)

test: unit staticcheck
