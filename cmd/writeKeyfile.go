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
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

// writeKeyfileCmd represents the writeKeyfile command
var writeKeyfileCmd = &cobra.Command{
	Use:   "writeKeyfile",
	Short: "command combined multiple keys into one keyring file",
	Long: `Convenience command to commbine multiple keys into one file. 
This command does same thing as "gpg --export --armor".

File can later be used with sigcheck command.`,
	Run: func(cmd *cobra.Command, args []string) {
		output := cmd.Flag("out").Value.String()

		_, err := os.Stat(output)

		combinedKeyring := make(openpgp.EntityList, 0)
		if err == nil {
			data, err := ioutil.ReadFile(output)
			if err != nil {
				log.Fatalf("Error reading exiting file %s: %v\n", output, err)
			}
			keyring, err := openpgp.ReadArmoredKeyRing(bytes.NewReader(data))
			if err != nil {
				log.Fatalf("bad PGP key in file %s: %v", output, err)
			}
			combinedKeyring = append(combinedKeyring, keyring...)
		}

		key := cmd.Flag("key")
		if key.Changed {
			data, err := ioutil.ReadFile(key.Value.String())
			if err != nil {
				log.Fatalf("Error reading public key file: %v\n", err)
			}
			publicKey = string(data)
		}

		newKeyring, err := openpgp.ReadArmoredKeyRing(strings.NewReader(publicKey))
		if err != nil {
			log.Fatalf("bad PGP key: %v", err)
		}
		combinedKeyring = append(combinedKeyring, newKeyring...)

		f, err := os.Create(output)
		defer f.Close()

		buf := bytes.NewBuffer(nil)
		headers := map[string]string{"purpose": "public keys for commit verification"}
		w, err := armor.Encode(buf, openpgp.PublicKeyType, headers)
		if err != nil {
			log.Fatalf("err: %v", err)
		}

		for _, e := range combinedKeyring {
			err = e.Serialize(w)
			if err != nil {
				log.Fatalf("err: %v", err)
			}
		}

		w.Close()

		f.Write(buf.Bytes())
		f.Sync()

		fmt.Printf("Written key file: %s\n", output)

	},
}

func init() {
	rootCmd.AddCommand(writeKeyfileCmd)

	writeKeyfileCmd.Flags().StringP("out", "o", "keys.asc", "Path to file, where key will be added to")
	writeKeyfileCmd.Flags().StringP("key", "k", "", "Path to armored PGP public key you want to add")
}
