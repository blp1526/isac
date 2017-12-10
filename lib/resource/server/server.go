package server

import (
	"fmt"
)

type ServerCollection struct {
	Count   int      `json:"Count,omitempty"`
	Servers []Server `json:"Servers,omitempty"`
}

type Server struct {
	Availability string   `json:"Availability,omitempty"`
	Description  string   `json:"Description,omitempty"`
	ID           string   `json:"ID,omitempty"`
	Instance     Instance `json:"Instance,omitempty"`
	Name         string   `json:"Name,omitempty"`
	ServiceClass string   `json:"ServiceClass,omitempty"`
	Zone         Zone     `json:"Zone,omitempty"`
}

type Instance struct {
	Status string `json:"Status,omitempty"`
}

type Zone struct {
	Name string `json:"Name,omitempty"`
}

func NewCollection(zone string) *ServerCollection {
	s := &ServerCollection{}
	return s
}

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
