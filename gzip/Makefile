.PHONY: build
build:
	mkdir -p ./out
	go build -ldflags="-w -s" -o ./out/ptxtrecov ./cmd/ptxtrecov
	go build -ldflags="-w -s" -o ./out/indcpa ./cmd/indcpa

.PHONY: test
test:
	stdout=$$(gofmt -l . 2>&1); \
	if [ "$$stdout" ]; then \
		exit 1; \
	fi
	go vet ./...
	gocyclo -over 10 $(shell find . -iname '*.go' -type f)
	staticcheck ./...
	go test -v -cover ./...

.PHONY: setup
setup:
	go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest

.PHONY: profile
profile:
	mkdir -p ./prof
	go test -tags profile -cpuprofile ./prof/cpu.prof -memprofile ./prof/mem.prof -run ^TestRun$$ github.com/eikendev/bridgefypwn/internal/commands -count=1
