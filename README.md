# isac

Interactive SAKURA Cloud

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

|Name|Description|
|-|-|
|ESC, C-c|exit|
|Arrow Up, C-p|move current row up|
|Arrow Down, C-n|move current row down|
|C-u|power on current row's server|
|C-r|refresh rows.|
|BackSpace, C-h|delete a filter character|
|C-s|sort rows|
|Enter|show current row's detail|

## Options

|Name|Description|
|-|-|
|--unanonymize|unanonymize personal data|
|--zones ZONES|set ZONES (separated by ",", example: "is1a,is1b,tk1a")|
|--help, -h|show help|
|--version, -v|print the version|
