package poiloader

import (
	"encoding/json"
	"io/ioutil"

	"github.com/golang/glog"
)

// LoadData accepts json file with points data and returns deserialized result
func LoadData(filePath string) ([]Element, error) {
	var elements []Element

	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		glog.Errorf("While load file %s, met error %v\n", filePath, err)
		return elements, err
	}

	err = json.Unmarshal(file, &elements)
	if err != nil {
		glog.Errorf("While unmarshal json file %s, met error %v\n", filePath, err)
		return elements, err
	}

	glog.Infof("Finished loading %d items from json file %s", len(elements), filePath)
	return elements, nil
}
