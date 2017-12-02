package server

import "fmt"

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

func NewCollection() *ServerCollection {
	s := &ServerCollection{}
	return s
}

func New() *Server {
	s := &Server{}
	return s
}

func (s *Server) String() string {
	status := fmt.Sprintf("%8v", s.Instance.Status)
	return fmt.Sprintf("%s %s %s %s", s.Zone.Name, s.ID, status, s.Name)
}