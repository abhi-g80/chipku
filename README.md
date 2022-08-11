![build](https://github.com/abhi-g80/chipku/actions/workflows/build.yml/badge.svg)
![test](https://github.com/abhi-g80/chipku/actions/workflows/test.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/abhi-g80/chipku)](https://goreportcard.com/report/github.com/abhi-g80/chipku)

Chipku - a no frill pastebin 🗑️
==============================
Paste bin in Golang.

This is an in-memory paste bin which tries to be extremely minimal and doesn't get in the way. Simply a tool to quickly share snippets. No backend databases, no code formatting or code commenting, no self-destructing messages.

If you restart the application, you **will lose** your snippets.


Installation
------------

You may download the standalone binary and run it as,

    ./chipku

Or you may download the project and build (or run) from source,

    go run .
    # or
    go build -o chipku


Usage
-----

Simply visit the homepage and paste your text that you would like to share.

Or you can simply use command-line to `PUT` your text, using `httpie`

![httpie](docs/httpie.png "httpie")

To get your snippet in command-line, set the `No-Html` HTTP header,

![httpie-get](docs/httpie-get.png "httpie-get")