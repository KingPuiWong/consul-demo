package consul

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/consul/api"
)

type ConsulClient interface {
	Register(port, checkPoint int, serviceName, serviceID, ip string, tags []string) error
	Deregister(serviceID string) error
	DiscoverService(serviceID string) (*api.AgentService, error)
	GetKV(key string) (string, error)
	PutKV(key string, value any) error
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

func (c *ConsulClientImpl) Register(port, checkPoint int, serviceName, serviceID, ip string, tags []string) error {
	srv := api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Tags:    tags,
		Port:    port,
		Address: ip,
	}

	checkHttpRoute := fmt.Sprintf("http://%s:%d%s", "127.0.0.1", checkPoint, "/check")
	log.Println("checkHttpRoute:", checkHttpRoute)
	srv.Check = &api.AgentServiceCheck{
		HTTP: checkHttpRoute,
		//Status:                         "passing",
		Timeout:                        "3s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "30s",
		// GRPC:     fmt.Sprintf("%v:%v/%v", IP, r.Port, r.Service),// grpc 支持，执行健康检查的地址，service 会传到 Health.Check 函数中
	}

	return c.client.Agent().ServiceRegister(&srv)
}

func (c *ConsulClientImpl) Deregister(serviceID string) error {
	return c.client.Agent().ServiceDeregister(serviceID)
}

func (c *ConsulClientImpl) DiscoverService(serviceID string) (*api.AgentService, error) {
	service, _, err := c.client.Agent().Service(serviceID, nil)
	if err != nil {
		return nil, err
	}
	return service, nil
}

func (c *ConsulClientImpl) GetKV(key string) (string, error) {
	kv := c.client.KV()
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return string(pair.Value), nil
}

func (c *ConsulClientImpl) PutKV(key string, value any) error {
	valueJson, err := json.Marshal(value)
	if err != nil {
		return err
	}

	kv := c.client.KV()
	p := api.KVPair{Key: key, Value: valueJson}
	if _, err := kv.Put(&p, nil); err != nil {
		return err
	}

	return nil
}
