[![Build Status](https://travis-ci.org/blp1526/isac.svg?branch=travis)](https://travis-ci.org/blp1526/isac)
[![Go Report Card](https://goreportcard.com/badge/github.com/blp1526/isac)](https://goreportcard.com/report/github.com/blp1526/isac)
[![GoDoc](https://godoc.org/github.com/blp1526/isac?status.svg)](https://godoc.org/github.com/blp1526/isac)

# isac

Interactive [SAKURA Cloud](https://cloud.sakura.ad.jp/)

## Installation

Download a binary from [here](https://github.com/blp1526/isac/releases).

## Usage

1. Create your `$HOME/.usacloud/default/config.json` by below command.

```
$ isac init
```

2. Write your ACCESS TOKEN, ACCESS TOKEN SECRET and default ZONE to `$HOME/.usacloud/default/config.json`.

3. Run below command.

```
$ isac
```

![isac](https://user-images.githubusercontent.com/1040576/33887076-e12c7de8-df8b-11e7-9466-5af9b6af8904.gif)

## Keybindings

|Keys|Description|
|---|---|
{{ range .Keybindings -}}
|{{.Keys}}|{{.Desc}}|
{{end}}

## Options

|Option|Description|
|---|---|
|help, h|show help|
|version, v|print the version|
{{ range .Options -}}
|{{.Name}}|{{.Usage}}|
{{end}}
