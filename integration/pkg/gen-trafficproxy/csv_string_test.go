package proxy

import "testing"

func TestFlowCSVString(t *testing.T) {

	cases := []struct {
		f Flow
		s string
	}{
		{Flow{WayId: 829733412, Speed: 20.280001, TrafficLevel: 7}, "829733412,20.280001,7"},
		{Flow{WayId: -129639168, Speed: 31.389999, TrafficLevel: 7}, "-129639168,31.389999,7"},
	}

	for _, c := range cases {
		s := c.f.CSVString()
		if s != c.s {
			t.Errorf("flow: %v, expect csv string %s, but got %s", c.f, c.s, s)
		}
	}

}

func TestIncidentCSVString(t *testing.T) {

	cases := []struct {
		incident Incident
		s        string
	}{
		{
			Incident{
				IncidentId:            "TTI-f47b8dba-59a3-372d-9cec-549eb252e2d5-TTR46312939215361-1",
				AffectedWayIds:        []int64{100663296, -1204020275, 100663296, -1204020274, 100663296, -916744017, 100663296, -1204020245, 100663296, -1194204646, 100663296, -1204394608, 100663296, -1194204647, 100663296, -129639168, 100663296, -1194204645},
				IncidentType:          IncidentType_MISCELLANEOUS,
				IncidentSeverity:      IncidentSeverity_CRITICAL,
				IncidentLocation:      &Location{Lat: 44.181220, Lon: -117.135840},
				Description:           "Construction on I-84 EB near MP 359, Drive with caution.",
				FirstCrossStreet:      "",
				SecondCrossStreet:     "",
				StreetName:            "I-84 E",
				EventCode:             500,
				AlertCEventQuantifier: 0,
				IsBlocking:            false,
			},
			"TTI-f47b8dba-59a3-372d-9cec-549eb252e2d5-TTR46312939215361-1,\"100663296,-1204020275,100663296,-1204020274,100663296,-916744017,100663296,-1204020245,100663296,-1194204646,100663296,-1204394608,100663296,-1194204647,100663296,-129639168,100663296,-1194204645\",MISCELLANEOUS,CRITICAL,44.181220,-117.135840,\"Construction on I-84 EB near MP 359, Drive with caution.\",,,I-84 E,500,0,false",
		},
	}

	for _, c := range cases {
		s := c.incident.CSVString()
		if s != c.s {
			t.Errorf("incident: %v, expect csv string %s, but got %s", c.incident, c.s, s)
		}
	}

}
