package consul_client

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"log"
)

type ConsulClient interface {
	Register(port, checkPoint int, serviceName, serviceID, ip string, tags []string) error
	Deregister(serviceID string) error
	GetKV()
	PutKV()
	GetService()
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

	//TODO:这里的地址应该换成 ip,但是不清楚为什么我本地启动这个服务，出口ip不能访问到
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

func (c *ConsulClientImpl) GetKV() {
	//TODO implement me
	panic("implement me")
}

func (c *ConsulClientImpl) PutKV() {
	//TODO implement me
	panic("implement me")
}

func (c *ConsulClientImpl) GetService() {
	panic("implement me")
}
