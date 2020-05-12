package oasis

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Telenav/osrm-backend/integration/api/oasis"
	"github.com/Telenav/osrm-backend/integration/service/oasis/osrmhelper"
	"github.com/Telenav/osrm-backend/integration/service/oasis/selectionstrategy"
	"github.com/golang/glog"
)

// Handler handles oasis request and provide response
type Handler struct {
	resourceMgr *selectionstrategy.ResourceMgr
}

// New creates new Handler object
func New(osrmBackend, finderType, searchEndpoint, apiKey, apiSignature, dataFolderPath string) (*Handler, error) {
	resourceMgr, err := selectionstrategy.NewResourceMgr(osrmBackend, finderType, searchEndpoint, apiKey, apiSignature, dataFolderPath)
	if err != nil {
		glog.Errorf("Failed to create Handler due to error %+v.\n", err)
		return nil, err
	}

	return &Handler{
		resourceMgr: resourceMgr,
	}, nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	glog.Infof("Handle incoming request %s from remote addr %s", req.RequestURI, req.RemoteAddr)

	// parse oasis request
	oasisReq, err := oasis.ParseRequestURL(req.URL)
	if err != nil || len(oasisReq.Coordinates) != 2 {
		glog.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%v", err)
		return
	}

	// generate route response based on given oasis's orig/destination
	routeResp, err := osrmhelper.RequestRoute4InputOrigDest(oasisReq, h.resourceMgr.OSRMConnector())
	if err != nil {
		glog.Error(err)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%v", err)
		return
	}

	// check whether orig and dest is reachable
	if len(routeResp.Routes) == 0 {
		info := "Orig and destination is not reachable for request " + oasisReq.RequestURI() + "."
		glog.Info(info)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, info)
		return
	}

	// check whether has enough energy
	b, remainRange, err := selectionstrategy.HasEnoughEnergy(oasisReq.CurrRange, oasisReq.SafeLevel, routeResp)
	if err != nil {
		glog.Error(err)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%v", err)
		return
	}
	if b {
		selectionstrategy.GenerateOASISResponse4NoChargeNeeded(w, routeResp, remainRange)
		return
	}

	// check whether could achieve by single charge
	overlap := selectionstrategy.GetOverlapChargeStations4OrigDest(oasisReq, routeResp.Routes[0].Distance, h.resourceMgr)
	if len(overlap) > 0 {
		selectionstrategy.GenerateResponse4SingleChargeStation(w, oasisReq, overlap, h.resourceMgr)
		return
	}

	// generate result for multiple charge
	selectionstrategy.GenerateSolutions4MultipleCharge(w, oasisReq, routeResp, h.resourceMgr)
	return
}

func generateFakeOASISResponse(w http.ResponseWriter, req *oasis.Request) {
	w.WriteHeader(http.StatusOK)

	r := oasis.GenerateFakeResponse(req)

	json.NewEncoder(w).Encode(r)
}
