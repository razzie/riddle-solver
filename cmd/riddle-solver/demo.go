package main

import (
	"github.com/razzie/riddle-solver/riddle"
)

// Demo represents Einstein's 5 house riddle
type Demo struct {
	*riddle.Riddle
}

// NewDemo returns a new Demo
func NewDemo() *Demo {
	demo := &Demo{Riddle: riddle.NewRiddle()}

	demo.items("house", "1", "2", "3", "4", "5")
	demo.items("nationality", "norwegian", "brit", "swede", "dane", "german")
	demo.items("color", "red", "green", "white", "yellow", "blue")
	demo.items("beverage", "tea", "coffee", "milk", "beer", "water")
	demo.items("cigar", "PallMall", "Dunhill", "blends", "BlueMaster", "Prince")
	demo.items("pet", "dogs", "birds", "cats", "horses", "fish")

	// the Brit lives in the red house
	demo.hint("nationality:brit", "color:red")
	// the Swede keeps dogs as pets
	demo.hint("nationality:swede", "pet:dogs")
	// the Dane drinks tea
	demo.hint("nationality:dane", "beverage:tea")
	// the green house is on the left of the white house
	demo.leftneighbor("color:green", "color:white")
	// the green house's owner drinks coffee
	demo.hint("color:green", "beverage:coffee")
	// the person who smokes Pall Mall rears birds
	demo.hint("cigar:PallMall", "pet:birds")
	// the owner of the yellow house smokes Dunhill
	demo.hint("color:yellow", "cigar:Dunhill")
	// the man living in the center house drinks milk
	demo.hint("house:3", "beverage:milk")
	// the Norwegian lives in the first house
	demo.hint("nationality:norwegian", "house:1")
	// the man who smokes blends lives next to the one who keeps cats
	demo.anyneighbor("cigar:blends", "pet:cats")
	// the man who keeps horses lives next to the man who smokes Dunhill
	demo.anyneighbor("pet:horses", "cigar:Dunhill")
	// the owner who smokes BlueMaster drinks beer
	demo.hint("cigar:BlueMaster", "beverage:beer")
	// the German smokes Prince
	demo.hint("nationality:german", "cigar:Prince")
	// the Norwegian lives next to the blue house
	demo.anyneighbor("nationality:norwegian", "color:blue")
	// the man who smokes blend has a neighbor who drinks water
	demo.anyneighbor("cigar:blends", "beverage:water")

	return demo
}

func (demo *Demo) items(itemType string, items ...string) {
	demo.Setup[itemType] = items
}

func (demo *Demo) hint(itemA, itemB string) {
	demo.Rules = append(demo.Rules, riddle.Rule{
		ItemA:    riddle.Item(itemA),
		ItemB:    riddle.Item(itemB),
		Relation: riddle.RelAssociated,
	})
}

func (demo *Demo) leftneighbor(itemA, itemB string) {
	demo.Rules = append(demo.Rules, riddle.Rule{
		ItemA:             riddle.Item(itemA),
		ItemB:             riddle.Item(itemB),
		Relation:          riddle.RelAssociated,
		ConditionItemType: "house",
		Condition:         "A == B - 1",
		IsReversible:      false,
	})
	demo.Rules = append(demo.Rules, riddle.Rule{
		ItemA:             riddle.Item(itemB),
		ItemB:             riddle.Item(itemA),
		Relation:          riddle.RelAssociated,
		ConditionItemType: "house",
		Condition:         "A == B + 1",
		IsReversible:      false,
	})
}

func (demo *Demo) anyneighbor(itemA, itemB string) {
	demo.Rules = append(demo.Rules, riddle.Rule{
		ItemA:             riddle.Item(itemA),
		ItemB:             riddle.Item(itemB),
		Relation:          riddle.RelAssociated,
		ConditionItemType: "house",
		Condition:         "(A == B - 1) || (A == B + 1)",
		IsReversible:      true,
	})
}
