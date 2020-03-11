.PHONY: clean build

build:
	mkdir -p build
	GOOS=darwin go build -o build/macos-verify-commit .
	GOOS=linux go build -o build/verify-commit .

clean:
	rm -rf build