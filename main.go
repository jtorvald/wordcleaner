package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"os"
)

func main() {
	var err error
	if len(os.Args) == 4 {
		err = run(os.Args[1], os.Args[2], os.Args[3])
	} else if len(os.Args) == 3 {
		err = run("", os.Args[1], os.Args[2])
	} else {
		err = errors.New("invalid number of arguments: templateFile inputFile outputFile")
	}
	if err != nil {
		_, err = fmt.Fprint(os.Stderr, err.Error(), "\n")
		if err != nil {
			panic(err)
		}
	}
}

func run(templateFile, inputFile, outputFile string) error {

	tpl := []byte("{{ .Contents }}")

	var err error

	if templateFile != "" {
		tpl, err = os.ReadFile(templateFile)
		if err != nil {
			return err
		}
	}

	var inputBytes []byte
	if inputFile != "" {
		inputBytes, err = os.ReadFile(inputFile)
		if err != nil {
			return err
		}
	}
	p := bluemonday.NewPolicy()

	// Allow URLs to be parseable by net/url.Parse and either:
	//   mailto: http:// or https://
	p.AllowStandardURLs()

	// Allow lists <ul> <ol>
	p.AllowLists()

	// We only allow <p> and <a href="">
	p.AllowAttrs("href").OnElements("a")
	p.AllowAttrs("name").OnElements("a")
	p.AllowAttrs("class").OnElements("p", "span")
	p.AllowElements("p", "h1", "h2", "h3", "h4", "h5", "h6", "span")

	html := p.SanitizeBytes(
		inputBytes,
	)

	html = bytes.Replace(html, []byte("<p class=\"Standard\"><span> </span></p>"), []byte(""), -1)

	html = bytes.TrimSpace(html)
	tpl = bytes.Replace(tpl, []byte("{{ .Contents }}"), html, -1)

	err = os.WriteFile(outputFile, tpl, 0644)
	return err
}
