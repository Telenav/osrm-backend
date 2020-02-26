package nametable

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/Telenav/osrm-backend/integration/osrmfiles/osrmtype"
)

func TestIndexedData(t *testing.T) {
	blocksBytes, err := ioutil.ReadFile("testdata/blocks")
	if err != nil {
		t.Error(err)
	}
	valuesBytes, err := ioutil.ReadFile("testdata/values")
	if err != nil {
		t.Error(err)
	}

	indexedData := IndexedData{
		BlocksMeta:   14953,
		ValuesMeta:   862764,
		BlocksBuffer: *bytes.NewBuffer(blocksBytes),
		ValuesBuffer: *bytes.NewBuffer(valuesBytes),
	}
	if err := indexedData.Assemble(); err != nil {
		t.Error(err)
	}
	if err := indexedData.Validate(); err != nil {
		t.Error(err)
	}

	invalidCases := []osrmtype.NameID{
		1, 2, 3, 4,
		101, 102, 103, 104,
		14953 * blockContainingSize, // more than indexedData stored
		osrmtype.InvalidNameID,
	}

	for _, c := range invalidCases {
		if n, err := indexedData.GetNamesForID(c); err == nil {
			t.Errorf("GetNamesForID(%d), expect failed but got %#v", c, n)
		}
	}

	cases := []struct {
		osrmtype.NameID
		Names
	}{
		{0, Names{}},
		{5, Names{Name: "Bruce Woodbury Beltway", Destinations: "", Pronuciation: "", Ref: "CR 215", Exits: ""}},
		{45, Names{Name: "Gulch Road", Destinations: "", Pronuciation: "", Ref: "", Exits: ""}},
		{730, Names{Name: "Summerlin Parkway", Destinations: "NV 613 West: Summerlin South", Pronuciation: "", Ref: "NV 613", Exits: "81A"}},
		{200260, Names{Name: "Summerlin Parkway", Destinations: "US 95 South: Downtown Las Vegas", Pronuciation: "", Ref: "NV 613", Exits: ""}},
		{244625, Names{Name: "Prater Way", Destinations: "Reno", Pronuciation: "", Ref: "I 80 BUS; NV 647", Exits: ""}},
		{246045, Names{Name: "North Rancho Drive", Destinations: "US 95 North: Tonopah, Reno", Pronuciation: "", Ref: "US 95 BUS; NV 599", Exits: ""}},
	}

	for _, c := range cases {
		n, err := indexedData.GetNamesForID(c.NameID)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(n, c.Names) {
			t.Errorf("GetNamesForID(%d), expect %#v but got %#v", c.NameID, c.Names, n)
		}
	}

}
