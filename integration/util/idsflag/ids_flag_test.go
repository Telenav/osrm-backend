package idsflag

import (
	"fmt"
	"reflect"
	"testing"
)

func TestIDsFlagString(t *testing.T) {
	cases := []struct {
		idsFlag IDs
		s       string
	}{
		{IDs{}, ""},
		{IDs{829733412}, "829733412"},
		{IDs{-104489539}, "-104489539"},
		{IDs{829733412, -104489539}, "829733412,-104489539"},
	}

	for _, c := range cases {
		s := c.idsFlag.String()
		if s != c.s {
			t.Errorf("idsFlag: %v, expect string %s, but got %s", c.idsFlag, c.s, s)
		}
	}
}

func TestIDsFlagSet(t *testing.T) {
	cases := []struct {
		idsFlag IDs
		s       string
		err     error
	}{
		{IDs{}, "", nil},
		{IDs{829733412}, "+829733412", nil},
		{IDs{829733412}, "829733412", nil},
		{IDs{-104489539}, "-104489539", nil},
		{IDs{829733412, -104489539}, "829733412,-104489539", nil},
		{IDs{}, "a", fmt.Errorf("error")},
	}

	for _, c := range cases {
		idsFlagValue := IDs{}
		err := idsFlagValue.Set(c.s)
		if err != nil && c.err != nil {
			continue
		} else if err == nil && c.err == nil {
			// compare wayIDs slice
			if !reflect.DeepEqual(idsFlagValue, c.idsFlag) {
				t.Errorf("string %s, expect %v, but got %v", c.s, c.idsFlag, idsFlagValue)
			}
		} else {
			t.Errorf("string %s, expect err %v, but got %v", c.s, c.err, err)
		}
	}
}
