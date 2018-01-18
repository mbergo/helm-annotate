# Helm Annotate Plugin

This is a Helm plugin to help chart developers to tag their releases

[![Build Status](https://travis-ci.org/Tradeshift/helm-annotate.svg?branch=master)](https://travis-ci.org/Tradeshift/helm-annotate)

## Usage

Tag releases with custom information, and retrieve them again.

```shell
Usage:
  annotate [command]

Available Commands:
  get         gets annotation on a release
  help        Help about any command
  set         sets annotation on a release

Flags:
  -h, --help   help for annotate

Use "annotate [command] --help" for more information about a command.
```

## Install

```
$ helm plugin install https://github.com/Tradeshift/helm-annotate
```

The above will fetch the latest binary release of `helm annotate` and install it.

### Developer (From Source) Install

If you would like to handle the build yourself, instead of fetching a binary,
this is how recommend doing it.

First, set up your environment:

- You need to have [Go](http://golang.org) installed. Make sure to set `$GOPATH`

Clone this repo into your `$GOPATH`. You can use `go get -d github.com/Tradeshift/helm-annotate`
for that.

```
$ cd $GOPATH/src/github.com/Tradeshift/helm-annotate
$ go build
$ SKIP_BIN_INSTALL=1 helm plugin install $GOPATH/src/github.com/Tradeshift/helm-annotate
```

That last command will skip fetching the binary and use the one you build.
