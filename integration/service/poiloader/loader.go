package poiloader

import (
	"encoding/json"
	"io/ioutil"

	"github.com/golang/glog"
)

// LoadPOI accepts json file recorded with poi data and returns deserialized result
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
	return elements, nil
}
