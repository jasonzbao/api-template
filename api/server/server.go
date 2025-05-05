package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/jasonzbao/api-template/api/dao/db"
	"github.com/jasonzbao/api-template/api/utils/config"
	dbUtils "github.com/jasonzbao/api-template/api/utils/db"
)

type Config struct {
	Env            string   `json:"env"`
	Version        string   `json:"version"`
	AllowedOrigins []string `json:"allowed_origins"`

	DBConnectionMain string `json:"db_connection_main"`
}

type Server struct {
	config   *Config
	dbClient *db.Client
}

func NewServer(ctx context.Context, configFile string, version string) *Server {
	cfg, err := config.NewConfig[Config](configFile)
	if err != nil {
		panic(err)
	}

	cfg.Version = version

	dbConn, err := dbUtils.NewDBConnection(cfg.DBConnectionMain, "api")
	if err != nil {
		panic(err)
	}

	dbClient := db.NewClient(dbConn)

	return &Server{
		config:   cfg,
		dbClient: dbClient,
	}
}

// Run is blocking
func (s *Server) Run(ctx context.Context, port string) error {
	router := s.NewRouter()
	srv := &http.Server{
		Addr:    port,
		Handler: router,
	}

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		log.Println("Got graceful shutdown message")
		if err := srv.Shutdown(ctx); err != nil {
			log.Println("Error shutting down server", err)
		}
	}()

	fmt.Println("Starting server")
	// Start the server
	if err := srv.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (s *Server) Stop() error {
	// noop for now
	return nil
}
