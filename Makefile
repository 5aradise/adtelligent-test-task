build:
	go build -C cmd/app/ -o ../../bin/

run: build
	./bin/app
