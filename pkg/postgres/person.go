package postgres

import (
	"context"
	"fmt"
	"volleybot/pkg/domain/person"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PersonPgRepository struct {
	dbpool            *pgxpool.Pool
	TableName         string
	RolesTableName    string
	SettingsTableName string
}

func NewPersonPgRepository(dbpool *pgxpool.Pool) (pgrep PersonPgRepository, err error) {
	pgrep.TableName = "persons"
	pgrep.RolesTableName = "person_roles"
	pgrep.SettingsTableName = "person_params"
	if err != nil {
		return
	}

	pgrep.dbpool = dbpool

	return
}

func (rep *PersonPgRepository) Get(pid uuid.UUID) (p person.Person, err error) {
	p = person.NewPerson("")
	sql := "SELECT person_id, telegram_id, firstname, lastname, fullname, sex " +
		"FROM %s " +
		"WHERE person_id = $1"
	row := rep.dbpool.QueryRow(context.Background(), fmt.Sprintf(sql, rep.TableName), pid)
	err = row.Scan(&p.Id, &p.TelegramId, &p.Firstname, &p.Lastname, &p.Fullname, &p.Sex)
	if err != nil {
		return
	}
	p.LocationRoles, err = rep.GetRoles(p.Id)
	if err != nil {
		return
	}
	p.Settings, err = rep.GetSettings(p.Id)
	return
}

func (rep *PersonPgRepository) GetByTelegramId(tid int) (p person.Person, err error) {
	p = person.NewPerson("")
	sql := "SELECT person_id, telegram_id, firstname, lastname, fullname, sex " +
		"FROM %s " +
		"WHERE telegram_id = $1"
	row := rep.dbpool.QueryRow(context.Background(), fmt.Sprintf(sql, rep.TableName), tid)

	err = row.Scan(&p.Id, &p.TelegramId, &p.Firstname, &p.Lastname, &p.Fullname, &p.Sex)
	if err != nil {
		if err.Error() == "no rows in result set" {
			err = person.ErrPersonNotFound
			return
		}
		return
	}
	p.LocationRoles, err = rep.GetRoles(p.Id)
	if err != nil {
		return
	}
	p.Settings, err = rep.GetSettings(p.Id)
	return
}

func (rep *PersonPgRepository) Add(p person.Person) (per person.Person, err error) {
	sql := "INSERT INTO %s " +
		"(person_id, telegram_id, firstname, lastname, fullname, sex) " +
		"VALUES ($1, $2, $3, $4, $5, $6) " +
		"RETURNING person_id"
	sql = fmt.Sprintf(sql, rep.TableName)

	row := rep.dbpool.QueryRow(context.Background(), sql,
		p.Id, p.TelegramId, p.Firstname, p.Lastname, p.Fullname, p.Sex)

	err = row.Scan(&p.Id)
	return p, err
}

func (rep *PersonPgRepository) Update(p person.Person) (err error) {
	sql := "UPDATE %s SET " +
		"telegram_id = $1, firstname = $2, lastname = $3, fullname = $4, sex = $5 " +
		"WHERE person_id = $6"
	sql = fmt.Sprintf(sql, rep.TableName)

	rows, err := rep.dbpool.Query(context.Background(), sql,
		p.TelegramId, p.Firstname, p.Lastname, p.Fullname, p.Sex, p.Id)

	if err != nil {
		return
	}
	err = rep.UpdateParams(p)
	defer rows.Close()
	return
}

func (rep *PersonPgRepository) UpdateParams(p person.Person) (err error) {
	usql := "UPDATE %s SET param_value = $3 " +
		"WHERE person_id = $1 AND param_name =$2"
	usql = fmt.Sprintf(usql, rep.SettingsTableName)

	isql := "INSERT INTO %s " +
		"(person_id, param_name, param_value) " +
		"VALUES ($1, $2, $3) "
	isql = fmt.Sprintf(isql, rep.SettingsTableName)

	for param, val := range p.Settings {
		rres, _ := rep.dbpool.Exec(context.Background(), usql, p.Id, param, val)

		if rres.RowsAffected() < 1 {
			_, err = rep.dbpool.Exec(context.Background(), isql, p.Id, param, val)
		}
	}
	return
}

func (rep *PersonPgRepository) UpdateDB() (err error) {
	sql := "CREATE TABLE IF NOT EXISTS %s " +
		"(person_id UUID PRIMARY KEY, telegram_id BIGINT, " +
		"firstname VARCHAR(20), lastname VARCHAR(20), fullname VARCHAR(60), " +
		"sex INT, level INT, " +
		"roles varchar(250));"
	sql += "CREATE TABLE IF NOT EXISTS %s " +
		"(person_id UUID, location_id UUID, role VARCHAR(30));"
	sql += "CREATE TABLE IF NOT EXISTS %s " +
		"(person_id UUID, param_name VARCHAR(20), param_value VARCHAR(250));"
	sql = fmt.Sprintf(sql, rep.TableName, rep.RolesTableName, rep.SettingsTableName)
	_, err = rep.dbpool.Exec(context.Background(), sql)

	if err != nil {
		return
	}

	return
}

func (rep *PersonPgRepository) GetRoles(pid uuid.UUID) (pmap map[uuid.UUID][]string, err error) {
	sql := "SELECT roles.person_id, roles.location_id, roles.role " +
		"FROM %s AS roles " +
		"INNER JOIN %s AS p ON roles.person_id = p.person_id " +
		"WHERE roles.person_id = $1;"
	sql = fmt.Sprintf(sql, rep.RolesTableName, rep.TableName)
	rows, err := rep.dbpool.Query(context.Background(), sql, pid)
	pmap = make(map[uuid.UUID][]string)
	var (
		PersonId, LocationId uuid.UUID
		Role                 string
	)

	for rows.Next() {
		rows.Scan(&PersonId, &LocationId, &Role)
		pmap[LocationId] = append(pmap[LocationId], Role)
	}
	return
}

func (rep *PersonPgRepository) GetSettings(pid uuid.UUID) (pmap map[string]string, err error) {
	sql := "SELECT params.param_name, params.param_value " +
		"FROM %s AS params " +
		"INNER JOIN %s AS p ON params.person_id = p.person_id " +
		"WHERE params.person_id = $1;"
	sql = fmt.Sprintf(sql, rep.SettingsTableName, rep.TableName)
	rows, err := rep.dbpool.Query(context.Background(), sql, pid)
	pmap = make(map[string]string)
	var Param, Value string

	for rows.Next() {
		rows.Scan(&Param, &Value)
		pmap[Param] = Value
	}
	return
}
