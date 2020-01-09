package oasis

import (
	"encoding/json"
	"net/http"

	"github.com/Telenav/osrm-backend/integration/pkg/api/oasis"
)

type Handler struct {
	osrmBackend string
}

func New(osrmBackend string) *Handler {
	return &Handler{
		osrmBackend,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(generateFakeOasisResponse())
}

func generateFakeOasisResponse() *oasis.Response {
	fakeSolution1 := new(oasis.Solution)
	fakeSolution1.Distance = 90.0
	fakeSolution1.Duration = 300.0
	fakeSolution1.Weight = 300.0
	fakeSolution1.WeightName = "duration"

	fakeStation1 := new(oasis.ChargeStation)
	fakeStation1.Location[0] = 13.39677
	fakeStation1.Location[1] = 52.54366
	fakeStation1.WaitTime = 30.0
	fakeStation1.ChargeTime = 100.0
	fakeStation1.ChargeRange = 100.0
	fakeSolution1.ChargeStations = append(fakeSolution1.ChargeStations, fakeStation1)

	fakeStation2 := new(oasis.ChargeStation)
	fakeStation2.Location[0] = 13.40677
	fakeStation2.Location[1] = 52.53333
	fakeStation2.WaitTime = 100.0
	fakeStation2.ChargeTime = 100.0
	fakeStation2.ChargeRange = 100.0
	fakeSolution1.ChargeStations = append(fakeSolution1.ChargeStations, fakeStation2)

	r := new(oasis.Response)
	r.Code = "200"
	r.Message = "Optimized charge station selection result:"
	r.Solutions = append(r.Solutions, fakeSolution1)

	return r
}
