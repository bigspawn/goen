package main

import (
	"errors"
	"testing"
)

var ErrSome = errors.New("some error")

func Test_defer(t *testing.T) {
	err := do()
	if !errors.Is(err, ErrSome) {
		t.Logf("%v", err)
		t.Fail()
	}
}

func do() (err error) {

	defer func() {
		err = ErrSome
	}()

	return nil
}
