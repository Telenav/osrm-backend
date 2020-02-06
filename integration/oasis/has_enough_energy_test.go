package oasis

import (
	"strconv"
	"testing"

	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/route"
)

func TesthasEnoughEnergyPositive1(t *testing.T) {
	response := &route.Response{
		Routes: []*route.Route{&route.Route{Distance: 10000.0}},
	}

	currRange := 20000.0
	destRange := 5000.0
	b, remainRange, err := hasEnoughEnergy(currRange, destRange, response)
	if !b || err != nil {
		t.Errorf("Incorrect result generated for TesthasEnoughEnergyPositive1, return value is (%t, %v)", b, err)
	}

	expect := 10000.0
	if !floatEquals(remainRange, expect) {
		t.Errorf("Incorrect remaining range calculated, expect %s while actual value is %s", strconv.FormatFloat(expect, 'f', -1, 64), strconv.FormatFloat(remainRange, 'f', -1, 64))
	}

}

func TesthasEnoughEnergyPositive2(t *testing.T) {
	response := &route.Response{
		Routes: []*route.Route{&route.Route{Distance: 10000.0}},
	}

	currRange := 10000.0
	destRange := 5000.0
	b, remainRange, err := hasEnoughEnergy(currRange, destRange, response)
	if b || err != nil {
		t.Errorf("Incorrect result generated for TesthasEnoughEnergyPositive1, return value is (%t, %v)", b, err)
	}

	expect := 0.0
	if !floatEquals(remainRange, expect) {
		t.Errorf("Incorrect remaining range calculated, expect %s while actual value is %s", strconv.FormatFloat(expect, 'f', -1, 64), strconv.FormatFloat(remainRange, 'f', -1, 64))
	}
}

var epsilon float64 = 0.00000001

func floatEquals(a, b float64) bool {
	if (a-b) < epsilon && (b-a) < epsilon {
		return true
	}
	return false
}
