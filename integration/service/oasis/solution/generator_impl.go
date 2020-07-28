package solution

import (
	"fmt"

	"github.com/Telenav/osrm-backend/integration/api/oasis"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/osrmhelper"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/resourcemanager"
	"github.com/Telenav/osrm-backend/integration/service/oasis/solution/selectionstrategy"
)

type generatorImpl struct {
	resourceMgr *resourcemanager.ResourceMgr
}

// NewGeneratorImpl creates implementation of interface Generator
func NewGeneratorImpl(resourceMgr *resourcemanager.ResourceMgr) *generatorImpl {
	return &generatorImpl{
		resourceMgr: resourceMgr,
	}
}

// Generate calculates oasis solutions based on oasis request
func (g *generatorImpl) Generate(oasisReq *oasis.Request) (int, []*oasis.Solution, error) {
	// generate route response based on given oasis's orig/destination
	routeResp, err := osrmhelper.RequestRoute4InputOrigDest(oasisReq, g.resourceMgr.OSRMConnector())
	if err != nil {
		return StatusFailedToCalculateRoute, nil, err
	}

	// check whether orig and dest is reachable
	if len(routeResp.Routes) == 0 {
		return StatusOrigAndDestIsNotReachable, nil, nil
	}

	// check whether has enough energy
	b, energyLeft, err := selectionstrategy.HasEnoughEnergy(oasisReq.CurrRange, oasisReq.SafeLevel, routeResp)
	if err != nil {
		return StatusFailedToGenerateChargeResult, nil, err
	}
	if b {
		solutions, err := selectionstrategy.GenerateSolution4NoChargeNeeded(routeResp, energyLeft)
		return StatusNoNeedCharge, solutions, err
	}

	// check whether orig and dest reachable range covers shared stations
	overlap := selectionstrategy.GetOverlapChargeStations4OrigDest(oasisReq, routeResp.Routes[0].Distance, g.resourceMgr)
	if len(overlap) > 0 {
		solutions, err := selectionstrategy.GenerateResponse4SingleChargeStation(nil, oasisReq, overlap, g.resourceMgr)
		return StatusChargeForSingleTime, solutions, err
	}

	// generate result for multiple charge
	solutions, err := selectionstrategy.GenerateSolutions4ChargeStationBasedRoute(oasisReq, g.resourceMgr)
	if len(solutions) == 0 && err != nil {
		err = fmt.Errorf("Failed to generate charge solution, internal error.\n")
	}
	return StatusChargeForMultipleTime, solutions, err
}
