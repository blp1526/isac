package keybinding

// Keybinding shows Keys and Desc.
type Keybinding struct {
	Keys string
	Desc string
}

// Keybindings returns this TUI tool's keybindings.
func Keybindings() []Keybinding {
	return []Keybinding{
		{Keys: "C-c                ", Desc: "exit"},
		{Keys: "Arrow Up, C-p      ", Desc: "move current row up"},
		{Keys: "Arrow Down, C-n    ", Desc: "move current row down"},
		{Keys: "C-u                ", Desc: "power on current row's server"},
		{Keys: "C-r                ", Desc: "refresh rows"},
		{Keys: "BackSpace, C-b, C-h", Desc: "delete a filter character"},
		{Keys: "C-a                ", Desc: "delete all filter characters"},
		{Keys: "C-s                ", Desc: "sort rows"},
		{Keys: "C-/                ", Desc: "show help"},
		{Keys: "Enter              ", Desc: "show current row's detail"},
	}
}
