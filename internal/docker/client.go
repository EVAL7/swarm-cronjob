package docker

import (
	"bytes"
	"context"
	"strings"

	"github.com/crazy-max/swarm-cronjob/internal/model"
	"github.com/docker/cli/cli/command"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

// Client for Swarm
type Client interface {
	DistributionInspect(ctx context.Context, image, encodedAuth string) (registry.DistributionInspect, error)
	RetrieveAuthTokenFromImage(ctx context.Context, image string) (string, error)
	ServiceUpdate(ctx context.Context, serviceID string, version swarm.Version, service swarm.ServiceSpec, options types.ServiceUpdateOptions) (types.ServiceUpdateResponse, error)
	ServiceInspectWithRaw(ctx context.Context, serviceID string, opts types.ServiceInspectOptions) (swarm.Service, []byte, error)
	Events(ctx context.Context, options types.EventsOptions) (<-chan events.Message, <-chan error)

	ServiceList(args *model.ServiceListArgs) ([]*model.ServiceInfo, error)
	Service(name string) (*model.ServiceInfo, error)
	TaskList(service string) ([]*model.TaskInfo, error)
	TaskLogs(ctx context.Context, taskid string) (string, error)
}

type dockerClient struct {
	api *client.Client
	cli command.Cli
}

// NewEnvClient initializes a new Docker API client based on environment variables
func NewEnvClient() (*dockerClient, error) {
	dockerApiCli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.12"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize Docker API client")
	}

	_, err = dockerApiCli.ServerVersion(context.Background())
	if err != nil {
		return nil, err
	}

	dockerCli, err := command.NewDockerCli()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Docker cli")
	}

	return &dockerClient{
		api: dockerApiCli,
		cli: dockerCli,
	}, err
}

// DistributionInspect returns the image digest with full Manifest
func (c *dockerClient) DistributionInspect(ctx context.Context, image, encodedAuth string) (registry.DistributionInspect, error) {
	return c.api.DistributionInspect(ctx, image, encodedAuth)
}

// RetrieveAuthTokenFromImage retrieves an encoded auth token given a complete image
func (c *dockerClient) RetrieveAuthTokenFromImage(ctx context.Context, image string) (string, error) {
	return command.RetrieveAuthTokenFromImage(ctx, c.cli, image)
}

// ServiceUpdate updates a Service. The version number is required to avoid conflicting writes.
// It should be the value as set *before* the update. You can find this value in the Meta field
// of swarm.Service, which can be found using ServiceInspectWithRaw.
func (c *dockerClient) ServiceUpdate(ctx context.Context, serviceID string, version swarm.Version, service swarm.ServiceSpec, options types.ServiceUpdateOptions) (types.ServiceUpdateResponse, error) {
	return c.api.ServiceUpdate(ctx, serviceID, version, service, options)
}

// ServiceInspectWithRaw returns the service information and the raw data.
func (c *dockerClient) ServiceInspectWithRaw(ctx context.Context, serviceID string, opts types.ServiceInspectOptions) (swarm.Service, []byte, error) {
	return c.api.ServiceInspectWithRaw(ctx, serviceID, opts)
}

// Events returns a stream of events in the daemon. It's up to the caller to close the stream
// by cancelling the context. Once the stream has been completely read an io.EOF error will
// be sent over the error channel. If an error is sent all processing will be stopped. It's up
// to the caller to reopen the stream in the event of an error by reinvoking this method.
func (c *dockerClient) Events(ctx context.Context, options types.EventsOptions) (<-chan events.Message, <-chan error) {
	return c.api.Events(ctx, options)
}

func normalizeImage(image string) string {
	if i := strings.Index(image, "@sha256:"); i > 0 {
		image = image[:i]
	}
	return image
}

func (c *dockerClient) TaskLogs(ctx context.Context, taskid string) (string, error) {
	r, err := c.api.TaskLogs(ctx, taskid, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Details:    true,
	})

	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(r)
	if err != nil {
		return "", err
	}
	return buf.String(), nil

}
