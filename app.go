package main

import (
	"context"
	"file-service/api"
	"file-service/config"
	dbfile "file-service/database/file"
	"file-service/database/object"
	"file-service/service/file"
	"fmt"
	"io/fs"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sunshineOfficial/golib/db"
	"github.com/sunshineOfficial/golib/gohttp/goserver"
	"github.com/sunshineOfficial/golib/golog"
)

const (
	serviceName = "file-service"
	dbTimeout   = 15 * time.Second
)

type App struct {
	/* main */
	mainCtx  context.Context
	log      golog.Logger
	settings config.Settings

	/* http */
	server goserver.Server

	/* db */
	postgres *sqlx.DB
	minio    *minio.Client

	/* services */
	fileService *file.Service
}

func NewApp(mainCtx context.Context, log golog.Logger, settings config.Settings) *App {
	return &App{
		mainCtx:  mainCtx,
		log:      log,
		settings: settings,
	}
}

func (a *App) InitDatabases(fs fs.FS, path string) (err error) {
	postgresCtx, cancelPostgresCtx := context.WithTimeout(a.mainCtx, dbTimeout)
	defer cancelPostgresCtx()

	a.postgres, err = db.NewPgx(postgresCtx, a.settings.Databases.Postgres)
	if err != nil {
		return fmt.Errorf("init postgres: %w", err)
	}

	err = db.Migrate(fs, a.log, a.postgres, path)
	if err != nil {
		return fmt.Errorf("migrate postgres: %w", err)
	}

	a.minio, err = minio.New(
		a.settings.Databases.Minio.Endpoint,
		&minio.Options{
			Creds:  credentials.NewStaticV4(a.settings.Databases.Minio.User, a.settings.Databases.Minio.Password, ""),
			Secure: a.settings.Databases.Minio.UseSSL,
		},
	)
	if err != nil {
		return fmt.Errorf("init minio: %w", err)
	}

	return nil
}

func (a *App) InitServices() error {
	documentCtx, cancelDocumentCtx := context.WithTimeout(a.mainCtx, dbTimeout)
	defer cancelDocumentCtx()

	documentStorage, err := object.NewMinio(documentCtx, a.minio, string(file.BucketDocuments), a.settings.Databases.Minio.Host)
	if err != nil {
		return fmt.Errorf("init document storage: %w", err)
	}

	imageCtx, cancelImageCtx := context.WithTimeout(a.mainCtx, dbTimeout)
	defer cancelImageCtx()

	imageStorage, err := object.NewMinio(imageCtx, a.minio, string(file.BucketImages), a.settings.Databases.Minio.Host)
	if err != nil {
		return fmt.Errorf("init image storage: %w", err)
	}

	bucketStorages := map[file.Bucket]*object.Minio{
		file.BucketDocuments: documentStorage,
		file.BucketImages:    imageStorage,
	}

	fileRepository := dbfile.NewPostgres(a.postgres)

	a.fileService = file.NewService(bucketStorages, fileRepository)

	return nil
}

func (a *App) InitServer() {
	sb := api.NewServerBuilder(a.mainCtx, a.log, a.settings)
	sb.AddDebug()
	sb.AddFiles(a.fileService)

	a.server = sb.Build()
}

func (a *App) Start() {
	a.server.Start()
}

func (a *App) Stop(_ context.Context) {
	a.server.Stop()

	err := a.postgres.Close()
	if err != nil {
		a.log.Error("failed to close postgres connection")
	}
}
