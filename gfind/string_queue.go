package gfind

import (
	"fmt"
)

type NoSuchElementError string

func (e NoSuchElementError) Error() string {
	return fmt.Sprintf("%s", e)
}

type StringQueue struct {
	values []string
}

func (sq *StringQueue) Front() (string, error) {
	if len(sq.values) < 1 {
		return "", NoSuchElementError("Cannot call Front() on an empty queue")
	}
	return sq.values[0], nil
}

func (sq *StringQueue) Push(s string) {
	sq.values = append(sq.values, s)
}

func (sq *StringQueue) Pop() (popped string, err error) {
	popped, err = sq.Front()
	if err != nil {
		return "", NoSuchElementError("Cannot call Pop() on an empty queue")
	}
	sq.values[0] = ""
	sq.values = sq.values[1:]
	return
}

func (sq *StringQueue) Size() int {
	return len(sq.values)
}

func (sq *StringQueue) Empty() bool {
	return sq.Size() == 0
}
