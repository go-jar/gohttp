package gracehttp

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const (
	GRACEFUL_ENVIRON_KEY    = "IS_GRACEFUL"
	GRACEFUL_ENVIRON_STRING = GRACEFUL_ENVIRON_KEY + "=1"
	GRACEFUL_LISTENER_FD    = 3
)

type Server struct {
	httpServer *http.Server
	listener   net.Listener

	isChild      bool
	signalChan   chan os.Signal
	shutdownChan chan bool
}

func NewServer(addr string, handler http.Handler, readTimeout, writeTimeout time.Duration) *Server {
	isChild := false
	if os.Getenv(GRACEFUL_ENVIRON_KEY) != "" {
		isChild = true
	}

	return &Server{
		httpServer: &http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		},
		isChild:      isChild,
		signalChan:   make(chan os.Signal),
		shutdownChan: make(chan bool),
	}
}

func (s *Server) ListenAndServe() error {
	addr := s.httpServer.Addr
	if addr == "" {
		addr = ":http"
	}

	ln, err := s.Listen(addr)
	if err != nil {
		return err
	}
	s.listener = ln

	return s.Serve(ln)
}

func (s *Server) ListenAndServeTLS(certFile, keyFile string) error {
	addr := s.httpServer.Addr
	if addr == "" {
		addr = ":https"
	}

	config := &tls.Config{}
	if s.httpServer.TLSConfig != nil {
		*config = *s.httpServer.TLSConfig
	}
	if config.NextProtos == nil {
		config.NextProtos = []string{"http/1.1"}
	}

	var err error
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}

	ln, err := s.Listen(addr)
	if err != nil {
		return err
	}

	s.listener = tls.NewListener(ln, config)
	return s.Serve(s.listener)
}

func (s *Server) Listen(addr string) (net.Listener, error) {
	var ln net.Listener
	var err error

	if s.isChild {
		file := os.NewFile(GRACEFUL_LISTENER_FD, "")
		ln, err = net.FileListener(file)
		if err != nil {
			s.logf("net.FileListener error: %v", err)
			return nil, err
		}
	} else {
		ln, err = net.Listen("tcp", addr)
		if err != nil {
			s.logf("net.Listen error: %v", err)
			return nil, err
		}
	}
	return ln, nil
}

func (s *Server) Serve(l net.Listener) error {
	go s.handleSignals()
	err := s.httpServer.Serve(l)

	s.logf("Waiting for connections to be closed.")
	<-s.shutdownChan
	s.logf("All connections has been closed.")

	return err
}

func (s *Server) handleSignals() {
	var sig os.Signal

	signal.Notify(
		s.signalChan,
		syscall.SIGTERM,
		syscall.SIGUSR2,
	)

	for {
		sig = <-s.signalChan
		switch sig {
		case syscall.SIGTERM:
			s.logf("Receive signal SIGTERM. To shut down the server gracefully.")
			s.Shutdown()
		case syscall.SIGUSR2:
			s.logf("Receive signal SIGUSR2. To restart the server gracefully.")
			pid, err := s.startNewProcess()
			if err != nil {
				s.logf("Failed to start a new process: %v. Continue serving.", err)
			} else {
				s.logf("A new process %v has been started successfully. To shut down the server gracefully.", pid)
				s.Shutdown()
			}
		}
	}
}

func (s *Server) Shutdown() {
	if err := s.httpServer.Shutdown(context.Background()); err != nil { // Shutdown gracefully
		s.logf("httpServer.Shutdown error: %v", err)
	} else {
		s.logf("The http server has been shut down successfully.")
		s.shutdownChan <- true
	}
}

func (s *Server) startNewProcess() (uintptr, error) {
	listenerFd, err := s.getListenerFd()
	if err != nil {
		s.logf("getListenerFd error: %v", err)
		return 0, err
	}

	envs := os.Environ()
	if os.Getenv(GRACEFUL_ENVIRON_KEY) == "" {
		envs = append(envs, GRACEFUL_ENVIRON_STRING)
	}

	procAttr := &syscall.ProcAttr{
		Env:   envs,
		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd(), listenerFd},
	}

	child, err := syscall.ForkExec(os.Args[0], os.Args, procAttr)
	if err != nil {
		s.logf("syscall.ForkExec error: %v", err)
		return 0, err
	}
	return uintptr(child), nil
}

func (s *Server) getListenerFd() (uintptr, error) {
	file, err := s.listener.(*net.TCPListener).File()
	if err != nil {
		return 0, err
	} else {
		return file.Fd(), nil
	}
}

func (s *Server) logf(pattern string, args ...interface{}) {
	pid := strconv.Itoa(os.Getpid())
	pattern = "[pid: " + pid + "] " + pattern

	if s.httpServer.ErrorLog != nil {
		if args != nil {
			s.httpServer.ErrorLog.Printf(pattern, args)
		} else {
			s.httpServer.ErrorLog.Println(pattern, args)
		}
	} else {
		if args != nil {
			log.Printf(pattern, args)
		} else {
			log.Println(pattern)
		}
	}
}
