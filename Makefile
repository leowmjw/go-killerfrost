run:
	@cd ./cmd/teleport && go run main.go

test:
	@gotest ./...

tools:
	@echo "Download all the tools ... like teleport? .."
