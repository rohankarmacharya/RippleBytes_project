package journal

import "github.com/rohankarmacharya/TigIntegration/pkg/client"

type Service struct {
	client *client.TiggClient
}

func NewService(c *client.TiggClient) *Service {
	return &Service{client: c}
}
