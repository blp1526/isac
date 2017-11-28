package resource

type ServerCollection struct {
	Count   int      `json:"Count,omitempty"`
	Servers []Server `json:"Servers,omitempty"`
}

type Server struct {
	ID   string `json:"ID,omitempty"`
	Name string `json:"Name,omitempty"`
}
