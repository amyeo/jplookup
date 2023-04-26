build:
	mkdir -p  "bin"
	go build -o bin/jisho2 .
run:
	go run .
index:
	go run . --init