package s2indexer

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Telenav/osrm-backend/integration/api"
	"github.com/golang/glog"
)

const cellID2PointIDsFileName = "cellID2PointIDs.gob"
const pointID2LocationFileName = "pointID2Location.gob"

func serializeS2Indexer(indexer *S2Indexer, folderPath string) error {
	if !strings.HasSuffix(folderPath, api.Slash) {
		folderPath += "/"
	}

	if err := serializeCellID2PointIDs(indexer, folderPath); err != nil {
		return err
	}

	if err := serializePointID2Location(indexer, folderPath); err != nil {
		return err
	}

	glog.Infof("Successfully serialize S2Indexer to folder %s. len(indexer.cellID2PointIDs) = %d, len(indexer.pointID2Location) = %d\n",
		folderPath, len(indexer.cellID2PointIDs), len(indexer.pointID2Location))

	return nil
}

func serializeCellID2PointIDs(indexer *S2Indexer, folderPath string) error {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(indexer.cellID2PointIDs)

	if err != nil {
		glog.Errorf("During encode S2Indexer's cellID2PointIDs met error %v", err)
		return err
	}

	if err = ioutil.WriteFile(folderPath+cellID2PointIDsFileName, buf.Bytes(), 0644); err != nil {
		glog.Errorf("During dump S2Indexer's cellID2PointIDs met error %v", err)
		return err
	}

	return nil
}

func serializePointID2Location(indexer *S2Indexer, folderPath string) error {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(indexer.pointID2Location)

	if err != nil {
		glog.Errorf("During encode S2Indexer's pointID2Location met error %v", err)
		return err
	}

	if err = ioutil.WriteFile(folderPath+pointID2LocationFileName, buf.Bytes(), 0644); err != nil {
		glog.Errorf("During dump S2Indexer's pointID2Location met error %v", err)
		return err
	}

	return nil
}

func deSerializeS2Indexer(indexer *S2Indexer, folderPath string) error {
	if !strings.HasSuffix(folderPath, api.Slash) {
		folderPath += "/"
	}

	if err := deSerializeCellID2PointIDs(indexer, folderPath); err != nil {
		return err
	}

	if err := deSerializePointID2Location(indexer, folderPath); err != nil {
		return err
	}

	glog.Infof("Successfully deserialize S2Indexer from folder %s. len(indexer.cellID2PointIDs) = %d, len(indexer.pointID2Location) = %d\n",
		folderPath, len(indexer.cellID2PointIDs), len(indexer.pointID2Location))

	return nil
}

func deSerializeCellID2PointIDs(indexer *S2Indexer, folderPath string) error {
	byteArray, err := ioutil.ReadFile(folderPath + cellID2PointIDsFileName)
	if err != nil {
		glog.Errorf("During load S2Indexer's cellID2PointIDs from %s met error %v", folderPath, err)
		return err
	}

	buf := bytes.NewBuffer(byteArray)
	decoder := gob.NewDecoder(buf)
	err = decoder.Decode(&indexer.cellID2PointIDs)

	if err != nil {
		glog.Errorf("During decode S2Indexer's cellID2PointIDs from %s met error %v", folderPath, err)
		return err
	}

	return nil
}

func deSerializePointID2Location(indexer *S2Indexer, folderPath string) error {
	byteArray, err := ioutil.ReadFile(folderPath + pointID2LocationFileName)
	if err != nil {
		glog.Errorf("During load S2Indexer's pointID2Location from %s met error %v", folderPath, err)
		return err
	}

	buf := bytes.NewBuffer(byteArray)
	decoder := gob.NewDecoder(buf)
	err = decoder.Decode(&indexer.pointID2Location)

	if err != nil {
		glog.Errorf("During decode S2Indexer's pointID2Location from %s met error %v", folderPath, err)
		return err
	}

	return nil
}

func removeAllDumpFiles(folderPath string) error {
	if !strings.HasSuffix(folderPath, api.Slash) {
		folderPath += "/"
	}

	_, err := os.Stat(folderPath + cellID2PointIDsFileName)
	if !os.IsNotExist(err) {
		err = os.Remove(folderPath + cellID2PointIDsFileName)
		if err != nil {
			glog.Errorf("Remove file failed %s\n", folderPath+cellID2PointIDsFileName)
			return err
		}
	} else {
		glog.Warningf("There is no %s file in folder %s\n", cellID2PointIDsFileName, folderPath)
	}

	_, err = os.Stat(folderPath + pointID2LocationFileName)
	if !os.IsNotExist(err) {
		err = os.Remove(folderPath + pointID2LocationFileName)
		if err != nil {
			glog.Errorf("Remove file failed %s\n", folderPath+pointID2LocationFileName)
			return err
		}
	} else {
		glog.Warningf("There is no %s file in folder %s\n", pointID2LocationFileName, folderPath)
	}

	return nil
}
