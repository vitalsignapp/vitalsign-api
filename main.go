package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/firestore"
	"github.com/alexandrevicenzi/go-sse"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/b3"
	"google.golang.org/api/option"

	"github.com/vitalsignapp/vitalsign-api/auth"
	"github.com/vitalsignapp/vitalsign-api/patient"
	"github.com/vitalsignapp/vitalsign-api/ward"
)

var projectID string

func init() {
	initConfig()

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
	var fsClient *firestore.Client
	var err error
	var ctx context.Context
	{
		opt := option.WithCredentialsFile("configs/firebase-credentials.json")
		ctx = context.Background()
		fsClient, err = firestore.NewClient(ctx, projectID, opt)
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}
	}
	defer fsClient.Close()

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

	r.HandleFunc("/auth", auth.Authen).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/login", auth.Login(fsClient)).Methods(http.MethodPost, http.MethodOptions)

	r.HandleFunc("/health_check", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	}).Methods(http.MethodGet)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	}).Methods(http.MethodGet)

	
	secure := r.NewRoute().Subrouter()
	secure.Use(auth.Authorization)

	secure.HandleFunc("/patient/scheduler/{patientID}", patient.NewScheduler(fsClient))
	secure.HandleFunc("/patient/{patientID}", patient.ByIDHandler(patient.NewRepoByID(fsClient)))
	secure.HandleFunc("/patient/hospital/{hospitalID}", patient.ByHospital(patient.NewRepoByHospital(fsClient)))

	secure.HandleFunc("/ward/{hospitalKey}", ward.Rooms(ward.NewRepository(fsClient)))
	secure.HandleFunc("/ward/{patientRoomKey}/patients", patient.ByRoomKeyHandler(patient.NewRepoByRoomKey(fsClient)))


	s := sse.NewServer(nil)
	defer s.Shutdown()
	r.Handle("/events/patient/{channel}", s)
	go watchingPatientData(context.Background(), fsClient, s)

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

func initConfig() {
	viper.SetDefault("port", "1323")
	viper.SetDefault("cors.allow_origin", "*")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func watchingPatientData(ctx context.Context, fsClient *firestore.Client, s *sse.Server) {
	iter := fsClient.Collection("patientData").Snapshots(ctx)
	defer iter.Stop()
	for {
		docsnap, err := iter.Next()
		if err != nil {
			log.Println("error", err)
		}
		docs, _ := docsnap.Documents.GetAll()
		patients := []patient.PatientData{}
		for _, doc := range docs {
			patient := patient.PatientData{}
			_ = doc.DataTo(&patient)
			patients = append(patients, patient)
		}

		// get all connection clients 
		clients := s.Channels()
		for _, c := range clients {
			channelName := strings.ReplaceAll(c, "/events/patient/", "")

			patientsByHospital := []patient.PatientData{}
			for _, patient := range patients {
				if channelName == patient.HospitalKey {
					patientsByHospital = append(patientsByHospital, patient)
				}
			}

			if len(patientsByHospital) > 0 {
				payload, err := json.Marshal(&patientsByHospital)
				if err != nil {
					log.Println(err)
					continue
				}

				s.SendMessage(c, sse.SimpleMessage(string(payload)))
			}
		}
	}
}