package s2indexer

import (
	"os"
	"reflect"
	"testing"

	"github.com/Telenav/osrm-backend/integration/service/oasis/spatialindexer"
	"github.com/golang/geo/s2"
)

func TestDumper(t *testing.T) {
	cases := []S2Indexer{
		S2Indexer{
			cellID2PointIDs: map[s2.CellID][]spatialindexer.PointID{
				9263834064756932608: {1, 2, 3}, // 4/0010133
				9263851656942977024: {1},       // 4/00101332
				9263847258896465920: {2},       // 4/001013321
				9263843960361582592: {3},       // 4/0010133210
			},
			pointID2Location: map[spatialindexer.PointID]spatialindexer.Location{
				1: spatialindexer.Location{
					Lat: 11.11,
					Lon: 11.11,
				},
				2: spatialindexer.Location{
					Lat: 22.22,
					Lon: 22.22,
				},
				3: spatialindexer.Location{
					Lat: 33.33,
					Lon: 33.33,
				},
			},
		},
	}
	// check whether curent folder is writeable
	path, _ := os.Getwd()
	_, err := os.OpenFile(path+"/tmp", os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	err = os.Remove(path + "/tmp")

	removeAllDumpFiles(path)
	for _, c := range cases {
		err := serializeS2Indexer(&c, path)
		if err != nil {
			t.Errorf("During running serializeS2Indexer for case %v, met error %v", c, err)
		}

		actual := NewS2Indexer()
		err = deSerializeS2Indexer(actual, path)
		if err != nil {
			t.Errorf("During running deSerializeS2Indexer for case %v, met error %v", c, err)
		}

		if !reflect.DeepEqual(actual, &c) {
			t.Errorf("Expect result \n%v but got \n%v\n", &c, actual)
		}

		removeAllDumpFiles(path)
	}

}
