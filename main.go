package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	consulapi "github.com/hashicorp/consul/api"
)

func main() {
	serviceRegistryWithConsul()
	log.Println("Starting frontend World Server...")
	http.HandleFunc("/frontend", helloworld)
	http.HandleFunc("/api", apiWorld)

	http.HandleFunc("/check", check)
	http.ListenAndServe(getPort(), nil)

}

func serviceRegistryWithConsul() {
	config := consulapi.DefaultConfig()
	consul, err := consulapi.NewClient(config)
	if err != nil {
		log.Println(err)
	}

	port, _ := strconv.Atoi(getPort()[1:len(getPort())])
	address := getPodIP()
	serviceID := fmt.Sprintf("frontend-server-%s:%v", address, port)

	tags := []string{"urlprefix-/frontend host=test.netslum.dev", "urlprefix-/api host=test-api.netslum.dev/api"}

	registration := &consulapi.AgentServiceRegistration{
		ID:      serviceID,
		Name:    "frontend-server",
		Port:    port,
		Address: address,
		Tags:    tags,
		Check: &consulapi.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%v/check", address, port),
			Interval: "10s",
			Timeout:  "30s",
		},
	}

	regiErr := consul.Agent().ServiceRegister(registration)

	if regiErr != nil {
		log.Printf("Failed to register service: %s:%v ", address, port)
	} else {
		log.Printf("successfully register service: %s:%v", address, port)
	}
}

func helloworld(w http.ResponseWriter, r *http.Request) {
	log.Println("helloworld service is called.")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello world.")
}

func apiWorld(w http.ResponseWriter, r *http.Request) {
	log.Println("apiWorld service is called.")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "apiWorld.")
}

func check(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Consul check")
}

func getPort() (port string) {
	port = os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	port = ":" + port
	return
}

func getHostname() (hostname string) {
	hostname, _ = os.Hostname()
	return
}

func getPodIP() (podIP string) {
	podIP = os.Getenv("POD_IP")
	return
}
