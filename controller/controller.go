package controller

import (
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	api "github.com/embarkstudios/cassini/api"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/examples/exporter"
	"go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"go.opencensus.io/zpages"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Controller struct {
}

var (
	command        = kingpin.Command("controller", "Controller")
	grpcAddress    = command.Flag("grpc-address", "GRPC address").Default(":3237").String()
	metricsAddress = command.Flag("metrics-address", "Metrics address").Default(":8888").String()
)

func FullCommand() string {
	return command.FullCommand()
}

func NewController() *Controller {
	controller := &Controller{}

	return controller
}

func (c *Controller) Serve() error {
	errCh := make(chan error)

	pe, err := prometheus.NewExporter(prometheus.Options{
		Namespace: "cassini",
	})
	if err != nil {
		return err
	}
	view.RegisterExporter(pe)

	exporter := &exporter.PrintExporter{}
	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	view.SetReportingPeriod(1 * time.Second)

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", pe)
		zpages.Handle(mux, "/")
		log.WithField("address", *metricsAddress).Info("Starting metrics server")
		if err := http.ListenAndServe(*metricsAddress, mux); err != nil {
			errCh <- err
		}
	}()

	go func() {
		log.WithField("address", *grpcAddress).Info("Starting GRPC server")
		listen, err := net.Listen("tcp", *grpcAddress)
		if err != nil {
			errCh <- err
			return
		}
		gs := grpc.NewServer()
		api.RegisterCassiniServer(gs, c)
		reflection.Register(gs)
		err = gs.Serve(listen)
		if err != nil {
			errCh <- err
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigCh:
		log.Warn("Received SIGTERM, exiting gracefully...")
	case err := <-errCh:
		log.WithError(err).Error("Got an error from errCh, exiting gracefully")
		return err
	}

	return nil
}

func RunController() {
	c := NewController()
	err := c.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
