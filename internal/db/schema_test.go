package db

import (
	"reflect"
	"testing"
)

type getPartitionBordersTest struct {
	now, timeRange, offset, start, end int64
}

var getPartitionBordersTests = []getPartitionBordersTest{
	getPartitionBordersTest{0, 10, 0, 0, 10},
	getPartitionBordersTest{0, 10, 1, 10, 20},
	getPartitionBordersTest{0, 10, 2, 20, 30},
}

func TestGetPartitionBorders(t *testing.T) {
	for _, test := range getPartitionBordersTests {
		start, end := getPartitionBorders(test.now, test.timeRange, test.offset)

		if start != test.start {
			t.Errorf("Expected start value %d but got %d", test.start, start)
		}

		if end != test.end {
			t.Errorf("Expected end value %d but got %d", test.start, end)
		}
	}
}

type getPartitionNameTest struct {
	table         string
	start         int64
	end           int64
	partitionName string
}

var getPartitionNameTests = []getPartitionNameTest{
	getPartitionNameTest{"table", 0, 1, "table_0_1"},
	getPartitionNameTest{"table", 1, 2, "table_1_2"},
}

func TestGetPartitionName(t *testing.T) {
	for _, test := range getPartitionNameTests {
		name := getPartitionName(test.table, partitionBorders{test.start, test.end})

		if test.partitionName != name {
			t.Errorf("Expected partition name '%s' but got '%s'", test.partitionName, name)
		}
	}
}

type getOldPartitionsTest struct {
	all    []string
	actual []string
	old    []string
}

var getOldPartitionsTests = []getOldPartitionsTest{
	getOldPartitionsTest{[]string{}, []string{}, []string{}},
	getOldPartitionsTest{[]string{"p1"}, []string{"p1"}, []string{}},
	getOldPartitionsTest{[]string{"p1"}, []string{}, []string{"p1"}},
	getOldPartitionsTest{[]string{"p1", "p2"}, []string{"p1", "p2"}, []string{}},
	getOldPartitionsTest{[]string{"p1", "p2"}, []string{}, []string{"p1", "p2"}},
	getOldPartitionsTest{[]string{"p1", "p2"}, []string{"p1"}, []string{"p2"}},
	getOldPartitionsTest{[]string{"p1", "p2"}, []string{"p2"}, []string{"p1"}},
	getOldPartitionsTest{[]string{"p1", "p2"}, []string{"p1", "p2"}, []string{}},
	getOldPartitionsTest{[]string{"p1", "p2", "p3"}, []string{"p2"}, []string{"p1", "p3"}},
}

func TestGetOldPartitions(t *testing.T) {
	for _, test := range getOldPartitionsTests {
		old := filterPartitionsGetOld(test.all, test.actual)

		if !reflect.DeepEqual(test.old, old) {
			t.Errorf("Expected old partitions '%+v' but got '%+v'", test.old, old)
		}
	}
}
