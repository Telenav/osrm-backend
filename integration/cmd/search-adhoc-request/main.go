package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/Telenav/osrm-backend/integration/pkg/api/search/nearbychargestation"
	"github.com/Telenav/osrm-backend/integration/pkg/api/search/searchcoordinate"
	"github.com/Telenav/osrm-backend/integration/pkg/backend"
)

func main() {
	flag.Parse()

	if len(flags.entityEndpoint) == 0 {
		fmt.Println("[ERROR] -entity should not be empty!")
		return
	}

	if len(flags.apiKey) == 0 {
		fmt.Println("[ERROR] -apikey should not be empty.")
		return
	}
	if len(flags.apiSignature) == 0 {
		fmt.Println("[ERROR] -apisignature should not be empty.")
		return
	}
	req := nearbychargestation.NewRequest()
	req.APIKey = flags.apiKey
	req.APISignature = flags.apiSignature
	req.Location = searchcoordinate.Coordinate{Lat: 37.78509, Lon: -122.41988}

	clt := http.Client{Timeout: backend.Timeout()}
	requestURL := flags.entityEndpoint + req.RequestURI()
	fmt.Println(requestURL)
	resp, err := clt.Get(requestURL)
	if err != nil {
		fmt.Printf("route request %s against search failed, err %v", requestURL, err)
		return
	}
	defer resp.Body.Close()

	var response nearbychargestation.Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	jsonResult, _ := json.Marshal(response)
	fmt.Println(string(jsonResult))
}
