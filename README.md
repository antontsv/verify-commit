# Commit signature check

Small utility to check that HEAD commit on provided git repository was signed using known PGP key.

## Build your own

* Get your public key:
    `gpg -a --export yourname@youremail.com`
* Paste value in `key.go`
* Compile: `go build -o verify-commit .`

## Using

```bash
verify-commit -dir ~/some/git-repo
```