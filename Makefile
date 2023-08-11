BIN := $(subst cmd/,,$(wildcard cmd/*))

.PHONY: all clean test run build install upgrade help

all: $(SUBDIR) 		# default action
	@[ -f .git/hooks/pre-commit ] || pre-commit install --install-hooks
	@git config commit.template .git-commit-template

clean: $(SUBDIR)	# clean-up environment
	@find . -name '*.sw[po]' -delete
	rm -f $(BIN)

test:				# run test
	gofmt -s -w .
	go test -v ./...

run:				# run in the local environment
	go run cmd/$(BIN)/main.go

build: $(BIN)		# build the binary/library
	@go mod tidy

install: test		# install the binary tool
	go install ./...

upgrade:			# upgrade all the necessary packages
	pre-commit autoupdate

help:				# show this message
	@printf "Usage: make [OPTION]\n"
	@printf "\n"
	@perl -nle 'print $$& if m{^[\w-]+:.*?#.*$$}' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?#"} {printf "    %-18s %s\n", $$1, $$2}'

%: cmd/%/main.go
	go build -ldflags "-w -s" -o $@ $<
