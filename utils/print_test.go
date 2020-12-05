package utils

import "testing"

type Tmp0 struct {
	Name  string
	Age   int
	extra string
}

type Tmp struct {
	Name  string
	Age   int
	extra string
	T     Tmp0
}

func TestBeautifyPrint(t *testing.T) {
	tmp0 := Tmp0{
		Name:  "happyxhw",
		Age:   22,
		extra: "extra",
	}
	tmp := &Tmp{
		Name:  "happyxhw",
		Age:   22,
		extra: "extra",
		T:     tmp0,
	}
	PrintStruct(tmp, false)
	PrintStruct(tmp, true)
}
