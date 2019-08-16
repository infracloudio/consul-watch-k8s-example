package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

const (
	backend     = "consul"
	backendAddr = "myconsul-consul-server.default.svc.cluster.local:8500"
	backendPath = "mykey1"
	configType  = "json"
)

var done chan bool

func init() {

	GetConfig()
	fmt.Println("Initializing Config values...")
	fmt.Println(ReadConfig())

}

func GetConfig() {
	var runtime_viper = viper.New()
	err := runtime_viper.AddRemoteProvider(backend, backendAddr, backendPath)
	if err != nil {
		fmt.Println("Unable to connect to Consul", err)
	}
	runtime_viper.SetConfigType(configType)

	// read from remote config the first time.
	err = runtime_viper.ReadRemoteConfig()
	if err != nil {
		fmt.Println("Unable to read config", err)
	}
	runtime_viper.WriteConfigAs("/var/log/config.json")
}

func ReadConfig() *bytes.Buffer {
	f, err := ioutil.ReadFile("/var/log/config.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	b := bytes.NewBuffer(f)

	return b
}

func CliPrint(done chan bool) {
	ticker := time.NewTicker(2 * time.Second)
	for range ticker.C {
		select {
		case <-done:
			return
		default:
			fmt.Println(ReadConfig())
		}

	}
}

func handler(w http.ResponseWriter, r *http.Request) {

	done <- true
	fmt.Println("Dropping old values and refetching new values ...")
	GetConfig()
	go CliPrint(done)

	w.Header().Set("Content-type", "application/json")

	if _, err := ReadConfig().WriteTo(w); err != nil {
		fmt.Fprintf(w, "%s", err)
	}

}

func main() {
	done = make(chan bool)
	go CliPrint(done)
	http.HandleFunc("/", handler)
	http.ListenAndServe(":7001", nil)

}
