package parser

import "lammps-file-parser/structs"

func Parse(content string) (structs.LammpsStruct, error) {
	loader := structs.LammpsLoader{}
	return loader.Load(content)
}
