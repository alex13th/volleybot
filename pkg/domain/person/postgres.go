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

func NewPgRepository(url string) (pgrep PgRepository, err error) {
	pgrep.TableName = "persons"
	dbpool, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		return
	}

	pgrep.dbpool = dbpool

	return
}

func (rep *PgRepository) Get(pid uuid.UUID) (p Person, err error) {
	sql := "SELECT person_id, telegram_id, firstname, lastname, fullname " +
		"FROM %s " +
		"WHERE person_id = $1"
	row := rep.dbpool.QueryRow(context.Background(), fmt.Sprintf(sql, rep.TableName), pid)

	var PersonId uuid.UUID
	var TelegramId int
	var Firstname, Lastname, Fullname string
	err = row.Scan(&PersonId, &TelegramId, &Firstname, &Lastname, &Fullname)

	if err != nil {
		return
	}

	p = Person{
		Id:         PersonId,
		TelegramId: TelegramId,
		Firstname:  Firstname,
		Lastname:   Lastname,
		Fullname:   Fullname,
	}

	return
}

func (rep *PgRepository) GetByTelegramId(tid int) (p Person, err error) {
	sql := "SELECT person_id, telegram_id, firstname, lastname, fullname " +
		"FROM %s " +
		"WHERE telegram_id = $1"
	row := rep.dbpool.QueryRow(context.Background(), fmt.Sprintf(sql, rep.TableName), tid)

	var PersonId uuid.UUID
	var TelegramId int
	var Firstname, Lastname, Fullname string
	err = row.Scan(&PersonId, &TelegramId, &Firstname, &Lastname, &Fullname)

	if err != nil {
		return
	}

	p = Person{
		Id:         PersonId,
		TelegramId: TelegramId,
		Firstname:  Firstname,
		Lastname:   Lastname,
		Fullname:   Fullname,
	}

	return
}

func (rep *PgRepository) Add(p Person) (per Person, err error) {
	sql := "INSERT INTO %s " +
		"(person_id, telegram_id, firstname, lastname, fullname) " +
		"VALUES ($1, $2, $3, $4, $5) " +
		"RETURNING person_id"
	sql = fmt.Sprintf(sql, rep.TableName)

	row := rep.dbpool.QueryRow(context.Background(), sql, p.Id, p.TelegramId, p.Firstname, p.Lastname, p.Fullname)

	var PersonId uuid.UUID
	err = row.Scan(&PersonId)

	if err != nil {
		return
	}

	p.Id = PersonId
	per = p

	return
}

func (rep *PgRepository) Update(p Person) (err error) {
	sql := "UPDATE %s " +
		"SET telegram_id = '%d', firstname = '%s', lastname = '%s', fullname = '%s' " +
		"WHERE person_id = '%d'"
	sql = fmt.Sprintf(sql, rep.TableName, p.TelegramId, p.Firstname, p.Lastname, p.Fullname, p.Id)

	_, err = rep.dbpool.Exec(context.Background(), sql)

	return
}

func (rep *PgRepository) UpdateDB() (err error) {
	sql := "CREATE TABLE IF NOT EXISTS %s " +
		"(person_id UUID PRIMARY KEY, telegram_id int, " +
		"firstname varchar(20), lastname varchar(20), fullname varchar(60))"
	rows, err := rep.dbpool.Query(context.Background(), fmt.Sprintf(sql, rep.TableName))

	if err != nil {
		return err
	}
	defer rows.Close()

	return err
}
