package test

import (
	"testing"

	"chat/pkg/queue"
)

var (
	testCaseQueuePushPop = []string{
		"test1",
		"test2",
		"test3",
	}

	testCaseQueueFromHeadToTail = []string{
		"test1",
		"test2",
		"test3",
		"test4",
		"test5",
	}
)

func TestQueuePush(t *testing.T) {
	queue := queue.CreateQueue(3)

	for _, value := range testCaseQueuePushPop {
		queue.Push(value)
	}

	for index, value := range testCaseQueuePushPop {
		if queue.Buf[index] != value {
			t.FailNow()
		}
	}

	queue.Push("test4")
	if queue.Buf[0] != "test4" {
		t.FailNow()
	}
	for _, value := range testCaseQueuePushPop {
		queue.Push(value)
	}
	if queue.Buf[0] != "test3" {
		t.FailNow()
	}
}

func TestQueuePop(t *testing.T) {
	queue := queue.CreateQueue(3)

	for _, value := range testCaseQueuePushPop {
		queue.Push(value)
	}

	for _, testCase := range testCaseQueuePushPop {
		if value, err := queue.Pop(); err != nil || testCase != value {
			t.FailNow()
		}
	}

	if value, err := queue.Pop(); err == nil || value != "" {
		t.FailNow()
	}
}

func TestQueueFromHeadToTail(t *testing.T) {
	queue := queue.CreateQueue(3)
	for _, value := range testCaseQueueFromHeadToTail {
		queue.Push(value)
	}

	result := queue.FromHeadToTail()
	expected := []string{"test3", "test4", "test5"}

	for index, expectedValue := range expected {
		if result[index] != expectedValue {
			t.FailNow()
		}
	}

	for queue.Len != 0 {
		_, _ = queue.Pop()
	}

	if elem := queue.FromHeadToTail(); elem != nil {
		t.FailNow()
	}
}
