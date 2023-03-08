package queue

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type QueueTest struct {
	suite.Suite
	queue *Queue
}

func (suite *QueueTest) SetupTest() {
	suite.queue = NewQueue(5, "test")
}

func TestQueueTestSuite(t *testing.T) {
	suite.Run(t, new(QueueTest))
}

func (suite *QueueTest) TestEnqueueLock() {
	//test enqueue without lock
	suite.NoError(suite.queue.Enqueue(1), "Unlocked queue allows to enqueue elements")

	suite.queue.Lock()
	err := suite.queue.Enqueue(1)
	suite.Error(err, "Locked queue does not allow to enqueue elements")

	//test queue lock error message
	suite.Equal(err.Error(), "the queue is locked")
}

func (suite *QueueTest) TestEnqueueFullCapacity() {
	//test enqueue without lock for 5 entries
	suite.NoError(suite.queue.Enqueue(1), "Unlocked queue allows to enqueue elements")
	suite.NoError(suite.queue.Enqueue(1), "Unlocked queue allows to enqueue elements")
	suite.NoError(suite.queue.Enqueue(1), "Unlocked queue allows to enqueue elements")
	suite.NoError(suite.queue.Enqueue(1), "Unlocked queue allows to enqueue elements")
	suite.NoError(suite.queue.Enqueue(1), "Unlocked queue allows to enqueue elements")

	err := suite.queue.Enqueue(1)
	suite.Error(err, "Queue at full capacity does not allow new elements")

	//test queue lock error message
	suite.Equal(err.Error(), "Queue is at full capacity")
}

func (suite *QueueTest) TestIsValidIdentity() {
	suite.True(suite.queue.IsValidIdentity("test"))
}

func (suite *QueueTest) TestDequeueLock() {
	suite.NoError(suite.queue.Enqueue(1), "Unlocked queue allows to enqueue elements")
	suite.NoError(suite.queue.Enqueue(1), "Unlocked queue allows to enqueue elements")

	val, err := suite.queue.Dequeue()
	suite.NoError(err, "Unlocked queue allows to dequeue elements")
	suite.Equal(1, val)

	suite.queue.Lock()
	val, err = suite.queue.Dequeue()
	suite.Nil(val, "dequeue after locking returns nil value with error")
	suite.Error(err, "Locked queue does not allow to dequeue elements")

	//test queue lock error message
	suite.Equal(err.Error(), "the queue is locked")
}

func (suite *QueueTest) TestDequeueEmptyQueue() {

	val, err := suite.queue.Dequeue()
	suite.Error(err, "dequeue on empty queue is not allowed")
	suite.Nil(val, "dequeue on empty queue return nil with error")

	suite.Equal(err.Error(), "empty queue")
}