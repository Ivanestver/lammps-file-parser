package serialize

import "github.com/Ivanestver/lammps-file-parser/structs"

func Serialize(lammpsStruct *structs.LammpsStruct) (string, error) {
	serializer := _NewSerializer(lammpsStruct)
	return serializer.Serialize()
}
