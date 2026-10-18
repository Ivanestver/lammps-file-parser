package structs

type AtomCoords struct{ X, Y, Z float64 }
type PeriodicImageFlags struct{ NX, NY, NZ int }

func (crds *AtomCoords) equals(others *AtomCoords) bool {
	return crds.X == others.X &&
		crds.Y == others.Y &&
		crds.Z == others.Z
}

type Atom struct {
	Label      string
	AtomID     int
	MoleculeID int
	AtomType   int
	Q          float64
	AtomCoords
	PeriodicImageFlags
}

func NewAtom(label string, atomID, moleculeID, atomType int, q, x, y, z float64, nx, ny, nz int) *Atom {
	return &Atom{
		Label:      label,
		AtomID:     atomID,
		MoleculeID: moleculeID,
		AtomType:   atomType,
		Q:          q,
		AtomCoords: AtomCoords{
			X: x,
			Y: y,
			Z: z,
		},
		PeriodicImageFlags: PeriodicImageFlags{
			NX: nx,
			NY: ny,
			NZ: nz,
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
