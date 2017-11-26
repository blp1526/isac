package server

type Server struct {
	ID     string
	Zone   string
	Status string
	Name   string
}

func All(zones []string) (servers []Server, err error) {
	for _, zone := range zones {
		zonedServers, err := zonedAll(zone)
		if err != nil {
			return servers, err
		}

		for _, zonedServer := range zonedServers {
			servers = append(servers, zonedServer)
		}
	}

	return servers, err
}

func zonedAll(zone string) (servers []Server, err error) {
	return servers, err
}
