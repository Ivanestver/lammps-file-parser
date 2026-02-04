package structs

type AtomCoords struct{ X, Y, Z float64 }

func (crds *AtomCoords) equals(others *AtomCoords) bool {
	return crds.X == others.X &&
		crds.Y == others.Y &&
		crds.Z == others.Z
}

type Atom struct {
	Label      string
	AtomID     int
	MoleculeID int
	Q          float64
	AtomCoords
}

func NewAtom(label string, atomId, moleculeId int, q, x, y, z float64) *Atom {
	return &Atom{
		Label:      label,
		AtomID:     atomId,
		MoleculeID: moleculeId,
		Q:          q,
		AtomCoords: AtomCoords{
			X: x,
			Y: y,
			Z: z,
		},
	}
}

func (atom *Atom) Equals(other *Atom) bool {
	if other == nil {
		return false
	}

	return atom.Label == other.Label &&
		atom.AtomID == other.AtomID &&
		atom.MoleculeID == other.MoleculeID &&
		atom.AtomCoords.equals(&other.AtomCoords)
}
