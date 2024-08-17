tools: tools/lint

tools/lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1
	go install github.com/segmentio/golines@v0.12.2
	go install github.com/daixiang0/gci@v0.13.4
	go install github.com/securego/gosec/v2/cmd/gosec@v2.20.0

lint: tools/lint
	golangci-lint run

lint/fix: tools/lint
	golines . -m 88 -w
	golangci-lint run --fix
	gci write --skip-generated -s standard -s default .

test:
	go test ./...

test/v:
	go test ./... -v 

test/race:
	go test ./... -race

test/pprof:
	go test ./... -cpuprofile cpu.prof -memprofile mem.prof -bench .

test/cov/report:
	go test ./... -coverpkg=./... -covermode=atomic -coverprofile=coverage.out -race

test/cov: test/cov/report
	go tool cover -html=coverage.out -o=coverage.html
	make test/cov/open

test/cov/open:
	open ./coverage.html
