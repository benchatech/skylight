package skylight

type ObserverHandler func(*Event)
type ObserverCondition func(*Event) bool

type Observer struct {
	id      string
	cond    ObserverCondition
	handler ObserverHandler
}

func WildcardObserver(handler ObserverHandler) *Observer {
	return &Observer{
		cond:    func(e *Event) bool { return true },
		handler: handler,
	}
}

func TopicObserver(topic string, handler ObserverHandler) *Observer {
	return &Observer{
		cond:    func(e *Event) bool { return e.topic == topic },
		handler: handler,
	}
}

func LevelObserver(level Level, handler ObserverHandler) *Observer {
	return &Observer{
		cond:    func(e *Event) bool { return e.level <= level },
		handler: handler,
	}
}
