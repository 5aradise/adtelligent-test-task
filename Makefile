build-auction:
	go build -C cmd/auction/ -o ../../bin/

build-stitching:
	go build -C cmd/stitching/ -o ../../bin/

run-auction: build-auction
	./bin/auction

run-stitching: build-stitching
	./bin/stitching