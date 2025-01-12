package ui

import (
	"slices"

	"github.com/janmalch/argus/internal/handler"
)

// TODO: move to timeline package and rethink models
type timeline struct {
	order        []uint64
	data         map[uint64]*handler.Exchange
	reqBodySizes map[uint64]int64
	resBodySizes map[uint64]int64
}

func newTimeline() *timeline {
	return &timeline{
		order:        make([]uint64, 0),
		data:         make(map[uint64]*handler.Exchange),
		reqBodySizes: make(map[uint64]int64),
		resBodySizes: make(map[uint64]int64),
	}
}

func (t *timeline) add(r *handler.Exchange) {
	if _, isKnown := t.data[r.Id]; !isKnown {
		t.data[r.Id] = r
		i, _ := slices.BinarySearch(t.order, r.Id)
		t.order = slices.Insert(t.order, i, r.Id)
	}
}

func (t *timeline) setReqBodySize(id uint64, size int64) {
	if size > 0 {
		t.reqBodySizes[id] = size
	}
}

func (t *timeline) setResBodySize(id uint64, size int64) {
	if size > 0 {
		t.resBodySizes[id] = size
	}
}

// Returns -1 if error occurred while determining size
func (t *timeline) getReqBodySize(id uint64) int64 {
	size, found := t.reqBodySizes[id]
	if found {
		return size
	} else {
		return 0
	}
}

// Returns -1 if error occurred while determining size
func (t *timeline) getResBodySize(id uint64) int64 {
	size, found := t.resBodySizes[id]
	if found {
		return size
	} else {
		return 0
	}
}

func (t *timeline) len() int {
	return len(t.order)
}

func (t *timeline) at(i int) *handler.Exchange {
	if i >= len(t.order) {
		return nil
	}
	id := t.order[i]
	return t.data[id]
}
