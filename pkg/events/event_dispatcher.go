package events

import (
	"errors"
	"sync"
)

var ErrHandlerAlreadyRegistered = errors.New("handler already registered")

type EventDispatcher struct {
	handlers map[string][]EventHandlerInterface
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[string][]EventHandlerInterface),
	}
}

func (ed *EventDispatcher) Dispatch(event EventInterface) error {
	if _, ok := ed.handlers[event.GetName()]; !ok {
		return nil
	}
	wg := &sync.WaitGroup{}
	for _, handler := range ed.handlers[event.GetName()] {
		wg.Add(1)
		go handler.Handle(event, wg) // rodando de forma assincrona
	}
	wg.Wait()
	return nil
}

func (ed *EventDispatcher) Register(eventName string, handler EventHandlerInterface) error {
	if _, ok := ed.handlers[eventName]; ok {
		for _, h := range ed.handlers[eventName] {
			if h == handler {
				return ErrHandlerAlreadyRegistered
			}
		}
	}
	ed.handlers[eventName] = append(ed.handlers[eventName], handler)
	return nil
}

func (ed *EventDispatcher) Has(eventName string, handler EventHandlerInterface) bool {
	if _, ok := ed.handlers[eventName]; !ok {
		return false
	}
	for _, h := range ed.handlers[eventName] {
		if h == handler {
			return true
		}
	}
	return false
}

func (ed *EventDispatcher) Clear() error {
	ed.handlers = make(map[string][]EventHandlerInterface)
	return nil
}

func (ed *EventDispatcher) Remove(eventName string, handler EventHandlerInterface) error {
	if _, ok := ed.handlers[eventName]; !ok {
		return nil
	}
	for i, h := range ed.handlers[eventName] {
		if h == handler {
			ed.handlers[eventName] = append(ed.handlers[eventName][:i], ed.handlers[eventName][i+1:]...)
			return nil
		}
	}
	return nil
}