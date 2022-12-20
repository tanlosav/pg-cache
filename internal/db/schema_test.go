package db

import (
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
