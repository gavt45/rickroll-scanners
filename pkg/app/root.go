package app

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gavt45/rickroll-scanners/pkg/config"
	"github.com/gavt45/rickroll-scanners/pkg/rickrolls"
	"github.com/gorilla/mux"
	"golang.org/x/sync/errgroup"
)

const (
	base10             = 10
	readHeaderTimeout  = time.Second
	writeTimeout       = 15 * time.Second
	shutdownPeriod     = 3 * time.Second
	hardShutdownPeriod = 2 * time.Second
)

type rickrollApp struct {
	cfg *config.AppConfig
	mux *mux.Router
}

func New(cfg *config.AppConfig) *rickrollApp {
	app := &rickrollApp{
		cfg: cfg,
		mux: mux.NewRouter(),
	}

	app.initHandlers()

	return app
}

func (r *rickrollApp) defaultHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Server", r.cfg.Server.Useragent)

	rr := rand.IntN(len(rickrolls.RickRolls))

	rickrolls.RickRolls[rr](w, req)
}

func (r *rickrollApp) loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s %s %s\n", r.RemoteAddr, r.Method, r.URL.Path, r.Proto, r.Header.Get("User-Agent"))

		next.ServeHTTP(w, r)
	})
}

func (r *rickrollApp) initHandlers() {
	r.mux.Use(r.loggerMiddleware)

	for _, pattern := range r.cfg.BadPatterns {
		r.mux.HandleFunc(pattern, r.defaultHandler)
	}

	for _, prefix := range r.cfg.PrefixPatterns {
		r.mux.PathPrefix(prefix).HandlerFunc(r.defaultHandler)
	}
}

func (r *rickrollApp) runServer(ctx context.Context) error {
	errchan := make(chan error, 1)

	// Context for all incoming connections
	ongoingCtx, stopOngoingGracefully := context.WithCancel(context.Background())

	srv := &http.Server{
		Addr:              net.JoinHostPort(r.cfg.Server.Host, strconv.FormatUint(uint64(r.cfg.Server.Port), base10)),
		Handler:           r.mux,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readHeaderTimeout,
		WriteTimeout:      writeTimeout,
		BaseContext: func(_ net.Listener) context.Context {
			return ongoingCtx
		},
	}

	go func() {
		log.Println("Listening on ", srv.Addr)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errchan <- err
		}
	}()

	select {
	case err := <-errchan:
		stopOngoingGracefully()

		return err
	case <-ctx.Done():
		log.Println("Stopping accepting connections and shutting down ongoing requests")

		// Now shutdown our connections
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			time.Duration(shutdownPeriod)*time.Second,
		)
		defer cancel()

		err := srv.Shutdown(shutdownCtx)

		// Stop handlers operating right now
		stopOngoingGracefully()

		if err != nil {
			log.Println("Error shutting down, waiting hard shutdown period: ", err.Error())

			time.Sleep(time.Duration(hardShutdownPeriod) * time.Second)
		}

		return nil
	}
}

func (p *rickrollApp) Start(_ctx context.Context) error {
	grp, ctx := errgroup.WithContext(_ctx)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	grp.Go(func() error {
		select {
		case sig := <-sigs:
			return fmt.Errorf("stopping due to signal %d", sig)
		case <-ctx.Done():
			return nil
		}
	})

	grp.Go(func() error {
		return p.runServer(ctx)
	})

	return grp.Wait()
}
