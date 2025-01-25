package timeline

import (
	"slices"

	"github.com/janmalch/argus/internal/handler"
)

// TODO: rethink models
type Timeline struct {
	order        []uint64
	data         map[uint64]*handler.Exchange
	reqBodySizes map[uint64]int64
	resBodySizes map[uint64]int64
}

func NewTimeline() *Timeline {
	return &Timeline{
		order:        make([]uint64, 0),
		data:         make(map[uint64]*handler.Exchange),
		reqBodySizes: make(map[uint64]int64),
		resBodySizes: make(map[uint64]int64),
	}
}

func (t *Timeline) Add(r *handler.Exchange) {
	if _, isKnown := t.data[r.Id]; !isKnown {
		t.data[r.Id] = r
		i, _ := slices.BinarySearch(t.order, r.Id)
		t.order = slices.Insert(t.order, i, r.Id)
	}
}

func (t *Timeline) SetReqBodySize(id uint64, size int64) {
	if size > 0 {
		t.reqBodySizes[id] = size
	}
}

func (t *Timeline) SetResBodySize(id uint64, size int64) {
	if size > 0 {
		t.resBodySizes[id] = size
	}
}

// Returns -1 if error occurred while determining size
func (t *Timeline) GetReqBodySize(id uint64) int64 {
	size, found := t.reqBodySizes[id]
	if found {
		return size
	} else {
		return 0
	}
}

// Returns -1 if error occurred while determining size
func (t *Timeline) GetResBodySize(id uint64) int64 {
	size, found := t.resBodySizes[id]
	if found {
		return size
	} else {
		return 0
	}
}

func (t *Timeline) Len() int {
	return len(t.order)
}

func (t *Timeline) At(i int) *handler.Exchange {
	if i >= len(t.order) {
		return nil
	}
	id := t.order[i]
	return t.data[id]
}

func (t *Timeline) Clear() {
	t.order = make([]uint64, 0)
	t.data = make(map[uint64]*handler.Exchange)
	t.reqBodySizes = make(map[uint64]int64)
	t.resBodySizes = make(map[uint64]int64)
}

func (t *Timeline) Data() []*handler.Exchange {
	d := make([]*handler.Exchange, len(t.order))
	for i, id := range t.order {
		d[i] = t.data[id]
	}
	return d
}
