OBJECT=api-tester

all: test

run: $(OBJECT)
	./$(OBJECT)	

$(OBJECT): */*.go *.go

$(OBJECT):
	go build .

test: $(OBJECT)
	go test -cover ./...

bench: $(OBJECT)
	go test -benchmem -bench=. ./...

ci:
	ci -l -m"." *.go */*.go Makefile

clean:
	rm -f $(OBJECT)
