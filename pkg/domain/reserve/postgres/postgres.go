package postgres

import (
	"context"
	"fmt"
	"time"
	"volleybot/pkg/domain/reserve"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PgReserve struct {
	dbpool    *pgxpool.Pool
	TableName string
	Reserve   *reserve.Reserve
}

func NewPgReserve(url string, ReserveId uuid.UUID) (pgr PgReserve, err error) {
	dbpool, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		return
	}

	pgr = PgReserve{dbpool: dbpool, Reserve: &reserve.Reserve{Id: ReserveId}}

	if ReserveId != uuid.Nil {
		err = pgr.Get()
	}
	return
}

func (pgr *PgReserve) Get() (err error) {
	if err != nil {
		return
	}

	sql := "SELECT reserve_id, start_time, end_time " +
		"FROM %s " +
		"WHERE reserve_id = '%d'"
	row := pgr.dbpool.QueryRow(context.Background(), fmt.Sprintf(sql, pgr.TableName, pgr.Reserve.Id))

	var ReserveId uuid.UUID
	var StartTime, EndTime time.Time
	err = row.Scan(&ReserveId, &StartTime, &EndTime)
	if err != nil {
		return
	}

	pgr.Reserve.Id = ReserveId
	pgr.Reserve.StartTime = StartTime
	pgr.Reserve.EndTime = EndTime

	return
}

func (pgr *PgReserve) UpdateDB() (err error) {
	if err != nil {
		return
	}

	sql := "CREATE TABLE IF NOT EXISTS %s " +
		"(reserve_id UUID PRIMARY KEY, " +
		"start_time timestamp, end_time timestamp)"
	rows, err := pgr.dbpool.Query(context.Background(), fmt.Sprintf(sql, pgr.TableName))

	if err != nil {
		return
	}
	defer rows.Close()

	return
}
