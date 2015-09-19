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

// The Object type is the collection that holds all the items choose from.
type Object struct {
	pool        itemPool
	totalWeight int
	inverse     bool
}

// NewObject creates and returns a new Object to work with. If inverse is true then
// smaller weights are more likely.
func NewObject(inverse bool) *Object {
	return &Object{make(itemPool, 0), 0, inverse}
}

// NewItem adds a new Item to the Object with ithe given value and weight.
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

// RandomItem returns a random Item out of the ones that have been added via NewItem
// taking into account the weights of each item.
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

type Selectable interface {
	Weight() float64
}

// Select returns a random item from the given list with a probability corresponding to
// the relative weights of each item. Each item's Wegiht function will be called once.
func Select(items []Selectable) Selectable {
	cumWeights := make([]float64, len(items))
	cumWeights[0] = items[0].Weight()
	for i, item := range items {
		if i > 0 {
			cumWeights[i] = cumWeights[i - 1] + item.Weight()
		}
	}

	// rand.Float64() is strictly less than 1.0 so rnd will never be equal to the total weight.
	// Therefore SearchFloat64s will always return a valid index.
	rnd := rand.Float64() * cumWeights[len(items)-1]
	return items[sort.SearchFloat64s(cumWeights, rnd)]
}

// SelectIndex takes a list of weights and returns an index with a probability corresponding
// to the relative weight of each index. Behavior is undefined if len(weights) == 0. A weight
// of 0 will never be selected unless all are 0, in which case any index may be selected.
func SelectIndex(weights []float64) int {
	cumWeights := make([]float64, len(weights))
	cumWeights[0] = weights[0]
	for i, w := range weights {
		if i > 0 {
			cumWeights[i] = cumWeights[i - 1] + w
		}
	}

	if cumWeights[len(weights)-1] == 0.0 {
		return rand.Intn(len(weights))
	}

	rnd := rand.Float64() * cumWeights[len(weights)-1]
	return sort.SearchFloat64s(cumWeights, rnd)
}
