LAMBDA_ENTRY := ./cmd/lambda
LAMBDA_DIR   := .
ARCH        ?= arm64
BUILD_TAGS  ?= lambda.norpc

.PHONY: run build bootstrap zip package clean deploy

run_dev:
	air ./..

bootstrap:
	@echo ">> Building bootstrap for linux/$(ARCH)"
	cd $(LAMBDA_ENTRY) && \
		GOOS=linux GOARCH=$(ARCH) CGO_ENABLED=0 \
		go build -trimpath -ldflags="-s -w" -tags "$(BUILD_TAGS)" \
		-o ../../bootstrap .
	@ls -lh $(LAMBDA_DIR)/bootstrap

zip: bootstrap
	cd $(LAMBDA_DIR) && zip -9r function.zip bootstrap
	@ls -lh $(LAMBDA_DIR)/function.zip

package: zip

clean:
	@rm -f $(LAMBDA_DIR)/bootstrap $(LAMBDA_DIR)/function.zip

deploy: bootstrap
	cd ../infra && npx cdk deploy

test:
	go test ./...
