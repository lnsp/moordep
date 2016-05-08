package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

var (
	certPath   = flag.String("cert", "server.pem", "Certificate file for HTTPS")
	keyPath    = flag.String("key", "server.key", "Key file for HTTPS")
	configPath = flag.String("config", "config.json", "Configuration file")
	debugMode  = flag.Bool("debug", false, "Enables debug mode (no HTTPS)")
	cfg        = Config{
		Host:  "localhost",
		Port:  8080,
		Hooks: map[string]string{},
	}
)

type Config struct {
	Host  string            `json:"host"`
	Port  int               `json:"port"`
	Hooks map[string]string `json:"hooks"`
	Token string            `json:"token"`
}

func HookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Not implemented", http.StatusNotImplemented)
		return
	}

	query := r.URL.Query()
	token := query.Get("token")
	hook := query.Get("hook")

	if token != cfg.Token {
		http.Error(w, "Unauthorized request", http.StatusUnauthorized)
		return
	}

	if action, ok := cfg.Hooks[hook]; ok {
		cmd := exec.Command(action)
		cmd.Stdout = os.Stdout
		err := cmd.Start()

		if err != nil {
			log.Fatal(err)
			http.Error(w, "Failed to run hook command", http.StatusInternalServerError)
			return
		}

		log.Printf("Running %s: %s\n", hook, action)
		fmt.Fprintf(w, "Running command %s: %s\n", hook, action)
	} else {
		http.Error(w, "Hook not found", http.StatusNotFound)
		return
	}
}

func main() {
	flag.Parse()

	cfgBytes, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = json.Unmarshal(cfgBytes, &cfg)
	if err != nil {
		log.Fatal(err)
		return
	}
	hostport := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	http.HandleFunc("/", HookHandler)
	if *debugMode {
		log.Fatal(http.ListenAndServe(hostport, nil))
	} else {
		log.Fatal(http.ListenAndServeTLS(hostport, *certPath, *keyPath, nil))
	}
}
