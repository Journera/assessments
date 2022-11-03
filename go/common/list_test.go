package common

import "testing"

func TestSort(t *testing.T) {
	list := NewLinkedList[string]()
	list.AddLast("ccc")
	list.AddLast("aaa")
	list.AddLast("zzz")
	list.AddLast("jjj")
	list.Sort(func(i, j string) int {
		if i < j {
			return -1
		}
		if i > j {
			return 1
		}
		return 0
	})

	c := list.Iter()
	if <-c != "aaa" {
		t.Error("first value does not equal aaa")
	}
	if <-c != "ccc" {
		t.Error("second value does not equal ccc")
	}
	if <-c != "jjj" {
		t.Error("third value does not equal jjj")
	}
	if <-c != "zzz" {
		t.Error("last value does not equal zzz")
	}
}
