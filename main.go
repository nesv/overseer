/*
   Copyright 2014 Nick Saika

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package main

import (
	"flag"
	"net/http"
	"path/filepath"

	"github.com/golang/glog"
)

var (
	ConfigDir = flag.String("config-dir", "/usr/local/etc/overseer", "Path to the directory where configurations are stored")
)

func main() {
	flag.Parse()
	glog.Infoln("Loading configuration file")
	config, err := LoadConfig(filepath.Join(*ConfigDir, "config.yml"))
	if err != nil {
		glog.Fatalln(err)
	}

	glog.Infoln("Loading process configuration files")
	processConfigDir := filepath.Join(*ConfigDir, "process.d")
	processConfigs, err := filepath.Glob(filepath.Join(processConfigDir, "*.yml"))
	if err != nil {
		glog.Fatalln(err)
	}
	processes := make([]*Process, 0)
	for _, pconfig := range processConfigs {
		if proc, err := LoadProcess(pconfig, *config); err != nil {
			glog.Fatalln(err)
		} else {
			processes = append(processes, proc)
		}
	}

	glog.Infoln("Starting processes")
	for _, proc := range processes {
		go func() {
			if err := proc.Start(); err != nil {
				glog.Warningln(err)
			}
		}()
	}

	if config.HTTP.Enabled {
		glog.Infoln("Starting HTTP server; listening on", config.HTTP.ListenAddr)
		http.HandleFunc("/", HTTPStatus)
		if err := http.ListenAndServe(config.HTTP.ListenAddr, nil); err != nil {
			glog.Fatalln(err)
		}
	}
	glog.Infoln("Shutting down")
}
