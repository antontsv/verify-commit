.PHONY: clean build

build:
	mkdir -p build
	GOOS=darwin go build -o build/macos-verify-commit .
	GOOS=darwin GOARCH=arm64 go build -o build/macos-verify-commit-arm .
	GOOS=linux go build -o build/verify-commit .
	GOOS=linux GOARCH=arm64 go build -o build/verify-commit-arm .

clean:
	rm -rf build