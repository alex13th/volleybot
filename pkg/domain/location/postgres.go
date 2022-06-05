package location

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PgRepository struct {
	dbpool    *pgxpool.Pool
	TableName string
}

func NewPgRepository(url string) (pgrep PgRepository, err error) {
	pgrep.TableName = "locations"
	dbpool, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		return
	}

	pgrep.dbpool = dbpool

	return
}

func (rep *PgRepository) Get(id uuid.UUID) (loc Location, err error) {
	sql := "SELECT location_id, location_name, location_descr, location_chat_id " +
		"FROM %s " +
		"WHERE location_id = $1"
	row := rep.dbpool.QueryRow(context.Background(), fmt.Sprintf(sql, rep.TableName), id)

	err = row.Scan(&loc.Id, &loc.Name, &loc.Description, &loc.ChatId)

	if err != nil {
		return
	}
	return
}

func (rep *PgRepository) GetByName(name string) (loc Location, err error) {
	sql := "SELECT location_id, location_name, location_descr, location_chat_id " +
		"FROM %s " +
		"WHERE location_name = $1"
	row := rep.dbpool.QueryRow(context.Background(), fmt.Sprintf(sql, rep.TableName), name)

	err = row.Scan(&loc.Id, &loc.Name, &loc.Description, &loc.ChatId)

	if err != nil {
		return
	}
	return
}

func (rep *PgRepository) Add(l Location) (loc Location, err error) {
	sql := "INSERT INTO %s " +
		"(location_id, location_name, location_descr, location_chat_id) " +
		"VALUES ($1, $2, $3, $4) " +
		"RETURNING location_id"
	sql = fmt.Sprintf(sql, rep.TableName)

	row := rep.dbpool.QueryRow(context.Background(), sql, l.Id, l.Name, l.Description, l.ChatId)

	var LocationId uuid.UUID
	err = row.Scan(&LocationId)

	if err != nil {
		return
	}
	l.Id = LocationId
	loc = l
	return
}

func (rep *PgRepository) Update(loc Location) (err error) {
	sql := "UPDATE %s SET " +
		"location_name = $1, location_descr = $2, location_chat_id = $3" +
		"WHERE location_id = $4"
	sql = fmt.Sprintf(sql, rep.TableName)

	_, err = rep.dbpool.Exec(context.Background(), sql, loc.Name, loc.Description, loc.ChatId, loc.Id)

	return
}

func (rep *PgRepository) UpdateDB() (err error) {
	sql := "CREATE TABLE IF NOT EXISTS %s " +
		"(location_id UUID PRIMARY KEY, location_name varchar(20), location_descr varchar(100), location_chat_id int) "
	rows, err := rep.dbpool.Query(context.Background(), fmt.Sprintf(sql, rep.TableName))

	if err != nil {
		return err
	}
	defer rows.Close()

	return err
}
