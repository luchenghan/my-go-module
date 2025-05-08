package service

import "mymodule/sql"

type Service struct {
	*sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db}
}
