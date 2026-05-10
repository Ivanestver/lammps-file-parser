package deserialize

import "github.com/Ivanestver/lammps-file-parser/structs"

/*
Deserialize function converts a valid LAMMPS file into a set of objects
that represents the file.

Params:
  - content: a LAMMPS file content
  - filename: a LAMMPS file name

Returns:
  - LammpsStruct: the file contents' representation
  - error: any error occured
*/
func Deserialize(content, fileName string) (*structs.LammpsStruct, error) {
	loader := structs.LammpsLoader{}
	if result, err := loader.Load(content); err != nil {
		return nil, err
	} else {
		result.FileName = fileName
		return result, nil
	}
}
