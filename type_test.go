package nomad

import "testing"

func TestRandomListingGen(t *testing.T) {
	L := NewListingRand()

	if L.Arrive.DateTime.Before(L.Depart.DateTime) {
		t.Fail()
	}
	if L.Arrive.Location == L.Depart.Location {
		t.Fail()
	}

}
