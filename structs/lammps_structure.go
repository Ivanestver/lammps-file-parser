package structs

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type LammpsStruct struct {
	FileName string
	Atoms    []Atom
	Bonds    []Bond
}

type _LammpsMetadata struct {
	atomsCount     int
	atomTypesCount int
	bondsCount     int
	bondTypesCount int
	atomTypes      map[string]struct {
		Mass  int
		Label string
	}
	bondTypes map[string]struct {
		sth1 int
		sth2 float64
	}
	atoms map[int64]struct {
		Atom        *Atom
		ChainNumber int
	}
	bonds []*Bond
}

type LammpsLoader struct {
	_LammpsMetadata
	builtGlobula LammpsStruct
	scanner      *bufio.Scanner
}

func (loader *LammpsLoader) Load(content string) (LammpsStruct, error) {
	loader.scanner = bufio.NewScanner(strings.NewReader(content))
	if err := loader.load(); err != nil {
		return LammpsStruct{}, err
	}
	loader.scanner = nil
	return loader.builtGlobula, nil
}

func (loader *LammpsLoader) load() error {
	if err := loader.loadMetadata(); err != nil {
		return err
	}
	if err := loader.loadMasses(); err != nil {
		return err
	}
	if err := loader.loadBondTypes(); err != nil {
		return err
	}
	if err := loader.loadAtoms(); err != nil {
		return err
	}
	if err := loader.loadBonds(); err != nil {
		return err
	}
	if err := loader.constructLammpsStruct(); err != nil {
		return err
	}
	return nil
}

func (loader *LammpsLoader) loadMetadata() error {
	for loader.scanner.Scan() {
		line := loader.scanner.Text()
		if len(line) == 0 {
			continue
		}
		if line[0] < '0' || line[0] > '9' {
			continue
		}
		// read atoms count section
		if count, err := getNumber(line); err == nil {
			loader.atomsCount = count
			break
		}
	}
	loader.atoms = make(map[int64]struct {
		Atom        *Atom
		ChainNumber int
	})

	// read atoms types section
	if err := writeMetadata(loader, &loader.atomTypesCount, "atom types"); err != nil {
		return err
	}
	loader.atomTypes = make(map[string]struct {
		Mass  int
		Label string
	})

	// read bonds section
	if err := writeMetadata(loader, &loader.bondsCount, "bonds counts"); err != nil {
		return err
	}
	loader.bondTypes = make(map[string]struct {
		sth1 int
		sth2 float64
	})

	// read bonds types section
	if err := writeMetadata(loader, &loader.bondTypesCount, "bonds types"); err != nil {
		return err
	}

	return nil
}

func getNumber(s string) (int, error) {
	if len(s) == 0 {
		return 0, errors.New("string is empty")
	}

	isNumber := func(b byte) bool { return '0' <= b && b <= '9' }

	var builder strings.Builder
	currentSymbol := 0
	for isNumber(s[currentSymbol]) && currentSymbol < len(s) {
		builder.WriteByte(s[currentSymbol])
		currentSymbol++
	}
	return strconv.Atoi(builder.String())
}

func writeMetadata(loader *LammpsLoader, metadata *int, sectionName string) error {
	if loader.scanner.Scan() {
		line := loader.scanner.Text()
		if len(line) == 0 {
			return errors.New("Could not find " + sectionName + " section")
		}
		if line[0] < '0' || line[0] > '9' {
			return errors.New("Could not find " + sectionName + " section")
		}
		if count, err := getNumber(line); err == nil {
			*metadata = count
		}
	}
	return nil
}

func (loader *LammpsLoader) loadMasses() error {
	for loader.scanner.Scan() {
		if loader.scanner.Text() == "Masses" {
			break
		}
	}
	loader.scanner.Scan()

	for atomTypeLineNumber := 0; atomTypeLineNumber < loader.atomTypesCount && loader.scanner.Scan(); atomTypeLineNumber++ {
		line := loader.scanner.Text()
		if len(line) == 0 {
			break
		}
		parts := strings.Split(line, " ")
		if len(parts) < 2 || len(parts) > 4 {
			return fmt.Errorf("wrong line in the Masses section (line number in there: %d)", atomTypeLineNumber+1)
		}
		number := parts[0]
		mass, err := strconv.Atoi(parts[1])
		if err != nil {
			return err
		}
		literals := make(map[int]string)
		literals[0] = "O"
		literals[1] = "N"
		literals[2] = "C"
		literals[3] = "S"
		var label string
		if len(parts) > 2 {
			label = parts[3]
		} else {
			for t, l := range literals {
				n, _ := strconv.Atoi(number)
				if int(t) == n {
					label = l
					break
				}
			}
			if len(label) == 0 {
				label = "C"
			}
		}
		loader.atomTypes[number] = struct {
			Mass  int
			Label string
		}{
			Mass:  mass,
			Label: label,
		}
	}

	return nil
}

func (loader *LammpsLoader) loadBondTypes() error {
	for loader.scanner.Scan() && loader.scanner.Text() != "Bond Coeffs # harmonic" {
	}
	loader.scanner.Scan()

	for bondTypeLineNumber := 0; bondTypeLineNumber < loader.bondTypesCount && loader.scanner.Scan(); bondTypeLineNumber++ {
		parts := strings.Split(loader.scanner.Text(), " ")
		if len(parts) != 3 {
			return fmt.Errorf("wrong line in the Bond Coeffs section (line number in there: %d)", bondTypeLineNumber+1)
		}
		number := parts[0]
		sth1, err := strconv.Atoi(parts[1])
		if err != nil {
			return err
		}
		sth2, err := strconv.ParseFloat(parts[2], 64)
		if err != nil {
			return err
		}
		loader.bondTypes[number] = struct {
			sth1 int
			sth2 float64
		}{
			sth1: sth1,
			sth2: sth2,
		}
	}

	return nil
}

func (loader *LammpsLoader) loadAtoms() error {
	for loader.scanner.Scan() && loader.scanner.Text() != "Atoms # full" {
	}
	loader.scanner.Scan()

	for atomLineNumber := 0; atomLineNumber < loader.atomsCount && loader.scanner.Scan(); atomLineNumber++ {
		parts := strings.Split(loader.scanner.Text(), " ")
		if len(parts) != 10 {
			return fmt.Errorf("wrong line in the Atoms section (line number in there: %d)", atomLineNumber+1)
		}

		atomID, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return err
		}

		polymerID, err := strconv.Atoi(parts[1])
		if err != nil {
			return err
		}

		atomTypeNumber := parts[2]
		atomType, _ := strconv.Atoi(atomTypeNumber)

		// Miss the charge because of not needing it

		x, err := strconv.ParseFloat(parts[4], 64)
		if err != nil {
			return err
		}

		y, err := strconv.ParseFloat(parts[5], 64)
		if err != nil {
			return err
		}

		z, err := strconv.ParseFloat(parts[6], 64)
		if err != nil {
			return err
		}

		monomer := NewAtom(loader.atomTypes[atomTypeNumber].Label,
			int(atomID),
			polymerID,
			atomType,
			0.0,
			x,
			y,
			z)

		loader.atoms[atomID] = struct {
			Atom        *Atom
			ChainNumber int
		}{Atom: monomer, ChainNumber: polymerID}
	}
	return nil
}

func (loader *LammpsLoader) loadBonds() error {
	for loader.scanner.Scan() && loader.scanner.Text() != "Bonds" {
	}
	loader.scanner.Scan()

	loader.bonds = make([]*Bond, 0)
	for bondLineNumber := 0; bondLineNumber < loader.bondsCount && loader.scanner.Scan(); bondLineNumber++ {
		parts := strings.Split(loader.scanner.Text(), " ")
		if len(parts) != 4 {
			return fmt.Errorf("wrong line in the Bonds section (line number in there: %d)", bondLineNumber+1)
		}
		bondID, err := strconv.Atoi(parts[0])
		if err != nil {
			return err
		}

		connectionType, err := strconv.Atoi(parts[1])
		if err != nil {
			return err
		}

		firstAtomID, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			return err
		}

		secondAtomID, err := strconv.ParseInt(parts[3], 10, 64)
		if err != nil {
			return err
		}

		firstAtom := loader.atoms[firstAtomID]
		secondAtom := loader.atoms[secondAtomID]

		bond := NewBond(bondID, connectionType, [2]*Atom{firstAtom.Atom, secondAtom.Atom})
		loader.bonds = append(loader.bonds, bond)
	}
	return nil
}

func (loader *LammpsLoader) constructLammpsStruct() error {
	loader.builtGlobula = LammpsStruct{}
	for i := range loader.atoms {
		loader.builtGlobula.Atoms = append(loader.builtGlobula.Atoms, *loader.atoms[i].Atom)
	}
	for i := range loader.bonds {
		loader.builtGlobula.Bonds = append(loader.builtGlobula.Bonds, *loader.bonds[i])
	}
	return nil
}
