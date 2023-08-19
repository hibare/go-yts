SHELL=/bin/bash
MAKEFLAGS += -s

clean: 
	go mod tidy

test: 
	go test ./... -cover

.PHONY = clean test 