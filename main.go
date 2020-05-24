package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v6"
	"github.com/nuttmeister/go-shadowsocks2/core"
)

// Constants.
const (
	minPortNo       = 10000 // Min port to be allowed
	maxPortNo       = 19999 // Max port to be allowed
	ssPortIncrement = 10000 // ssPort = port + ssPortIncrement

	listen  = "127.0.0.1"             // Ip that v2ray will listen to.
	forward = "127.0.0.1"             // Ip that v2ray will relay to.
	v2ray   = "/usr/bin/v2ray-plugin" // The location of v2ray-plugin.

	// v2ray-plugin options.
	v2rayOpts = "server;fast-open;tls=true;cert=/ssl/cert.pem;key=/ssl/private.key;host=%s;path=/prxy/%d;loglevel=%s"
)

var cfg *config

// Config contains the programs config.
type config struct {
	// General
	FullHost string `env:"FULL_HOST,required"`
	Cipher   string `env:"CIPHER,required"`
	Password string `env:"PASSWORD,required"`
	Port     int    `env:"PORT,required"`
	ssPort   int

	// Debug
	Verbose bool `env:"VERBOSE" envDefault:"false"`
}

func main() {
	// Configure the service.
	if err := configure(); err != nil {
		log.Fatal(err)
	}

	saltstack.

	// Start v2ray.
	if err := cfg.startV2Ray(); err != nil {
		log.Fatal(err)
	}

	// Create the Cipher.
	ciph, err := core.PickCipher(cfg.Cipher, cfg.Password)
	if err != nil {
		log.Fatal(err)
	}

	// Start the tcp listener.
	go tcpRemote(cfg.ssPort, ciph.StreamConn)

	// Listen for exit and error to quit.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	killPlugin()
}

// configure will configure the program from env vars.
// Returns *config and error.
func configure() error {
	cfg = &config{}
	if err := env.Parse(cfg); err != nil {
		return err
	}

	// Check that the ports are within range.
	if cfg.Port < minPortNo || cfg.Port > maxPortNo {
		return fmt.Errorf("port is not within the allowed range")
	}
	cfg.ssPort = cfg.Port + ssPortIncrement

	return nil
}

// plugin will start any configured plugins.
// Returns error.
func (cfg *config) startV2Ray() error {
	// Set loglevel to none or debug.
	loglevel := "none"
	if cfg.Verbose {
		loglevel = "debug"
	}

	// Start the plugin.
	opts := fmt.Sprintf(v2rayOpts, cfg.FullHost, cfg.Port, loglevel)
	return startPlugin(v2ray, opts, cfg.Port, cfg.ssPort)
}
