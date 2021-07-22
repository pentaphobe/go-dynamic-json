VERSION=$(shell cat VERSION)
MODULES=$(shell go list ./...)

all: build

build: $(MODULES)
	go build $(MODULES)

.PHONY: update_version_tag
update_version_tag: DEFAULT_BRANCH:=$(shell git remote show origin | awk '/HEAD branch/ {print $$NF}')
update_version_tag: CURRENT_BRANCH:=$(shell git rev-parse --abbrev-ref HEAD)
update_version_tag: 		
	@if ! [ "$(CURRENT_BRANCH)" == "$(DEFAULT_BRANCH)" ]; then \
		echo Not on default branch; \
		false; \
	fi
	@echo Updating ${VERSION} tag
	git push origin :refs/tags/${VERSION}
	git tag -f ${VERSION}
	git push origin --tags

.PHONY: test
test: 
	go test -v $(MODULES)

.PHONY: coverage
coverage: OUTFILE:=$(shell mktemp -t "XXXXXX.out")
coverage:
	go test -v -coverprofile=${OUTFILE} $(MODULES)
	go tool cover -html=${OUTFILE}
	rm ${OUTFILE}