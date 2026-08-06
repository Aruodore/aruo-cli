.PHONY: check build test

check: test

build:
	npx tsc

test: build
	node --test dist/
