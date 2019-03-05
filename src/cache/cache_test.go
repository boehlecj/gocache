package main

import (
	"testing"
)

func TestGetAction(t *testing.T) {
	t.Log("expecting an action of 0")
	action := getAction("get")

	if action != 0 {
		t.Errorf("Expected 0, but it was %d instead.", action)
	}
}

func TestPutAction(t *testing.T) {
	t.Log("expecting an action of 1")
	action := getAction("put")

	if action != 1 {
		t.Errorf("Expected 1, but it was %d instead.", action)
	}
}

func TestExitAction(t *testing.T) {
	t.Log("expecting an action of 2")
	action := getAction("quit")

	if action != 2 {
		t.Errorf("Expected 2, but it was %d instead.", action)
	}
}

func TestCache(t *testing.T) {
	t.Log("expecting cached value to be 12345")
  cache := NewCache()
	cache.put("test1","12345")

  response, found := cache.get("test1")
  expected := string(response.data[:])

	if found && expected != "12345" {
		t.Errorf("Expected 12345, but it was %s instead.", expected)
	}
}
