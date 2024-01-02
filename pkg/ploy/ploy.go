package ploy

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Params map[string]string

type Server struct {
	Host       string `json:"host"`
	Port       int    `json:"port"`
	User       string `json:"user"`
	PrivateKey string `json:"private-key"`
}

type Servers map[string]Server

type Tasks map[string][]string

type Pipeline struct {
	Servers []string `json:"servers"`
	Tasks   []string `json:"tasks"`
}

type Pipelines map[string]Pipeline

type Config struct {
	Params    Params    `json:"params"`
	Servers   Servers   `json:"servers"`
	Tasks     Tasks     `json:"tasks"`
	Pipelines Pipelines `json:"pipelines"`
}

type PipelineExecutor interface {
	Execute(pipeline string) (string, error)
}

type RemotePipelineExecutor struct {
	Config Config
}

func (r *RemotePipelineExecutor) Execute(pipeline string) (string, error) {
	pl, exists := r.Config.Pipelines[pipeline]
	if !exists {
		return "", fmt.Errorf("pipeline %s is not defined", pipeline)
	}

	var out bytes.Buffer

	if len(pl.Servers) == 0 {
		return "", fmt.Errorf("pipeline %s has no servers", pipeline)
	}

	for _, s := range pl.Servers {
		out.Write([]byte(fmt.Sprintf("Running pipeline %s on server %s\n", pipeline, s)))

		server, exists := r.Config.Servers[s]
		if !exists {
			return "", fmt.Errorf("server %s does not exist", s)
		}

		key, err := os.ReadFile(server.PrivateKey)
		if err != nil {
			return "", fmt.Errorf("error reading private key (%s): %v", server.PrivateKey, err)
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return "", fmt.Errorf("error parsing private key (%s): %v", server.PrivateKey, err)
		}

		sshCfg := ssh.ClientConfig{
			User: server.User,
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(), //Update for production
		}

		client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", server.Host, server.Port), &sshCfg)
		if err != nil {
			return "", fmt.Errorf("error dialing SSH server (%s): %v", s, err)
		}

		for _, t := range pl.Tasks {
			commands, exists := r.Config.Tasks[t]
			if !exists {
				return "", fmt.Errorf("task %s is not defined", t)
			}

			for _, c := range commands {
				populatePlaceholders(&c, r.Config.Params)
				session, err := client.NewSession()
				if err != nil {
					return "", fmt.Errorf("error creating SSH session: %v", err)
				}

				op, err := session.CombinedOutput(c)
				if err != nil {
					return "", fmt.Errorf("error running command (%s): %v", c, err)
				}

				out.Write(op)

				if err := session.Close(); err != nil && err != io.EOF {
					return "", fmt.Errorf("error closing SSH session: %v", err)
				}
			}
		}

		if err := client.Close(); err != nil {
			return "", fmt.Errorf("error closing SSH client: %v", err)
		}
	}

	return out.String(), nil
}

func populatePlaceholders(s *string, params Params) {
	timestamp := fmt.Sprintf("%s", time.Now().Format("20060102150405"))
	*s = strings.ReplaceAll(*s, "{{timestamp}}", timestamp)

	for k, v := range params {
		*s = strings.ReplaceAll(*s, fmt.Sprintf("{{%s}}", k), v)
	}
}

type LocalPipelineExecutor struct {
	Config Config
}

func (l *LocalPipelineExecutor) Execute(pipeline string) (string, error) {
	var out bytes.Buffer

	pl, exists := l.Config.Pipelines[pipeline]
	if !exists {
		return "", fmt.Errorf("pipeline %s is not defined", pipeline)
	}

	for _, t := range pl.Tasks {
		task, exists := l.Config.Tasks[t]
		if !exists {
			return "", fmt.Errorf("task %s is not defined", t)
		}

		for _, c := range task {
			populatePlaceholders(&c, l.Config.Params)
			cmd := exec.Command("sh", "-c", c)
			cmd.Stdout = &out
			cmd.Stderr = &out

			err := cmd.Run()
			if err != nil {
				return "", err
			}
		}
	}

	return out.String(), nil
}

type Ploy struct {
	Config Config
}
