package testdata

import (
	"fmt"
	"time"
)

type Foo struct {
	ID   int
	Name string
}

type Bar struct {
	ID        int
	Age       int
	name      string
	createdAt time.Time
}

type Buzz struct {
	id   int
	name string
}

func CheckFoo() {
	f := Foo{ID: 1, Name: "knsh14"}
	f = Foo{2, "knsh14"}
	ff := &Foo{ID: 2, Name: "knsh14"}
	ff = &Foo{3, "knsh14"}
	fmt.Println(f)
	fmt.Println(ff)
}
