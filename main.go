package main

import (
	"consul-demo/consul"
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	addr = flag.String("listen", ":8080", "The address to listen on for HTTP requests.")
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
	flag.Parse()

	temp := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "host_temperature_celsius",
		Help: "The current temperature is in degress Celsius",
	})

	prometheus.MustRegister(temp)
	temp.Set(39)

	consulClient, err := consul.NewConsulClient("127.0.0.1:8500")
	if err != nil {
		return
	}

	ip, err := getOutIp()
	if err != nil {
		panic(err)
	}

	log.Println(ip.String())

	port := 20001
	checkPoint := 8080
	serviceName := "hello-consul"
	serviceNameID := "hello-consul-20001"
	tags := []string{"consul-demo"}
	if err := consulClient.Register(port, checkPoint, serviceName, serviceNameID, ip.String(), tags); err != nil {
		panic(err)
	}

	//discovery service
	service, err := consulClient.DiscoverService(serviceNameID)
	if err != nil {
		return
	}

	//put key
	if err := consulClient.PutKV("test", "test_value"); err != nil {
		log.Println(err)
	}

	//get key
	value, err := consulClient.GetKV("test")
	if err != nil {
		log.Println(err)
	}
	fmt.Println("get value:", value)

	fmt.Printf("service:host:%s port:%d\n", service.Address, service.Port)

	var count int64

	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		s := "consul check" + fmt.Sprint(count) + "remote:" + r.RemoteAddr + " " + r.URL.String()
		fmt.Println(s)
		fmt.Fprintf(w, s)
		count++
	})

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		h := promhttp.Handler()
		h.ServeHTTP(w, r)
	})

	err = http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	err = consulClient.Deregister(serviceNameID)
	if err != nil {
		log.Fatal(err)
		return
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
