/*
Copyright © 2020 Anton Tsviatkou

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// sigcheckCmd represents the sigcheck command
var sigcheckCmd = &cobra.Command{
	Use:   "sigcheck",
	Short: "Verifies PGP signature on commits",
	Long: `Sometimes important information is cloned using git. 
	
This command will verify signature on the commit before you decide to use files from that commit`,
	Run: func(cmd *cobra.Command, args []string) {
		key := cmd.Flag("key")
		reference := cmd.Flag("ref").Value.String()
		if key.Changed {
			fmt.Println("here is the key:", key.Value)
		}

		dir := cmd.Flag("path").Value.String()

		r, err := git.PlainOpen(dir)
		if err != nil {
			log.Fatalf("Cannot open repo at dir '%s': %v", dir, err)
		}

		hash, err := r.ResolveRevision(plumbing.Revision(reference))
		if err != nil {
			log.Fatalf("Cannot read %s information: %v\n", reference, err)
		}

		commit, err := r.CommitObject(*hash)
		if err != nil {
			log.Fatalf("Cannot get commit information: %v", err)
		}
		entity, err := commit.Verify(publicKey)

		if err != nil {
			log.Fatalf("Cannot run PGP verify: %v", err)
		}

		for _, v := range entity.Identities {
			if v.UserId != nil {
				fmt.Printf("%s commit in %s has verified signature by %s\n", commit.Hash, dir, v.UserId.Name)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(sigcheckCmd)

	sigcheckCmd.Flags().StringP("path", "p", ".", "Path to git repository where to check for commit")
	sigcheckCmd.Flags().StringP("ref", "r", "HEAD", "Commit sha/reference to check for the signature")
	sigcheckCmd.Flags().StringP("key", "k", "", "Path to armored PGP public key file to use for the verification")
}
