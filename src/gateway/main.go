package main

import (
	"context"
	"fmt"
	"os"
	"time"

	_ "cloud.google.com/go/profiler"
	_ "contrib.go.opencensus.io/exporter/stackdriver"
	_ "github.com/gorilla/mux"
	_ "github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	_ "go.opencensus.io/plugin/ocgrpc"
	_ "go.opencensus.io/plugin/ochttp"
	_ "go.opencensus.io/plugin/ochttp/propagation/b3"
	_ "go.opencensus.io/stats/view"
	_ "go.opencensus.io/trace"
	_ "google.golang.org/grpc"
)

const (
	serviceName = "users-service"
	port        = "8080"
)

type ctxKeySessionId struct {
}

func main() {
	ctx := context.Background()
	// Create a new logger configuration
	log := logrus.New()
	log.Level = logrus.DebugLevel
	log.Formatter = &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
		TimestampFormat: time.RFC3339Nano,
	}
	log.Out = os.Stdout
}

// Function defines the necessary tracing configurations for Jaeger framework.
func initJaegerTracing(log logrus.FieldLogger) {
	svcAddr := os.Getenv("JAEGER_SERVICE_ADDR")
	if svcAddr == "" {
		log.Info("jaeger initialization disabled.")
		return
	}

	// Register Jaeger exporter to retrieve colleted spans
	exporter, err := jaeger.NewExporter(jaeger.Options{
		Endpoint: fmt.Sprintf("http://%s", svcAddr),
		Process: jaeger.Process{
			ServiceName: serviceName,
		},
	})

	if err != nil {
		log.Fatal(err)
	}
	trace.RegisterExporter(exporter)
	log.Info("jaeger initialization completed.")
}
