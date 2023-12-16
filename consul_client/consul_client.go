package consul_client

import (
	"github.com/hashicorp/consul/api"
)

type ConsulClient interface {
	Register(port int, serviceName, serviceID, ip string, tags []string) error
	Deregister()
	HealthCheck()
	GetKV()
	PutKV()
}

type ConsulClientImpl struct {
	client *api.Client
}

var _ ConsulClient = (*ConsulClientImpl)(nil)

func NewConsulClient(addr string) (*ConsulClientImpl, error) {
	cfg := api.DefaultConfig()
	cfg.Address = addr
	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &ConsulClientImpl{client: client}, nil
}

func (c *ConsulClientImpl) Register(port int, serviceName, serviceID, ip string, tags []string) error {
	srv := api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Tags:    tags,
		Port:    port,
		Address: ip,
	}

	return c.client.Agent().ServiceRegister(&srv)
	
}

func (c *ConsulClientImpl) Deregister() {
	//TODO implement me
	panic("implement me")
}

func (c *ConsulClientImpl) HealthCheck() {
	//TODO implement me
	panic("implement me")

}

func (c *ConsulClientImpl) GetKV() {
	//TODO implement me
	panic("implement me")
}

func (c *ConsulClientImpl) PutKV() {
	//TODO implement me
	panic("implement me")
}
