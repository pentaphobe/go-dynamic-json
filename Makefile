MODULES=$(shell go mod list ./...)

.PHONY: test
test: 
	go test -v $(MODULES)

coverage: OUTFILE:=$(shell mktemp -t "XXXXXX.out")
coverage: $(MODULES)
	go test -v -coverprofile=${OUTFILE} ./...
	go tool cover -html=${OUTFILE}
	rm ${OUTFILE}