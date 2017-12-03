package server

import (
	"fmt"
)

type ServerCollection struct {
	Count   int      `json:"Count,omitempty"`
	Servers []Server `json:"Servers,omitempty"`
}

type Server struct {
	ID       string   `json:"ID,omitempty"`
	Name     string   `json:"Name,omitempty"`
	Instance Instance `json:"Instance,omitempty"`
	Zone     Zone     `json:"Zone,omitempty"`
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

func (s *Server) String(showServerID bool) string {
	id := "************"
	if showServerID {
		id = s.ID
	}

	status := fmt.Sprintf("%6v", s.Instance.Status)
	return fmt.Sprintf("%v %v %v %v", s.Zone.Name, id, status, s.Name)
}
