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

func NewStateRepository(dbpool *pgxpool.Pool) (pgrep StatePgRepository, err error) {
	pgrep.TableName = "tg_states"
	if err != nil {
		return
	}

	pgrep.dbpool = dbpool

	return
}

func (rep *StatePgRepository) UpdateDB() (err error) {
	sql := "CREATE TABLE IF NOT EXISTS %s " +
		"(chat_id bigint, message_id bigint, prefix varchar(63), state varchar(63)," +
		" action varchar(63), data varchar(240));"
	sql = fmt.Sprintf(sql, rep.TableName)
	_, err = rep.dbpool.Exec(context.Background(), sql)

	if err != nil {
		return
	}

	return
}

func (rep *StatePgRepository) Get(ChatId int) (slist []telegram.State, err error) {
	sql := "SELECT chat_id, message_id, prefix, state, action, data " +
		"FROM %s " +
		"WHERE chat_id = $1 and message_id <= 0" +
		"ORDER BY message_id DESC"
	sql = fmt.Sprintf(sql, rep.TableName)
	rows, err := rep.dbpool.Query(context.Background(), sql, ChatId)
	if err == nil {
		defer rows.Close()
		st := telegram.NewState()
		for rows.Next() {
			rows.Scan(&st.ChatId, &st.MessageId, &st.Prefix, &st.State, &st.Action, &st.Data)
			slist = append(slist, st)
		}
	}
	return
}

func (rep *StatePgRepository) GetByData(Data string) (slist []telegram.State, err error) {
	sql := "SELECT chat_id, message_id, prefix, state, action, data " +
		"FROM %s " +
		"WHERE data = $1 "
	sql = fmt.Sprintf(sql, rep.TableName)
	if rows, err := rep.dbpool.Query(context.Background(), sql, Data); err == nil {
		defer rows.Close()
		st := telegram.NewState()
		for rows.Next() {
			rows.Scan(&st.ChatId, &st.MessageId, &st.Prefix, &st.State, &st.Action, &st.Data)
			slist = append(slist, st)
		}
	}
	return
}

func (rep *StatePgRepository) GetByMessage(msg telegram.Message) (state telegram.State, err error) {
	sql := "SELECT chat_id, message_id, prefix, state, action, data " +
		"FROM %s " +
		"WHERE chat_id = $1 AND message_id = $2 "
	sql = fmt.Sprintf(sql, rep.TableName)
	rows, err := rep.dbpool.Query(context.Background(), sql, msg.Chat.Id, msg.MessageId)
	if err == nil {
		defer rows.Close()
		st := telegram.NewState()
		for rows.Next() {
			rows.Scan(&st.ChatId, &st.MessageId, &st.Prefix, &st.State, &st.Action, &st.Data)
			state = st
		}
	}
	return
}

func (rep *StatePgRepository) Set(st telegram.State) (err error) {
	sql := "UPDATE %s SET " +
		"prefix =$1, state = $2, action =$3, data = $4 " +
		"WHERE (chat_id = $5) AND (message_id = $6)"
	sql = fmt.Sprintf(sql, rep.TableName)

	rres, err := rep.dbpool.Exec(context.Background(), sql,
		st.Prefix, st.State, st.Action, st.Data, st.ChatId, st.MessageId)
	if rres.RowsAffected() < 1 {
		rep.Add(st)
	}
	return
}

func (rep *StatePgRepository) Add(st telegram.State) error {
	sql := "INSERT INTO %s " +
		"(chat_id, message_id, prefix, state, action, data) " +
		"VALUES ($1, $2, $3, $4, $5, $6);"
	sql = fmt.Sprintf(sql, rep.TableName)

	row := rep.dbpool.QueryRow(context.Background(), sql,
		st.ChatId, st.MessageId, st.Prefix, st.State, st.Action, st.Data)
	return row.Scan()
}

func (rep *StatePgRepository) Clear(st telegram.State) error {
	sql := "DELETE FROM %s " +
		"WHERE (chat_id = $1) AND (message_id = $2);"
	sql = fmt.Sprintf(sql, rep.TableName)

	_, err := rep.dbpool.Exec(context.Background(), sql,
		st.ChatId, st.MessageId)
	return err
}
