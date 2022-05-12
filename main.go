package main

import (
	"net/http"

	"github.com/docopt/docopt-go"
	"github.com/reconquest/pkg/log"
)

var version = "[manual build]"

var usage = `mihome - Xiaomi based smart home automation

Usage:
  mihome [options]

Options:
  -c --config <path>  Read specified config [default: mihome.yaml]
  -h --help           Show this help.
`

type Opts struct {
	ValueConfig string `docopt:"--config"`
}

func main() {
	args, err := docopt.ParseArgs(usage, nil, "mihome "+version)
	if err != nil {
		panic(err)
	}

	var opts Opts
	err = args.Bind(&opts)
	if err != nil {
		log.Fatal(err)
	}

	log.SetLevel(log.LevelDebug)

	config, err := LoadConfig(opts.ValueConfig)
	if err != nil {
		log.Fatalf(err, "load config")
	}

	log.Infof(nil, "requesting list of xiaomi devies")

	devices, err := getDevices(config)
	if err != nil {
		log.Fatalf(err, "get devices")
	}

	log.Infof(nil, "found %d xiaomi devices", len(devices))

	server := Server{
		config:  config,
		devices: devices,
	}

	log.Infof(nil, "listening and serving on %s", config.Listen)
	err = http.ListenAndServe(config.Listen, &server)
	if err != nil {
		log.Fatal(err)
	}
}
