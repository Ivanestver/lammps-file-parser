package serialize

import (
	"fmt"
	"strings"

	"github.com/Ivanestver/lammps-file-parser/structs"
)

type _Serializer struct {
	lammpsStruct *structs.LammpsStruct
	builder      strings.Builder
}

func (serializer *_Serializer) writeString(line string) (int, error) {
	return serializer.builder.WriteString(line)
}

func (serializer *_Serializer) writeStringf(format string, a ...any) (int, error) {
	return serializer.builder.WriteString(fmt.Sprintf(format, a...))
}

func (serializer *_Serializer) writeLine(line string) (int, error) {
	return serializer.builder.WriteString(fmt.Sprintf("%s\n", line))
}

func (serializer *_Serializer) writeLinef(format string, a ...any) (int, error) {
	return serializer.builder.WriteString(fmt.Sprintf(format+"\n", a...))
}

func _NewSerializer(lammpsStruct *structs.LammpsStruct) *_Serializer {
	return &_Serializer{
		lammpsStruct: lammpsStruct,
		builder:      strings.Builder{},
	}
}

func (serializer *_Serializer) Serialize() (string, error) {
	if err := serializer.serializeMetadata(); err != nil {
		return "", err
	}
	serializer.writeLine("")
	if err := serializer.serializeMasses(); err != nil {
		return "", err
	}
	serializer.writeLine("")
	if err := serializer.serializeBondCoeffs(); err != nil {
		return "", err
	}
	serializer.writeLine("")
	if err := serializer.serializeAtoms(); err != nil {
		return "", err
	}
	serializer.writeLine("")
	if err := serializer.serializeBonds(); err != nil {
		return "", err
	}
	return serializer.builder.String(), nil
}

func (serializer *_Serializer) serializeMetadata() error {
	if err := serializer.serializeHeader(); err != nil {
		return err
	}
	serializer.writeString("")
	if err := serializer.serializeAtomsCount(); err != nil {
		return err
	}
	if err := serializer.serializeAtomsTypesCount(); err != nil {
		return err
	}
	if err := serializer.serializeBondsCount(); err != nil {
		return err
	}
	if err := serializer.serializeBondsTypesCount(); err != nil {
		return err
	}
	serializer.writeLine("")
	if err := serializer.serializeSpaceMeasures(); err != nil {
		return err
	}
	return nil
}

// ================== Metadata ==================
func (serializer *_Serializer) serializeHeader() error {
	_, err := serializer.writeLine("LAMMPS data file via write_data, version 24 Dec 2020, timestep = 40000000")
	return err
}

func (serializer *_Serializer) serializeAtomsCount() error {
	_, err := serializer.writeLinef("%d atoms", len(serializer.lammpsStruct.Atoms))
	return err
}

func (serializer *_Serializer) serializeAtomsTypesCount() error {
	_, err := serializer.writeLinef("%d atom types", len(serializer.lammpsStruct.AtomTypes))
	return err
}

func (serializer *_Serializer) serializeBondsCount() error {
	_, err := serializer.writeLinef("%d bonds", len(serializer.lammpsStruct.Bonds))
	return err
}

func (serializer *_Serializer) serializeBondsTypesCount() error {
	_, err := serializer.writeLinef("%d bond types", len(serializer.lammpsStruct.BondTypes))
	return err
}

func (serializer *_Serializer) serializeSpaceMeasures() error {
	axes := [3]rune{'x', 'y', 'z'}
	for i, axis := range axes {
		if err := serializer.serializeSpaceMeasure(
			serializer.lammpsStruct.SpaceDimention[i][0],
			serializer.lammpsStruct.SpaceDimention[i][1],
			axis); err != nil {
			return err
		}
	}
	return nil
}

func (serializer *_Serializer) serializeSpaceMeasure(lower, higher float64, axis rune) error {
	_, err := serializer.writeLinef("%d %d %slo %shi", lower, higher, axis, axis)
	return err
}

// ================== Masses ==================

func (serializer *_Serializer) serializeMasses() error {
	for _, atomType := range serializer.lammpsStruct.AtomTypes {
		if _, err := serializer.writeLinef("%d %f", atomType.Item1, atomType.Item2); err != nil {
			return err
		}
	}
	return nil
}

// ================== Bond coeffs ==================

func (serializer *_Serializer) serializeBondCoeffs() error {
	for _, bondType := range serializer.lammpsStruct.BondTypes {
		if _, err := serializer.writeLinef("%d %f %f", bondType.Item1, bondType.Item2, bondType.Item3); err != nil {
			return err
		}
	}
	return nil
}

func (serializer *_Serializer) serializeAtoms() error {
	for _, atom := range serializer.lammpsStruct.Atoms {
		if _, err := serializer.writeLinef("%d %d %d %f %f %f %f 0 0 0",
			atom.AtomID, atom.MoleculeID, atom.AtomType, atom.Q,
			atom.X, atom.Y, atom.Z,
		); err != nil {
			return err
		}
	}
	return nil
}

func (serializer *_Serializer) serializeBonds() error {
	for _, bond := range serializer.lammpsStruct.Bonds {
		if _, err := serializer.writeLinef("%d, %d, %d, %d",
			bond.BondID, bond.ConnectionType, bond.Ends[0], bond.Ends[1],
		); err != nil {
			return err
		}
	}
	return nil
}
