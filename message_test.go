package tailor

import (
	"encoding/json"
	"reflect"
	"testing"
)

func Test_NewMessage(t *testing.T) {

	m := NewMessage("foo.log", "foo")

	// returns pointer
	if reflect.ValueOf(m).Kind() != reflect.Ptr {
		t.Errorf("expected *Message, got %v", reflect.ValueOf(m).Kind())
	}

	// is not nil
	if m == nil {
		t.Errorf("expected non-nil pointer, got %v", m)
	}

	// returns expected source and body
	if m.Source != "foo.log" {
		t.Errorf("expected source to equal 'foo.log', got %v", m.Source)
	}

	if m.Body != "foo" {
		t.Errorf("expected body to equal 'foo', got %v", m.Body)
	}

}

func Test_Marshal(t *testing.T) {

	m := NewMessage("foo.log", "foo")

	_, err := json.Marshal(m)

	if err != nil {
		t.Errorf("expect no errors, got %v", err)
	}

}
