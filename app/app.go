package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

const (
	backend     = "consul"
	backendAddr = "localhost:8500"
	backendPath = "default/global/app/cache/shinto"
	configType  = "json"
)

var config string
var done chan bool

func init() {
	GetConfigCache()
	fmt.Println("Initializing Config values...")
	fmt.Println(config)
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

func GetRandomInt() int {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	num := r1.Intn(100)
	fmt.Println("new config... %d", num)
	return num
}

func GetConfigCache() {
	fmt.Println("Print from cache...")
	// config = strconv.Itoa(GetRandomInt())

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
	// runtime_viper.WriteConfigAs("/var/log/config.json")
	config = runtime_viper.GetString("name") + " : " + runtime_viper.GetString("age")
}

func ReadConfigCache() *bytes.Buffer {
	b := bytes.NewBufferString(config)
	return b
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
			fmt.Println(ReadConfigCache())
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	done <- true
	fmt.Println("Dropping old values and refetching new values ...")
	GetConfigCache()
	go CliPrint(done)

	w.Header().Set("Content-type", "application/json")

	if _, err := ReadConfigCache().WriteTo(w); err != nil {
		fmt.Fprintf(w, "%s", err)
	}
}

func main() {
	done = make(chan bool)
	go CliPrint(done)
	http.HandleFunc("/", handler)
	http.ListenAndServe(":7001", nil)
}
