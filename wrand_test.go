package wrand

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

const count = 100000

func checkFrequency(o *Object, text string) {
	counts := make(map[string]int)
	for i := 0; i < count; i++ {
		item := o.RandomItem().(*testItem)
		value := item.value
		counts[value] += 1
	}

	fmt.Println("Check", text)
	for v, c := range counts {
		fmt.Printf(" %s: %f\n", v, float64(c)/count)
	}
}

type testItem struct {
	value     string
	weight    int
	cumWeight int
}

func (i *testItem) Weight() int {
	return i.weight
}

func (i *testItem) WeightIs(w int) {
	i.weight = w
}

func (i *testItem) CumWeight() int {
	return i.cumWeight
}

func (i *testItem) CumWeightIs(w int) {
	i.cumWeight = w
}

func nonInverse(t *testing.T) {
	o := NewObject(false)

	a := testItem{"a", 1, 0}
	o.NewItem(&a)
	b := testItem{"b", 1, 0}
	o.NewItem(&b)
	fmt.Println(a, b, o.pool.Len())
	checkFrequency(o, "a=1, b=1")

	fmt.Println(a, b, o.pool.Len())
	o.UpdateItemWeight(&b, 2)
	checkFrequency(o, "a=1, b=2")

	c := testItem{"c", 1, 0}
	o.NewItem(&c)
	fmt.Println(a, b, c, o.pool.Len())
	checkFrequency(o, "a=1, b=2, c=1")

	o.UpdateItemWeight(&c, 4)
	fmt.Println(a, b, c, o.pool.Len())
	checkFrequency(o, "a=1, b=2, c=4")
}

func inverse(t *testing.T) {
	o := NewObject(true)

	a := testItem{"a", 1, 0}
	o.NewItem(&a)
	b := testItem{"b", 1, 0}
	o.NewItem(&b)
	fmt.Println(a, b, o.pool.Len())
	checkFrequency(o, "a=1, b=1")

	fmt.Println(a, b, o.pool.Len())
	o.UpdateItemWeight(&b, 2)
	checkFrequency(o, "a=1, b=2")

	c := testItem{"c", 1, 0}
	o.NewItem(&c)
	fmt.Println(a, b, c, o.pool.Len())
	checkFrequency(o, "a=1, b=2, c=1")

	o.UpdateItemWeight(&c, 4)
	fmt.Println(a, b, c, o.pool.Len())
	checkFrequency(o, "a=1, b=2, c=4")
}

func TestWrand(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	nonInverse(t)
	inverse(t)
}
