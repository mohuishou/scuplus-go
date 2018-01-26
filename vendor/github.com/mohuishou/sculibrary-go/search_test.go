package sculibrary

import (
	"testing"
)

func TestLibrary_Search(t *testing.T) {
	t.Log(Search("php", "WRD", "http://opac.scu.edu.cn:8080/F/GEL6R9V1M2BAEFITQVS6B8LV61EJ2BBVPFQU4LNIUT36RD2FV3-09149?func=short-jump&jump=41&pag=now"))
}
