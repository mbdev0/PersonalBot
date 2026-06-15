package models

type ProgramRepository struct {
	Id      int    `db:"id"`
	Program string `db:"program"`
}
