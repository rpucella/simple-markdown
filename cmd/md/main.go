package main

import (
	"github.com/russross/blackfriday/v2"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

/*************************************************************
 *
 * A simple wrapper around the blackfriday markdown library
 *
 */

var rep *log.Logger = log.New(os.Stdout, "" /* log.Ldate| */, log.Ltime)

func main() {
	args := os.Args[1:]

	if len(args) == 1 {
		ProcessMarkdown(args[0], "")
	} else if len(args) == 2 {
		ProcessMarkdown(args[0], args[1])
	} else {
		Usage()
	}
}

func Usage() {
	rep.Println("USAGE: md input.md [template.html]")
}

func ProcessMarkdown(inputMD string, inputTpl string) {
	w := os.Stdout
	md, err := ioutil.ReadFile(inputMD)
	if err != nil {
		rep.Fatalf("ERROR: %s\n", err)
	}
	// Do nothing with the metadata for now.
	_, restmd, err := ExtractMetadata(md)
	if err != nil {
		rep.Fatalf("ERROR: %s\n", err)
	}
	output := blackfriday.Run(restmd, blackfriday.WithNoExtensions())
	if inputTpl != "" {
		result, err := ProcessMarkdownTemplate(inputTpl, template.HTML(output))
		if err != nil {
			rep.Fatalf("ERROR: %s\n", err)
		}
		output = []byte(result)
	}
	if _, err := w.Write(output); err != nil {
		rep.Fatalf("ERROR: %s\n", err)
	}
}

func ExtractMetadata(md []byte) (map[string]string, []byte, error) {
	lines := strings.Split(string(md), "\n")
	foundMetadata := false
	metadata := make(map[string]string)
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			if line == "---" {
				if foundMetadata {
					// We're done.
					rest := []byte(strings.Join(lines[idx+1:], "\n"))
					return metadata, rest, nil
				}
				foundMetadata = true
			} else if foundMetadata {
				fields := strings.Split(line, ":")
				if len(fields) == 2 {
					fieldname := strings.TrimSpace(fields[0])
					fieldvalue := strings.TrimSpace(fields[1])
					metadata[fieldname] = fieldvalue
				}
			}
		}
	}
	return metadata, md, nil
}

func ProcessMarkdownTemplate(inputTpl string, content template.HTML) (template.HTML, error) {
	tpl, err := template.ParseFiles(inputTpl)
	if err != nil {
		return template.HTML(""), err
	}
	var b strings.Builder
	if err := tpl.Execute(&b, content); err != nil {
		return template.HTML(""), err
	}
	result := template.HTML(b.String())
	return result, nil
}
