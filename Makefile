.PHONY: all test clean golog

all: golog

test: golog
	./golog

golog:
	go build github.com/goj/golog

clean:
	rm -f golog
