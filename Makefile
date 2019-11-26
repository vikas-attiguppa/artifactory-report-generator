REPO			:= artifactory-report-generator
NAME			?= report-generator
BUILD_IMAGE		:= golang:1.11.2
CURRENT			:= ${CURDIR}
DEST			:= /root/$(REPO)
COVFILE		    ?= coverage.out
PUBFILE		    ?= coverage.html


ifeq ($(BUILD_ARTIFACT),)
BUILD_ARTIFACT := $(NAME)
endif

ifeq ($(RELEASE_VERSION),)
	RELEASE_VERSION := 0.0.0
endif

OS ?= linux
ifeq ($(shell uname), Darwin)
	OS := darwin
endif

PKGS=$$(go list ./... | grep -v 'vendor')

.PHONY: build docker_build docker_login docker_test test_report release publish clean fmt vet lint

default: docker_build

lint: ## Golint
	golint $(PKGS)

fmt: ## Go Fmt
	go fmt $(PKGS)

vet: ## Go Vet
	go vet -shadow=true $(PKGS)

build: fmt vet ## Build Project
	@printf ">>> building for %s\n" $(OS)
	@# Builds statically linked binary for target OS and Architecture (amd64).
	@# Uses linker flags to set main.Version to RELEASE_VERSION.
	@# Outputs BUILD_ARTIFACT
	GOOS=darwin GOARCH=amd64 go build -ldflags "-s" -o ./bin/$(BUILD_ARTIFACT)
	@printf ">>> fixing ownership (will error on osx, that's fine)\n"
	@# Fixes ownerships of files generated as a result of building in a container.
	@# This command will error on OSX hosts due to differences between command implementations.
	-find ./bin -not -uid $$(stat -c "%u" .) -exec chown --reference=. {} \;

docker_build: ## Build using docker image
	@# Runs build target in docker container.
	docker run --rm -e "CGO_ENABLED=0" -e "GOPATH=/go" -e "BUILD_ARTIFACT=$(BUILD_ARTIFACT)" -v "$(CURRENT):$(DEST)" -w "$(DEST)" $(BUILD_IMAGE) make build

docker_login: ## Login to registry
	docker login t

docker_test: ## run coverage reports
	docker run --rm -e "CGO_ENABLED=1" -e "GOPATH=/go" -v "$(CURRENT):$(DEST)" -w "$(DEST)" $(BUILD_IMAGE) \
		 make test_report COVFILE=$(COVFILE) PUBFILE=$(PUBFILE)

release: ## Build and push image
	docker build -t $(IMAGE):$(TAG) .
	docker push $(IMAGE):$(TAG)

publish: docker_login release

test_report: clean ## Run coverage report
	@printf ">>> running tests\n"
	@# Runs go test and outputs coverage profile.
	@# Generates HTML formatted coverage report.
	mkdir -p coverage
	go test -race -coverprofile=coverage/$(COVFILE) ./...
	go tool cover -html=coverage/$(COVFILE) -o coverage/$(PUBFILE)
