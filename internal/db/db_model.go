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
	if s.Ctx == nil {
		return nil, fmt.Errorf("context is nil")
	}
	if s.Tx != nil {
		return s.Tx.PrepareContext(s.Ctx, sql)
	}
	return s.DB.PrepareContext(s.Ctx, sql)
}

func (s *DbCtx) Rollback() error {
	if s.Tx != nil {
		return s.Tx.Rollback()
	}

	return fmt.Errorf("no transaction to rollback")
}

func (s *DbCtx) Commit() error {
	if s.Tx != nil {
		return s.Tx.Commit()
	}

	return fmt.Errorf("no transaction to commit")
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

func (s *DbCtx) Transactional(fn func(*DbCtx) error) (err error) {
	err = s.BeginTxn()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			s.Rollback()
			panic(p)
		} else if err != nil {
			s.Rollback()
		} else {
			err = s.Commit()
		}
	}()

	err = fn(s)
	return err
}
