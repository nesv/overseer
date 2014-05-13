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
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	yaml "gopkg.in/yaml.v1"
)

func LoadYAML(path string, v interface{}) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, v)
}

type Configuration struct {
	HTTP struct {
		Enabled    bool   `json:"enabled" yaml:"enabled"`
		ListenAddr string `json:"listen" yaml:"listen"`
	} `json:"http" yaml:"http,flow"`
	LogDirectory string `yaml:"log_directory"`
}

func LoadConfig(path string) (*Configuration, error) {
	b, err := ioutil.ReadAll(path)
	if err != nil {
		return nil, err
	}

	var c Configuration
	switch ext := filepath.Ext(path); ext {
	case ".yaml", ".yml":
		err = yaml.Unmarshal(b, &c)

	case ".json":
		err = json.Unmarshal(b, &c)
	}

	return &c, err
}
