package parser

import "github.com/Ivanestver/lammps-file-parser/structs"

func Parse(content, fileName string) (*structs.LammpsStruct, error) {
	loader := structs.LammpsLoader{}
	if result, err := loader.Load(content); err != nil {
		return nil, err
	} else {
		result.FileName = fileName
		return result, nil
	}
}
