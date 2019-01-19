package tempserver

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"text/template"
)

type Config struct {
	Path           string
	Arguments      func(string) []string
	ConfigTemplate *template.Template
	Config         interface{}
	WaitFor        string
}

type Server struct {
	config    *Config
	dir       string
	cmd       *exec.Cmd
	stdout    io.Reader
	stdoutBuf bytes.Buffer
	stderr    io.Reader
}

func Start(config *Config) (server *Server, err error) {
	if config == nil {
		config = &Config{}
	}

	dir, err := ioutil.TempDir(os.TempDir(), "tempserver")
	if err != nil {
		return nil, err
	}

	server = &Server{
		dir:    dir,
		config: config,
	}

	err = server.start()

	if err != nil {
		return server, err
	}

	if config.WaitFor != "" {
		err = server.waitFor(config.WaitFor)
	}

	return server, err
}

func (s *Server) start() (err error) {

	if s.cmd != nil {
		return fmt.Errorf("Server has already been started")
	}

	configPath := s.dir + "/config"

	if s.config.ConfigTemplate != nil && s.config.Config != nil {

		f, err := os.Create(configPath)

		if err != nil {
			return err
		}

		err = s.config.ConfigTemplate.Execute(f, s.config.Config)

		if err != nil {
			return err
		}

	}

	s.cmd = exec.Command(s.config.Path, s.config.Arguments(configPath)...)

	s.stdout, _ = s.cmd.StdoutPipe()
	s.stderr, _ = s.cmd.StderrPipe()

	return s.cmd.Start()

}

func (s *Server) waitFor(search string) (err error) {

	fmt.Println("Starting waitFor")

	var line string

	scanner := bufio.NewScanner(s.stdout)
	for scanner.Scan() {
		line = scanner.Text()
		fmt.Fprintf(&s.stdoutBuf, "%s\n", line)
		if strings.Contains(line, search) {
			return nil
		}
	}
	fmt.Println("After scan")
	err = scanner.Err()
	if err == nil {
		err = io.EOF
	}
	return err
}

func (s *Server) Term() (err error) {
	return s.signalAndCleanup(syscall.SIGTERM)
}

func (s *Server) Kill() (err error) {
	return s.signalAndCleanup(syscall.SIGKILL)
}

func (s *Server) signalAndCleanup(sig syscall.Signal) error {
	s.cmd.Process.Signal(sig)
	_, err := s.cmd.Process.Wait()
	os.RemoveAll(s.dir)
	return err
}
