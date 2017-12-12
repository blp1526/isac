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
|ESC, C-c|Exit.|
|Arrow Up, C-p|Move current row up.|
|Arrow Down, C-n|Move current row down.|
|C-u|Power on current row's server.|
|C-r|Refresh rows.|
|BackSpace, C-h|Delete a filter character.|
|C-s|Sort rows.|
|Enter|Show current row's detail.|
