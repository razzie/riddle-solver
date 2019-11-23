package riddle

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Demo riddles
var (
	EinsteinRiddle *Riddle
	JindoshRiddle  *Riddle
)

// Riddle contains the setup and rules of the riddle
type Riddle struct {
	Setup           Setup
	Rules           []Rule
	PrimaryItemType string `json:",omitempty"`
}

// NewRiddle returns a new Riddle
func NewRiddle() *Riddle {
	return &Riddle{Setup: make(Setup)}
}

// Check returns an error if the riddle is invalid
func (r *Riddle) Check() error {
	if err := r.Setup.Check(); err != nil {
		return err
	}

	for i, rule := range r.Rules {
		if err := rule.Check(r.Setup); err != nil {
			return fmt.Errorf("rule#%d error: %v", i+1, err)
		}
	}

	if len(r.PrimaryItemType) > 0 {
		if _, ok := r.Setup[r.PrimaryItemType]; !ok {
			return fmt.Errorf("Primary item type %q not found", r.PrimaryItemType)
		}
	}

	return nil
}

// Solve solves the riddle and returns the entries
func (r *Riddle) Solve() ([]SolverEntry, bool, error) {
	solver := NewSolver(r.Setup, r.Rules)
	primaryItemType := r.PrimaryItemType
	if len(primaryItemType) == 0 {
		primaryItemType = solver.GuessPrimaryItemType()
	}

	_, err := solver.Solve(primaryItemType)
	return solver.Entries, solver.IsSolved(), err
}

// LoadRiddle loads the riddle from a byte slice in JSON format
func LoadRiddle(data []byte) (*Riddle, error) {
	var r Riddle

	err := json.Unmarshal(data, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

// LoadRiddleFromFile loads the riddle from a JSON file
func LoadRiddleFromFile(file string) (*Riddle, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return LoadRiddle(data)
}

// Save saves the riddle to a byte slice in JSON format
func (r *Riddle) Save() ([]byte, error) {
	return json.MarshalIndent(r, "", "\t")
}

// SaveToFile saves the riddle to a JSON file
func (r *Riddle) SaveToFile(file string) error {
	data, err := r.Save()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, data, 0777)
}

func (r *Riddle) items(itemType string, items ...string) {
	r.Setup[itemType] = items
}

func (r *Riddle) hint(itemA, itemB string) {
	r.Rules = append(r.Rules, Rule{
		ItemA:    Item(itemA),
		ItemB:    Item(itemB),
		Relation: RelAssociated,
	})
}

func (r *Riddle) leftneighbor(itemA, itemB string) {
	r.Rules = append(r.Rules, Rule{
		ItemA:             Item(itemA),
		ItemB:             Item(itemB),
		Relation:          RelAssociated,
		ConditionItemType: "house",
		Condition:         "A == B - 1",
		IsReversible:      false,
	})
	r.Rules = append(r.Rules, Rule{
		ItemA:             Item(itemB),
		ItemB:             Item(itemA),
		Relation:          RelAssociated,
		ConditionItemType: "house",
		Condition:         "A == B + 1",
		IsReversible:      false,
	})
}

func (r *Riddle) anyneighbor(itemA, itemB string) {
	r.Rules = append(r.Rules, Rule{
		ItemA:             Item(itemA),
		ItemB:             Item(itemB),
		Relation:          RelAssociated,
		ConditionItemType: "house",
		Condition:         "(A == B - 1) || (A == B + 1)",
		IsReversible:      true,
	})
}

func newEinsteinRiddle() *Riddle {
	r := NewRiddle()

	r.PrimaryItemType = "house"
	r.items("house", "1", "2", "3", "4", "5")
	r.items("nationality", "norwegian", "brit", "swede", "dane", "german")
	r.items("color", "red", "green", "white", "yellow", "blue")
	r.items("beverage", "tea", "coffee", "milk", "beer", "water")
	r.items("cigar", "PallMall", "Dunhill", "blends", "BlueMaster", "Prince")
	r.items("pet", "dogs", "birds", "cats", "horses", "fish")

	// the Brit lives in the red house
	r.hint("nationality:brit", "color:red")
	// the Swede keeps dogs as pets
	r.hint("nationality:swede", "pet:dogs")
	// the Dane drinks tea
	r.hint("nationality:dane", "beverage:tea")
	// the green house is on the left of the white house
	r.leftneighbor("color:green", "color:white")
	// the green house's owner drinks coffee
	r.hint("color:green", "beverage:coffee")
	// the person who smokes Pall Mall rears birds
	r.hint("cigar:PallMall", "pet:birds")
	// the owner of the yellow house smokes Dunhill
	r.hint("color:yellow", "cigar:Dunhill")
	// the man living in the center house drinks milk
	r.hint("house:3", "beverage:milk")
	// the Norwegian lives in the first house
	r.hint("nationality:norwegian", "house:1")
	// the man who smokes blends lives next to the one who keeps cats
	r.anyneighbor("cigar:blends", "pet:cats")
	// the man who keeps horses lives next to the man who smokes Dunhill
	r.anyneighbor("pet:horses", "cigar:Dunhill")
	// the owner who smokes BlueMaster drinks beer
	r.hint("cigar:BlueMaster", "beverage:beer")
	// the German smokes Prince
	r.hint("nationality:german", "cigar:Prince")
	// the Norwegian lives next to the blue house
	r.anyneighbor("nationality:norwegian", "color:blue")
	// the man who smokes blend has a neighbor who drinks water
	r.anyneighbor("cigar:blends", "beverage:water")

	return r
}

func newJindoshRiddle() *Riddle {
	r := NewRiddle()

	r.PrimaryItemType = "seat"
	r.items("seat", "1", "2", "3", "4", "5")
	r.items("name", "Winslow", "Marcolla", "Natsiou", "Finch")
	r.items("color", "purple", "blue", "red", "green", "white")
	r.items("drink", "whiskey", "beer", "absinthe", "rum", "wine")
	r.items("jewel", "SnuffTin", "WarMedal", "Ring", "BirdPendant", "Diamond")
	r.items("place", "Dabokva", "Fraeport", "Dunwall", "Karnaca", "Baleton")

	// Lady Winslow wore a jaunty purple hat.
	r.hint("name:Winslow", "color:purple")
	// Doctor Marcolla was at the far left, next to the guest wearing a blue jacket
	r.hint("name:Marcolla", "seat:1")
	r.hint("color:blue", "seat:2")
	// The lady in the red sat left of someone in green.
	r.leftneighbor("color:red", "color:green")
	// I remember that red outfit because the woman spilled her whiskey all over it.
	r.hint("color:red", "drink:whiskey")
	// The traveler from Dabokva was dressed entirely in white.
	r.hint("place:Dabokva", "color:white")
	// When one of the dinner guests bragged about her Snuff Tin, the woman next to her said
	// they were finer in Dabokva, where she lived
	r.anyneighbor("jewel:SnuffTin", "place:Dabokva")
	// So Countess Contee showed off a prized War Medal, at which the lady from Fraeport scoffed,
	// saying it was no match for her Ring.
	r.hint("name:Contee", "jewel:WarMedal")
	r.hint("place:Fraeport", "jewel:Ring")
	// Someone else carried a valuable Bird Pendant and when she saw it,
	// the visitor from Dunwall next to her almost spilled her neighbor's absinthe.
	r.anyneighbor("jewel:BirdPendant", "place:Dunwall")
	r.hint("jewel:BirdPendant", "drink:absinthe")
	// Baroness Finch raised her rum in toast.
	r.hint("name:Finch", "drink:rum")
	// The lady from Karnaca, full of wine, jumped up onto the table, falling onto
	// the guest in the center seat, spilling the poor woman's beer.
	r.hint("place:Karnaca", "drink:wine")
	r.hint("seat:3", "drink:beer")
	// Then Madam Natsiou captivated them all with a story about her wild youth in Baleton.
	r.hint("name:Natsiou", "place:Baleton")

	return r
}

func init() {
	EinsteinRiddle = newEinsteinRiddle()
	JindoshRiddle = newJindoshRiddle()
}
