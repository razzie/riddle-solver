package riddle

// Relation describes how items affect each other
type Relation int

// Relation types
const (
	RelAssociated Relation = iota
	RelDisassociated
)

func (rel Relation) String() string {
	switch rel {
	case RelAssociated:
		return "associated"
	case RelDisassociated:
		return "disassociated"
	default:
		return "unknown"
	}
}
