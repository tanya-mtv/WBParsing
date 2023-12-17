package agent

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"parsingWB/internal/api"
	"parsingWB/internal/config"
	"parsingWB/internal/db"
	"parsingWB/internal/logger"
	"syscall"
	"time"
)

type Agent struct {
	ps  *api.ProductService
	cfg *config.Config
	log logger.Logger
}

func NewAgent(cfg *config.Config, log logger.Logger) *Agent {
	return &Agent{
		cfg: cfg,
		log: log,
	}
}

func (a *Agent) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	sqlClient := db.NewMSSQLDB(a.cfg.MSSQL)
	defer sqlClient.Close() // nolint: errcheck
	err := sqlClient.PingContext(ctx)
	if err != nil {
		a.log.Fatal("Error creating connection pool: " + err.Error())
	}
	a.log.Infof("MSSQL connected: %+v", a.cfg.MSSQL.DSN)

	sqlRepo := db.NewSQLStorage(a.log, a.cfg, sqlClient)
	a.ps = api.NewProductService(a.cfg, a.log, sqlRepo)

	// reportIntervalTicker := time.NewTicker(time.Duration(a.cfg.ReportInterval) * time.Hour)
	reportIntervalTicker := time.NewTicker(time.Duration(a.cfg.ReportInterval) * time.Second)
	defer reportIntervalTicker.Stop()

	a.log.Infof("Report interval set to: %d seconds", a.cfg.ReportInterval)

	for {
		select {
		case <-ctx.Done():
			stop()

			return nil
		case <-reportIntervalTicker.C:

			res, err := a.ps.PostPagination(ctx)

			if err != nil {
				a.log.Errorf("Error post query")
			}
			for _, elem := range res {
				products := elem.Data.Cards
				for _, val := range products {
					product := a.ps.ParsePage(ctx, val)
					if len(product.SellerPrice) != 0 {
						bProduct, err := json.Marshal(&product)
						if err != nil {
							a.log.Errorf(err.Error())
						}
						a.log.Info(string(bProduct))
					}

				}
			}

			stop()

		}
	}

}
