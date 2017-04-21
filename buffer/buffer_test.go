package buffer

import (
	"testing"
	"time"
)

func TestAddTime(t *testing.T) {
	b := New(3)
	if len(b.ring) != 3 {
		t.Fatal("Length of ring should be 3")
	}

	t1, err := time.Parse(time.RFC3339, "2017-04-21T22:03:00+00:00")
	if err != nil {
		t.Fatal(err)
	}
	t2, err := time.Parse(time.RFC3339, "2017-04-22T22:03:00+00:00")
	if err != nil {
		t.Fatal(err)
	}
	t3, err := time.Parse(time.RFC3339, "2017-04-23T22:03:00+00:00")
	if err != nil {
		t.Fatal(err)
	}
	t4, err := time.Parse(time.RFC3339, "2017-04-24T22:03:00+00:00")
	if err != nil {
		t.Fatal(err)
	}
	t5, err := time.Parse(time.RFC3339, "2017-04-25T22:03:00+00:00")
	if err != nil {
		t.Fatal(err)
	}
	b.AddTime(t1)
	b.AddTime(t2)
	b.AddTime(t3)
	b.AddTime(t4)
	b.AddTime(t5)
	if b.ring[0] != t4 {
		t.Fatalf("Expected %v, got %v", t4, b.ring[0])
	}
	if b.ring[1] != t5 {
		t.Fatalf("Expected %v, got %v", t5, b.ring[1])
	}

	b2 := New(0)
	if len(b2.ring) != 1 {
		t.Fatalf("Ring should be 1 if lower than 1 is given as size. Got %v", len(b2.ring))
	}
}

func TestIsOkay(t *testing.T) {
	t1, err := time.Parse(time.RFC3339, "2017-04-21T22:00:00+00:00")
	if err != nil {
		t.Fatal(err)
	}
	notOkay, err := time.Parse(time.RFC3339, "2017-04-21T23:00:00+00:00")
	if err != nil {
		t.Fatal(err)
	}
	okay, err := time.Parse(time.RFC3339, "2017-04-21T23:00:01+00:00")
	if err != nil {
		t.Fatal(err)
	}
	b := New(2)
	b.AddTime(t1)
	if !b.IsOkay(t1.Add(time.Minute), time.Hour) {
		t.Fatal("If ring is not full any time should be okay")
	}
	b.AddNow()
	if b.IsOkay(notOkay, time.Hour) {
		t.Fatal("Exactly one hour should not be not okay.")
	}
	if !b.IsOkay(okay, time.Hour) {
		t.Fatal("One hour and one second should be cokay.")
	}
}
