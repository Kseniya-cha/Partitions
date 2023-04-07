package cmd

import (
	"context"
	"os"

	"github.com/Kseniya-cha/Partitions/internal/devicetestingresults"
	repoRes "github.com/Kseniya-cha/Partitions/internal/devicetestingresults/repository"
	"github.com/Kseniya-cha/Partitions/internal/monitoringcycle"
	repoGlob "github.com/Kseniya-cha/Partitions/internal/monitoringcycle/repository"
	"github.com/Kseniya-cha/Partitions/pkg/config"
	"github.com/Kseniya-cha/Partitions/pkg/database/postgresql"
	"github.com/Kseniya-cha/Partitions/pkg/logger"
	"go.uber.org/zap"
)

type App interface {
	Run(context.Context)
	GracefulShutdown(cancel context.CancelFunc)
}

// app - прототип приложения
type app struct {
	cfg *config.Config
	log *zap.Logger
	db  postgresql.IDB

	sigChan chan os.Signal

	repoGlobal monitoringcycle.Repository
	repoResult devicetestingresults.Repository
}

// NewApp инициализирует прототип приложения
func NewApp(ctx context.Context, cfg *config.Config) (App, error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	log := logger.NewLogger(cfg)

	db, err := postgresql.NewDB(ctx, cfg.Database, log)
	if err != nil {
		return nil, err
	}

	sigChan := make(chan os.Signal, 1)

	return &app{
		cfg: cfg,
		db:  db,
		log: log,

		sigChan: sigChan,

		// repoResult: repos,
		repoResult: repoRes.NewRepository(db, log),
		repoGlobal: repoGlob.NewRepository(db, log),
	}, nil
}
