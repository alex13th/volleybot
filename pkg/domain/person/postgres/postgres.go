package postgres

import (
	"context"
	"fmt"
	"volleybot/pkg/domain/person"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PgPerson struct {
	dbpool    *pgxpool.Pool
	TableName string
	Person    *person.Person
}

func NewPgPerson(url string, PersonId uuid.UUID) (pgp PgPerson, err error) {
	dbpool, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		return
	}

	pgp = PgPerson{dbpool: dbpool, Person: &person.Person{Id: PersonId}}

	if PersonId != uuid.Nil {
		err = pgp.Get()
	}
	return
}

func (pgp *PgPerson) Get() (err error) {
	if err != nil {
		return err
	}

	sql := "SELECT Person_id, firstname, lastname, fullname " +
		"FROM %s " +
		"WHERE Person_id = '%d'"
	row := pgp.dbpool.QueryRow(context.Background(), fmt.Sprintf(sql, pgp.TableName, pgp.Person.Id))

	var PersonId uuid.UUID
	var Firstname, Lastname, Fullname string
	err = row.Scan(&PersonId, &Firstname, &Lastname, &Fullname)

	if err != nil {
		return
	}

	pgp.Person.Id = PersonId
	pgp.Person.Firstname = Firstname
	pgp.Person.Lastname = Lastname
	pgp.Person.Fullname = Fullname

	return
}

func (pgp *PgPerson) Add() (err error) {
	sql := "INSERT INTO %s " +
		"(firstname, lastname, fullname) " +
		"VALUES ('%s',' %s', '%s') " +
		"RETURNING person_id"
	sql = fmt.Sprintf(sql,
		pgp.TableName,
		pgp.Person.Firstname,
		pgp.Person.Lastname,
		pgp.Person.Fullname)

	row := pgp.dbpool.QueryRow(context.Background(), sql)

	var PersonId uuid.UUID
	err = row.Scan(&PersonId)

	if err != nil {
		return
	}

	pgp.Person.Id = PersonId

	return
}

func (pgp *PgPerson) Update() (err error) {
	sql := "UPDATE %s " +
		"SET firstname = '%s', lastname = '%s', fullname = '%s' " +
		"WHERE person_id = '%d'"
	sql = fmt.Sprintf(sql,
		pgp.TableName,
		pgp.Person.Firstname,
		pgp.Person.Lastname,
		pgp.Person.Fullname,
		pgp.Person.Id)

	_, err = pgp.dbpool.Exec(context.Background(), sql)

	return
}

func (pgp *PgPerson) UpdateDB() (err error) {
	if err != nil {
		return err
	}

	sql := "CREATE TABLE IF NOT EXISTS %s " +
		"(person_id UUID PRIMARY KEY, " +
		"firstname varchar(20), lastname varchar(20), fullname varchar(60))"
	rows, err := pgp.dbpool.Query(context.Background(), fmt.Sprintf(sql, pgp.TableName))

	if err != nil {
		return err
	}
	defer rows.Close()

	return err
}
