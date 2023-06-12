package dbstore

import "context"

type Storage struct { // TODO
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Init(ctx context.Context) error {
	// TODO
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	// TODO
	return nil
}
