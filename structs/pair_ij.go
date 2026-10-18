package structs

import (
	"fmt"
	"strconv"
	"strings"
)

type IPairIJFactory interface {
	CreateFromLine(line string) (IPairIJ, error)
}

func NewPairIJFactory(pairIJType string) IPairIJFactory {
	if strings.Contains(pairIJType, "dpd") {
		return &PairIJDPDFactory{}
	} else {
		return nil
	}
}

type IPairIJ interface {
	ToString() string
}

type PairIJDPDFactory struct {
}

func (pair *PairIJDPDFactory) CreateFromLine(line string) (IPairIJ, error) {
	parts := strings.Split(line, " ")
	if len(parts) < 4 {
		return nil, fmt.Errorf("в строке \"%s\" секции \"PairIJ Coeffs\" не хватает параметров", line)
	}
	I, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, err
	}
	J, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, err
	}
	Aij, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return nil, err
	}
	Gammaij, err := strconv.ParseFloat(parts[3], 64)
	if err != nil {
		return nil, err
	}
	cutoff := 0.0
	cutoffExists := len(parts) == 5
	if cutoffExists {
		cutoff, err = strconv.ParseFloat(parts[4], 64)
		if err != nil {
			return nil, err
		}
	}
	return &PairIJDPD{
		I:            I,
		J:            J,
		Aij:          Aij,
		Gammaij:      Gammaij,
		Cutoff:       cutoff,
		CutoffExists: cutoffExists,
	}, nil
}

type PairIJDPD struct {
	I            int     // Тип частицы
	J            int     // Тип частицы
	Aij          float64 // Коэффициент консервативной силы
	Gammaij      float64 // Коэффициент силы трения
	Cutoff       float64 // [Опционально] Радиус отсечения взаимодействия
	CutoffExists bool    // Определяет, задани ли радиус отсечения взаимодействия
}

func (pair *PairIJDPD) ToString() string {
	if pair.CutoffExists {
		return fmt.Sprintf("%d %d %s %s %s",
			pair.I, pair.J,
			strconv.FormatFloat(pair.Aij, 'f', -1, 64),
			strconv.FormatFloat(pair.Gammaij, 'f', -1, 64),
			strconv.FormatFloat(pair.Cutoff, 'f', -1, 64))
	} else {
		return fmt.Sprintf("%d %d %s %s",
			pair.I, pair.J,
			strconv.FormatFloat(pair.Aij, 'f', -1, 64),
			strconv.FormatFloat(pair.Gammaij, 'f', -1, 64))
	}
}
