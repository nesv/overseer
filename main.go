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
	"os"
	"path/filepath"

	"github.com/golang/glog"
)

var (
	ConfigDir = flag.String("config-dir", "/usr/local/etc/overseer", "Path to the directory where configurations are stored")
)

func main() {
	flag.Parse()
	glog.Infoln("Loading configuration file")
	configPaths := []string{
		filepath.Join(*ConfigDir, "config.yml"),
		filepath.Join(*ConfigDir, "config.json")}
	var config *Configuration
	var err error
	for _, cfgpth := range configPaths {
		config, err = LoadConfig(cfgpth)
		if err != nil && os.IsNotExist(err) {
			// The configuration file we just attempted to load
			// does not exist. No issue; we will just move on
			// to the next one.
			continue
		} else if err != nil {
			glog.Fatalln(err)
		} else if err == nil {
			break
		}
	}
	if config == nil {
		// *WELP*, it looks like we couldn't find a configuration file.
		glog.Fatalln("No config.json or config.yml file found")
	}

	glog.Infoln("Loading process configuration files")
	processConfigDir := filepath.Join(*ConfigDir, "process.d")
	processConfigs := make([]string, 0)
	globs := []string{
		filepath.Join(processConfigDir, "*.yml"),
		filepath.Join(processConfigDir, "*.yaml"),
		filepath.Join(processConfigDir, "*.json")}
	for _, pat := range globs {
		matches, err := filepath.Glob(pat)
		if err != nil {
			glog.Fatalln(err)
		}
		processConfigs = append(processConfigs, matches...)
	}
	processes := make([]*Process, len(processConfigs))
	for i, pconfig := range processConfigs {
		if proc, err := LoadProcess(pconfig, *config); err != nil {
			glog.Fatalln(err)
		} else {
			processes[i] = proc
		}
	}

	glog.Infoln("Starting processes")
	for _, proc := range processes {
		glog.V(1).Infoln("|-", proc.Name)
		go func(p *Process) {
			if err := p.Start(); err != nil {
				glog.Warningf("%s: %v", p.Name, err)
			}
		}(proc)
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
