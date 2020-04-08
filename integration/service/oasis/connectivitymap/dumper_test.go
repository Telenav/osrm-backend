package connectivitymap

import (
	"os"
	"reflect"
	"testing"

	"github.com/golang/glog"
)

func TestDumpGivenObjectThenLoadAndThenCompareWithOriginalObject(t *testing.T) {
	cases := []ConnectivityMap{
		ConnectivityMap{
			id2nearByIDs: fakeID2NearByIDsMap1,
			maxRange:     fakeDistanceLimit,
			statistic:    &fakeStatisticResult1,
		},
	}

	// check whether curent folder is writeable
	path, _ := os.Getwd()
	_, err := os.OpenFile(path+"/tmp", os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	if err := os.Remove(path + "/tmp"); err != nil {
		glog.Errorf("During dumper test, remove path %s failed.\n", path+"/tmp")
		return
	}

	if err := removeAllDumpFiles(path); err != nil {
		t.Errorf("During running removeAllDumpFiles met error %v", err)
	}
	for _, c := range cases {
		if err := serializeConnectivityMap(&c, path); err != nil {
			t.Errorf("During running serializeConnectivityMap for case %v, met error %v", c, err)
		}

		actual := NewConnectivityMap(0.0)
		if err := deSerializeConnectivityMap(actual, path); err != nil {
			t.Errorf("During running deSerializeConnectivityMap for case %v, met error %v", c, err)
		}

		if !reflect.DeepEqual(actual.id2nearByIDs, c.id2nearByIDs) {
			t.Errorf("Expect result \n%+v but got \n%+v\n", c.id2nearByIDs, actual.id2nearByIDs)
		}

		if !reflect.DeepEqual(actual.maxRange, c.maxRange) {
			t.Errorf("Expect result \n%+v but got \n%+v\n", c.maxRange, actual.maxRange)
		}

		if !reflect.DeepEqual(actual.statistic, c.statistic) {
			t.Errorf("Expect result \n%+v but got \n%+v\n", c.statistic, actual.statistic)
		}

		if err := removeAllDumpFiles(path); err != nil {
			t.Errorf("During running removeAllDumpFiles met error %v", err)
		}
	}
}
