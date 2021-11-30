// Code generated by truss. DO NOT EDIT.
// Rerunning truss will overwrite this file.
// Version: fcd9ff140d
// Version Date: 2021-07-14T06:36:40Z

package main

import (
	"flag"

	// This Service
	"github.com/mises-id/socialsvc/handlers"
	"github.com/mises-id/socialsvc/svc/server"
)

func main() {
	// Update addresses if they have been overwritten by flags
	flag.Parse()

	cfg := server.DefaultConfig
	cfg = handlers.SetConfig(cfg)

	server.Run(cfg)
}
