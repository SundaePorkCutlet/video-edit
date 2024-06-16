package db

import (
	"context"
	"database/sql"
	"fmt"
)

type DbCtx struct {
	DB  *sql.DB
	Tx  *sql.Tx
	Ctx context.Context
}

func (s *DbCtx) CreatePrepareStmt(sql string) (*sql.Stmt, error) {
	if s.DB == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	return s.DB.PrepareContext(s.Ctx, sql)
}

func (s *DbCtx) Rollback() error {
	if s.Tx != nil {
		return s.Tx.Rollback()
	}

	return nil
}

func (s *DbCtx) Commit() error {
	if s.Tx != nil {
		return s.Tx.Commit()
	}

	return nil
}

func (s *DbCtx) BeginTxn() error {
	if s.DB == nil {
		return fmt.Errorf("database connection is nil")
	}
	tx, err := s.DB.BeginTx(s.Ctx, nil)
	if err != nil {
		return err
	}
	s.Tx = tx
	return nil
}
