package tests

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/server/common/persistence"
	"go.temporal.io/server/common/persistence/sql/sqlplugin"
	"go.temporal.io/server/common/shuffle"
)

const (
	testQueueMessageEncoding = "random encoding"
)

var (
	testQueueMessageData = []byte("random queue data")
)

type (
	queueMessageSuite struct {
		suite.Suite
		*require.Assertions

		store sqlplugin.QueueMessage
	}
)

func NewQueueMessageSuite(
	t *testing.T,
	store sqlplugin.QueueMessage,
) *queueMessageSuite {
	return &queueMessageSuite{
		Assertions: require.New(t),
		store:      store,
	}
}

func (s *queueMessageSuite) SetupSuite() {

}

func (s *queueMessageSuite) TearDownSuite() {

}

func (s *queueMessageSuite) SetupTest() {
	s.Assertions = require.New(s.T())
}

func (s *queueMessageSuite) TearDownTest() {

}

func (s *queueMessageSuite) TestInsert_Single_Success() {
	queueType := persistence.NamespaceReplicationQueueType
	messageID := rand.Int63()

	message := s.newRandomQueueMessageRow(queueType, messageID)
	result, err := s.store.InsertIntoMessages(newExecutionContext(), []sqlplugin.QueueMessageRow{message})
	s.NoError(err)
	rowsAffected, err := result.RowsAffected()
	s.NoError(err)
	s.Equal(1, int(rowsAffected))
}

func (s *queueMessageSuite) TestInsert_Multiple_Success() {
	queueType := persistence.NamespaceReplicationQueueType
	messageID := rand.Int63()

	message1 := s.newRandomQueueMessageRow(queueType, messageID)
	messageID++
	message2 := s.newRandomQueueMessageRow(queueType, messageID)
	result, err := s.store.InsertIntoMessages(newExecutionContext(), []sqlplugin.QueueMessageRow{message1, message2})
	s.NoError(err)
	rowsAffected, err := result.RowsAffected()
	s.NoError(err)
	s.Equal(2, int(rowsAffected))
}

func (s *queueMessageSuite) TestInsert_Single_Fail_Duplicate() {
	queueType := persistence.NamespaceReplicationQueueType
	messageID := rand.Int63()

	message := s.newRandomQueueMessageRow(queueType, messageID)
	result, err := s.store.InsertIntoMessages(newExecutionContext(), []sqlplugin.QueueMessageRow{message})
	s.NoError(err)
	rowsAffected, err := result.RowsAffected()
	s.NoError(err)
	s.Equal(1, int(rowsAffected))

	message = s.newRandomQueueMessageRow(queueType, messageID)
	_, err = s.store.InsertIntoMessages(newExecutionContext(), []sqlplugin.QueueMessageRow{message})
	s.Error(err) // TODO persistence layer should do proper error translation
}

func (s *queueMessageSuite) TestInsert_Multiple_Fail_Duplicate() {
	queueType := persistence.NamespaceReplicationQueueType
	messageID := rand.Int63()

	message1 := s.newRandomQueueMessageRow(queueType, messageID)
	messageID++
	message2 := s.newRandomQueueMessageRow(queueType, messageID)
	result, err := s.store.InsertIntoMessages(newExecutionContext(), []sqlplugin.QueueMessageRow{message1, message2})
	s.NoError(err)
	rowsAffected, err := result.RowsAffected()
	s.NoError(err)
	s.Equal(2, int(rowsAffected))

	message2 = s.newRandomQueueMessageRow(queueType, messageID)
	messageID++
	message3 := s.newRandomQueueMessageRow(queueType, messageID)
	_, err = s.store.InsertIntoMessages(newExecutionContext(), []sqlplugin.QueueMessageRow{message2, message3})
	s.Error(err) // TODO persistence layer should do proper error translation
}

func (s *queueMessageSuite) TestInsertSelect() {
	queueType := persistence.NamespaceReplicationQueueType
	messageID := rand.Int63()

	message := s.newRandomQueueMessageRow(queueType, messageID)
	result, err := s.store.InsertIntoMessages(newExecutionContext(), []sqlplugin.QueueMessageRow{message})
	s.NoError(err)
	rowsAffected, err := result.RowsAffected()
	s.NoError(err)
	s.Equal(1, int(rowsAffected))

	filter := sqlplugin.QueueMessagesFilter{
		QueueType: queueType,
		MessageID: messageID,
	}
	rows, err := s.store.SelectFromMessages(newExecutionContext(), filter)
	s.NoError(err)
	for index := range rows {
		rows[index].QueueType = queueType
	}
	s.Equal([]sqlplugin.QueueMessageRow{message}, rows)
}

func (s *queueMessageSuite) TestInsertSelect_Multiple() {
	numMessages := 20

	queueType := persistence.NamespaceReplicationQueueType
	minMessageID := rand.Int63()
	messageID := minMessageID + 1
	maxMessageID := messageID + int64(numMessages)

	var messages []sqlplugin.QueueMessageRow
	for i := 0; i < numMessages; i++ {
		message := s.newRandomQueueMessageRow(queueType, messageID)
		messageID++
		messages = append(messages, message)
	}
	result, err := s.store.InsertIntoMessages(newExecutionContext(), messages)
	s.NoError(err)
	rowsAffected, err := result.RowsAffected()
	s.NoError(err)
	s.Equal(numMessages, int(rowsAffected))

	filter := sqlplugin.QueueMessagesRangeFilter{
		QueueType:    queueType,
		MinMessageID: minMessageID,
		MaxMessageID: maxMessageID,
		PageSize:     numMessages,
	}
	rows, err := s.store.RangeSelectFromMessages(newExecutionContext(), filter)
	s.NoError(err)
	for index := range rows {
		rows[index].QueueType = queueType
	}
	s.Equal(messages, rows)
}

func (s *queueMessageSuite) TestDeleteSelect_Single() {
	queueType := persistence.NamespaceReplicationQueueType
	messageID := rand.Int63()

	filter := sqlplugin.QueueMessagesFilter{
		QueueType: queueType,
		MessageID: messageID,
	}
	result, err := s.store.DeleteFromMessages(newExecutionContext(), filter)
	s.NoError(err)
	rowsAffected, err := result.RowsAffected()
	s.NoError(err)
	s.Equal(0, int(rowsAffected))

	rows, err := s.store.SelectFromMessages(newExecutionContext(), filter)
	s.NoError(err)
	for index := range rows {
		rows[index].QueueType = queueType
	}
	s.Equal([]sqlplugin.QueueMessageRow(nil), rows)
}

func (s *queueMessageSuite) TestDeleteSelect_Multiple() {
	pageSize := 100

	queueType := persistence.NamespaceReplicationQueueType
	minMessageID := rand.Int63()
	maxMessageID := minMessageID + int64(20)

	filter := sqlplugin.QueueMessagesRangeFilter{
		QueueType:    queueType,
		MinMessageID: minMessageID,
		MaxMessageID: maxMessageID,
		PageSize:     0,
	}
	result, err := s.store.RangeDeleteFromMessages(newExecutionContext(), filter)
	s.NoError(err)
	rowsAffected, err := result.RowsAffected()
	s.NoError(err)
	s.Equal(0, int(rowsAffected))

	filter.PageSize = pageSize
	rows, err := s.store.RangeSelectFromMessages(newExecutionContext(), filter)
	s.NoError(err)
	for index := range rows {
		rows[index].QueueType = queueType
	}
	s.Equal([]sqlplugin.QueueMessageRow(nil), rows)
}

func (s *queueMessageSuite) TestInsertDeleteSelect_Single() {
	queueType := persistence.NamespaceReplicationQueueType
	messageID := rand.Int63()

	message := s.newRandomQueueMessageRow(queueType, messageID)
	result, err := s.store.InsertIntoMessages(newExecutionContext(), []sqlplugin.QueueMessageRow{message})
	s.NoError(err)
	rowsAffected, err := result.RowsAffected()
	s.NoError(err)
	s.Equal(1, int(rowsAffected))

	filter := sqlplugin.QueueMessagesFilter{
		QueueType: queueType,
		MessageID: messageID,
	}
	result, err = s.store.DeleteFromMessages(newExecutionContext(), filter)
	s.NoError(err)
	rowsAffected, err = result.RowsAffected()
	s.NoError(err)
	s.Equal(1, int(rowsAffected))

	rows, err := s.store.SelectFromMessages(newExecutionContext(), filter)
	s.NoError(err)
	for index := range rows {
		rows[index].QueueType = queueType
	}
	s.Equal([]sqlplugin.QueueMessageRow(nil), rows)
}

func (s *queueMessageSuite) TestInsertDeleteSelect_Multiple() {
	numMessages := 20
	pageSize := numMessages

	queueType := persistence.NamespaceReplicationQueueType
	minMessageID := rand.Int63()
	messageID := minMessageID + 1
	maxMessageID := messageID + int64(numMessages)

	var messages []sqlplugin.QueueMessageRow
	for i := 0; i < numMessages; i++ {
		message := s.newRandomQueueMessageRow(queueType, messageID)
		messageID++
		messages = append(messages, message)
	}
	result, err := s.store.InsertIntoMessages(newExecutionContext(), messages)
	s.NoError(err)
	rowsAffected, err := result.RowsAffected()
	s.NoError(err)
	s.Equal(numMessages, int(rowsAffected))

	filter := sqlplugin.QueueMessagesRangeFilter{
		QueueType:    queueType,
		MinMessageID: minMessageID,
		MaxMessageID: maxMessageID,
		PageSize:     0,
	}
	result, err = s.store.RangeDeleteFromMessages(newExecutionContext(), filter)
	s.NoError(err)
	rowsAffected, err = result.RowsAffected()
	s.NoError(err)
	s.Equal(numMessages, int(rowsAffected))

	filter.PageSize = pageSize
	rows, err := s.store.RangeSelectFromMessages(newExecutionContext(), filter)
	s.NoError(err)
	for index := range rows {
		rows[index].QueueType = queueType
	}
	s.Equal([]sqlplugin.QueueMessageRow(nil), rows)
}

func (s *queueMessageSuite) newRandomQueueMessageRow(
	queueType persistence.QueueType,
	messageID int64,
) sqlplugin.QueueMessageRow {
	return sqlplugin.QueueMessageRow{
		QueueType:       queueType,
		MessageID:       messageID,
		MessagePayload:  shuffle.Bytes(testQueueMessageData),
		MessageEncoding: testQueueMessageEncoding,
	}
}
