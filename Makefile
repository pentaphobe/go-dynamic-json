MODULES=$(shell go list ./...)

.PHONY: test
test: 
	go test -v $(MODULES)

.PHONY: coverage
coverage: OUTFILE:=$(shell mktemp -t "XXXXXX.out")
coverage:
	go test -v -coverprofile=${OUTFILE} $(MODULES)
	go tool cover -html=${OUTFILE}
	rm ${OUTFILE}