package postgres

import (
	"context"
	"fmt"
	"volleybot/pkg/domain/location"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type LocationPgRepository struct {
	dbpool    *pgxpool.Pool
	TableName string
}

func NewLocationRepository(dbpool *pgxpool.Pool) (pgrep LocationPgRepository, err error) {
	pgrep.TableName = "locations"
	if err != nil {
		return
	}

	pgrep.dbpool = dbpool

	return
}

func (rep *LocationPgRepository) Get(id uuid.UUID) (loc location.Location, err error) {
	sql := "SELECT location_id, location_name, location_descr, location_chat_id, location_court_count " +
		"FROM %s " +
		"WHERE location_id = $1"
	row := rep.dbpool.QueryRow(context.Background(), fmt.Sprintf(sql, rep.TableName), id)

	err = row.Scan(&loc.Id, &loc.Name, &loc.Description, &loc.ChatId, &loc.CourtCount)

	if err != nil {
		return
	}
	return
}

func (rep *LocationPgRepository) GetByName(name string) (loc location.Location, err error) {
	sql := "SELECT location_id, location_name, location_descr, location_chat_id, location_court_count " +
		"FROM %s " +
		"WHERE location_name = $1"
	row := rep.dbpool.QueryRow(context.Background(), fmt.Sprintf(sql, rep.TableName), name)

	err = row.Scan(&loc.Id, &loc.Name, &loc.Description, &loc.ChatId, &loc.CourtCount)

	if err != nil {
		return
	}
	return
}

func (rep *LocationPgRepository) Add(l location.Location) (loc location.Location, err error) {
	sql := "INSERT INTO %s " +
		"(location_id, location_name, location_descr, location_chat_id, location_court_count) " +
		"VALUES ($1, $2, $3, $4, $5) " +
		"RETURNING location_id"
	sql = fmt.Sprintf(sql, rep.TableName)

	row := rep.dbpool.QueryRow(context.Background(), sql, l.Id, l.Name, l.Description, l.ChatId, l.CourtCount)

	var LocationId uuid.UUID
	err = row.Scan(&LocationId)

	if err != nil {
		return
	}
	l.Id = LocationId
	loc = l
	return
}

func (rep *LocationPgRepository) Update(loc location.Location) (err error) {
	sql := "UPDATE %s SET " +
		"location_name = $1, location_descr = $2, location_chat_id = $3, location_court_count = $4" +
		"WHERE location_id = $5"
	sql = fmt.Sprintf(sql, rep.TableName)

	_, err = rep.dbpool.Exec(context.Background(), sql, loc.Name, loc.Description, loc.ChatId, loc.CourtCount, loc.Id)

	return
}

func (rep *LocationPgRepository) UpdateDB() (err error) {
	sql := "CREATE TABLE IF NOT EXISTS %s (" +
		"location_id UUID PRIMARY KEY, location_name VARCHAR(20), location_descr VARCHAR(100), " +
		"location_chat_id BIGINT, location_court_count INT) "
	rows, err := rep.dbpool.Query(context.Background(), fmt.Sprintf(sql, rep.TableName))

	if err != nil {
		return err
	}
	defer rows.Close()

	return err
}

type LocationConfigPgRepository struct {
	dbpool    *pgxpool.Pool
	TableName string
}

func NewLocationConfigRepository(dbpool *pgxpool.Pool) (pgrep LocationConfigPgRepository, err error) {
	pgrep.TableName = "locations_configs"
	if err != nil {
		return
	}

	pgrep.dbpool = dbpool

	return
}

func (rep *LocationConfigPgRepository) UpdateDB() (err error) {
	sql := "CREATE TABLE IF NOT EXISTS %s (" +
		"location_id UUID, service_name VARCHAR(20), location_config JSONB, " +
		"location_chat_id BIGINT, location_court_count INT) "
	rows, err := rep.dbpool.Query(context.Background(), fmt.Sprintf(sql, rep.TableName))

	if err != nil {
		return err
	}
	defer rows.Close()

	return err
}

func (rep *LocationConfigPgRepository) Add(loc location.Location, service string, config interface{}) error {
	sql := "INSERT INTO %s (location_id, service_name, location_config) VALUES ($1, $2, $3)"
	sql = fmt.Sprintf(sql, rep.TableName)

	row := rep.dbpool.QueryRow(context.Background(), sql, loc.Id, service, config)

	return row.Scan()
}

func (rep *LocationConfigPgRepository) Get(loc location.Location, service string, config interface{}) (err error) {
	sql := "SELECT location_config " +
		"FROM %s " +
		"WHERE location_id = $1 AND service_name = $2"
	row := rep.dbpool.QueryRow(context.Background(), fmt.Sprintf(sql, rep.TableName), loc.Id, service)

	err = row.Scan(config)

	if err != nil {
		return
	}
	return
}

func (rep *LocationConfigPgRepository) Update(loc location.Location, service string, config interface{}) error {
	sql := "UPDATE %s SET " +
		"location_config = $1 " +
		"WHERE location_id = $2 AND service_name = $3"
	sql = fmt.Sprintf(sql, rep.TableName)

	_, err := rep.dbpool.Exec(context.Background(), sql, config, loc.Id, service)

	return err
}
