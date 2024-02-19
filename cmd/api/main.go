package main

import (
	"context"
	"database/sql"
	"expvar"
	"flag"
	"log"
	"log/slog"
	"os"
	"runtime"
	"sync"
	"time"

	db "github.com/Hopertz/rmgmt/db/sqlc"
	"github.com/Hopertz/rmgmt/pkg/mailer"
	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}

	mailergun struct {
		domain     string
		privateKey string
		sender     string
	}
}

type application struct {
	config config
	mailer mailer.Mailer
	wg     sync.WaitGroup
	store  db.Store
}

func init() {

	var programLevel = new(slog.LevelVar)

	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel})
	slog.SetDefault(slog.New(h))

}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 5040, "API server port")
	flag.StringVar(&cfg.env, "env", os.Getenv("ENV_STAGE"), "Environment (development|Staging|production")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("DB_DSN"), "PostgreSQL DSN")

	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max ilde connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection  connections")

	flag.StringVar(&cfg.mailergun.domain, "mg-domain", os.Getenv("MAILGUN_DOMAIN"), "Mailgun-domain")
	flag.StringVar(&cfg.mailergun.privateKey, "mg-privatekey", os.Getenv("MAILGUN_PRIVATEKEY"), "Mailgun-privatekey")
	flag.StringVar(&cfg.mailergun.sender, "mg-sender", os.Getenv("MAILGUN_SENDER"), "Mailgun-sender")

	flag.Parse()

	dbConn, err := openDB(cfg)
	if err != nil {
		log.Fatal("error opening db", err)
	}

	defer dbConn.Close()

	slog.Info("database connection established")

	expvar.NewString("version").Set(version)

	expvar.Publish("goroutines", expvar.Func(func() interface{} {
		return runtime.NumGoroutine()
	}))

	expvar.Publish("database", expvar.Func(func() interface{} {
		return dbConn.Stats()
	}))

	expvar.Publish("timestamp", expvar.Func(func() interface{} {
		return time.Now().Unix()
	}))

	app := &application{
		config: cfg,
		mailer: mailer.New(cfg.mailergun.domain, cfg.mailergun.privateKey, cfg.mailergun.sender),
		store:  db.NewStore(dbConn),
	}

	err = app.serve()
	if err != nil {
		log.Fatal("error starting server", err)
	}

}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)

	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)

	if err != nil {
		return nil, err
	}

	return db, nil
}
