# simple-markdown

A bare bones CLI wrapper around the [Blackfriday](https://github.com/russross/blackfriday) markdown library.

## To build

Clone and run `make`. The binary `md` can be found in `bin`.

## To use

Run `bin/md` with a Markdown file as a first argument and an optional HTML template as a second argument. The resulting HTML output will be sent to standard out.

A template is an HTML file with an `{{.}}` in it that will get replaced by the output of the Markdown processor.

