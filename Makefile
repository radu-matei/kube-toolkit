# Runs unit tests for the kube-toolkit examples(example/*) and the library(pkg/*)
.PHONY: test
test:
	./scripts/test

# TODO: Build the ktk command which is used to generate a project
.PHONY: build
build:
	echo building ktk <VERSION>...
