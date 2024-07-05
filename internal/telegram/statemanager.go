package telegram

import "gopkg.in/telebot.v3"

type StateKey string

type State func(c telebot.Context) (StateKey, error)

type StateManager struct {
	States map[StateKey]State
}

func (s *StateManager) Register(key StateKey, state State) {
	s.States[key] = state
}
