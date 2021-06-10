package gfind

import (
	"log"
	"testing"
)

func Test_Push(t *testing.T) {
	sq := &StringQueue{}
	expected := "foobar"
	sq.Push(expected)
	got, err := sq.Front()
	if err != nil {
		log.Fatal(err)
	}
	if got != expected {
		log.Fatalf("Front() returned %s, expected %s", got, expected)
	}
}

func Test_FrontEmpty(t *testing.T) {
	sq := &StringQueue{}
	_, err := sq.Front()
	if err == nil {
		log.Fatal("Front() on an empty queue does not error")
	}
}

func Test_FrontNonEmpty(t *testing.T) {
	sq := &StringQueue{}
	actual := "foobar"
	sq.Push(actual)
	got, err := sq.Front()
	if err != nil {
		log.Fatal(err)
	}
	if got != actual {
		log.Fatalf("Front() returned %s, expected %s", got, actual)
	}
}

func Test_PopEmpty(t *testing.T) {
	sq := &StringQueue{}
	_, err := sq.Pop()
	if err == nil {
		log.Fatal("Pop() on an empty queue does not error")
	}
}

func Test_Pop(t *testing.T) {
	sq := &StringQueue{}
	actual := "foobar"
	sq.Push(actual)
	got, err := sq.Pop()
	if err != nil {
		log.Fatal(err)
	}
	if got != actual {
		log.Fatalf("Pop() returned %s, expected %s", got, actual)
	}
}

func Test_SizeEmpty(t *testing.T) {
	sq := &StringQueue{}
	got := sq.Size()
	expected := 0
	if got != expected {
		log.Fatalf("Size() returned %d, expected %d", got, expected)
	}
}

func Test_SizeNonEmpty(t *testing.T) {
	sq := &StringQueue{}
	sq.Push("foo")
	got := sq.Size()
	expected := 1
	if got != expected {
		log.Fatalf("Size() returned %d, expected %d", got, expected)
	}
}

func Test_EmptyEmpty(t *testing.T) {
	sq := &StringQueue{}
	expected := true
	got := sq.Empty()
	if got != expected {
		log.Fatalf("Empty() returned %t, expected %t", got, expected)
	}
}

func Test_EmptyNonEmpty(t *testing.T) {
	sq := &StringQueue{}
	sq.Push("foo")
	expected := false
	got := sq.Empty()
	if got != expected {
		log.Fatalf("Empty() returned %t, expected %t", got, expected)
	}
}
