package db

import (
	"database/sql"
	"parsingWB/internal/config"
	"parsingWB/internal/logger"
	"parsingWB/internal/models"
	"time"
)

type MSSQLStorage struct {
	log       logger.Logger
	cfg       *config.Config
	sqlClient *sql.DB
}

func NewSQLStorage(log logger.Logger, cfg *config.Config, sqlClient *sql.DB) *MSSQLStorage {
	return &MSSQLStorage{log: log, cfg: cfg, sqlClient: sqlClient}
}

func (s *MSSQLStorage) InsertData(product models.Product) error {
	tx, err := s.sqlClient.Begin()
	if err != nil {
		return err
	}
	defer func() {
		err = tx.Rollback()
		if err != nil {
			return
		}
	}()

	timeNow := time.Now().Format("2006-01-02")
	stmtProd, err := tx.Prepare("INSERT INTO wbProduct (ModifiedDate, nmID, name, price, barcode) OUTPUT Inserted.ID" +
		" VALUES(?,?,?,?,?)")
	if err != nil {
		return err
	}
	defer stmtProd.Close()

	stmtPrice, err := tx.Prepare("INSERT INTO wbSellerPrice (ModifiedDate, nmID , seller, price)" +
		" VALUES(?,?,?,?)")
	if err != nil {
		return err
	}
	defer stmtPrice.Close()

	var id int
	err = stmtProd.QueryRow(timeNow, product.NmID, product.Name, product.Price, product.Barcode).Scan(&id)
	if err != nil {
		s.log.Errorf("Cant insert data to  wbProduct", err.Error())
		return err
	}

	for _, v := range product.SellerPrice {
		_, err := stmtPrice.Exec(timeNow, id, v.Saller, v.Price)
		if err != nil {
			s.log.Errorf("Cant insert data to  SellerPrice", err.Error())
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		s.log.Errorf("Transaction fail ", err.Error())
		return err
	}

	return nil
}
