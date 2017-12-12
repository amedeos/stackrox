# Set these variables only if not set by the including Makefile.
IMAGE ?= $(ROX_PROJECT)
IMAGE_TAG ?= latest
PROJECT_SUBDIR ?= $(ROX_PROJECT)
BINARY ?= $(ROX_PROJECT)
BASE_PATH ?= $(CURDIR)/..
GO_BASE_PATH ?= /go/src/bitbucket.org/stack-rox/apollo

.DEFAULT_GOAL = all

# VERSION_FILE file contains the version string of the platform
# XXX/e: Find a better way to determine this location.
VERSION_FILE ?= $(BASE_PATH)"/ROX_VERSION"
VERSION ?= `cat $(VERSION_FILE) | xargs echo -n`

.PHONY: version
version:
	echo "$(VERSION)"

###########
## Clean ##
###########
.PHONY: clean
clean: preclean clean-common postclean

.PHONY: preclean
preclean:

.PHONY: postclean
postclean:

.PHONY: clean-common
clean-common:
	@echo "+ $@"
	@go clean -i
	@rm -rf $(GOPATH)/src/github.com/grpc-ecosystem
	@rm -rf $(GOPATH)/src/github.com/golang/protobuf
	@rm -f $(GOPATH)/bin/protoc-gen-grpc-gateway
	@rm -f $(GOPATH)/bin/protoc-gen-go
	@rm -rf $(PROTOC_TMP)
	@rm -f $(PROTOC_FILE)
	@test -n "$(GENERATED_API_PATH)" && rm -rf "$(GENERATED_API_PATH)" || true


###########################
## Developer Environment ##
###########################
.PHONY: dev
dev:
	@echo "+ $@"
	@go get -u github.com/golang/lint/golint
	@go get -u golang.org/x/tools/cmd/goimports
	@go get -u github.com/jstemmer/go-junit-report
	@curl https://glide.sh/get | sh

.PHONY: docs
docs: generated-srcs
	@echo "+ $@"
	@echo
	@echo 'Access your docs at http://localhost:6061/pkg/bitbucket.org/stack-rox/apollo/$(ROX_PROJECT)/'
	@echo 'Hit CTRL-C to quit.'
	@godoc -http=:6061


###########
## Style ##
###########
.PHONY: style
style: fmt imports lint vet

.PHONY: fmt
fmt:
	@echo "+ $@"
ifdef CIRCLECI
		@echo "The environment indicates we are in CI; checking gofmt."
		@$(eval FMT=`find . -name vendor -prune -name generated -prune -o -name '*.go' -print | xargs gofmt -l`)
		@echo "gofmt problems in the following files, if any:"
		@echo $(FMT)
		@test -z "$(FMT)"
endif
	@find . -name vendor -prune -name generated -prune -o -name '*.go' -print | xargs gofmt -l -w

.PHONY: imports
imports:
	@echo "+ $@"
ifdef CIRCLECI
		@echo "The environment indicates we are in CI; checking goimports."
		@$(eval IMPORTS=`find . -name vendor -prune -name generated -prune -o -name '*.go' -print | xargs goimports -l`)
		@echo "goimports problems in the following files, if any:"
		@echo $(IMPORTS)
		@test -z "$(IMPORTS)"
endif
	@find . -name vendor -prune -name generated -prune -o -name '*.go' -print | xargs goimports -w

.PHONY: lint
lint:
	@echo "+ $@"
	@set -e; for pkg in $(shell go list -e ./... | grep -v generated | grep -v vendor); do golint -set_exit_status $$pkg; done

.PHONY: vet
vet:
	@echo "+ $@"
	@go vet $(shell go list -e ./... | grep -v generated | grep -v vendor)


#######################
## Local Compilation ##
#######################

generated-srcs: $(GENERATED_SRCS)

.PHONY: build
build: generated-srcs
	bazel run //:gazelle
	bazel build --cpu=k8 \
		//$(PROJECT_SUBDIR)


#############
## Testing ##
#############
.PHONY: test
test: pretest test-common posttest

.PHONY: pretest
pretest:

.PHONY: posttest
posttest:

.PHONY: test-common
test-common:
	@echo "+ $@"
	@go test -cover $(TESTFLAGS) -v $(shell go list -e ./... | grep -v generated | grep -v integration-tests | grep -v vendor) 2>&1 | tee test.log

.PHONY: test-integration
test-integration:
	@echo "+ $@"
	@go test -cover -tags integration -v $(shell go list -e ./... | grep -v generated | grep -v integration-tests | grep -v vendor) 2>&1 | tee test.log

.PHONY: test-all
test-all: test-integration

.PHONY: report
report:
	@echo "+ $@"
	@cat test.log | go-junit-report > report.xml
	@echo
	@echo "Test coverage summary:"
	@grep "^coverage: " -A1 test.log | grep -v -e '--' | paste -d " "  - -
	@echo
	@echo "Test pass/fail summary:"
	@grep failures report.xml
	@echo
	@echo "`grep 'FAIL	bitbucket.org/stack-rox/apollo' test.log | wc -l` package(s) detected with compilation or test failures."
	@-grep 'FAIL	bitbucket.org/stack-rox/apollo' test.log || true
	@echo
	@testerror="$$(grep -e 'can.t load package' -e '^# bitbucket.org/stack-rox/apollo/' -e 'FAIL	bitbucket.org/stack-rox/apollo' test.log | wc -l)" && test $$testerror -eq 0
