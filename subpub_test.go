package subpub

import (
	"context"
	"sync"
	"testing"
)

func TestNewSubPub(t *testing.T) {
	sp := NewSubPub()
	if sp == nil {
		t.Errorf("NewSubPub() returned nil")
	}
}

func TestQueueSize(t *testing.T) {
	sp := NewSubPub()
	if sp.GetQueueSize() != 100 {
		t.Errorf("Wrong queue size")
	}

	spWithQueueSize := sp.WithQueueSize(10)
	if spWithQueueSize == nil {
		t.Errorf("WithQueueSize() returned nil")
	} else if spWithQueueSize != sp {
		t.Errorf("spWithQueueSize is not equals to sp")
	} else if spWithQueueSize.GetQueueSize() != 10 {
		t.Errorf("Wrong queue size")
	}
}

func TestSubscribe(t *testing.T) {
	sp := NewSubPub()
	subscription, err := sp.Subscribe("subject", func(msg interface{}) {})

	if err != nil {
		t.Errorf("Failed to subscribe: %v", err)
	}
	if subscription == nil {
		t.Errorf("subscription is nil")
	}
}

func TestNumberOfSubjects(t *testing.T) {
	sp := NewSubPub()

	_, _ = sp.Subscribe("subject1", func(msg interface{}) {})
	if sp.GetNumberOfSubjects() != 1 {
		t.Errorf("Wrong number of subjects")
	}

	_, _ = sp.Subscribe("subject1", func(msg interface{}) {})
	if sp.GetNumberOfSubjects() != 1 {
		t.Errorf("Wrong number of subjects")
	}

	_, _ = sp.Subscribe("subject2", func(msg interface{}) {})
	if sp.GetNumberOfSubjects() != 2 {
		t.Errorf("Wrong number of subjects")
	}
}

func TestNumberOfSubscriptions(t *testing.T) {
	sp := NewSubPub()

	subscription1, _ := sp.Subscribe("subject", func(msg interface{}) {})
	if sp.GetNumberOfSubscriptions("subject") != 1 {
		t.Errorf("Wrong number of subscriptions")
	}

	subscription2, _ := sp.Subscribe("subject", func(msg interface{}) {})
	if sp.GetNumberOfSubscriptions("subject") != 2 {
		t.Errorf("Wrong number of subscriptions")
	}

	subscription2.Unsubscribe()
	if sp.GetNumberOfSubscriptions("subject") != 1 {
		t.Errorf("Wrong number of subscriptions")
	}

	sp.Unsubscribe(subscription1)
	if sp.GetNumberOfSubscriptions("subject") != 0 {
		t.Errorf("Wrong number of subscriptions")
	}
}

func TestPublish(t *testing.T) {
	sp := NewSubPub()

	str := ""

	wg := sync.WaitGroup{}
	wg.Add(10)

	subscription, _ := sp.Subscribe("subject", func(msg interface{}) {
		str += msg.(string)
		wg.Done()
	})

	for i := 0; i < 10; i++ {
		if err := sp.Publish("subject", "1"); err != nil {
			t.Errorf("Failed to publish: %v", err)
		}
	}

	go subscription.HandleMessages()
	wg.Wait()

	if str != "1111111111" {
		t.Errorf("Wrong publish result")
	}
}

func TestIsClosed(t *testing.T) {
	sp := NewSubPub()

	err := sp.Close(context.Background())
	if err != nil {
		t.Errorf("Failed to close: %v", err)
	}
}

func TestIsNotClosed(t *testing.T) {
	sp := NewSubPub()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := sp.Close(ctx)
	if err == nil {
		t.Errorf("Close should have failed")
	}
	if sp.IsClosed() {
		t.Errorf("IsClosed() returned true")
	}
}
