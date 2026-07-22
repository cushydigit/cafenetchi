package otp

import "context"

type fakeStore struct {
	values map[string]string

	setErr error
	getErr error
	delErr error
}

func newFakeStore() *fakeStore {
	return &fakeStore{
		values: make(map[string]string),
	}
}

func (f *fakeStore) Set(ctx context.Context, phone, otp string) error {

	if f.setErr != nil {
		return f.setErr
	}

	f.values[phone] = otp
	return nil
}

func (f *fakeStore) Get(ctx context.Context, phone string) (string, error) {
	if f.getErr != nil {
		return "", f.getErr
	}

	code, ok := f.values[phone]
	if !ok {
		return "", ErrNotFound
	}

	return code, nil
}

func (f *fakeStore) Del(ctx context.Context, phone string) error {
	if f.delErr != nil {
		return f.delErr
	}

	delete(f.values, phone)
	return nil
}
