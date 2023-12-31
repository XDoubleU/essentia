init:
	go install github.com/segmentio/golines@v0.11.0
	go install github.com/daixiang0/gci@v0.11.2
	go install github.com/securego/gosec/v2/cmd/gosec@v2.17.0

lint:
	golangci-lint run

lint/fix:
	golines . -m 88 -w
	golangci-lint run --fix
	gci write --skip-generated -s standard -s default -s "prefix(check-in/api)" .

test:
	go test ./cmd/api

test/v:
	go test -v ./cmd/api

test/cov/report:
	go test -covermode=set -coverprofile=coverage.out ./cmd/api -coverpkg ./cmd/api,./internal/...

test/cov:
	go test -covermode=set -coverprofile=coverage.out ./cmd/api -coverpkg ./cmd/api,./internal/...
	go tool cover -html=coverage.out -o=coverage.html
	make test/cov/open

test/cov/open:
	CMD /C start chrome /new-tab %CD%/coverage.html
