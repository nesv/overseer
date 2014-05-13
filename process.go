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
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/golang/glog"
)

var (
	ErrProcessAlreadyStarted = errors.New("process has already been started")
	ErrProcessNotRunning     = errors.New("process is not running")
	ErrNoProcessName         = errors.New("process name cannot be empty")
)

type (
	Process struct {
		Name           string            `yaml:"name"`
		Command        string            `yaml:"command"`
		Env            map[string]string `yaml:"env,flow"`
		WorkingDir     string            `yaml:"working_dir,omitempty"`
		RedirectStderr bool              `yaml:"redirect_stderr,omitempty"`
		StdoutLogfile  string            `yaml:"stdout_logfile,omitempty"`
		StderrLogfile  string            `yaml:"stderr_logilfe,omitempty"`

		running bool
		stop    chan struct{}
		pid     int
	}
)

func LoadProcess(configPath string, globalConfig Configuration) (*Process, error) {
	var proc Process
	err := LoadYAML(configPath, &proc)

	if proc.Name == "" {
		return nil, ErrNoProcessName
	}

	now := time.Now()
	if proc.StdoutLogfile == "" {
		lf := fmt.Sprintf("%s-%d.out.log", proc.Name, now.Unix())
		proc.StdoutLogfile = filepath.Join(globalConfig.LogDirectory, lf)
	}
	if proc.StderrLogfile == "" && !proc.RedirectStderr {
		lf := fmt.Sprintf("%s-%d.err.log", proc.Name, now.Unix())
		proc.StderrLogfile = filepath.Join(globalConfig.LogDirectory, lf)
	}

	return &proc, err
}

func (p *Process) Start() error {
	if p.running {
		return ErrProcessAlreadyStarted
	}

	cmdParts := strings.Split(p.Command, " ")
	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	if p.RedirectStderr {
		// Combine STDOUT and STDERR into the same stream.
		lf, err := os.OpenFile(p.StdoutLogfile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
		if err != nil {
			return fmt.Errorf("error opening logfile %q: %v", p.StdoutLogfile, err)
		}
		defer lf.Close()
		cmd.Stdout = lf
		cmd.Stderr = lf
	} else {
		of, err := os.OpenFile(p.StdoutLogfile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
		if err != nil {
			return fmt.Errorf("error opening logfile %q: %v", p.StdoutLogfile, err)
		}
		defer of.Close()
		cmd.Stdout = of

		ef, err := os.OpenFile(p.StderrLogfile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
		if err != nil {
			return fmt.Errorf("error opening logfile %q: %v", p.StderrLogfile, err)
		}
		defer ef.Close()
		cmd.Stderr = ef
	}

	// Set the environment for the command.
	var env = make([]string, 0)
	for k, v := range p.Env {
		env = append(env, fmt.Sprintf("%s=%q", k, v))
	}
	if len(env) > 0 {
		cmd.Env = env
	}

	// Set the working directory for the process (if specified in the
	// config file).
	if p.WorkingDir != "" {
		cmd.Dir = p.WorkingDir
	}

	// Start the command!
	if err := cmd.Start(); err != nil {
		return err
	}
	p.running = true
	p.pid = cmd.Process.Pid

	var done = make(chan *os.ProcessState)
	var errors = make(chan error)
	go func() {
		if err := cmd.Wait(); err != nil {
			errors <- err
		}
		done <- cmd.ProcessState
	}()
	defer close(done)
	defer close(errors)

	p.stop = make(chan struct{}, 1)
	for {
		select {
		case pstate := <-done:
			// The process has exited.
			glog.Infoln("Process %d has stopped", pstate.Pid())
			break

		case <-p.stop:
			// Stop the process.
			if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
				p.running = false
			}
			break

		case err := <-errors:
			p.running = false
			return err
		}
	}
	return nil
}

func (p *Process) Stop() error {
	if !p.running {
		return ErrProcessNotRunning
	}
	p.stop <- struct{}{}
	return nil
}
