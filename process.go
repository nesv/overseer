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

import "errors"

var (
	ErrProcessAlreadyStarted = errors.New("process has already been started")
	ErrProcessNotRunning     = errors.New("process is not running")
)

type (
	Process struct {
		Name           string            `yaml:"name"`
		Command        string            `yaml:"command"`
		Env            map[string]string `yaml:"env,flow"`
		RedirectStderr bool              `yaml:"redirect_stderr,omitempty"`
		StdoutLogfile  string            `yaml:"stdout_logfile,omitempty"`
		StderrLogfile  string            `yaml:"stderr_logilfe,omitempty"`

		running bool
		stop    chan struct{}
		pid     int
	}
)

func LoadProcess(configPath string) (*Process, error) {
	var proc Process
	err := LoadYAML(configPath, &proc)
	return &proc, err
}

func (p *Process) Start() error {
	if p.running {
		return ErrProcessAlreadyStarted
	}
	return nil
}

func (p *Process) Stop() error {
	if !p.running {
		return ErrProcessNotRunning
	}
	return nil
}
