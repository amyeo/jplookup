build:
	mkdir -p  "bin"
	go build -o bin/jplookup .
run:
	go run .
index:
	go run . --init