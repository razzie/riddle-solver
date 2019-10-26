package main

// Relation describes how items affect each other
type Relation int

// Relation types
const (
	RelAssociated Relation = iota
	RelDisassociated
	RelUnknown
)

// Rule is used by the riddle solver algorithm
type Rule struct {
	ItemA             string
	ItemB             string
	Relation          Relation
	ConditionItemType string
	Condition         string
}
