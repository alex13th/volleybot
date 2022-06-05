package person

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

func NewPgRepository(dbpool *pgxpool.Pool) (pgrep PgRepository, err error) {
	pgrep.TableName = "persons"
	if err != nil {
		return
	}

	pgrep.dbpool = dbpool

	return
}

func (rep *PgRepository) Get(pid uuid.UUID) (p Person, err error) {
	sql := "SELECT person_id, telegram_id, firstname, lastname, fullname, roles " +
		"FROM %s " +
		"WHERE person_id = $1"
	row := rep.dbpool.QueryRow(context.Background(), fmt.Sprintf(sql, rep.TableName), pid)

	err = row.Scan(&p.Id, &p.TelegramId, &p.Firstname, &p.Lastname, &p.Fullname, &p.Roles)
	return
}

func (rep *PgRepository) GetByTelegramId(tid int) (p Person, err error) {
	sql := "SELECT person_id, telegram_id, firstname, lastname, fullname, roles " +
		"FROM %s " +
		"WHERE telegram_id = $1"
	row := rep.dbpool.QueryRow(context.Background(), fmt.Sprintf(sql, rep.TableName), tid)

	err = row.Scan(&p.Id, &p.TelegramId, &p.Firstname, &p.Lastname, &p.Fullname, &p.Roles)
	return
}

func (rep *PgRepository) Add(p Person) (per Person, err error) {
	sql := "INSERT INTO %s " +
		"(person_id, telegram_id, firstname, lastname, fullname, roles) " +
		"VALUES ($1, $2, $3, $4, $5, $6) " +
		"RETURNING person_id"
	sql = fmt.Sprintf(sql, rep.TableName)

	row := rep.dbpool.QueryRow(context.Background(), sql,
		p.Id, p.TelegramId, p.Firstname, p.Lastname, p.Fullname, p.Roles)

	err = row.Scan(&p.Id)
	return p, err
}

func (rep *PgRepository) Update(p Person) (err error) {
	sql := "UPDATE %s SET " +
		"telegram_id = $1, firstname = $2, lastname = $3, fullname = $4 " +
		"roles = $5 " +
		"WHERE person_id = $6"
	sql = fmt.Sprintf(sql, rep.TableName)

	rows, err := rep.dbpool.Query(context.Background(), sql,
		p.TelegramId, p.Firstname, p.Lastname, p.Fullname, p.Roles, p.Id)

	if err != nil {
		return
	}
	defer rows.Close()
	return
}

func (rep *PgRepository) UpdateDB() (err error) {
	sql := "CREATE TABLE IF NOT EXISTS %s " +
		"(person_id UUID PRIMARY KEY, telegram_id int, " +
		"firstname varchar(20), lastname varchar(20), fullname varchar(60), " +
		"roles varchar(250))"
	rows, err := rep.dbpool.Query(context.Background(), fmt.Sprintf(sql, rep.TableName))

	if err != nil {
		return
	}
	defer rows.Close()

	return
}
