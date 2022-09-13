package postgres

import (
	"context"
	"fmt"
	"strconv"
	"volleybot/pkg/domain/location"
	"volleybot/pkg/domain/person"
	"volleybot/pkg/domain/volley"

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

func NewVolleyPgRepository(dbpool *pgxpool.Pool, PersonRepository person.PersonRepository, LocationRepository location.LocationRepository) (rep VolleyPgRepository, err error) {
	rep.PersonRepository = PersonRepository
	rep.LocationRepository = LocationRepository
	rep.TableName = "reserves"
	rep.PlayersTableName = "reserve_players"
	rep.PlayerSpName = "sp_reserve_player_update"

	if err != nil {
		return
	}

	rep.dbpool = dbpool
	return
}

type VolleyPgRepository struct {
	dbpool             *pgxpool.Pool
	PersonRepository   person.PersonRepository
	LocationRepository location.LocationRepository
	TableName          string
	PlayersTableName   string
	PlayerSpName       string
}

func (rep *VolleyPgRepository) UpdateDB() (err error) {
	sql := "CREATE TABLE IF NOT EXISTS %[1]s " +
		"(reserve_id UUID PRIMARY KEY, person_id UUID, location_id UUID, " +
		"start_time TIMESTAMP, end_time TIMESTAMP, price INT, " +
		"min_level INT, court_count INT, max_players INT, " +
		"ordered BOOL, approved BOOL, canceled BOOL, description varchar(4000), activity INT);"

	pl_sql := "CREATE TABLE IF NOT EXISTS %[2]s (player_id serial, reserve_id UUID, person_id UUID, count INT, arrive_time TIMESTAMP);"
	sp_sql := "CREATE OR REPLACE PROCEDURE " +
		"%[3]s(res_id UUID, per_id UUID, c INT, at TIMESTAMP) " +
		"LANGUAGE plpgsql AS $$ " +
		"DECLARE cur_count INT;\n" +
		"BEGIN\n" +
		"SELECT SUM(count) INTO cur_count " +
		"FROM %[2]s WHERE reserve_id = res_id AND person_id = per_id; " +
		"IF cur_count = 0 THEN\n" +
		"DELETE FROM %[2]s WHERE reserve_id = res_id AND person_id = per_id;\n" +
		"END IF;\n" +
		"CASE\n" +
		"WHEN cur_count > 0 THEN\n" +
		"UPDATE %[2]s SET count = c, arrive_time = at WHERE reserve_id = res_id AND person_id = per_id;\n" +
		"ELSE\n" +
		"INSERT INTO %[2]s (reserve_id, person_id, count, arrive_time) VALUES (res_id, per_id, c, at);\n" +
		"END CASE;\n" +
		"END;$$;"
	sql = fmt.Sprintf(sql+pl_sql+sp_sql, rep.TableName, rep.PlayersTableName, rep.PlayerSpName)
	_, err = rep.dbpool.Exec(context.Background(), sql)

	if err != nil {
		return
	}

	return
}

func (rep *VolleyPgRepository) GetPlayers(rid uuid.UUID) (plist []person.Player, err error) {
	sql := "SELECT player_id, count, arrive_time, person_id " +
		"FROM %s " +
		"WHERE reserve_id = $1 " +
		"ORDER BY player_id "
	sql = fmt.Sprintf(sql, rep.PlayersTableName)
	rows, err := rep.dbpool.Query(context.Background(), sql, rid)
	pl := person.Player{}
	for rows.Next() {
		rows.Scan(&pl.PlayerId, &pl.Count, &pl.ArriveTime, &pl.Id)
		pl.Person, _ = rep.PersonRepository.Get(pl.Id)
		plist = append(plist, pl)
	}
	return
}

func (rep *VolleyPgRepository) Get(rid uuid.UUID) (res volley.Volley, err error) {
	sql_str := "SELECT reserve_id, person_id, location_id, start_time, end_time, price, " +
		"min_level, court_count, max_players, approved, canceled, description, activity " +
		"FROM %s " +
		"WHERE reserve_id = $1"
	sql_str = fmt.Sprintf(sql_str, rep.TableName)
	row := rep.dbpool.QueryRow(context.Background(), sql_str, rid)

	err = row.Scan(&res.Id, &res.Person.Id, &res.Location.Id, &res.StartTime, &res.EndTime, &res.Price,
		&res.MinLevel, &res.CourtCount, &res.MaxPlayers, &res.Approved, &res.Canceled, &res.Description, &res.Activity)
	if err != nil {
		return
	}
	res.Person, _ = rep.PersonRepository.Get(res.Person.Id)
	res.Location, _ = rep.LocationRepository.Get(res.Location.Id)
	plist, err := rep.GetPlayers(res.Id)
	res.Players = plist
	return
}

func (rep *VolleyPgRepository) GetByFilter(filter volley.Volley, oredered bool, sorted bool) (rmap []volley.Volley, err error) {
	sql_str := "SELECT reserve_id, person_id, start_time, end_time, price, " +
		"min_level, court_count, max_players, approved, canceled, description, activity " +
		"FROM %s "
	sql_str = fmt.Sprintf(sql_str, rep.TableName)
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
	if sorted {
		squery += " ORDER BY start_time"
	}
	rows, err := rep.dbpool.Query(context.Background(), squery, params...)
	if err != nil {
		return rmap, err
	}
	rmap = []volley.Volley{}

	for rows.Next() {
		res := volley.Volley{}
		err = rows.Scan(&res.Id, &res.Person.Id, &res.StartTime, &res.EndTime, &res.Price,
			&res.MinLevel, &res.CourtCount, &res.MaxPlayers, &res.Approved, &res.Canceled,
			&res.Description, &res.Activity)
		if err != nil {
			return
		}
		res.Person, _ = rep.PersonRepository.Get(res.Person.Id)
		res.Location, _ = rep.LocationRepository.Get(res.Location.Id)
		res.Players, err = rep.GetPlayers(res.Id)
		rmap = append(rmap, res)
	}
	return
}

func (rep *VolleyPgRepository) Add(r volley.Volley) (res volley.Volley, err error) {
	sql := "INSERT INTO %s " +
		"(reserve_id, person_id, location_id, start_time, end_time, price, " +
		"min_level, court_count, max_players, approved, ordered, canceled, description, activity) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) " +
		"RETURNING reserve_id"
	sql = fmt.Sprintf(sql, rep.TableName)

	row := rep.dbpool.QueryRow(context.Background(), sql,
		r.Id, r.Person.Id, r.Location.Id, r.StartTime, r.GetEndTime(), r.Price, r.MinLevel,
		r.CourtCount, r.MaxPlayers, r.Approved, r.Ordered(), r.Canceled, r.Description, r.Activity)

	var ReserveId uuid.UUID
	err = row.Scan(&ReserveId)

	if err != nil {
		return
	}
	res, err = rep.Get(ReserveId)

	return
}

func (rep *VolleyPgRepository) Update(r volley.Volley) (err error) {
	sql := "UPDATE %s SET " +
		"person_id = $1, location_id = $2, start_time = $3, end_time = $4, " +
		"price = $5, min_level = $6, court_count = $7, max_players = $8, " +
		"approved = $9, ordered = $10, canceled = $11, description = $12, activity = $13  " +
		"WHERE reserve_id = $14"
	sql = fmt.Sprintf(sql, rep.TableName)

	rows, err := rep.dbpool.Query(context.Background(), sql,
		r.Person.Id, r.Location.Id, r.StartTime, r.GetEndTime(), r.Price, r.MinLevel,
		r.CourtCount, r.MaxPlayers, r.Approved, r.Ordered(), r.Canceled, r.Description, r.Activity, r.Id)
	if err != nil {
		return
	}
	defer rows.Close()
	for _, pl := range r.Players {
		rep.UpdatePlayer(r, pl)
	}
	return
}

func (rep *VolleyPgRepository) AddPlayer(r volley.Volley, pl person.Player) (res volley.Volley, err error) {
	sql := "INSERT INTO %s (reserve_id, person_id, count, arrive_time) " +
		"VALUES ($1, $2, $3, $4)"
	sql = fmt.Sprintf(sql, rep.PlayersTableName)

	rows, err := rep.dbpool.Query(context.Background(), sql, r.Id, pl.Id, pl.Count, pl.ArriveTime)
	if err != nil {
		return
	}
	defer rows.Close()
	res, err = rep.Get(r.Id)
	res.Person, _ = rep.PersonRepository.Get(res.Person.Id)
	return
}

func (rep *VolleyPgRepository) UpdatePlayer(r volley.Volley, pl person.Player) (res volley.Volley, err error) {
	sql := "call " + rep.PlayerSpName + " ($1, $2, $3, $4);"
	_, err = rep.dbpool.Exec(context.Background(), sql, r.Id, pl.Id, pl.Count, pl.ArriveTime)
	return
}
