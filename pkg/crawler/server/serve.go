package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-logr/logr"

	"github.com/morlay/crawler/pkg/crawler"
)

func Serve(ctx context.Context, c crawler.Crawler) error {
	srv := &http.Server{Addr: ":7666", Handler: http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		q := req.URL.Query()
		op := ""
		params := map[string]string{}

		for name := range q {
			if name == "operation" {
				op = q.Get(name)
				continue
			}
			params[name] = q.Get(name)
		}

		ctx := req.Context()

		ret, err := c.Do(ctx, op, params)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintf(rw, "%#v", err)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		if err := ret.Scan(ctx, rw); err != nil {
			_, _ = fmt.Fprintf(rw, "%#v", err)
			return
		}
	})}

	l := logr.FromContextOrDiscard(ctx)

	go func() {
		l.Info(fmt.Sprintf("listen on %s", srv.Addr))

		if err := srv.ListenAndServe(); err != nil {
			l.Error(err, "")
			if err != http.ErrServerClosed {
				panic(err)
			}
		}
	}()

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)
	<-stopCh

	timeout := 10 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	l.Info("shutdowning in %s", timeout)

	return srv.Shutdown(ctx)
}
