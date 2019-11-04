package proxy

import (
	"bytes"
	"encoding/csv"
	fmt "fmt"
	"strconv"
	"strings"
)

// CSVString represents Flow as defined CSV format.
// I.e. 'wayID,Speed,int(TrafficLevel)'
func (f *Flow) CSVString() string {
	return fmt.Sprintf("%d,%f,%d", f.WayId, f.Speed, int32(f.TrafficLevel))
}

// CSVString represents Incident as defined CSV format.
func (i *Incident) CSVString() string {

	records := []string{}
	records = append(records, i.IncidentId)

	affectedWayIDsString := []string{}
	for _, wayID := range i.AffectedWayIds {
		affectedWayIDsString = append(affectedWayIDsString, strconv.FormatInt(wayID, 10))
	}
	records = append(records, strings.Join(affectedWayIDsString, ","))

	records = append(records, i.IncidentType.String(), i.IncidentSeverity.String())
	records = append(records, fmt.Sprintf("%f", i.IncidentLocation.Lat), fmt.Sprintf("%f", i.IncidentLocation.Lon))
	records = append(records, i.Description, i.FirstCrossStreet, i.SecondCrossStreet, i.StreetName)
	records = append(records, strconv.Itoa(int(i.EventCode)), strconv.Itoa(int(i.AlertCEventQuantifier)))
	records = append(records, strconv.FormatBool(i.IsBlocking))

	var buff bytes.Buffer
	w := csv.NewWriter(&buff)
	w.UseCRLF = false
	w.Write(records) // string only operation, don't need to handle error
	w.Flush()
	s := strings.TrimRight(buff.String(), "\n")

	return s
}
