// Package wrand provides a way to choose random items where each item has it's own
// weight. For example if item "a" has weight 1 and item "b" has weight 2 then "b
// will be selected about twice as often as "a".
package wrand

import (
	"math/rand"
	"sort"
)

// The Item type holds the user's value
type Item interface {
	Weight() int
	WeightIs(int)
	CumWeight() int
	CumWeightIs(int)
}

// The Object type is the collection that holds all the items choos from.
type Object struct {
	pool itemPool
	totalWeight int
	inverse bool
}

// NewObject creates and returns a new Object to work with. If inverse is true then
// smaller weights are more likely.
func NewObject(inverse bool) *Object {
	return &Object{make(itemPool, 0),  0, inverse}
}

// NewItem adds a new Item to the Object with ithe given value and weight. It returns
// a pointer to the newly created Item.
func (o *Object) NewItem(item Item) {
	// O(n)
	o.pool = append(o.pool, item)
	o.update()
}

// updateItemWeight, as the name suggests, sets the given Item's weight to the value
// provided. You should use this instead of setting the Item's weight yourself.
func (o *Object) UpdateItemWeight(item Item, weight int) {
	// O(n)
	item.WeightIs(weight)
	o.update()
}

func (o *Object) update() {
	maxWeight := 0
	for _, item := range o.pool {
		if item.Weight() > maxWeight {
			maxWeight = item.Weight()
		}
	}
	cumWeight := 0
	for _, item := range o.pool {
		w := item.Weight()
		if o.inverse {
			w = maxWeight - w + 1
		}
		cumWeight += w
		item.CumWeightIs(cumWeight)
	}
	o.totalWeight = cumWeight
	sort.Sort(o.pool)
}

// RandomItem returns a printer to a random Item out of the ones that have been added
// via NewItem taking into account the weights of each item.
func (o *Object) RandomItem() Item {
	// O(log n)
	rnd := int(rand.Float64() * float64(o.totalWeight))
	i := sort.Search(o.pool.Len(), func(i int) bool { return o.pool[i].CumWeight() > rnd })
	return o.pool[i]
}

// itemPool is a sortable list of Items
type itemPool []Item

func (p itemPool) Len() int {
	return len(p)
}

func (p itemPool) Less(i, j int) bool {
	return p[i].CumWeight() < p[j].CumWeight()
}

func (p itemPool) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}