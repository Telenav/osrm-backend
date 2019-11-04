package proxy

import fmt "fmt"

// CSVString represents Flow as defined CSV format.
// I.e. 'wayID,Speed,(int)TrafficLevel'
func (f *Flow) CSVString() string {
	return fmt.Sprintf("%d,%f,%d", f.WayId, f.Speed, int32(f.TrafficLevel))
}

//TODO: func (i *Incident)CSVString() string
