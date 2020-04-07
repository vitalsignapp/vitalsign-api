package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"cloud.google.com/go/compute/metadata"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/b3"
)

var projectID string

func init() {
	viper.SetDefault("port", "1323")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	log.SetFlags(0)

	projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		projectID, _ = metadata.ProjectID()
	}
	if projectID == "" {
		log.Println("Could not determine Google Cloud Project. Running without log correlation. For local use set the GOOGLE_CLOUD_PROJECT environment variable.")
	}
}

func main() {

	r := mux.NewRouter()
	r.Use(mux.CORSMethodMiddleware(r))

	// all routes required headers
	r.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", viper.GetString("cors.allow_origin"))
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Header().Set("X-XSS-Protection", "1; mode=block")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("Strict-Transport-Security", "max-age=604800; includeSubDomains; preload")
			handler.ServeHTTP(w, r)
		})
	})

	// additional headers for preflight request
	r.PathPrefix("/").Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Max-Age", viper.GetString("cors.max_age"))
		w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
	})

	srv := &http.Server{
		Handler: &ochttp.Handler{
			Handler:     r,
			Propagation: &b3.HTTPFormat{},
		},
		Addr:         "0.0.0.0:" + viper.GetString("port"),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		log.Printf("serve on %s\n", ":"+viper.GetString("port"))
		log.Printf("%s", srv.ListenAndServe())
	}()

	gracefulshutdown(srv)
}

func gracefulshutdown(srv *http.Server) {
	sigterm := make(chan os.Signal)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	<-sigterm

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("%s", err.Error())
	}
}
