package nodes2wayblotdb

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func testQuery(db *DB) error {
	if db == nil {
		return fmt.Errorf("invalid db")
	}

	queryWayCases := []struct {
		fromNodeID int64
		toNodeID   int64
		wayID      int64
	}{
		{84760891102, 19496208102, 24418325},
		{19496208102, 84760891102, -24418325},
		{84762609102, 244183320001101, 24418332},
		{244183320001101, 84762607102, 24418332},
		{244183320001101, 84762609102, -24418332},
		{84762607102, 244183320001101, -24418332},
	}
	for _, c := range queryWayCases {
		wayID, err := db.QueryWay(c.fromNodeID, c.toNodeID)
		if err != nil {
			return err
		}
		if wayID != c.wayID {
			return fmt.Errorf("query %d,%d in db, expect %d but got %d", c.fromNodeID, c.toNodeID, c.wayID, wayID)
		}
	}

	queryWayExpectFailCases := []struct {
		fromNodeID int64
		toNodeID   int64
	}{
		{0, 0},
		{84760891102, 84760891102},
		{84762609102, 84762607102},
	}
	for _, c := range queryWayExpectFailCases {
		wayID, err := db.QueryWay(c.fromNodeID, c.toNodeID)
		if err == nil {
			return fmt.Errorf("query %d,%d in db, expect fail but got %d", c.fromNodeID, c.toNodeID, wayID)
		}
	}

	queryWaysCases := []struct {
		nodeIDs []int64
		wayIDs  []int64
	}{
		{
			[]int64{84760891102, 19496208102}, []int64{24418325},
		},
		{
			[]int64{19496208102, 84760891102}, []int64{-24418325},
		},
		{
			[]int64{84762609102, 244183320001101, 84762607102}, []int64{24418332, 24418332},
		},
		{
			[]int64{84762607102, 244183320001101, 84762609102}, []int64{-24418332, -24418332},
		},
		{
			[]int64{84762609102, 244183320001101, 84762607199}, []int64{24418332, 24418330},
		},
		{
			[]int64{84762607102, 244183320001101, 84762607199}, []int64{-24418332, 24418330},
		},
		{
			[]int64{84762607199, 244183320001101, 84762607102}, []int64{-24418330, 24418332},
		},
		{
			[]int64{84762607199, 244183320001101, 84762609102}, []int64{-24418330, -24418332},
		},
	}
	for _, c := range queryWaysCases {
		wayIDs, err := db.QueryWays(c.nodeIDs)
		if err != nil {
			return err
		}
		if !reflect.DeepEqual(wayIDs, c.wayIDs) {
			return fmt.Errorf("query %v in db, expect %v but got %v", c.nodeIDs, c.wayIDs, wayIDs)
		}
	}

	queryWaysExpectFailCases := []struct {
		nodeIDs []int64
	}{
		{[]int64{0, 0}},
		{[]int64{0, 0, 0}},
		{[]int64{0, 1, 2}},
		{[]int64{84760891102, 84762609102}},
	}
	for _, c := range queryWaysExpectFailCases {
		wayIDs, err := db.QueryWays(c.nodeIDs)
		if err == nil {
			return fmt.Errorf("query %v in db, expect fail but got %v", c.nodeIDs, wayIDs)
		}
	}
	return nil
}

func TestDB(t *testing.T) {

	tempDBFile := "temp.db"

	// Open a new DB.
	db, err := Open(tempDBFile, false)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(db.db.Path()) // remove temp db file after test.

	// Insert data.
	if err := db.Write(24418325, []int64{84760891102, 19496208102}); err != nil {
		t.Error(err)
	}
	if err := db.Write(24418332, []int64{84762609102, 244183320001101, 84762607102}); err != nil {
		t.Error(err)
	}
	if err := db.Write(24418330, []int64{244183320001101, 84762607199}); err != nil {
		t.Error(err)
	}

	// Query data.
	if err := testQuery(db); err != nil {
		t.Error(err)
	}

	// Close database to release the file lock.
	if err := db.Close(); err != nil {
		t.Error(err)
	}
	db = nil

	// Open exist DB with readonly.
	dbRead, err := Open(tempDBFile, true)
	if err != nil {
		t.Error(err)
	}

	// Query data.
	if err := testQuery(dbRead); err != nil {
		t.Error(err)
	}

	// Close database to release the file lock.
	if err := dbRead.Close(); err != nil {
		t.Error(err)
	}
}
