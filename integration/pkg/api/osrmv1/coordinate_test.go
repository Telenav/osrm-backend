package osrmv1

import "testing"

func TestCoordinates(t *testing.T) {
	cases := []struct {
		s string
		c Coordinates
	}{
		{"13.388860,52.517037", Coordinates{Coordinate{52.517037, 13.388860}}},
		{
			"13.388860,52.517037;13.397634,52.529407;13.428555,52.523219",
			Coordinates{
				Coordinate{52.517037, 13.388860},
				Coordinate{52.529407, 13.397634},
				Coordinate{52.523219, 13.428555},
			},
		},
	}

	for _, c := range cases {
		s := c.c.String()
		if s != c.s {
			t.Errorf("%v String(), expect %s, but got %s", c.c, c.s, s)
		}
	}
}
