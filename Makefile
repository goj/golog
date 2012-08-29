.PHONY: all clean golog

all: golog

golog:
	go build github.com/goj/golog

clean:
	rm -f golog
