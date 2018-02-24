package server

import (
	"fmt"
)

// Collection shows Server's collection.
type Collection struct {
	Count   int      `json:"Count,omitempty"`
	Servers []Server `json:"Servers,omitempty"`
}

// Server shows GET /server/:id attributes.
type Server struct {
	Availability    string   `json:"Availability,omitempty"`
	CreatedAt       string   `json:"CreatedAt,omitempty"`
	Description     string   `json:"Description,omitempty"`
	ID              string   `json:"ID,omitempty"`
	Instance        Instance `json:"Instance,omitempty"`
	InterfaceDriver string   `json:"InterfaceDriver,omitempty"`
	ModifiedAt      string   `json:"ModifiedAt,omitempty"`
	Name            string   `json:"Name,omitempty"`
	ServiceClass    string   `json:"ServiceClass,omitempty"`
	Tags            []string `json:"Tags,omitempty"`
	Zone            Zone     `json:"Zone,omitempty"`
}

// Instance shows Server.Instance.
type Instance struct {
	Status string `json:"Status,omitempty"`
}

// Zone shows Server.Zone.
type Zone struct {
	Name string `json:"Name,omitempty"`
}

// NewCollection initializes *Collection.
func NewCollection(zone string) *Collection {
	s := &Collection{}
	return s
}

// New initializes *Server.
func New() *Server {
	s := &Server{}
	return s
}

func (s *Server) String(unanonymize bool) string {
	id := "************"
	if unanonymize {
		id = s.ID
	}

	status := fmt.Sprintf("%6v", s.Instance.Status)
	return fmt.Sprintf("%v %v %v %v", s.Zone.Name, id, status, s.Name)
}
