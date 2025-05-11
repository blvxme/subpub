package subpub

import "context"

type Subscription struct {
	subject string
	msgs    <-chan interface{}
	cb      MessageHandler
	sp      *SubPub
}

func (s *Subscription) HandleMessages() {
	for msg := range s.msgs {
		s.cb(msg)
	}
}

func (s *Subscription) HandleMessagesWithContext(ctx context.Context) {
	for {
		select {
		case msg, ok := <-s.msgs:
			if !ok {
				return
			}
			s.cb(msg)
		case <-ctx.Done():
			return
		}
	}
}

func (s *Subscription) Unsubscribe() {
	s.sp.Unsubscribe(s)
}
