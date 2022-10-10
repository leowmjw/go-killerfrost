run:
	@echo "Schema check .."
	@atlas schema inspect -u "postgres://foo:password@localhost:5432/myterraform?sslmode=disable"

apply:
	@eco "Applying .."
	@atlas schema apply -s s2 -u "postgres://foo:password@localhost:5432/myterraform?sslmode=disable" -f ./import_partition

tools:
	@echo "Download all the tools ... like Atlas? .."

deps:
	@echo "Installing the deps .."
	@brew install ariga/tap/atlas
