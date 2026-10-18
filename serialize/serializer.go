package serialize

import (
	"errors"
	"fmt"
	"strconv"
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
	if err := serializer.serializePairIJCoeffs(); err != nil {
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
	if err := serializer.serializeVelocities(); err != nil {
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
	serializer.writeLine("")
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
	_, err := serializer.writeLine("LAMMPS data file via write_data")
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
	_, err := serializer.writeLinef("%s %s %slo %shi",
		strconv.FormatFloat(lower, 'e', 16, 64), strconv.FormatFloat(higher, 'e', 16, 64),
		string(axis), string(axis))
	return err
}

// ================== Masses ==================

func (serializer *_Serializer) serializeMasses() error {
	if len(serializer.lammpsStruct.BondTypes) == 0 {
		return errors.New("секция Masses не может быть пустой")
	}
	serializer.writeLine("Masses\n")
	for _, atomType := range serializer.lammpsStruct.AtomTypes {
		if _, err := serializer.writeLinef(
			"%d %s # %s",
			atomType.AtomType,
			strconv.FormatFloat(atomType.AtomMass, 'f', -1, 64),
			atomType.AtomLabel); err != nil {
			return err
		}
	}
	return nil
}

// ================== PairIJ coeffs ==================

func (serializer *_Serializer) serializePairIJCoeffs() error {
	if len(serializer.lammpsStruct.PairIJCoeffs) == 0 {
		return nil
	}
	serializer.writeLinef("PairIJ Coeffs # %s\n", serializer.lammpsStruct.PairIJCoeffs[0].GetType())
	for _, pairij := range serializer.lammpsStruct.PairIJCoeffs {
		if _, err := serializer.writeLine(pairij.ToString()); err != nil {
			return err
		}
	}
	return nil
}

// ================== Bond coeffs ==================

func (serializer *_Serializer) serializeBondCoeffs() error {
	if len(serializer.lammpsStruct.BondTypes) == 0 {
		return errors.New("секция Bond Coeffs не может быть пустой")
	}
	serializer.writeLine("Bond Coeffs # harmonic\n")
	for _, bondType := range serializer.lammpsStruct.BondTypes {
		if _, err := serializer.writeLinef("%d %s %s",
			bondType.BondID,
			strconv.FormatFloat(bondType.Sth1, 'f', -1, 64),
			strconv.FormatFloat(bondType.Sth2, 'f', -1, 64),
		); err != nil {
			return err
		}
	}
	return nil
}

func (serializer *_Serializer) serializeAtoms() error {
	if len(serializer.lammpsStruct.Atoms) == 0 {
		return errors.New("секция Atoms не может быть пустой")
	}
	serializer.writeLine("Atoms # full\n")
	for _, atom := range serializer.lammpsStruct.Atoms {
		if _, err := serializer.writeLinef("%d %d %d %s %s %s %s %d %d %d",
			atom.AtomID, atom.MoleculeID, atom.AtomType, strconv.FormatFloat(atom.Q, 'e', 16, 64),
			strconv.FormatFloat(atom.X, 'e', 16, 64), strconv.FormatFloat(atom.Y, 'e', 16, 64), strconv.FormatFloat(atom.Z, 'e', 16, 64),
			atom.NX, atom.NY, atom.NZ,
		); err != nil {
			return err
		}
	}
	return nil
}

func (serializer *_Serializer) serializeBonds() error {
	if len(serializer.lammpsStruct.Bonds) == 0 {
		return nil
	}
	serializer.writeLine("Bonds\n")
	for _, bond := range serializer.lammpsStruct.Bonds {
		if _, err := serializer.writeLinef("%d %d %d %d",
			bond.BondID, bond.ConnectionType, bond.Ends[0], bond.Ends[1],
		); err != nil {
			return err
		}
	}
	return nil
}

// ================== Velocities ==================

func (serializer *_Serializer) serializeVelocities() error {
	if len(serializer.lammpsStruct.Velocities) == 0 {
		return nil
	}
	serializer.writeLine("Velocities\n")
	for _, velocity := range serializer.lammpsStruct.Velocities {
		if _, err := serializer.writeLinef("%d %s %s %s",
			velocity.AtomID,
			strconv.FormatFloat(velocity.VX, 'e', 16, 64),
			strconv.FormatFloat(velocity.VY, 'e', 16, 64),
			strconv.FormatFloat(velocity.VZ, 'e', 16, 64),
		); err != nil {
			return err
		}
	}
	return nil
}
