package reserve

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"volleybot/pkg/domain/person"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

func AddWhereParam(wsql *string, params *[]interface{}, param interface{}, cond string) {
	*params = append(*params, param)
	if len(*params) > 1 {
		*wsql += " AND"
	}
	*wsql += " " + cond + " $" + strconv.Itoa(len(*params))
}

func NewPgRepository(url string) (rep PgRepository, err error) {
	rep.TableName = "reserves"
	rep.PersonsTableName = "persons"
	rep.LocationsTableName = "locations"
	rep.PlayersTableName = "reserve_players"
	rep.ViewName = "vw_reserves"
	rep.PlayerSpName = "sp_reserve_player_update"

	dbpool, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		return
	}

	rep.dbpool = dbpool
	return
}

type PgRepository struct {
	dbpool             *pgxpool.Pool
	TableName          string
	PersonsTableName   string
	LocationsTableName string
	PlayersTableName   string
	ViewName           string
	PlayerSpName       string
}

func (rep *PgRepository) UpdateDB() (err error) {
	sql := "CREATE TABLE IF NOT EXISTS %[1]s " +
		"(reserve_id UUID PRIMARY KEY, person_id UUID, location_id UUID, " +
		"start_time TIMESTAMP, end_time TIMESTAMP, price INT, " +
		"min_level INT, court_count INT, max_players INT, " +
		"ordered BOOL, approved BOOL, canceled BOOL);"

	pl_sql := "CREATE TABLE IF NOT EXISTS %[2]s (reserve_id UUID, person_id UUID, count INT);"
	vw_sql := "CREATE OR REPLACE VIEW %[4]s AS " +
		"SELECT reserve_id, r.person_id AS person_id, start_time, end_time, " +
		"price, min_level, court_count, max_players, ordered, approved, canceled, " +
		"telegram_id, firstname, lastname, fullname, " +
		"l.location_id AS location_id, location_name " +
		"FROM %[1]s AS r " +
		"INNER JOIN %[3]s AS p ON r.person_id = p.person_id " +
		"LEFT OUTER JOIN %[6]s AS l ON r.location_id = l.location_id; "
	sp_sql := "CREATE OR REPLACE PROCEDURE " +
		"%[5]s(res_id UUID, per_id UUID, c INT) " +
		"LANGUAGE plpgsql AS $$ " +
		"DECLARE cur_count INT;\n" +
		"BEGIN\n" +
		"SELECT SUM(count) INTO cur_count " +
		"FROM %[2]s WHERE reserve_id = res_id AND person_id = per_id; " +
		"CASE\n" +
		"WHEN c = 0 THEN\n" +
		"DELETE FROM %[2]s WHERE reserve_id = res_id AND person_id = per_id;\n" +
		"WHEN cur_count > 0 THEN\n" +
		"UPDATE %[2]s SET count = c WHERE reserve_id = res_id AND person_id = per_id;\n" +
		"ELSE\n" +
		"INSERT INTO %[2]s (reserve_id, person_id, count) VALUES (res_id, per_id, c);\n" +
		"END CASE;\n" +
		"END;$$;"
	sql = fmt.Sprintf(sql+pl_sql+vw_sql+sp_sql, rep.TableName, rep.PlayersTableName, rep.PersonsTableName,
		rep.ViewName, rep.PlayerSpName, rep.LocationsTableName)
	_, err = rep.dbpool.Exec(context.Background(), sql)

	if err != nil {
		return
	}

	return
}

func (rep *PgRepository) GetPlayers(rid uuid.UUID) (pmap map[uuid.UUID]Player, err error) {
	sql := "SELECT count, p.person_id, telegram_id, firstname, lastname, fullname " +
		"FROM %s AS pl " +
		"INNER JOIN %s AS p ON pl.person_id = p.person_id " +
		"WHERE reserve_id = $1;"
	sql = fmt.Sprintf(sql, rep.PlayersTableName, rep.PersonsTableName)
	rows, err := rep.dbpool.Query(context.Background(), sql, rid)
	pmap = make(map[uuid.UUID]Player)
	var (
		Count, TelegramId             int
		PersonId                      uuid.UUID
		FirstName, LastName, FullName string
	)

	for rows.Next() {
		rows.Scan(&Count, &PersonId, &TelegramId, &FirstName, &LastName, &FullName)
		pmap[PersonId] = Player{
			Person: person.Person{Id: PersonId, TelegramId: TelegramId,
				Firstname: FirstName, Lastname: LastName, Fullname: FullName},
			Count: Count}
	}
	return
}

func (rep *PgRepository) Get(rid uuid.UUID) (res Reserve, err error) {
	sql_str := "SELECT reserve_id, person_id, start_time, end_time, price, " +
		"min_level, court_count, max_players, approved, canceled, " +
		"telegram_id, firstname, lastname, fullname, " +
		"location_id, location_name " +
		"FROM %s " +
		"WHERE reserve_id = $1"
	sql_str = fmt.Sprintf(sql_str, rep.ViewName)
	row := rep.dbpool.QueryRow(context.Background(), sql_str, rid)

	var lname sql.NullString
	err = row.Scan(&res.Id, &res.Person.Id, &res.StartTime, &res.EndTime, &res.Price,
		&res.MinLevel, &res.CourtCount, &res.MaxPlayers, &res.Approved, &res.Canceled,
		&res.Person.TelegramId, &res.Person.Firstname, &res.Person.Lastname, &res.Person.Fullname,
		&res.Location.Id, &lname)
	if err != nil {
		return
	}
	if lname.Valid {
		res.Location.Name = lname.String
	}
	pmap, err := rep.GetPlayers(res.Id)
	res.Players = pmap
	return
}

func (rep *PgRepository) GetByFilter(filter Reserve, oredered bool) (rmap map[uuid.UUID]Reserve, err error) {
	sql_str := "SELECT reserve_id, person_id, start_time, end_time, price, " +
		"min_level, court_count, max_players, approved, canceled, " +
		"telegram_id, firstname, lastname, fullname, " +
		"location_id, location_name " +
		"FROM %s "
	sql_str = fmt.Sprintf(sql_str, rep.ViewName)
	wheresql := ""

	params := []interface{}{}
	if filter.Id != uuid.Nil {
		AddWhereParam(&wheresql, &params, filter.Id, "reserve_id =")
	}
	if filter.Person.Id != uuid.Nil {
		AddWhereParam(&wheresql, &params, filter.Person.Id, "person_id =")
	}
	if !filter.StartTime.IsZero() {
		AddWhereParam(&wheresql, &params, filter.StartTime, "start_time >=")
	}
	if !filter.EndTime.IsZero() {
		AddWhereParam(&wheresql, &params, filter.EndTime, "start_time <=")
	}
	if oredered {
		AddWhereParam(&wheresql, &params, oredered, "ordered =")
	}
	if len(params) > 0 {
		wheresql = " WHERE" + wheresql
	}
	squery := sql_str + wheresql
	rows, err := rep.dbpool.Query(context.Background(), squery, params...)
	if err != nil {
		return rmap, err
	}
	rmap = make(map[uuid.UUID]Reserve)

	for rows.Next() {
		res := Reserve{}
		var lname sql.NullString
		err = rows.Scan(&res.Id, &res.Person.Id, &res.StartTime, &res.EndTime, &res.Price,
			&res.MinLevel, &res.CourtCount, &res.MaxPlayers, &res.Approved, &res.Canceled,
			&res.Person.TelegramId, &res.Person.Firstname, &res.Person.Lastname, &res.Person.Fullname,
			&res.Location.Id, &lname)
		if err != nil {
			return
		}
		if lname.Valid {
			res.Location.Name = lname.String
		}
		res.Players, err = rep.GetPlayers(res.Id)
		rmap[res.Id] = res
	}
	return
}

func (rep *PgRepository) Add(r Reserve) (res Reserve, err error) {
	sql := "INSERT INTO %s " +
		"(reserve_id, person_id, location_id, start_time, end_time, price, " +
		"min_level, court_count, max_players, approved, ordered, canceled) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) " +
		"RETURNING reserve_id"
	sql = fmt.Sprintf(sql, rep.TableName)

	row := rep.dbpool.QueryRow(context.Background(), sql,
		r.Id, r.Person.Id, r.Location.Id, r.StartTime, r.GetEndTime(), r.Price, r.MinLevel,
		r.CourtCount, r.MaxPlayers, r.Approved, r.Orderd(), r.Canceled)

	var ReserveId uuid.UUID
	err = row.Scan(&ReserveId)

	if err != nil {
		return
	}
	res = r
	res.Id = ReserveId

	return
}

func (rep *PgRepository) Update(r Reserve) (err error) {
	sql := "UPDATE %s SET " +
		"person_id = $1, location_id = $2, start_time = $3, end_time = $4, " +
		"price = $5, min_level = $6, court_count = $7, max_players = $8, " +
		"approved = $9, ordered = $10, canceled = $11 " +
		"WHERE reserve_id = $12"
	sql = fmt.Sprintf(sql, rep.TableName)

	rows, err := rep.dbpool.Query(context.Background(), sql,
		r.Person.Id, r.Location.Id, r.StartTime, r.GetEndTime(), r.Price, r.MinLevel,
		r.CourtCount, r.MaxPlayers, r.Approved, r.Orderd(), r.Canceled, r.Id)
	if err != nil {
		return
	}
	defer rows.Close()
	for _, pl := range r.Players {
		rep.UpdatePlayer(r, pl.Person, pl.Count)
	}
	return
}

func (rep *PgRepository) AddPlayer(r Reserve, pl person.Person, count int) (res Reserve, err error) {
	sql := "INSERT INTO %s (reserve_id, person_id, count) " +
		"VALUES ($1, $2, $3)"
	sql = fmt.Sprintf(sql, rep.PlayersTableName)

	rows, err := rep.dbpool.Query(context.Background(), sql, r.Id, pl.Id, count)
	if err != nil {
		return
	}
	defer rows.Close()
	res, err = rep.Get(r.Id)
	return
}

func (rep *PgRepository) UpdatePlayer(r Reserve, pl person.Person, count int) (res Reserve, err error) {
	sql := "call " + rep.PlayerSpName + " ($1, $2, $3);"
	_, err = rep.dbpool.Exec(context.Background(), sql, r.Id, pl.Id, count)
	return
}
