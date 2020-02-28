package stationgraph

import (
	"reflect"
	"testing"
)

// node_0 -> node_1, duration = 30, distance = 30
// node_0 -> node_2, duration = 20, distance = 20
// node_0 -> node_3, duration = 10, distance = 10
func TestQueryHeap1(t *testing.T) {
	m := newQueryHeap()
	m.add(0, invalidNodeID, 0, 0)
	m.add(1, 0, 30, 30)
	m.add(2, 0, 20, 20)
	m.add(3, 0, 10, 10)

	expect := []nodeID{0, 3, 2, 1}
	for _, v := range expect {
		n := m.next()
		if n != v {
			t.Errorf("expect %d but got %d", v, n)
		}
	}

	id := m.next()
	if id != invalidNodeID {
		t.Errorf("expect %d but got %d", invalidNodeID, id)
	}
}

// node_0 -> node_1, duration = 30, distance = 30
// node_0 -> node_2, duration = 20, distance = 20
// node_1 -> node_3, duration = 10, distance = 10
// node_2 -> node_4, duration = 50, distance = 50
// node_2 -> node_3, duration = 50, distance = 50
func TestQueryHeap2(t *testing.T) {
	m := newQueryHeap()
	m.add(0, invalidNodeID, 0, 0)

	m.add(1, 0, 30, 30)
	m.add(2, 0, 20, 20)

	m.add(3, 1, 10, 10)

	m.add(4, 2, 50, 50)
	m.add(3, 2, 50, 50)

	expect := []nodeID{0, 2, 1, 3, 4}
	for _, v := range expect {
		n := m.next()
		if n != v {
			t.Errorf("expect %d but got %d", v, n)
		}
	}

	id := m.next()
	if id != invalidNodeID {
		t.Errorf("expect %d but got %d", invalidNodeID, id)
	}
}

func TestQueryEmptyHeap(t *testing.T) {
	m := newQueryHeap()
	id := m.next()
	if id != invalidNodeID {
		t.Errorf("expect %d but got %d", invalidNodeID, id)
	}
}

// node_0 -> node_1, duration = 30, distance = 30
// node_0 -> node_2, duration = 20, distance = 20
// node_1 -> node_3, duration = 10, distance = 10
// node_2 -> node_4, duration = 50, distance = 50
// node_2 -> node_3, duration = 50, distance = 50
func TestRetrieve1(t *testing.T) {
	m := newQueryHeap()
	m.add(0, invalidNodeID, 0, 0)

	m.add(1, 0, 30, 30)
	m.add(2, 0, 20, 20)

	m.add(3, 1, 10, 10)

	m.add(4, 2, 50, 50)
	m.add(3, 2, 50, 50)

	for {
		if v := m.next(); v == invalidNodeID {
			break
		}
	}

	expect1 := []nodeID{1}
	actual1 := m.retrieve(3)
	if !reflect.DeepEqual(expect1, actual1) {
		t.Errorf("expect %v but got %v", expect1, actual1)
	}

	expect2 := []nodeID{2}
	actual2 := m.retrieve(4)
	if !reflect.DeepEqual(expect2, actual2) {
		t.Errorf("expect %v but got %v", expect2, actual2)
	}

}

// node_0 -> node_1, duration = 30, distance = 30
// node_0 -> node_2, duration = 20, distance = 20
// node_1 -> node_3, duration = 10, distance = 10
// node_2 -> node_4, duration = 50, distance = 50
// node_2 -> node_3, duration = 50, distance = 50
// node_4 -> node_5, duration = 10, distance = 10
// node_5 -> node_6, duration = 10, distance = 10
func TestRetrieve2(t *testing.T) {
	m := newQueryHeap()
	m.add(0, invalidNodeID, 0, 0)

	m.add(1, 0, 30, 30)
	m.add(2, 0, 20, 20)

	m.add(3, 1, 10, 10)

	m.add(4, 2, 50, 50)
	m.add(3, 2, 50, 50)

	m.add(5, 4, 10, 10)

	m.add(6, 5, 10, 10)

	for {
		if v := m.next(); v == invalidNodeID {
			break
		}
	}

	expect1 := []nodeID{1}
	actual1 := m.retrieve(3)
	if !reflect.DeepEqual(expect1, actual1) {
		t.Errorf("expect %v but got %v", expect1, actual1)
	}

	expect2 := []nodeID{2, 4, 5}
	actual2 := m.retrieve(6)
	if !reflect.DeepEqual(expect2, actual2) {
		t.Errorf("expect %v but got %v", expect2, actual2)
	}

}
