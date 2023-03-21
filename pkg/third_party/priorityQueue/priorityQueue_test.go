package priorityQueue

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

type myInt int

func (a myInt) Less(b myInt) bool {
	return a > b
}

func TestNewPriorityQueue(t *testing.T) {
	que := NewPriorityQueue[myInt]()
	var Nums []myInt = []myInt{99, 67, 45, 22, 7, 84, 4, 4, 21, 2, 1}

	for _, v := range Nums {
		que.Push(v)
	}
	var elements []myInt
	for !que.IsEmpty() {
		elements = append(elements, que.Pop())
	}

	sort.Slice(Nums, func(i, j int) bool {
		return Nums[i] > Nums[j]
	})
	if !reflect.DeepEqual(elements, Nums) {
		t.Errorf("Expected %v, got %v", Nums, elements)
	}

}

type myInt2 int

func (a myInt2) Less(b myInt2) bool {
	return a < b
}
func TestNewPriorityQueue2(t *testing.T) {
	que := NewPriorityQueue[myInt2]()
	var Nums []myInt2 = []myInt2{99, 67, 45, 22, 7, 84, 4, 4, 21, 2, 1}

	for _, v := range Nums {
		que.Push(v)
	}
	var elements []myInt2
	for !que.IsEmpty() {
		elements = append(elements, que.Pop())
	}

	sort.Slice(Nums, func(i, j int) bool {
		return Nums[i] < Nums[j]
	})
	fmt.Println(Nums, elements)
	if !reflect.DeepEqual(elements, Nums) {
		t.Errorf("Expected %v, got %v", Nums, elements)
	}

}
