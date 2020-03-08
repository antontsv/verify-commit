package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"

	"golang.org/x/crypto/openpgp"
)

func command(ctx context.Context, dir, command string) (out, err string) {
	cmd := fmt.Sprintf("cd '%s' && git %s", dir, command)
	c := exec.CommandContext(ctx, "sh", "-c", cmd)
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	c.Stdout = writer
	c.Stderr = writer
	e := c.Run()
	writer.Flush()
	if e != nil {
		err = e.Error()
	}
	return b.String(), err
}

func main() {

	dir := flag.String("dir", "", "Directory where to check for commit signature")
	pubKeyFile := flag.String("pubKeyFile", "", "Location of armored PGP public key file you want to use. Note: this utility has a default key included.")

	mainCtx, mainCancel := context.WithCancel(context.Background())
	defer mainCancel()

	go func() {
		signl := make(chan os.Signal)
		signal.Notify(signl, os.Interrupt)
		s := <-signl
		fmt.Fprintf(os.Stderr, "Got signal: %v, canceling processes in flight...\n", s)
		mainCancel()
	}()

	flag.Parse()

	if len(*dir) <= 0 {
		log.Fatal("Please provide -dir for repository where to check for commit signature")
	}

	if len(*pubKeyFile) > 0 {
		data, err := ioutil.ReadFile(*pubKeyFile)
		if err != nil {
			log.Fatalf("Error reading public key file: %v\n", err)
		}
		publicKey = string(data)
	}

	if len(publicKey) <= 0 {
		log.Fatal("No public key was set during script compile time, and not prodived at runtime")
	}

	keyring, err := openpgp.ReadArmoredKeyRing(strings.NewReader(publicKey))
	if err != nil {
		log.Fatalf("bad PGP key: %v", err)
	}

	/*
		// Alternative using go-git
		r, err := git.PlainOpen(*dir)
		if err != nil {
			log.Fatalf("Cannot open repo at dir '%s': %v", *dir, err)
		}

		ref, err := r.Head()
		if err != nil {
			log.Fatalf("Cannot read HEAD information: %v", err)
		}

		commit, err := r.CommitObject(ref.Hash())
		if err != nil {
			log.Fatalf("Cannot get commit information: %v", err)
		}
		entity, err := commit.Verify(publicKey)

		fmt.Println(commit.Message)

		if err != nil {
			log.Fatalf("Cannot run PGP verify: %v", err)
		}

		for _, v := range entity.Identities {
			if v.UserId != nil {
				fmt.Printf("Verified signature by %s\n", v.UserId.Name)
			}
		}
	*/
	o, e := command(mainCtx, *dir, "cat-file commit HEAD")
	if e != "" {
		log.Fatalf("Cannot open repo at dir '%s': %s", *dir, e)
	}
	const (
		prefix      = "gpgsig "
		startMarker = "-----BEGIN PGP SIGNATURE-----\n"
		endMarker   = "-----END PGP SIGNATURE-----\n"
	)
	sigStart := strings.Index(o, fmt.Sprintf("%s%s", prefix, startMarker))
	sigEnd := strings.Index(o, endMarker)
	if sigStart < 0 || sigEnd <= sigStart {
		log.Fatalf("HEAD commit in %s does not have valid signature\n", *dir)
	}

	body := o[0:sigStart] + o[sigEnd+len(endMarker):]

	sigStart = sigStart + len(prefix)
	signature :=
		o[sigStart:sigStart+len(startMarker)] +
			// Here is actual signature content, but
			// we need to remove spaces (git cat-file has whitespace at the beninning of each line)
			strings.ReplaceAll(o[sigStart+len(startMarker):sigEnd], " ", "") +
			o[sigEnd:sigEnd+len(endMarker)]

	entity, err := openpgp.CheckArmoredDetachedSignature(
		keyring,
		strings.NewReader(body),
		strings.NewReader(signature),
	)

	if err != nil {
		log.Fatalf("Cannot run PGP verify: %v", err)
	}
	for _, v := range entity.Identities {
		if v.UserId != nil {
			fmt.Printf("HEAD commit in %s has verified signature by %s\n", *dir, v.UserId.Name)
		}
	}

}
