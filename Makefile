run:
	@cd ./cmd/teleport && go run main.go

test:
	@gotest ./...

tools:
	@echo "Download all the tools ... like teleport? .."

deps:
	@echo "Installing the deps .."
	@brew install ariga/tap/atlas