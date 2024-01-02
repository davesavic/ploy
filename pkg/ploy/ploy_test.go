package ploy_test

import (
	"fmt"
	"github.com/davesavic/ploy/pkg/ploy"
	"github.com/gliderlabs/ssh"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type SSHTestServer struct {
	server *ssh.Server
}

func NewSSHTestServer(address string) *SSHTestServer {
	return &SSHTestServer{
		server: &ssh.Server{
			Addr: address,
			PublicKeyHandler: func(ctx ssh.Context, key ssh.PublicKey) bool {
				return true
			},
		},
	}
}

func (s *SSHTestServer) ListenAndServe() error {
	return s.server.ListenAndServe()
}

func (s *SSHTestServer) Close() error {
	return s.server.Close()
}

func (s *SSHTestServer) SetOutputString(output string) {
	s.server.Handler = func(s ssh.Session) {
		_, _ = io.WriteString(s, output)
		_ = s.Exit(0)
	}
}

func TestRemotePipelineExecutor_Execute2(t *testing.T) {
	testCases := []struct {
		name           string
		prepareFunc    func() (*SSHTestServer, ploy.Config)
		pipeline       string
		expectedOutput string
		expectedError  string
	}{
		{
			name:           "valid execution",
			pipeline:       "test",
			expectedOutput: "Hello, World!",
			prepareFunc: func() (*SSHTestServer, ploy.Config) {
				s := NewSSHTestServer("localhost:2222")
				s.SetOutputString("Hello, World!")
				go func() {
					_ = s.ListenAndServe()
				}()

				cwd, _ := os.Getwd()
				testDataPath := filepath.Join(cwd, "testdata")

				cfg := ploy.Config{
					Servers: ploy.Servers{
						"test": ploy.Server{
							Host:       "localhost",
							Port:       2222,
							User:       "test",
							PrivateKey: fmt.Sprintf("%s/id_rsa.test", testDataPath),
						},
					},
					Pipelines: ploy.Pipelines{
						"test": ploy.Pipeline{
							Servers: []string{"test"},
							Tasks:   []string{"echo"},
						},
					},
					Tasks: map[string][]string{
						"echo": {"echo '{{message}}'"},
					},
					Params: map[string]string{
						"message": "Hello, World!",
					},
				}

				return s, cfg
			},
		},
		{
			name:          "invalid server",
			pipeline:      "test",
			expectedError: "server test does not exist",
			prepareFunc: func() (*SSHTestServer, ploy.Config) {
				cfg := ploy.Config{
					Pipelines: ploy.Pipelines{
						"test": ploy.Pipeline{
							Servers: []string{"test"},
							Tasks:   []string{"echo"},
						},
					},
					Tasks: map[string][]string{
						"echo": {"echo '{{message}}'"},
					},
					Params: map[string]string{
						"message": "Hello, World!",
					},
				}

				return nil, cfg
			},
		},
		{
			name:          "missing pipeline server",
			pipeline:      "test",
			expectedError: "pipeline test has no servers",
			prepareFunc: func() (*SSHTestServer, ploy.Config) {
				cfg := ploy.Config{
					Servers: map[string]ploy.Server{
						"test": {
							Host:       "localhost",
							Port:       2222,
							PrivateKey: "id_rsa.test",
							User:       "test",
						},
					},
					Pipelines: ploy.Pipelines{
						"test": ploy.Pipeline{
							Tasks: []string{"echo"},
						},
					},
					Tasks: map[string][]string{
						"echo": {"echo '{{message}}'"},
					},
					Params: map[string]string{
						"message": "Hello, World!",
					},
				}

				return nil, cfg
			},
		},
		{
			name:          "invalid pipeline",
			pipeline:      "invalid",
			expectedError: "pipeline invalid is not defined",
			prepareFunc: func() (*SSHTestServer, ploy.Config) {
				cfg := ploy.Config{
					Servers: map[string]ploy.Server{
						"test": {
							Host:       "localhost",
							Port:       2222,
							PrivateKey: "id_rsa.test",
							User:       "test",
						},
					},
					Pipelines: ploy.Pipelines{
						"test": ploy.Pipeline{
							Servers: []string{"test"},
							Tasks:   []string{"echo"},
						},
					},
					Tasks: map[string][]string{
						"echo": {"echo '{{message}}'"},
					},
					Params: map[string]string{
						"message": "Hello, World!",
					},
				}

				return nil, cfg
			},
		},
		{
			name:          "invalid task",
			pipeline:      "test",
			expectedError: "task invalid is not defined",
			prepareFunc: func() (*SSHTestServer, ploy.Config) {
				s := NewSSHTestServer("localhost:2222")
				s.SetOutputString("Hello, World!")
				go func() {
					_ = s.ListenAndServe()
				}()

				cwd, _ := os.Getwd()
				testDataPath := filepath.Join(cwd, "testdata")

				cfg := ploy.Config{
					Servers: ploy.Servers{
						"test": ploy.Server{
							Host:       "localhost",
							Port:       2222,
							User:       "test",
							PrivateKey: fmt.Sprintf("%s/id_rsa.test", testDataPath),
						},
					},
					Pipelines: ploy.Pipelines{
						"test": ploy.Pipeline{
							Servers: []string{"test"},
							Tasks:   []string{"invalid"},
						},
					},
					Tasks: map[string][]string{
						"echo": {"echo '{{message}}'"},
					},
					Params: map[string]string{
						"message": "Hello, World!",
					},
				}

				return s, cfg
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server, cfg := tc.prepareFunc()
			if server != nil {
				defer func() {
					_ = server.Close()
				}()
			}

			re := ploy.RemotePipelineExecutor{Config: cfg}
			out, err := re.Execute(tc.pipeline)

			if tc.expectedError != "" {
				if err == nil || !strings.Contains(err.Error(), tc.expectedError) {
					t.Errorf("expected error to contain '%s', got '%s'", tc.expectedError, err.Error())
				}

				return
			}

			if err != nil {
				t.Errorf("error executing pipeline: %v", err)
			}

			if !strings.Contains(out, tc.expectedOutput) {
				t.Errorf("expected output to contain '%s', got '%s'", tc.expectedOutput, out)
			}
		})
	}
}

func TestLocalPipelineExecutor_Execute(t *testing.T) {
	testCases := []struct {
		name           string
		prepareFunc    func() ploy.Config
		pipeline       string
		expectedOutput string
		expectedError  string
	}{
		{
			name:           "valid execution",
			pipeline:       "test",
			expectedOutput: "Hello, World!",
			prepareFunc: func() ploy.Config {
				cfg := ploy.Config{
					Pipelines: ploy.Pipelines{
						"test": ploy.Pipeline{
							Tasks: []string{"echo"},
						},
					},
					Tasks: map[string][]string{
						"echo": {fmt.Sprintf("echo '%s'", "Hello, World!")},
					},
				}

				return cfg
			},
		},
		{
			name:          "invalid pipeline",
			pipeline:      "invalid",
			expectedError: "pipeline invalid is not defined",
			prepareFunc: func() ploy.Config {
				cfg := ploy.Config{
					Pipelines: ploy.Pipelines{
						"test": ploy.Pipeline{
							Tasks: []string{"echo"},
						},
					},
					Tasks: map[string][]string{
						"echo": {fmt.Sprintf("echo '%s'", "Hello, World!")},
					},
				}

				return cfg
			},
		},
		{
			name:          "invalid task",
			pipeline:      "test",
			expectedError: "task invalid is not defined",
			prepareFunc: func() ploy.Config {
				cfg := ploy.Config{
					Pipelines: ploy.Pipelines{
						"test": ploy.Pipeline{
							Tasks: []string{"invalid"},
						},
					},
					Tasks: map[string][]string{
						"echo": {fmt.Sprintf("echo '%s'", "Hello, World!")},
					},
				}

				return cfg
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := tc.prepareFunc()
			le := ploy.LocalPipelineExecutor{Config: cfg}
			out, err := le.Execute(tc.pipeline)

			if tc.expectedError != "" {
				if err == nil || !strings.Contains(err.Error(), tc.expectedError) {
					t.Errorf("expected error to contain '%s', got '%s'", tc.expectedError, err.Error())
				}

				return
			}

			if err != nil {
				t.Errorf("error executing pipeline: %v", err)
			}

			if !strings.Contains(out, tc.expectedOutput) {
				t.Errorf("expected output to contain '%s', got '%s'", tc.expectedOutput, out)
			}
		})
	}
}
