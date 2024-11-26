package value

import (
	"reflect"
	"testing"

)

func TestRopeValueUpdate(t *testing.T) {

	t.Run("insert a ropeValue in the middle", func(t *testing.T) {
		rp := NewRopeValue("ropeBuffer", "h!")
		content := "i"
		want := "hi!"
		rp.Update(1, content, 0)
		got := rp.String()
		if !reflect.DeepEqual(want, got) {
			t.Errorf("Got %v, expected %v", got, want)
		}
	})
	t.Run("Appending a ropeValue", func(t *testing.T) {
		rp := NewRopeValue("ropeBuffer", "hi")
		content := "!"
		want := "hi!"
		rp.Update(5, content, 0)
		got := rp.String()
		if !reflect.DeepEqual(want, got) {
			t.Errorf("Got %v, expected %v", got, want)
		}
	})
	t.Run("Prepending a ropeValue", func(t *testing.T) {
		rp := NewRopeValue("ropeBuffer", "i!")
		content := "h"
		want := "hi!"
		rp.Update(0, content, 0)
		got := rp.String()
		if !reflect.DeepEqual(want, got) {
			t.Errorf("Got %v, expected %v", got, want)
		}
	})
	t.Run("append after empty space", func(t *testing.T) {
		rp := NewRopeValue("ropeBuffer", "Hi")
		rp.Update(999, " ", 0) //append
		rp.Update(999, "!", 0)
		want := NewRopeValue("ropeBuffer", "Hi !")
		got := rp
		if !reflect.DeepEqual(want, got) {
			t.Errorf("Expected %v, got %v", want.String(), got.String())
		}
	})
	t.Run("deleting a value", func(t *testing.T) {
		rp := NewRopeValue("ropeBuffer", "0123456789")
		_,err := rp.Update(1, "", 3)
		if err != nil {
			t.Fatal("Failed to update")
		}
		want := "0456789"
		got := rp.String()
		if want != got {
			t.Errorf("Expected %v, found %v", want, got)
		}
	})
	t.Run("deleting the whole buffer", func(t *testing.T) {
		rp := NewRopeValue("ropeBuffer", "0123456789")
		_,err := rp.Update(0, "", 10)
		if err != nil {
			t.Fatal("Failed to update")
		}
		want := ""
		got := rp.String()
		if want != got {
			t.Errorf("Expected %v, found %v", want, got)
		}
	})
	t.Run("inserting and deleting", func(t *testing.T) {
		rp := NewRopeValue("ropeBuffer", "0123456789")
		_,err := rp.Update(1, "000", 4)
		if err != nil {
			t.Fatal("Failed to update")
		}
		want := "000056789"
		got := rp.String()
		if want != got {
			t.Errorf("Expected %v, found %v", want, got)
		}
	})
}
