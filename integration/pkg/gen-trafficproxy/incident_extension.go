package proxy

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
)

// CSVString represents Incident as defined CSV format.
func (i *Incident) CSVString() string {
	return i.csvString(false)
}

// HumanFriendlyCSVString represents Incident as defined CSV format, but prefer human friendly string instead of integer/boolean.
func (i *Incident) HumanFriendlyCSVString() string {
	return i.csvString(true)
}

func (i *Incident) csvString(humanFriendly bool) string {
	records := []string{}
	records = append(records, i.IncidentId)

	affectedWayIDsString := []string{}
	for _, wayID := range i.AffectedWayIds {
		affectedWayIDsString = append(affectedWayIDsString, strconv.FormatInt(wayID, 10))
	}
	records = append(records, strings.Join(affectedWayIDsString, ","))

	if humanFriendly {
		records = append(records, i.IncidentType.String(), i.IncidentSeverity.String())
	} else {
		records = append(records, strconv.Itoa(int(i.IncidentType)), strconv.Itoa(int(i.IncidentSeverity)))
	}

	if i.IncidentLocation == nil {
		records = append(records, "", "")
	} else {
		records = append(records, fmt.Sprintf("%f", i.IncidentLocation.Lat), fmt.Sprintf("%f", i.IncidentLocation.Lon))
	}

	records = append(records, i.Description, i.FirstCrossStreet, i.SecondCrossStreet, i.StreetName)
	records = append(records, strconv.Itoa(int(i.EventCode)), strconv.Itoa(int(i.AlertCEventQuantifier)))

	if humanFriendly {
		records = append(records, strconv.FormatBool(i.IsBlocking))
	} else {
		isBlockingInteger := 0
		if i.IsBlocking {
			isBlockingInteger = 1
		}
		records = append(records, strconv.Itoa(isBlockingInteger))
	}

	var buff bytes.Buffer
	w := csv.NewWriter(&buff)
	w.UseCRLF = false
	w.Write(records) // string only operation, don't need to handle error
	w.Flush()
	s := strings.TrimRight(buff.String(), "\n")

	return s
}
