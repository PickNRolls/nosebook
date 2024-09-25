package rootconvservice

import (
	prometheusmetrics "nosebook/src/deps_root/worker"
	domainchat "nosebook/src/domain/chat"
	domainmessage "nosebook/src/domain/message"
	"nosebook/src/errors"
	querybuilder "nosebook/src/infra/query_builder"
	worker "nosebook/src/lib/worker"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type chatRepository struct {
	db           *sqlx.DB
	done         chan struct{}
	bufferInsert *worker.Buffer[*bufferInsertMessage, *errors.Error, time.Time]
	bufferFind   *worker.Buffer[*bufferFindMessage, *bufferFindOut, time.Time]
}

type bufferFindMessage struct {
	leftId  uuid.UUID
	rightId uuid.UUID
}

type bufferFindOut struct {
	data map[uuid.UUID]map[uuid.UUID]*chatDest
	err  *errors.Error
}

type bufferInsertMessage struct {
	message *domainmessage.Message
	chatId  uuid.UUID
}

func newChatRepository(db *sqlx.DB) *chatRepository {
	out := &chatRepository{
		db:   db,
		done: make(chan struct{}),
	}

	ticker := time.NewTicker(time.Millisecond * 10)
	out.bufferInsert = worker.NewBuffer(func(bufferedMessages []*bufferInsertMessage) *errors.Error {
		qb := querybuilder.New()
		query := qb.Insert("messages").
			Columns(
				"id",
				"author_id",
				"text",
				"reply_to",
				"chat_id",
				"created_at",
				"removed_at",
			)

		for _, buffered := range bufferedMessages {
			message := buffered.message
			chatId := buffered.chatId

			query = query.Values(
				message.Id,
				message.AuthorId,
				message.Text,
				message.ReplyTo,
				chatId,
				message.CreatedAt,
				message.RemovedAt,
			)
		}

		sql, args, _ := query.ToSql()
		_, err := db.Exec(sql, args...)
		return errors.From(err)
	}, ticker.C, out.done, 256, prometheusmetrics.UsePrometheusMetrics("message_insert"))

	qb := querybuilder.New(querybuilder.OmitPlaceholder)
	out.bufferFind = worker.NewBuffer(func(bufferedMessages []*bufferFindMessage) *bufferFindOut {
		out := &bufferFindOut{
			data: make(map[uuid.UUID]map[uuid.UUID]*chatDest),
		}

		leftUnique := map[uuid.UUID]struct{}{}
		rightUnique := map[uuid.UUID]struct{}{}
		for _, message := range bufferedMessages {
			if _, has := leftUnique[message.leftId]; !has {
				leftUnique[message.leftId] = struct{}{}
			}

			if _, has := rightUnique[message.rightId]; !has {
				rightUnique[message.rightId] = struct{}{}
			}
		}

		leftIds := []uuid.UUID{}
		for id := range leftUnique {
			leftIds = append(leftIds, id)
		}
		rightIds := []uuid.UUID{}
		for id := range rightUnique {
			rightIds = append(rightIds, id)
		}

		query := qb.Select("pc.chat_id, cm.user_id").
			From("private_chats as pc").
			Join("chat_members as cm on pc.chat_id = cm.chat_id")

		leftSql, leftArgs, _ := query.Where(squirrel.Eq{"user_id": leftIds}).ToSql()
		rightSql, rightArgs, _ := query.Where(squirrel.Eq{"user_id": rightIds}).ToSql()

		sql, args, _ := qb.Select("l.chat_id as id", "l.user_id as left_user_id", "r.user_id as right_user_id", "c.created_at").
			Suffix("from ("+leftSql, leftArgs...).
			Suffix(") as l join ("+rightSql, rightArgs...).
			Suffix(") as r on l.chat_id = r.chat_id join chats as c on l.chat_id = c.id").
			PlaceholderFormat(squirrel.Dollar).
			ToSql()

		dests := []*chatDest{}
		err := errors.From(db.Select(&dests, sql, args...))
		if err != nil {
			out.err = err
			return out
		}

		for _, dest := range dests {
			if _, has := out.data[dest.LeftUserId]; !has {
				out.data[dest.LeftUserId] = make(map[uuid.UUID]*chatDest)
			}

			out.data[dest.LeftUserId][dest.RightUserId] = dest
		}

		return out
	}, ticker.C, out.done, 256, prometheusmetrics.UsePrometheusMetrics("chat_find"))

	return out
}

func (this *chatRepository) Run() {
	go this.bufferInsert.Run()
	go this.bufferFind.Run()
	<-this.done
}

func (this *chatRepository) OnDone() {
	this.done <- struct{}{}
	close(this.done)
}

func (this *chatRepository) FindByMemberIds(leftId uuid.UUID, rightId uuid.UUID) (*domainchat.Chat, *errors.Error) {
	out := this.bufferFind.Send(&bufferFindMessage{
		leftId:  leftId,
		rightId: rightId,
	})
	if out.err != nil {
		return nil, out.err
	}

	dest := out.data[leftId][rightId]
	if dest == nil {
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

			err := this.bufferInsert.Send(&bufferInsertMessage{
				message: message,
				chatId:  chat.Id,
			})
			return err
		}
	}

	return nil
}
