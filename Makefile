### dev tools
devtools:
	@echo "Installing devtools"
	go install golang.org/x/tools/cmd/goimports@latest
	go install mvdan.cc/gofumpt@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install go.uber.org/mock/mockgen@latest
	go install github.com/bufbuild/buf/cmd/buf@latest

### mock
mock:
	mockgen -source=./pkg/client/interface.go      -destination=./pkg/client/mock.go      -package=client
	mockgen -source=./pkg/wallet/interface.go      -destination=./pkg/wallet/mock.go      -package=wallet
	mockgen -source=./pkg/mailer/interface.go      -destination=./pkg/mailer/mock.go      -package=mailer
	mockgen -source=./pkg/nowpayments/interface.go -destination=./pkg/nowpayments/mock.go -package=nowpayments

### proto file generate
proto:
	rm -rf pkg/proto/gen
	cd pkg/proto && buf generate --template buf.gen.yaml ../proto

### Formatting, linting, and vetting
fmt:
	gofumpt -l -w .
	go mod tidy

check:
	golangci-lint run --timeout=20m0s

### Testing
test:
	go test ./... -covermode=atomic

### building
release:
	go build -ldflags "-s -w" -trimpath -o build/pagu ./internal/cmd

build:
	go build -o build/pagu ./internal/cmd

### Generating commands
gen:
	go run ./internal/generator/main.go \
		"./internal/engine/command/crowdfund/crowdfund.yml" \
		"./internal/engine/command/voucher/voucher.yml" \
		"./internal/engine/command/market/market.yml" \
		"./internal/engine/command/calculator/calculator.yml" \
		"./internal/engine/command/network/network.yml" \
		"./internal/engine/command/phoenix/phoenix.yml" \

	find . -name "*.gen.go" -exec gofumpt -l -w {} +

###
.PHONY: devtools mock proto fmt check test release build  gen
