package extruct

import (
	"fmt"
	"testing"
)

type Foo struct {
	Bar *Bar
	Qux []*Bar
}

type Bar struct {
	Baz string
	Bub *string
}

func TestNestedStruct(t *testing.T) {
	t.Parallel()
	f := Foo{
		Bar: &Bar{
			Baz: "baz",
		},
	}

	path := "Bar.Baz"
	want := "baz"
	v, err := Extruct(f, path)
	if err != nil {
		t.Fatalf("unexpected error: %q", err)
	}
	if str, ok := v.(string); !ok {
		t.Fatalf("got %T, want string for path %q", v, path)
	} else if str != want {
		t.Fatalf("got %q, want %q for path %q", str, want, path)
	}

	path = "Baz"
	want = "foo"
	v, err = Extruct(&Bar{Baz: "foo"}, path)
	if err != nil {
		t.Fatalf("unexpected error: %q", err)
	}
	if str, ok := v.(string); !ok {
		t.Fatalf("got %T, want string for path %q", v, path)
	} else if str != want {
		t.Fatalf("got %q, want %q for path %q", str, want, path)
	}

	f = Foo{
		Qux: []*Bar{
			&Bar{Baz: "1"},
			&Bar{Baz: "2"},
			&Bar{Baz: "3"},
			&Bar{Baz: "4"},
		},
	}
	path = "Qux.Baz"
	v, err = Extruct(f, path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if intSlice, ok := v.([]interface{}); !ok {
		t.Fatalf("got %T, want []interface{} for path %q", v, path)
	} else {
		if len(intSlice) != 4 {
			t.Fatalf("got %d elements, want 4", len(intSlice))
		}
		for i, iv := range intSlice {
			want := fmt.Sprint(i + 1)
			if iv.(string) != want {
				t.Fatalf("got %q, want %q for offset %d", iv.(string), want, i)
			}
		}
	}

	f = Foo{
		Bar: &Bar{},
	}
	path = "Bar.Bub"
	v, err = Extruct(f, path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strPtr, ok := v.(*string); !ok {
		t.Fatalf("got %T/%v, expected *string", v, v)
	} else if strPtr != nil {
		t.Fatalf("got v -> %q, expected v -> nil", *strPtr)
	}
}
