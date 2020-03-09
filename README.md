# Commit signature check

Small utility to check that HEAD commit on provided git repository was signed using known PGP key.

## Build your own

* Get your public key:
    `gpg -a --export yourname@youremail.com`
* Paste value in `cmd/key.go`
```go
//in cmd/key.go
package cmd

var publicKey = `-----BEGIN PGP PUBLIC KEY BLOCK-----
// actual key should be present here is here
-----END PGP PUBLIC KEY BLOCK-----
`
```
* Ensure dependencies: `dep ensure`
* Compile: `go build -o verify-commit .`

## Using

```bash
verify-commit sigcheck -p ~/some/git-repo
```