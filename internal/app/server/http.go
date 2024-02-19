package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type (
	server struct {
		server  *http.Server
		onClose []func()
	}
)

var (
	svr     *server
	svrInit sync.Once
)

func New(cfg Config, handler http.Handler, closeFns ...func()) *server {
	svrInit.Do(func() {
		svr = &server{
			server: &http.Server{
				Handler:      handler,
				Addr:         fmt.Sprintf(":%d", cfg.Port),
				ReadTimeout:  cfg.ReadTimeout,
				WriteTimeout: cfg.WriteTimeout,
				IdleTimeout:  cfg.IdleTimeout,
			},
			onClose: closeFns,
		}
	})
	return svr
}

func (s *server) serveWithGracefulShutdown(timeout time.Duration, proc func() error) error {
	// handle close connection to data stores like mongo, redis, etc.
	defer func() {
		for _, fn := range s.onClose {
			fn()
		}
	}()

	// Service listener for the 'SIGTERM' from kernel
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Wait for 'SIGTERM' from kernel
	var errRunning error
	go func() {
		if err := proc(); err != nil && err != http.ErrServerClosed {
			errRunning = err
		}
	}()
	<-quit

	// Service the cancelable context for help cancel the halted shutdown process
	srvCtx, srvCancel := context.WithTimeout(context.Background(), timeout)
	defer srvCancel()

	// Perform shutdown then wait until the server finished the shutdown
	// process or the timeout had been reached
	if errShutdown := s.server.Shutdown(srvCtx); errShutdown != nil {
		return errShutdown
	}
	return errRunning
}

func (s *server) ListenAndServe(timeout time.Duration, withTLS TLSConfig) error {
	if withTLS.Enabled {
		s.server.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
			ClientAuth: tls.RequireAnyClientCert,
		}
		return s.serveWithGracefulShutdown(timeout, func() error {
			return s.server.ListenAndServeTLS(withTLS.CertFile, withTLS.KeyFile)
		})
	}
	return s.serveWithGracefulShutdown(timeout, func() error {
		return s.server.ListenAndServe()
	})
}
