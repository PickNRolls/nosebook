package rootconvservice

import (
	domainchat "nosebook/src/domain/chat"
	"nosebook/src/errors"
	querybuilder "nosebook/src/infra/query_builder"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type chatRepository struct {
	db *sqlx.DB
}

func newChatRepository(db *sqlx.DB) *chatRepository {
	return &chatRepository{
		db: db,
	}
}

func (this *chatRepository) FindByMemberIds(leftId uuid.UUID, rightId uuid.UUID) (*domainchat.Chat, *errors.Error) {
	qb := querybuilder.New()

	sql, args, _ := qb.Select("l.chat_id as id", "c.created_at").
		Suffix(
			`from (
			select pc.chat_id from private_chats as pc
			join chat_members as cm on pc.chat_id = cm.chat_id
			where user_id = ?
		) as l join (
			select pc.chat_id from private_chats as pc
			join chat_members as cm on pc.chat_id = cm.chat_id
			where user_id = ? and user_id != ?
		) as r on l.chat_id = r.chat_id
		join chats as c on l.chat_id = c.id`,
			leftId, rightId, leftId,
		).ToSql()

	dest := chatDest{}
	err := this.db.Get(&dest, sql, args...)
	if err != nil {
		return nil, nil
	}

	return domainchat.New(
		dest.Id,
		uuid.UUIDs{leftId, rightId},
		"",
		true,
		dest.CreatedAt,
		nil,
		false,
	)
}

func (this *chatRepository) Save(chat *domainchat.Chat) *errors.Error {
	qb := querybuilder.New()

	for _, event := range chat.Events() {
		switch event.Type() {
		case domainchat.CREATED:
			sql, args, _ := qb.Insert("chats").
				Columns("id", "created_at").
				Values(chat.Id, chat.CreatedAt).
				ToSql()

			_, err := this.db.Exec(sql, args...)
			if err != nil {
				return errors.From(err)
			}

			sql, args, _ = qb.Insert("private_chats").
				Columns("chat_id").Values(chat.Id).
				ToSql()

			_, err = this.db.Exec(sql, args...)
			if err != nil {
				return errors.From(err)
			}

			if len(chat.MemberIds) > 0 {
				query := qb.Insert("chat_members").
					Columns("chat_id", "user_id")

				for _, memberId := range chat.MemberIds {
					query = query.Values(chat.Id, memberId)
				}

				sql, args, _ := query.ToSql()
				_, err := this.db.Exec(sql, args...)
				if err != nil {
					return errors.From(err)
				}
			}

		case domainchat.MESSAGE_SENT:
			messageSent := event.(*domainchat.MessageSentEvent)
			message := messageSent.Message
			sql, args, _ := qb.Insert("messages").
				Columns(
					"id",
					"author_id",
					"text",
					"reply_to",
					"chat_id",
					"created_at",
					"removed_at",
				).
				Values(
					message.Id,
					message.AuthorId,
					message.Text,
					message.ReplyTo,
					chat.Id,
					message.CreatedAt,
					message.RemovedAt,
				).
				ToSql()

			_, err := this.db.Exec(sql, args...)
			if err != nil {
				return errors.From(err)
			}
		}
	}

	return nil
}
