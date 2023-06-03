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
	rep.TableName = "bvreserves"
	rep.MembersTableName = "bvreserve_members"
	rep.MembersSpName = "sp_bvreserve_member_update"
	rep.PlayersTableName = "bvplayers"
	rep.PlayersSpName = "sp_bvplayer_update"

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
	MembersTableName   string
	MembersSpName      string
	PlayersTableName   string
	PlayersSpName      string
}

func (rep *VolleyPgRepository) UpdateDB() (err error) {
	sql := "CREATE TABLE IF NOT EXISTS %[1]s " +
		"(reserve_id UUID PRIMARY KEY, person_id UUID, location_id UUID, " +
		"start_time TIMESTAMP, end_time TIMESTAMP, price INT, " +
		"min_level INT, court_count INT, max_players INT, net_type INT, " +
		"ordered BOOL, approved BOOL, canceled BOOL, description varchar(4000), activity INT);"

	mb_sql := "CREATE TABLE IF NOT EXISTS %[2]s "
	mb_sql += "(member_id serial, reserve_id UUID, person_id UUID, count INT, "
	mb_sql += "arrive_time TIMESTAMP, paid BOOL);"
	pl_sql := "CREATE TABLE IF NOT EXISTS %[3]s (person_id UUID PRIMARY KEY, level INT);"
	sp_sql := "CREATE OR REPLACE PROCEDURE " +
		"%[4]s(res_id UUID, per_id UUID, c INT, at TIMESTAMP, pd BOOL) " +
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
		"UPDATE %[2]s SET count = c, arrive_time = at, paid = pd WHERE reserve_id = res_id AND person_id = per_id;\n" +
		"ELSE\n" +
		"INSERT INTO %[2]s (reserve_id, person_id, count, arrive_time, paid) VALUES (res_id, per_id, c, at, pd);\n" +
		"END CASE;\n" +
		"END;$$;"
	sp_pl_sql := "CREATE OR REPLACE PROCEDURE " +
		"%[5]s(per_id UUID, lvl INT) " +
		"LANGUAGE plpgsql AS $$ " +
		"BEGIN\n" +
		"IF (SELECT COUNT(*) FROM %[3]s WHERE person_id = per_id) = 0 THEN\n" +
		"INSERT INTO %[3]s (person_id, level) VALUES (per_id, lvl);\n" +
		"ELSE\n" +
		"UPDATE %[3]s SET level = lvl WHERE person_id = per_id;\n" +
		"END IF;\n" +
		"END;$$;"
	sql = fmt.Sprintf(sql+mb_sql+pl_sql+sp_sql+sp_pl_sql, rep.TableName, rep.MembersTableName, rep.PlayersTableName,
		rep.MembersSpName, rep.PlayersSpName)
	_, err = rep.dbpool.Exec(context.Background(), sql)

	if err != nil {
		return
	}

	return
}

func (rep *VolleyPgRepository) GetMembers(rid uuid.UUID) (mlist []volley.Member, err error) {
	sql := "SELECT member_id, count, arrive_time, paid, person_id " +
		"FROM %s " +
		"WHERE reserve_id = $1 " +
		"ORDER BY paid DESC, member_id "
	sql = fmt.Sprintf(sql, rep.MembersTableName)
	rows, err := rep.dbpool.Query(context.Background(), sql, rid)
	var mb volley.Member
	for rows.Next() {
		var paid bool
		rows.Scan(&mb.MemberId, &mb.Count, &mb.ArriveTime, &paid, &mb.Id)
		mb.SetPaid(paid)
		p, _ := rep.PersonRepository.Get(mb.Id)
		mb.Player, _ = rep.GetPlayer(p)
		mlist = append(mlist, mb)
	}
	return
}

func (rep *VolleyPgRepository) Get(rid uuid.UUID) (res volley.Volley, err error) {
	sql_str := "SELECT reserve_id, person_id, location_id, start_time, end_time, price, " +
		"min_level, court_count, max_players, net_type, approved, canceled, description, activity " +
		"FROM %s " +
		"WHERE reserve_id = $1"
	sql_str = fmt.Sprintf(sql_str, rep.TableName)
	row := rep.dbpool.QueryRow(context.Background(), sql_str, rid)

	err = row.Scan(&res.Id, &res.Person.Id, &res.Location.Id, &res.StartTime, &res.EndTime, &res.Price,
		&res.MinLevel, &res.CourtCount, &res.MaxPlayers, &res.NetType, &res.Approved, &res.Canceled, &res.Description, &res.Activity)
	if err != nil {
		return
	}
	res.Person, _ = rep.PersonRepository.Get(res.Person.Id)
	res.Location, _ = rep.LocationRepository.Get(res.Location.Id)
	plist, err := rep.GetMembers(res.Id)
	res.Members = plist
	return
}

func (rep *VolleyPgRepository) GetByFilter(filter volley.Volley, oredered bool, sorted bool) (rmap []volley.Volley, err error) {
	sql_str := "SELECT reserve_id, person_id, start_time, end_time, price, " +
		"min_level, court_count, max_players, net_type, approved, canceled, description, activity " +
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
			&res.MinLevel, &res.CourtCount, &res.MaxPlayers, &res.NetType, &res.Approved, &res.Canceled,
			&res.Description, &res.Activity)
		if err != nil {
			return
		}
		res.Person, _ = rep.PersonRepository.Get(res.Person.Id)
		res.Location, _ = rep.LocationRepository.Get(res.Location.Id)
		res.Members, err = rep.GetMembers(res.Id)
		rmap = append(rmap, res)
	}
	return
}

func (rep *VolleyPgRepository) Add(r volley.Volley) (res volley.Volley, err error) {
	sql := "INSERT INTO %s " +
		"(reserve_id, person_id, location_id, start_time, end_time, price, " +
		"min_level, court_count, max_players, net_type, approved, ordered, canceled, description, activity) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) " +
		"RETURNING reserve_id"
	sql = fmt.Sprintf(sql, rep.TableName)

	row := rep.dbpool.QueryRow(context.Background(), sql,
		r.Id, r.Person.Id, r.Location.Id, r.StartTime, r.GetEndTime(), r.Price, r.MinLevel,
		r.CourtCount, r.MaxPlayers, r.NetType, r.Approved, r.Ordered(), r.Canceled, r.Description, r.Activity)

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
		"price = $5, min_level = $6, court_count = $7, max_players = $8, net_type = $9, " +
		"approved = $10, ordered = $11, canceled = $12, description = $13, activity = $14  " +
		"WHERE reserve_id = $15"
	sql = fmt.Sprintf(sql, rep.TableName)

	rows, err := rep.dbpool.Query(context.Background(), sql,
		r.Person.Id, r.Location.Id, r.StartTime, r.GetEndTime(), r.Price, r.MinLevel,
		r.CourtCount, r.MaxPlayers, r.NetType, r.Approved, r.Ordered(), r.Canceled, r.Description, r.Activity, r.Id)
	if err != nil {
		return
	}
	defer rows.Close()
	for _, mb := range r.Members {
		rep.UpdateMember(r, mb)
	}
	return
}

func (rep *VolleyPgRepository) AddMember(r volley.Volley, mb volley.Member) (res volley.Volley, err error) {
	sql := "INSERT INTO %s (reserve_id, person_id, count, arrive_time, paid) " +
		"VALUES ($1, $2, $3, $4, $5)"
	sql = fmt.Sprintf(sql, rep.MembersTableName)

	rows, err := rep.dbpool.Query(context.Background(), sql, r.Id, mb.Id, mb.Count, mb.ArriveTime, mb.GetPaid())
	if err != nil {
		return
	}
	defer rows.Close()
	res, err = rep.Get(r.Id)
	res.Person, _ = rep.PersonRepository.Get(res.Person.Id)
	return
}

func (rep *VolleyPgRepository) UpdateMember(r volley.Volley, mb volley.Member) (res volley.Volley, err error) {
	sql := "call " + rep.MembersSpName + " ($1, $2, $3, $4, $5);"
	_, err = rep.dbpool.Exec(context.Background(), sql, r.Id, mb.Id, mb.Count, mb.ArriveTime, mb.GetPaid())
	return
}

func (rep *VolleyPgRepository) AddPlayer(pl volley.Player) (res volley.Player, err error) {
	sql := "INSERT INTO %s (person_id, level) " +
		"VALUES ($1, $2, $3)"
	sql = fmt.Sprintf(sql, rep.PlayersTableName)

	rows, err := rep.dbpool.Query(context.Background(), sql, pl.Id, pl.Level)
	if err != nil {
		return
	}
	defer rows.Close()
	return
}

func (rep *VolleyPgRepository) GetPlayer(p person.Person) (pl volley.Player, err error) {
	sql := "SELECT level FROM %s WHERE person_id = $1"
	row := rep.dbpool.QueryRow(context.Background(), fmt.Sprintf(sql, rep.PlayersTableName), p.Id)
	row.Scan(&pl.Level)
	pl.Person = p
	if err != nil {
		return
	}
	return
}

func (rep *VolleyPgRepository) UpdatePlayer(pl volley.Player) (err error) {
	sql := "call " + rep.PlayersSpName + " ($1, $2);"
	_, err = rep.dbpool.Exec(context.Background(), sql, pl.Id, pl.Level)
	rep.PersonRepository.Update(pl.Person)
	return
}
