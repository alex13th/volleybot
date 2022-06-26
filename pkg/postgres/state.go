package postgres

import (
	"context"
	"fmt"
	"volleybot/pkg/telegram"

	"github.com/jackc/pgx/v4/pgxpool"
)

type StatePgRepository struct {
	dbpool    *pgxpool.Pool
	TableName string
}

func NewStatePgRepository(dbpool *pgxpool.Pool) (pgrep StatePgRepository, err error) {
	pgrep.TableName = "tg_states"
	if err != nil {
		return
	}

	pgrep.dbpool = dbpool

	return
}

func (rep *StatePgRepository) UpdateDB() (err error) {
	sql := "CREATE TABLE IF NOT EXISTS %s " +
		"(chat_id bigint, message_id bigint, state varchar(240), data varchar(240));"
	sql = fmt.Sprintf(sql, rep.TableName)
	_, err = rep.dbpool.Exec(context.Background(), sql)

	if err != nil {
		return
	}

	return
}

func (rep *StatePgRepository) Get(ChatId int) (st telegram.State, err error) {
	sql := "SELECT chat_id, message_id, state, data " +
		"FROM %s " +
		"WHERE chat_id = $1"
	row := rep.dbpool.QueryRow(context.Background(), fmt.Sprintf(sql, rep.TableName), ChatId)
	err = row.Scan(&st.ChatId, &st.MessageId, &st.State, &st.Data)
	if err != nil {
		return
	}
	return
}

func (rep *StatePgRepository) Set(st telegram.State) (err error) {
	sql := "UPDATE %s SET " +
		"state = $1, data = $2, " +
		"WHERE (chat_id = $3) AND (message_id = $4)"
	sql = fmt.Sprintf(sql, rep.TableName)

	rows, err := rep.dbpool.Query(context.Background(), sql,
		st.State, st.Data, st.ChatId, st.MessageId)
	if rows.CommandTag().RowsAffected() < 1 {
		rep.Add(st)
	}
	if err != nil {
		return
	}
	defer rows.Close()
	return
}

func (rep *StatePgRepository) Add(st telegram.State) error {
	sql := "INSERT INTO %s " +
		"(chat_id, message_id, state, data) " +
		"VALUES ($1, $2, $3, $4);"
	sql = fmt.Sprintf(sql, rep.TableName)

	row := rep.dbpool.QueryRow(context.Background(), sql,
		st.ChatId, st.MessageId, st.State, st.Data)
	return row.Scan()
}

func (rep *StatePgRepository) Clear(st telegram.State) error {
	sql := "DELETE FROM %s " +
		"(chat_id, message_id, state, data) " +
		"WHERE (chat_id = $3) AND (message_id = $4);"
	sql = fmt.Sprintf(sql, rep.TableName)

	row := rep.dbpool.QueryRow(context.Background(), sql,
		st.ChatId, st.MessageId)
	return row.Scan()
}
