package entrypoint

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Telenav/osrm-backend/integration/api/oasis"
	"github.com/Telenav/osrm-backend/integration/service/oasis/solution"
	"github.com/golang/glog"
)

func generateResponse4UnableReachDestination(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(http.StatusOK)
	r := new(oasis.Response)
	r.Message = "StatusCode = " + strconv.Itoa(statusCode) + "(" + solution.StatusText(statusCode) + ")"
	json.NewEncoder(w).Encode(r)
}

func generateResponse4BadRequest(w http.ResponseWriter, err error) {
	glog.Error(err)
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, "%v", err)
}

func generateResponseWhenMetError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	r := new(oasis.Response)
	r.Code = strconv.Itoa(statusCode)
	r.Message = "StatusCode = " + strconv.Itoa(statusCode) + "(" + solution.StatusText(statusCode) + "), met error " + err.Error()
	json.NewEncoder(w).Encode(r)
}

func generateResponse(w http.ResponseWriter, statusCode int, solutions []*oasis.Solution) {
	w.WriteHeader(http.StatusOK)
	r := new(oasis.Response)
	r.Code = strconv.Itoa(statusCode)
	r.Message = "StatusCode = " + strconv.Itoa(statusCode) + "(" + solution.StatusText(statusCode) + ")"
	r.Solutions = append(r.Solutions, solutions...)
	json.NewEncoder(w).Encode(r)
}

// generateFakeResponse generate a static oasis response for initial integration and testing
func generateFakeResponse(w http.ResponseWriter, req *oasis.Request) {
	w.WriteHeader(http.StatusOK)
	r := oasis.GenerateFakeResponse(req)
	json.NewEncoder(w).Encode(r)
}
