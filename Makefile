all:
	./scripts/build.sh

lint:
	./scripts/golint.sh

publish-docs:
	./scripts/publish_docs.sh

validate-notices:
	@echo "Validating open source dependencies in NOTICES.md..."
	@echo "Extracting direct dependencies from go.mod..."
	@go list -m -json all | jq -r 'select(.Main != true and .Indirect != true) | .Path' | grep -E "^github\.com/" | sed 's|github\.com/||' | sed 's|/v[0-9][0-9]*$$||' | sort -u > /tmp/go_deps_github.txt
	@echo "Extracting documented dependencies from NOTICES.md..."
	@grep -E "https://github\.com/[^[:space:]]+" NOTICES.md | \
		sed 's/.*https:\/\/github\.com\///' | \
		sed 's/[[:space:]].*//' | \
		sort -u > /tmp/notices_deps.txt
	@echo ""
	@missing_deps=$$(comm -23 /tmp/go_deps_github.txt /tmp/notices_deps.txt); \
	extra_deps=$$(comm -13 /tmp/go_deps_github.txt /tmp/notices_deps.txt); \
	if [ -n "$$missing_deps" ]; then \
		echo "❌ GitHub dependencies missing from NOTICES.md:"; \
		echo "$$missing_deps" | sed 's/^/  - github.com\//'; \
		echo ""; \
	fi; \
	if [ -n "$$extra_deps" ]; then \
		echo "⚠️  Dependencies in NOTICES.md but not in direct go.mod dependencies:"; \
		echo "$$extra_deps" | sed 's/^/  - github.com\//'; \
		echo ""; \
	fi; \
	if [ -n "$$missing_deps" ] || [ -n "$$extra_deps" ]; then \
		echo "❌ Validation failed: Dependency mismatches found"; \
		exit 1; \
	else \
		echo "✅ All GitHub dependencies are documented in NOTICES.md"; \
	fi
	@rm -f /tmp/go_deps_github.txt /tmp/notices_deps.txt;

validate: lint validate-notices
	@echo "All validation checks completed!"

unit-test:
	@echo "Running unit tests..."
	@go test -v -cover -coverprofile=coverage.out -covermode=atomic $$(go list ./... | grep -v -E '(examples|services|testutil)') > test_results.txt 2>&1 || true

unit-test-coverage:
	@echo "Generating coverage report..."
	@go tool cover -func=coverage.out
	@go tool cover -html=coverage.out -o coverage.html

unit-test-check:
	@echo "Checking test results..."
	@if grep -q "FAIL" test_results.txt; then \
		echo "Unit tests failed!"; \
		exit 1; \
	else \
		echo "All tests passed!"; \
	fi

unit-test-all: unit-test unit-test-coverage unit-test-check
	@echo "Full unit test suite completed successfully!"

clean:
	rm -f ark
	rm -rf bin
