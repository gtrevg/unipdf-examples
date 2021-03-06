/*
 * Outputs protection information about locked PDFs.
 *
 * Run as: go run pdf_security_info.go input.pdf
 */

package main

import (
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v3/common/license"
	pdf "github.com/unidoc/unipdf/v3/model"
)

const licenseKey = `
-----BEGIN UNIDOC LICENSE KEY-----
Free trial license keys are available at: https://unidoc.io/
-----END UNIDOC LICENSE KEY-----
`

func init() {
	// Enable debug-level logging.
	// unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))

	err := license.SetLicenseKey(licenseKey, `Company Name`)
	if err != nil {
		panic(err)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run pdf_security_info.go input.pdf\n")
		os.Exit(0)
	}

	for _, inputPath := range os.Args[1:] {
		err := printSecurityInfo(inputPath)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}

}

func printSecurityInfo(inputPath string) error {
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return err
	}

	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return err
	}

	fmt.Printf("Input file %s\n", inputPath)
	if !isEncrypted {
		fmt.Printf(" - is not encrypted\n")
		return nil
	}

	// Try decrypting both with given password and an empty one if that fails.
	if isEncrypted {
		auth, err := pdfReader.Decrypt([]byte(""))
		if err != nil {
			return err
		}
		if !auth {
			fmt.Printf(" - has an opening password\n")
		}
	}

	method := pdfReader.GetEncryptionMethod()
	fmt.Printf(" - Method: %s\n", method)

	return nil
}
