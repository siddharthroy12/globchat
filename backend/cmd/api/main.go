package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"globechat.live/internal/models"
)

type config struct {
	env            string
	port           int
	dsn            string
	googleClientId string
	mediaDir       string
}

type application struct {
	logger       *slog.Logger
	db           *sql.DB
	config       config
	userModel    models.UserModel
	sessionModel models.SessionModel
	threadModel  models.ThreadModel
	messageModel models.MessageModel
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.dsn)

	if err != nil {
		return nil, err
	}

	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()

	err = db.PingContext(ctx)

	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "the environment the api server is running on")
	flag.StringVar(&cfg.dsn, "dsn", "", "dsn string to connect to postgres DB")
	flag.StringVar(&cfg.googleClientId, "gclientid", "", "google client id for oauth")
	flag.StringVar(&cfg.mediaDir, "mediadir", "", "directory to store uploaded media files")
	flag.Parse()

	// Ensure the profile pictures directory exists
	if err := os.MkdirAll(cfg.mediaDir, 0755); err != nil {
		fmt.Printf("failed to create media directory: %s", err.Error())
		os.Exit(1)
	}

	if strings.TrimSpace(cfg.dsn) == "" {
		fmt.Println("no dsn provided")
		os.Exit(1)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(cfg)

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info("database connection pool establised")

	app := application{
		logger: logger,
		db:     db,
		config: cfg,
		userModel: models.UserModel{
			DB: db,
		},
		sessionModel: models.SessionModel{
			DB: db,
		},
		threadModel: models.ThreadModel{
			DB: db,
		},
		messageModel: models.MessageModel{
			DB: db,
		},
	}

	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}
	logger.Info("starting server", "addr", srv.Addr, "env", cfg.env)
	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)

	fmt.Print("Hello world")

}
