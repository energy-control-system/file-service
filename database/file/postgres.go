package file

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Postgres struct {
	db *sqlx.DB
}

func NewPostgres(db *sqlx.DB) *Postgres {
	return &Postgres{
		db: db,
	}
}

//go:embed sql/add.sql
var addSQL string

func (p *Postgres) Add(ctx context.Context, f File) (id int, err error) {
	rows, err := p.db.NamedQueryContext(ctx, addSQL, f)
	if err != nil {
		return 0, fmt.Errorf("p.db.NamedQueryContext: %w", err)
	}
	defer func() {
		rowsErr := rows.Close()
		if rowsErr != nil {
			err = errors.Join(err, rowsErr)
		}
	}()

	if !rows.Next() {
		return 0, errors.New("rows.Next == false")
	}

	err = rows.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("rows.Scan: %w", err)
	}

	return id, err
}

//go:embed sql/get_by_id.sql
var getByIDSQL string

func (p *Postgres) GetByID(ctx context.Context, id int) (f File, err error) {
	err = p.db.GetContext(ctx, &f, getByIDSQL, id)
	if err != nil {
		return File{}, fmt.Errorf("p.db.GetContext: %w", err)
	}

	return f, nil
}
