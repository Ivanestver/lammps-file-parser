package parser

import "github.com/Ivanestver/lammps-file-parser/structs"

func Parse(content string) (structs.LammpsStruct, error) {
	loader := structs.LammpsLoader{}
	return loader.Load(content)
}
