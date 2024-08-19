package main

import (
	"context"
	"embed"
	"net/http"
	"os"
	"os/signal"

	_ "github.com/a-h/templ"
	"github.com/fiatjaf/khatru"
	"github.com/kelseyhightower/envconfig"
	"github.com/nbd-wtf/go-nostr"
	"github.com/pelletier/go-toml"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
)

type Settings struct {
	Port           string `envconfig:"PORT" default:"3002"`
	NostrSecretKey string `envconfig:"NOSTR_SECRET_KEY"`
	NostrPublicKey string
}

var (
	s          Settings
	log        = zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
	apps       = make(map[string]AppDefinition)
	categories = []string{
		"audio",
		"career",
		"community",
		"curation",
		"direct-message",
		"file-sharing",
		"group-chat",
		"meatspace",
		"marketplaces",
		"money",
		"music",
		"privacy",
		"signers",
		"social",
		"streaming",
		"tools",
		"video",
		"writing",
		"gaming",
	}
	platforms = []string{
		"All",     // "all"
		"iOS",     // "ios"
		"Android", // "android"
		"Web",     // "web"
		"Desktop", // "desktop"
	}
)

//go:embed static/*
var static embed.FS

//go:embed apps.toml
var appsToml []byte

func main() {
	err := envconfig.Process("", &s)
	if err != nil {
		log.Fatal().Err(err).Msg("couldn't process envconfig")
		return
	}
	s.NostrPublicKey, err = nostr.GetPublicKey(s.NostrSecretKey)
	if err != nil {
		log.Fatal().Err(err).Msg("invalid secret key given")
		return
	}

	// load all apps
	if err := toml.Unmarshal(appsToml, &apps); err != nil {
		log.Fatal().Err(err).Msg("failed to load apps")
		return
	}

	// setup http handler
	mux := http.NewServeMux()
	var handler http.Handler = mux

	// nostr relay
	if s.NostrSecretKey != "" {
		relay := khatru.NewRelay()
		relay.RejectEvent = append(relay.RejectEvent,
			func(ctx context.Context, event *nostr.Event) (reject bool, msg string) {
				return true, "this relay is read-only"
			},
		)
		relay.QueryEvents = append(relay.QueryEvents, handleRelayQuery)
		mux = relay.Router()
		handler = relay
	}
	// ~

	// routes
	mux.Handle("/static/", http.FileServer(http.FS(static)))
	mux.HandleFunc("/apps/{app}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/"+r.PathValue("app"), 301)
	})
	mux.HandleFunc("/{app}", handleAppPage)
	mux.HandleFunc("/{$}", handleHomePage)

	log.Printf("listening at http://0.0.0.0:%s and signing with public key %s", s.Port, s.NostrPublicKey)
	server := &http.Server{Addr: "0.0.0.0:" + s.Port, Handler: cors.AllowAll().Handler(handler)}
	defer server.Close()
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("")
		}
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc

	log.Info().Msg("exiting")
}
