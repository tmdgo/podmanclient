package podmanclient

type Sock struct {
	Value string
}

type Container struct {
	Name         string
	Image        string
	Labels       map[string]string
	Volumes      []ContainerVolume
	PortMappings []ContainerPortMapping
}

type ContainerVolume struct {
	Source      string
	Destination string
}

type ContainerPortMapping struct {
	Source      uint16
	Destination uint16
	Protocol    string
}
