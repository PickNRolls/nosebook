package cache

import (
	"testing"
)

func toSlice(lru *LRU[string, string]) []string {
	out := []string{}

	elem := lru.l.Front()
	for elem != nil {
		node := elem.Value.(node[string])
		value := node.value.(string)
		out = append(out, value)

		elem = elem.Next()
	}

	return out
}

func match(values []string, lru *LRU[string, string]) bool {
	lruValues := toSlice(lru)

	if len(lruValues) != len(values) {
		return false
	}

	for i := 0; i < len(values); i++ {
		if values[i] != lruValues[i] {
			return false
		}
	}

	return true
}

func TestLRU(t *testing.T) {
	lru := NewLRU[string, string](2)

	lru.Set("1", "1")
	lru.Set("2", "2")

	if lru.l.Len() != 2 {
		t.Fatalf("list length is wrong, want 2, has %v\n", lru.l.Len())
	}

	{
		want := []string{"2", "1"}
		if !match(want, lru) {
			t.Fatalf("inner linked list does not match\nwant %v\nhas %v", want, toSlice(lru))
		}
	}

  lru.Set("3", "3")
  
	{
		want := []string{"3", "2"}
		if !match(want, lru) {
			t.Fatalf("inner linked list does not match\nwant %v\nhas %v", want, toSlice(lru))
		}
	}

  value, _ := lru.Get("2")
  if value != "2" {
		t.Fatalf("list get is wrong, want 2, has %v\n", value)
  }
  
	{
		want := []string{"2", "3"}
		if !match(want, lru) {
			t.Fatalf("inner linked list does not match\nwant %v\nhas %v", want, toSlice(lru))
		}
	}

  lru.Set("3", "4")
  
	{
		want := []string{"4", "2"}
		if !match(want, lru) {
			t.Fatalf("inner linked list does not match\nwant %v\nhas %v", want, toSlice(lru))
		}
	}
}
