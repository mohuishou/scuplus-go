package sculibrary

import (
	"testing"
)

func TestLibrary_getURL(t *testing.T) {
	t.Log(getURL())
}

func TestLibrary_GetLoan(t *testing.T) {
	lib, err := NewLibrary("x", "x")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(lib.GetLoan())
}

func TestLibrary_GetLoanAll(t *testing.T) {
	lib, err := NewLibrary("xx", "xx")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(lib.GetLoanAll())
}

func TestLibrary_Loan(t *testing.T) {
	lib, err := NewLibrary("xx", "xx")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(lib.Loan(""))
}
