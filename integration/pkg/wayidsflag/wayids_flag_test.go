package wayidsflag

import (
	"fmt"
	"reflect"
	"testing"
)

func TestWayIDsFlagString(t *testing.T) {
	cases := []struct {
		wayIDsFlag WayIDs
		s          string
	}{
		{WayIDs{}, ""},
		{WayIDs{829733412}, "829733412"},
		{WayIDs{-104489539}, "-104489539"},
		{WayIDs{829733412, -104489539}, "829733412,-104489539"},
	}

	for _, c := range cases {
		s := c.wayIDsFlag.String()
		if s != c.s {
			t.Errorf("wayIDsFlag: %v, expect string %s, but got %s", c.wayIDsFlag, c.s, s)
		}
	}
}

func TestWayIDsFlagSet(t *testing.T) {
	cases := []struct {
		wayIDsFlag WayIDs
		s          string
		err        error
	}{
		{WayIDs{}, "", nil},
		{WayIDs{829733412}, "+829733412", nil},
		{WayIDs{829733412}, "829733412", nil},
		{WayIDs{-104489539}, "-104489539", nil},
		{WayIDs{829733412, -104489539}, "829733412,-104489539", nil},
		{WayIDs{}, "a", fmt.Errorf("error")},
	}

	for _, c := range cases {
		wayIDsFlagValue := WayIDs{}
		err := wayIDsFlagValue.Set(c.s)
		if err != nil && c.err != nil {
			continue
		} else if err == nil && c.err == nil {
			// compare wayIDs slice
			if !reflect.DeepEqual(wayIDsFlagValue, c.wayIDsFlag) {
				t.Errorf("string %s, expect %v, but got %v", c.s, c.wayIDsFlag, wayIDsFlagValue)
			}
		} else {
			t.Errorf("string %s, expect err %v, but got %v", c.s, c.err, err)
		}
	}
}
