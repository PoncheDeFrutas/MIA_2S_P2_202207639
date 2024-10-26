package global

import (
	"backend/structures"
	"errors"
)

const Carnet string = "39"

var (
	MountedPartitions = make(map[string]string) // id -> path
)

func GetMountedPartition(id string) (*structures.Partition, string, error) {
	path := MountedPartitions[id]
	if path == "" {
		return nil, "", errors.New("partition not mounted with id: " + id)
	}

	mbr := &structures.MBR{}

	if err := mbr.ReadMBR(path); err != nil {
		return nil, "", err
	}

	partition, err := mbr.GetPartitionByID(id)
	if err != nil {
		return nil, "", err
	}

	return partition, path, nil
}

func PrintMountedPartitions() string {
	result := "Mounted Partitions:\n"

	for id, path := range MountedPartitions {
		result += id + " -> " + path + "\n"
	}

	return result
}
