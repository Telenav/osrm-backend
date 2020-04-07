package connectivitymap

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Telenav/osrm-backend/integration/pkg/api"
	"github.com/golang/glog"
)

const id2NearByIDsMapFileName = "id2nearbyidsmap.gob"

func serializeConnectivityMap(cm *ConnectivityMap, folderPath string) error {
	glog.Info("Start serializeConnectivityMap()\n")

	if !strings.HasSuffix(folderPath, api.Slash) {
		folderPath += api.Slash
	}

	if err := serializeID2NearByIDsMap(cm, folderPath); err != nil {
		return err
	}

	if err := cm.statistic.dump(folderPath); err != nil {
		return err
	}

	glog.Infof("Finished serializeConnectivityMap() to folder %s.\n", folderPath)

	return nil
}

func deSerializeConnectivityMap(cm *ConnectivityMap, folderPath string) error {
	glog.Info("Start deSerializeConnectivityMap()")

	if !strings.HasSuffix(folderPath, api.Slash) {
		folderPath += api.Slash
	}

	if err := deSerializeID2NearByIDsMap(cm, folderPath); err != nil {
		return err
	}

	if err := cm.statistic.load(folderPath); err != nil {
		return err
	}
	cm.maxRange = cm.statistic.MaxRange

	glog.Infof("Finished deSerializeConnectivityMap() to folder %s.\n", folderPath)

	return nil
}

func removeAllDumpFiles(folderPath string) error {
	if !strings.HasSuffix(folderPath, api.Slash) {
		folderPath += "/"
	}

	_, err := os.Stat(folderPath + id2NearByIDsMapFileName)
	if !os.IsNotExist(err) {
		err = os.Remove(folderPath + id2NearByIDsMapFileName)
		if err != nil {
			glog.Errorf("Remove file failed %s\n", folderPath+id2NearByIDsMapFileName)
			return err
		}
	} else {
		glog.Warningf("There is no %s file in folder %s\n", id2NearByIDsMapFileName, folderPath)
	}

	_, err = os.Stat(folderPath + statisticFileName)
	if !os.IsNotExist(err) {
		err = os.Remove(folderPath + statisticFileName)
		if err != nil {
			glog.Errorf("Remove file failed %s\n", folderPath+statisticFileName)
			return err
		}
	} else {
		glog.Warningf("There is no %s file in folder %s\n", statisticFileName, folderPath)
	}

	return nil
}

func serializeID2NearByIDsMap(cm *ConnectivityMap, folderPath string) error {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(cm.id2nearByIDs)

	if err != nil {
		glog.Errorf("During encode ConnectivityMap's ID2NearByIDsMap met error %v", err)
		return err
	}

	if err = ioutil.WriteFile(folderPath+id2NearByIDsMapFileName, buf.Bytes(), 0644); err != nil {
		glog.Errorf("During dump ConnectivityMap's ID2NearByIDsMap to %s met error %v", folderPath+id2NearByIDsMapFileName, err)
		return err
	}

	return nil
}

func deSerializeID2NearByIDsMap(cm *ConnectivityMap, folderPath string) error {
	byteArray, err := ioutil.ReadFile(folderPath + id2NearByIDsMapFileName)
	if err != nil {
		glog.Errorf("During load ConnectivityMap's ID2NearByIDsMap from %s met error %v",
			folderPath+id2NearByIDsMapFileName, err)
		return err
	}

	buf := bytes.NewBuffer(byteArray)
	decoder := gob.NewDecoder(buf)
	err = decoder.Decode(&cm.id2nearByIDs)

	if err != nil {
		glog.Errorf("During decode ConnectivityMap's ID2NearByIDsMap from %s met error %v",
			folderPath+id2NearByIDsMapFileName, err)
		return err
	}

	return nil
}
