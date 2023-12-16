package main

import (
	"consul-demo/consul_client"
	"net"
)

func main() {
	/*
		// Get a new client
		client, err := capi.NewClient(capi.DefaultConfig())
		if err != nil {
			panic(err)
		}

		// Get a handle to the KV API
		kv := client.KV()

		// PUT a new KV pair
		p := &capi.KVPair{Key: "REDIS_MAXCLIENTS", Value: []byte("1000")}
		_, err = kv.Put(p, nil)
		if err != nil {
			panic(err)
		}

		// Lookup the pair
		pair, _, err := kv.Get("REDIS_MAXCLIENTS", nil)
		if err != nil {
			panic(err)
		}
		fmt.Printf("KV: %v %s\n", pair.Key, pair.Value)
	*/
	consulClient, err := consul_client.NewConsulClient("127.0.0.1:8500")
	if err != nil {
		return
	}

	ip, err := getOutIp()
	if err != nil {
		panic(err)
	}

	port := 20001
	serviceName := "hello-consul"
	serviceNameID := "hello-consul-20001"
	tags := []string{"consul-demo"}
	if err := consulClient.Register(port, serviceName, serviceNameID, ip.String(), tags); err != nil {
		panic(err)

	}

}

func getOutIp() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}
