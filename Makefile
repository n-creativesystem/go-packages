.PHONY: unit-test-module
unit-test-module:
	@echo "Running tests for all Go modules under lib/"
	@find lib -name 'go.mod' | while read modfile; do \
		dir=$$(dirname "$$modfile"); \
		echo "=== Testing $$dir ==="; \
		(cd "$$dir" && go test -covermode atomic -coverprofile=covprofile ./...); \
	done
