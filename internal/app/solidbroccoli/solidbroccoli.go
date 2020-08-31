package solidbroccoli

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/dstdfx/solid-broccoli/internal/app/exporter"
	"github.com/dstdfx/solid-broccoli/internal/pkg/backend"
	"github.com/dstdfx/solid-broccoli/internal/pkg/config"
	public "github.com/dstdfx/solid-broccoli/internal/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

const (
	pprofIndexPath   = "/debug/pprof/"
	pprofCmdlinePath = "/debug/pprof/cmdline"
	pprofProfilePath = "/debug/pprof/profile"
	pprofSymbolPath  = "/debug/pprof/symbol"
	pprofTracePath   = "/debug/pprof/trace"

	metricsPath = "/metrics"

	gracefulShutdownTimeout = 5 * time.Second
)

// StartOpts represents options to be passed to main gorountine.
type StartOpts struct {
	Interrupt      chan os.Signal
	BuildGitCommit string
	BuildGitTag    string
	BuildDate      string
	BuildCompiler  string
}

// StartService runs main service's goroutine.
func StartService(log *zap.Logger, opts StartOpts) error {
	if err := config.CheckConfig(); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}

	log.Debug("start service")

	// Init backend
	b, err := backend.New(log)
	if err != nil {
		return fmt.Errorf("failed to init backend: %w", err)
	}
	defer b.Shutdown()

	// Register new Prometheus exporter
	if err := prometheus.Register(exporter.NewAPIExporter(&exporter.NewAPIExporterOpts{
		BuildGitCommit: opts.BuildGitCommit,
		BuildGitTag:    opts.BuildGitTag,
		BuildDate:      opts.BuildDate,
		BuildCompiler:  opts.BuildCompiler,
	})); err != nil {
		return fmt.Errorf("failed to register prometheus exporter: %w", err)
	}

	// Register service API handler
	httpMux := http.NewServeMux()

	// Register prometheus handler
	httpMux.Handle(metricsPath, promhttp.Handler())

	// Register pprof handlers
	httpMux.HandleFunc(pprofIndexPath, pprof.Index)
	httpMux.HandleFunc(pprofCmdlinePath, pprof.Cmdline)
	httpMux.HandleFunc(pprofProfilePath, pprof.Profile)
	httpMux.HandleFunc(pprofSymbolPath, pprof.Symbol)
	httpMux.HandleFunc(pprofTracePath, pprof.Trace)

	// Configure Service API server
	serviceAPIServer := &http.Server{
		Addr: strings.Join([]string{
			config.Config.ServiceAPI.ServerAddress,
			strconv.Itoa(config.Config.ServiceAPI.ServerPort),
		}, ":"),
		ReadTimeout:  time.Duration(config.Config.ServiceAPI.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.Config.ServiceAPI.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(config.Config.ServiceAPI.IdleTimeout) * time.Second,
		Handler:      httpMux,
	}

	// Configure Public API server
	publicAPIServer := &http.Server{
		Addr: strings.Join([]string{
			config.Config.PublicAPI.ServerAddress,
			strconv.Itoa(config.Config.PublicAPI.ServerPort),
		}, ":"),
		ReadTimeout:  time.Duration(config.Config.PublicAPI.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.Config.PublicAPI.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(config.Config.PublicAPI.IdleTimeout) * time.Second,
		Handler:      public.InitAPIRouter(log, b),
	}

	log.Debug("wait for shutdown signals")
	signal.Notify(opts.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(opts.Interrupt)

	// Serve service API
	go func() {
		log.Info("running service API server", zap.String("addr", serviceAPIServer.Addr))
		if err := serviceAPIServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("failed to serve service API", zap.Error(err))
		}
	}()

	// Serve public API
	go func() {
		log.Info("running public API server", zap.String("addr", publicAPIServer.Addr))
		if err := publicAPIServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("failed to serve public API", zap.Error(err))
		}
	}()

	sig := <-opts.Interrupt
	log.Debug("got a signal", zap.Stringer("sig", sig))

	go func() {
		// Context to shutdown service API-server
		ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
		defer cancel()

		// Shutdown service API-server
		if err := serviceAPIServer.Shutdown(ctx); err != nil {
			log.Warn("service API server shutdown failed", zap.Error(err))
		}
	}()

	// Context to shutdown public API-server
	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
	defer cancel()

	// Shutdown public API-server
	if err := publicAPIServer.Shutdown(ctx); err != nil {
		log.Warn("public API server shutdown failed", zap.Error(err))
	}

	return nil
}
