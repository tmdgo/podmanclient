package podmanclient

import (
	"context"
	"log"

	networktype "github.com/containers/podman/v3/libpod/network/types"
	"github.com/containers/podman/v3/pkg/bindings"
	"github.com/containers/podman/v3/pkg/bindings/containers"
	"github.com/containers/podman/v3/pkg/bindings/images"
	"github.com/containers/podman/v3/pkg/specgen"
	runtimespec "github.com/opencontainers/runtime-spec/specs-go"
)

type Client struct {
	connection context.Context
}

func (client *Client) Init(sock Sock) {
	var err error
	client.connection, err = bindings.NewConnection(context.Background(), sock.Value)
	if err != nil {
		log.Panic(err)
	}
}

func (client *Client) PullImage(image string) (err error) {
	_, err = images.Pull(client.connection, "image", nil)
	return
}

func (client *Client) CreateContainer(container Container) (err error) {
	spec := specgen.NewSpecGenerator(container.Image, false)
	spec.Name = container.Name
	spec.Labels = container.Labels

	mounts := make([]runtimespec.Mount, 0)
	for _, mount := range container.Volumes {
		mounts = append(mounts, runtimespec.Mount{Source: mount.Source, Destination: mount.Destination})
	}
	spec.Mounts = mounts

	portMappings := make([]networktype.PortMapping, 0)
	for _, portmapping := range container.PortMappings {
		portMappings = append(portMappings, networktype.PortMapping{
			HostPort:      portmapping.Source,
			ContainerPort: portmapping.Destination,
			Protocol:      portmapping.Protocol,
		})
	}
	spec.PortMappings = portMappings

	_, err = containers.CreateWithSpec(client.connection, spec, nil)
	return
}

func (client *Client) StartContainer(container Container) (err error) {
	err = containers.Start(client.connection, container.Name, nil)
	return
}

func (client *Client) RunContainer(container Container) (err error) {
	err = client.CreateContainer(container)
	if err != nil {
		return
	}
	err = client.StartContainer(container)
	return
}
