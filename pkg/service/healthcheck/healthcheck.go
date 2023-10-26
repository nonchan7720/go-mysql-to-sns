package healthcheck

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/nonchan7720/go-mysql-to-sns/pkg/logging"
)

type HealthChecker struct {
	addr   string
	server *http.Server
	pinger IPing
}

var (
	_ http.Handler = (*HealthChecker)(nil)
)

func New(pinger IPing, addr string) *HealthChecker {
	checker := &HealthChecker{
		addr:   addr,
		pinger: pinger,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", checker.healthCheck)
	checker.server = &http.Server{
		Handler: mux,
	}
	return checker
}

func (checker *HealthChecker) healthCheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	mp := map[string]string{}
	if err := checker.pinger.PingContext(ctx); err != nil {
		slog.With(logging.WithStack(err)).ErrorContext(ctx, "Health check error.")
		w.WriteHeader(http.StatusInternalServerError)
		mp["error"] = err.Error()
	} else {
		w.WriteHeader(http.StatusOK)
		mp["status"] = "OK"
	}
	_ = json.NewEncoder(w).Encode(&mp)
}

func (checker *HealthChecker) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", checker.addr)
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = checker.shutdown(timeoutCtx)
	}()

	slog.Info(fmt.Sprintf("Start health check addr %s", lis.Addr().String()))
	if err := checker.server.Serve(lis); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (checker *HealthChecker) shutdown(ctx context.Context) error {
	return checker.server.Shutdown(ctx)
}

func (checker *HealthChecker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	checker.server.Handler.ServeHTTP(w, r)
}
