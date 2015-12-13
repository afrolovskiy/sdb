package sdb

import (
	"reflect"
	"testing"
)

func TestStorage(t *testing.T) {
	s := NewStorage()

	if s.HasVariable("a") != false {
		t.Errorf("HasVariable must return false if variable does not exist")
	}

	if s.Get("a") != "" {
		t.Errorf("Get must return empty string if variable does not exist")
	}

	s.Set("a", "10")

	if s.HasVariable("a") != true {
		t.Errorf("HasVariable must return true if variable exists")
	}

	if s.Get("a") != "10" {
		t.Errorf("Get must return variable value if variable exists")
	}

	s.Unset("a")

	if s.HasVariable("a") != false {
		t.Errorf("HasVariable must return false if variable does not exist")
	}

	if s.Get("a") != "" {
		t.Errorf("Get must return empty string if variable does not exist")
	}
}

func TestStorageNumEqualTo(t *testing.T) {
	s := NewStorage()

	if s.NumEqualTo("10") != 0 {
		t.Errorf("NumEqualTo must return 0 if no variables equal that value")
	}

	s.Set("a", "10")
	s.Set("b", "10")
	s.Set("c", "20")

	if s.NumEqualTo("10") != 2 {
		t.Errorf("NumEqualTo must return the number of variables that are currently set to value")
	}
}

func TestStorageCopy(t *testing.T) {
	s := NewStorage()
	s.Set("a", "10")

	ns := s.Copy()

	if reflect.DeepEqual(s, ns) != true {
		t.Errorf("Copy must return full storage copy")
	}

	s.Set("b", "10")

	if reflect.DeepEqual(s, ns) != false {
		t.Errorf("Storage copy must not see changes of parent storage")
	}
}
