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
	"errors"
	"io/ioutil"
	"path/filepath"

	yaml "gopkg.in/yaml.v1"
)

var ErrUnsupportedConfig = errors.New("unsupported config file type")

func LoadYAML(path string, v interface{}) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, v)
}

func LoadJSON(path string, v interface{}) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

func LoadFile(path string, v interface{}) error {
	var err error
	switch ext := filepath.Ext(path); ext {
	case ".yaml", ".yml":
		err = LoadYAML(path, v)
	case ".json":
		err = LoadJSON(path, v)
	}
	return err
}

type Configuration struct {
	HTTP struct {
		Enabled    bool   `json:"enabled" yaml:"enabled"`
		ListenAddr string `json:"listen" yaml:"listen"`
	} `json:"http" yaml:"http,flow"`
	LogDirectory string `yaml:"log_directory"`
}

func LoadConfig(path string) (*Configuration, error) {
	var c Configuration
	err := LoadFile(path, &c)
	return &c, err
}
