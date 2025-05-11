# SUBPUB

Требуется реализовать пакет subpub. Нужно написать простую шину событий, работающую по принципу Publisher-Subscriber.
Требования к шине:

- На один subject может подписываться и отписываться множество подписчиков.
- Один медленный подписчик не должен тормозить остальных.
- Нельзя терять порядок сообщений (FIFO очередь).
- Метод `Close` должен учитывать переданный контекст. Если он отменен - выходим сразу, работающие хендлеры оставляем
  работать.
- Горутины (если они будут) течь не должны.

Ниже представлен API пакеты subpub:

```go
package subpub

import "context"

// MessageHandler is a callback function processes messages delivered to subscribers.
type MessageHandler func(msg interface{})

type Subscription interface {
	// Unsubscribe will remove interest in the current subject subscription is for.
	Unsubscribe()
}

type SubPub interface {
	// Subscribe creates an asynchronous queue subscriber on the given subject.
	Subscribe(subject string, cb MessageHandler) (Subscription, error)

	// Publish publishes the msg argument to the given subject.
	Publish(subject string, msg interface{}) error

	// Close will shutdown sub-pub system.
	// May be blocked by data delivery until the context is canceled.
	Close(ctx context.Context) error
}

func NewSubPub() SubPub {
	panic("Implement me")
}
```

К заданию рекомендуется написать unit-тесты.

## Использование

```shell
go get github.com/blvxme/subpub
```
