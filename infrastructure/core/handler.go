package core

import "context"

type Option interface{}

type dataStore struct {
	dataItem interface{}
	dataType interface{}
}

type commandHandler struct {
	name        string
	description string
	data        []dataStore
}

func (ch *commandHandler) Description(d string) Option {
	ch.description = d
	return Option{}
}

func CommandHandler(c context.Context, opts ...Option) error {

}
