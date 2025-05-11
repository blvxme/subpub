package subpub

import (
	"context"
	"errors"
	"sync"
)

type SubPub struct {
	mutex     sync.Mutex
	isClosed  bool
	queueSize int
	subjects  map[string]map[*Subscription]chan<- interface{}
}

func NewSubPub() *SubPub {
	return &SubPub{subjects: make(map[string]map[*Subscription]chan<- interface{}), queueSize: 100}
}

func (sp *SubPub) WithQueueSize(queueSize int) *SubPub {
	sp.queueSize = queueSize
	return sp
}

func (sp *SubPub) IsClosed() bool {
	sp.mutex.Lock()
	defer sp.mutex.Unlock()

	return sp.isClosed
}

func (sp *SubPub) GetNumberOfSubjects() int {
	sp.mutex.Lock()
	defer sp.mutex.Unlock()

	return len(sp.subjects)
}

func (sp *SubPub) GetNumberOfSubscriptions(subject string) int {
	sp.mutex.Lock()
	defer sp.mutex.Unlock()

	return len(sp.subjects[subject])
}

func (sp *SubPub) GetQueueSize() int {
	sp.mutex.Lock()
	defer sp.mutex.Unlock()

	return sp.queueSize
}

func (sp *SubPub) Subscribe(subject string, cb MessageHandler) (*Subscription, error) {
	sp.mutex.Lock()
	defer sp.mutex.Unlock()

	if sp.isClosed {
		return nil, errors.New("closed")
	}

	if _, ok := sp.subjects[subject]; !ok {
		sp.subjects[subject] = make(map[*Subscription]chan<- interface{})
	}

	msgs := make(chan interface{}, sp.queueSize)
	subscription := &Subscription{subject: subject, msgs: msgs, cb: cb, sp: sp}
	sp.subjects[subject][subscription] = msgs

	return subscription, nil
}

func (sp *SubPub) Unsubscribe(s *Subscription) {
	sp.mutex.Lock()
	defer sp.mutex.Unlock()

	if sp.isClosed {
		return
	}

	subscriptions, ok := sp.subjects[s.subject]
	if !ok {
		return
	}

	close(subscriptions[s])
	delete(subscriptions, s)
	if len(subscriptions) == 0 {
		delete(sp.subjects, s.subject)
	}
}

func (sp *SubPub) Publish(subject string, msg interface{}) error {
	sp.mutex.Lock()
	defer sp.mutex.Unlock()

	if sp.isClosed {
		return errors.New("closed")
	}

	subscriptions, ok := sp.subjects[subject]
	if !ok {
		return errors.New("no such subscription")
	}

	for _, msgs := range subscriptions {
		msgs <- msg
	}

	return nil
}

func (sp *SubPub) Close(ctx context.Context) error {
	sp.mutex.Lock()
	defer sp.mutex.Unlock()

	if sp.isClosed {
		return errors.New("closed")
	}

	done := make(chan error)

	go func() {
		for subject := range sp.subjects {
			for subscription, msgs := range sp.subjects[subject] {
				close(msgs)
				delete(sp.subjects[subject], subscription)
			}
			delete(sp.subjects, subject)
		}

		done <- nil
	}()

	select {
	case err := <-done:
		sp.isClosed = true
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
