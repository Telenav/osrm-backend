package selectionstrategy

import (
	"strconv"
	"testing"

	"github.com/Telenav/osrm-backend/integration/api/osrm/route"
	"github.com/Telenav/osrm-backend/integration/util"
)

func TestHasEnoughEnergyPositive1(t *testing.T) {
	response := &route.Response{
		Routes: []*route.Route{{Distance: 10000.0}},
	}

	currRange := 20000.0
	destRange := 5000.0
	b, remainRange, err := HasEnoughEnergy(currRange, destRange, response)
	if !b || err != nil {
		t.Errorf("Incorrect result generated for TesthasEnoughEnergyPositive1, return value is (%t, %v)", b, err)
	}

	expect := 10000.0
	if !util.Float64Equal(remainRange, expect) {
		t.Errorf("Incorrect remaining range calculated, expect %s while actual value is %s", strconv.FormatFloat(expect, 'f', -1, 64), strconv.FormatFloat(remainRange, 'f', -1, 64))
	}

}

func TestHasEnoughEnergyPositive2(t *testing.T) {
	response := &route.Response{
		Routes: []*route.Route{{Distance: 10000.0}},
	}

	currRange := 10000.0
	destRange := 5000.0
	b, remainRange, err := HasEnoughEnergy(currRange, destRange, response)
	if b || err != nil {
		t.Errorf("Incorrect result generated for TesthasEnoughEnergyPositive1, return value is (%t, %v)", b, err)
	}

	expect := 0.0
	if !util.Float64Equal(remainRange, expect) {
		t.Errorf("Incorrect remaining range calculated, expect %s while actual value is %s", strconv.FormatFloat(expect, 'f', -1, 64), strconv.FormatFloat(remainRange, 'f', -1, 64))
	}
}
