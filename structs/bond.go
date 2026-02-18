package structs

type Bond struct {
	BondID         int
	ConnectionType int
	Ends           [2]int
}

func NewBond(bondId, connectionType int, ends [2]int) *Bond {
	return &Bond{
		BondID:         bondId,
		ConnectionType: connectionType,
		Ends:           ends,
	}
}

func (bond *Bond) Equals(other *Bond) bool {
	if other == nil {
		return false
	}
	return bond.BondID == other.BondID &&
		bond.ConnectionType == other.ConnectionType &&
		bond.equalEnds(other)

}

func (bond *Bond) equalEnds(other *Bond) bool {
	return (bond.Ends[0] == other.Ends[0] && bond.Ends[1] == other.Ends[1]) ||
		(bond.Ends[0] == other.Ends[1] && bond.Ends[1] == other.Ends[0])
}
