.PHONY: list %

list:
	@echo "Targets:"
	@ls -1 examples | sed 's/.go$$//' | sed 's/^/  make run name=/'

run:
	@: $(if $(name),,$(error Usage: make run name=<file-without-.go>))
	@echo "Running examples/$(name).go"
	@go run ./examples/$(name).go

# Convenience targets
gowork:
	@go env GOPATH && go env GOROOT && go version
