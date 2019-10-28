package main

import (
	"github.com/razzie/riddle-solver/riddle"
	"github.com/razzie/riddle-solver/ui"
)

// SetupDemo sets up the application to solve Einstein's 5 house riddle
func SetupDemo(root *ui.RootElement) {
	items := func(itemType string, items ...string) {
		root.SetupForm.AddItemType(itemType, items...)
	}

	hint := func(itemA, itemB string) {
		root.RuleList.SaveRule(&riddle.Rule{
			ItemA:    riddle.Item(itemA),
			ItemB:    riddle.Item(itemB),
			Relation: riddle.RelAssociated})
	}

	neighbor := func(itemA, itemB string, cond string) {
		root.RuleList.SaveRule(&riddle.Rule{
			ItemA:             riddle.Item(itemA),
			ItemB:             riddle.Item(itemB),
			Relation:          riddle.RelAssociated,
			ConditionItemType: "house",
			Condition:         cond})
		root.RuleList.SaveRule(&riddle.Rule{
			ItemA:    riddle.Item(itemA),
			ItemB:    riddle.Item(itemB),
			Relation: riddle.RelDisassociated})
	}

	items("house", "1", "2", "3", "4", "5")
	items("nationality", "norwegian", "brit", "swede", "dane", "german")
	items("color", "red", "green", "white", "yellow", "blue")
	items("beverage", "tea", "coffee", "milk", "beer", "water")
	items("cigar", "PallMall", "Dunhill", "blends", "BlueMaster", "Prince")
	items("pet", "dogs", "birds", "cats", "horses", "fish")

	// the Brit lives in the red house
	hint("nationality:brit", "color:red")
	// the Swede keeps dogs as pets
	hint("nationality:swede", "pet:dogs")
	// the Dane drinks tea
	hint("nationality:dane", "beverage:tea")
	// the green house is on the left of the white house
	neighbor("color:green", "color:white", "A == B - 1")
	// the green house's owner drinks coffee
	hint("color:green", "beverage:coffee")
	// the person who smokes Pall Mall rears birds
	hint("cigar:PallMall", "pet:birds")
	// the owner of the yellow house smokes Dunhill
	hint("color:yellow", "cigar:Dunhill")
	// the man living in the center house drinks milk
	hint("house:3", "beverage:milk")
	// the Norwegian lives in the first house
	hint("nationality:norwegian", "house:1")
	// the man who smokes blends lives next to the one who keeps cats
	neighbor("cigar:blends", "pet:cats", "(A == B - 1) || (A == B + 1)")
	// the man who keeps horses lives next to the man who smokes Dunhill
	neighbor("pet:horses", "cigar:Dunhill", "(A == B - 1) || (A == B + 1)")
	// the owner who smokes BlueMaster drinks beer
	hint("cigar:BlueMaster", "beverage:beer")
	// the German smokes Prince
	hint("nationality:german", "cigar:Prince")
	// the Norwegian lives next to the blue house
	neighbor("nationality:norwegian", "color:blue", "(A == B - 1) || (A == B + 1)")
	// the man who smokes blend has a neighbor who drinks water
	neighbor("cigar:blends", "beverage:water", "(A == B - 1) || (A == B + 1)")

	root.SwitchToPage(2)
}
