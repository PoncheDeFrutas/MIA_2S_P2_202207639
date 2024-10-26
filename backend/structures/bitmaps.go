package structures

import (
	"encoding/binary"
	"os"
)

func (sb *SuperBlock) CreateBitMaps(path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	if _, err = file.Seek(int64(sb.SBMInodeStart), 0); err != nil {
		return err
	}

	buffer := make([]byte, sb.SFreeInodeCount)
	for i := 0; i < len(buffer); i++ {
		buffer[i] = '0'
	}

	if err = binary.Write(file, binary.BigEndian, buffer); err != nil {
		return err
	}

	if _, err = file.Seek(int64(sb.SBMBlockStart), 0); err != nil {
		return err
	}

	buffer = make([]byte, sb.SFreeBlockCount)
	for i := 0; i < len(buffer); i++ {
		buffer[i] = 'O'
	}

	if err = binary.Write(file, binary.BigEndian, buffer); err != nil {
		return err
	}

	return nil
}

func (sb *SuperBlock) UpdateBitmapInode(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	if _, err = file.Seek(int64(sb.SBMInodeStart)+int64(sb.SInodesCount), 0); err != nil {
		return err
	}

	if _, err := file.Write([]byte{'1'}); err != nil {
		return err
	}

	sb.SInodesCount++
	sb.SFreeInodeCount--
	sb.SFirstIno += sb.SInodeSize

	return nil
}

func (sb *SuperBlock) UpdateBitmapBlock(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	if _, err = file.Seek(int64(sb.SBMBlockStart)+int64(sb.SBlocksCount), 0); err != nil {
		return err
	}

	if _, err := file.Write([]byte{'X'}); err != nil {
		return err
	}

	sb.SBlocksCount++
	sb.SFreeBlockCount--
	sb.SFirstBlo += sb.SBlockSize

	return nil
}
