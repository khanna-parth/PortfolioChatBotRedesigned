package main

import (
	"log"
	"net/http"
	"path/filepath"
	"server/handlers"
	"server/helper"
	"server/link"
	"server/middleware"

	"github.com/go-co-op/gocron/v2"
	"github.com/joho/godotenv"
)

// var validate *validator.Validate

var Config map[string]string

func main() {
	godotenv.Load()
	// var config map[string]string

	readCfg, err := godotenv.Read("dev.env")
	if err != nil {
		log.Fatalf(".ENV file invalid/missing: %v\n", err.Error())
	}

	Config = readCfg

	connections := link.NewConnectionStore(Config)
	connections.Config = Config

	log.Printf("Running with config: %+v\n", connections.Config)

	script := filepath.Join(Config["CHAIN_PATH"], Config["GENFILE"])

	absScript, err := filepath.Abs(script)
	if err != nil || absScript == "" {
		log.Fatalf("Cannot find generator script. '%s' is not valid", script)
	}

	connections.Script = absScript
	connections.UploadsDir = Config["UPLOADS_PATH"]

	connections.Presets = helper.ExtractPresets(Config, connections.UploadsDir)

	s, err := gocron.NewScheduler()
	if err != nil {
		log.Fatal("Could not start scheduler.")
	}
	_, cronErr := s.NewJob(
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(
				gocron.NewAtTime(0,0,0),
			),
		),
		gocron.NewTask(
			func() {
				connections.Lock.Lock()
				defer connections.Lock.Unlock()
				connections.IPRequestCount = map[string]int{}
				connections.PurgeDemoes()
			},
		),
	)

	if cronErr != nil {
		log.Fatal("Could not create job on scheduler")
	}

	s.Start()
	log.Printf("Created scheduler resetter cron job")

	listen(connections)
}

func listen(connections *link.ConnectionStore) {
	mux := http.NewServeMux()

	mux.HandleFunc("/add-doc", func(w http.ResponseWriter, r *http.Request) {
		handlers.DocumentUploadHandler(w, r, connections)
	})

	mux.HandleFunc("/delete-doc", func(w http.ResponseWriter, r *http.Request) {
		handlers.DocumentDeleteHandler(w, r, connections)
	})
	
	mux.HandleFunc("/list-docs", func(w http.ResponseWriter, r *http.Request) {
		handlers.DocumentListHandler(w, r, connections)
	})
	
	mux.HandleFunc("/heartbeat", handlers.HeartbeatHandler)

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleWebSocket(w, r, connections)
	})

	mux.HandleFunc("/prompt", func(w http.ResponseWriter, r *http.Request) {
		handlers.UserPromptHandler(w, r, connections)
	})

	mux.HandleFunc("/demo", func(w http.ResponseWriter, r *http.Request) {
		handlers.PresentationPromptHandler(w, r, connections)
	})

	mux.HandleFunc("/presets", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Presets: %v\n", connections.Presets)
		presetMappings := connections.Presets
		var presets []string
		for preset := range presetMappings {
			presets = append(presets, preset)
		}

		helper.SendResponse(w, presets)
	})

	wrappedMux := middleware.ApplyCommonMiddleware(mux)
	log.Println("Starting on port 3000")
	http.ListenAndServe(":3000", wrappedMux)
}