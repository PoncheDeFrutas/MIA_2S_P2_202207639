package structures

import (
	"fmt"
	"strings"
	"time"
)

type FolderElement struct {
	Name string `json:"name"`
	ID   int32  `json:"id"`
	Type string `json:"type"`
}

func (sb *SuperBlock) GetFile(path string, index int32, filePath []string) string {
	inode := &Inode{}
	inodePath := int64(sb.SInodeStart + index*sb.SInodeSize)

	if err := inode.ReadInode(path, inodePath); err != nil {
		return ""
	}

	if inode.IType == '0' {
		inodeIndex := sb.findInodeInBlock(path, filePath[0], inode)
		if inodeIndex != -1 {
			return sb.GetFile(path, inodeIndex, filePath[1:])
		}
	} else if inode.IType == '1' {
		return sb.getFileContent(path, inode)
	}

	return "Error path not found"
}

func (sb *SuperBlock) getFileContent(path string, inode *Inode) string {
	var content strings.Builder

	for _, block := range inode.IBlock[:12] {
		if block == -1 {
			continue
		}
		content.WriteString(sb.getContentBlock(path, block))
	}

	if inode.IBlock[12] != -1 {
		content.WriteString(sb.getIndirectBlockContent(path, inode.IBlock[12], 1))
	}

	if inode.IBlock[13] != -1 {
		content.WriteString(sb.getIndirectBlockContent(path, inode.IBlock[13], 2))
	}

	if inode.IBlock[14] != -1 {
		content.WriteString(sb.getIndirectBlockContent(path, inode.IBlock[14], 3))
	}

	return content.String()
}

func (sb *SuperBlock) getContentBlock(path string, index int32) string {
	block := &FileBlock{}
	blockPath := int64(sb.SBlockStart + index*sb.SBlockSize)

	if err := block.ReadFileBlock(path, blockPath); err != nil {
		return ""
	}

	return strings.TrimRight(string(block.BContent[:]), "\x00")
}

func (sb *SuperBlock) getIndirectBlockContent(path string, blockIndex int32, level int) string {
	if level == 0 {
		return sb.getContentBlock(path, blockIndex)
	}

	block := &PointerBlock{}
	blockPath := int64(sb.SBlockStart + blockIndex*sb.SBlockSize)

	if err := block.ReadPointerBlock(path, blockPath); err != nil {
		return ""
	}

	var content strings.Builder
	for _, subBlockIndex := range block.PPointers {
		if subBlockIndex == -1 {
			continue
		}
		content.WriteString(sb.getIndirectBlockContent(path, subBlockIndex, level-1))
	}

	return content.String()
}

// WriteFile writes a file in the filesystem
func (sb *SuperBlock) WriteFile(path string, index int32, filePath []string, content string) (int, error) {
	inode := &Inode{}
	if err := inode.ReadInode(path, int64(sb.SInodeStart+index*sb.SInodeSize)); err != nil {
		return 0, err
	}

	if inode.IType == '0' {
		inodeIndex := sb.findInodeInBlock(path, filePath[0], inode)
		if inodeIndex != -1 {
			writtenBytes, err := sb.WriteFile(path, inodeIndex, filePath[1:], content)
			if err != nil {
				return 0, err
			}

			if writtenBytes > 0 {
				return writtenBytes, err
			}
		}
	} else if inode.IType == '1' {
		return sb.writeFileContent(path, inode, content, index)
	}

	return 0, nil
}

func (sb *SuperBlock) writeFileContent(path string, inode *Inode, content string, index int32) (int, error) {
	inode.IMTime = float32(time.Now().Unix())
	var err error
	for i, blockIndex := range inode.IBlock[:12] {
		if content == "" {
			break
		}
		if blockIndex == -1 {
			inode.IBlock[i] = sb.SBlocksCount
			inode.IMTime = float32(time.Now().Unix())

			content, err = sb.CreateFileBlock(path, content)
			if err != nil {
				return 0, err
			}
		} else {
			content, err = sb.WriteFileBlock(path, blockIndex, content)
			if err != nil {
				return 0, err
			}
		}
	}

	for i, blockIndex := range inode.IBlock[12:] {
		if content == "" {
			break
		}

		if blockIndex == -1 {
			inode.IBlock[i] = sb.SBlocksCount
			inode.IMTime = float32(time.Now().Unix())

			if err := sb.CreatePointerBlock(path, i); err != nil {
				return 0, err
			}

			content, err = sb.CreateFileBlock(path, content)
			if err != nil {
				return 0, err
			}
		} else {

		}
	}

	if err := inode.WriteInode(path, int64(sb.SInodeStart+index*sb.SInodeSize),
		int64(sb.SInodeStart+(index+1)*sb.SInodeSize)); err != nil {
		return 0, err
	}

	return 0, nil
}

func (sb *SuperBlock) writePointerContent(path string, blockIndex, level int32, content string) (string, error) {
	block := &PointerBlock{}
	blockPath := int64(sb.SBlockStart + blockIndex*sb.SBlockSize)
	var err error
	err = block.ReadPointerBlock(path, blockPath)
	if err != nil {
		return "", err
	}

	for i, pointer := range block.PPointers {
		if content == "" {
			break
		}
		if pointer == -1 {
			block.PPointers[i] = sb.SBlocksCount
			if err := block.WritePointerBlock(path, blockPath, blockPath+int64(sb.SBlockSize)); err != nil {
				return "", err
			}

			if level != 0 {
				if err := sb.CreatePointerBlock(path, int(level-1)); err != nil {
					return "", err
				}
			}

			content, err = sb.CreateFileBlock(path, content)
			if err != nil {
				return "", err
			}
		} else {
			if level == 0 {
				content, err = sb.WriteFileBlock(path, pointer, content)
				if err != nil {
					return "", err
				}
			} else {
				content, err = sb.writePointerContent(path, pointer, level-1, content)
				if err != nil {
					return "", err
				}
			}
		}
	}

	return content, nil
}

// CreateInode creates a new inode in the filesystem
func (sb *SuperBlock) CreateInode(path string, isFile bool) error {
	if sb.SFreeInodeCount == 0 {
		return fmt.Errorf("no free inodes")
	}

	newInode := &Inode{}
	newInode.DefaultValue(-1)
	if isFile {
		newInode.IType = '1'
	} else {
		newInode.IType = '0'
	}
	newInode.IPerm = [3]byte{'6', '6', '4'}

	if err := newInode.WriteInode(path, int64(sb.SFirstIno), int64(sb.SFirstIno+sb.SInodeSize)); err != nil {
		return err
	}

	if err := sb.UpdateBitmapInode(path); err != nil {
		return err
	}

	return nil
}

// CreateFolderBlock creates a new folder block in the filesystem
func (sb *SuperBlock) CreateFolderBlock(path, name string, indexInode int32) error {
	if sb.SFreeBlockCount == 0 {
		return fmt.Errorf("no free blocks")
	}

	newBlock := &FolderBlock{}
	newBlock.DefaultValue()
	newBlock.BContent[0].BInode = indexInode // FIX CURRENT
	newBlock.BContent[1].BInode = indexInode // FIX FATHER

	newBlock.BContent[2].BInode = sb.SInodesCount
	copy(newBlock.BContent[2].BName[:], name)

	if err := newBlock.WriteFolderBlock(path, int64(sb.SFirstBlo), int64(sb.SFirstBlo+sb.SBlockSize)); err != nil {
		return err
	}

	if err := sb.UpdateBitmapBlock(path); err != nil {
		return err
	}

	return nil
}

// CreateFileBlock creates a new file block in the filesystem
func (sb *SuperBlock) CreateFileBlock(path string, content string) (string, error) {
	if sb.SFreeBlockCount == 0 {
		return "", fmt.Errorf("no free blocks")
	}

	newBlock := &FileBlock{}
	toWrite := min(len(content), 64)
	copy(newBlock.BContent[:], content[:toWrite])

	if err := newBlock.WriteFileBlock(path, int64(sb.SFirstBlo), int64(sb.SFirstBlo+sb.SBlockSize)); err != nil {
		return "", err
	}

	if err := sb.UpdateBitmapBlock(path); err != nil {
		return "", err
	}

	return content[toWrite:], nil
}

func (sb *SuperBlock) WriteFileBlock(path string, index int32, content string) (string, error) {
	block := &FileBlock{}
	blockPath := int64(sb.SBlockStart + index*sb.SBlockSize)

	if err := block.ReadFileBlock(path, blockPath); err != nil {
		return "", err
	}

	toWrite := min(len(content), 64)
	copy(block.BContent[:], content[:toWrite])

	if err := block.WriteFileBlock(path, blockPath, blockPath+int64(sb.SBlockSize)); err != nil {
		return "", err
	}

	return content[toWrite:], nil
}

// CreatePointerBlock creates a new pointer block in the filesystem
func (sb *SuperBlock) CreatePointerBlock(path string, level int) error {
	if sb.SFreeBlockCount == 0 {
		return fmt.Errorf("no free blocks")
	}

	newBlock := &PointerBlock{}
	newBlock.DefaultValue()
	newBlock.PPointers[0] = sb.SBlocksCount + 1

	if err := newBlock.WritePointerBlock(path, int64(sb.SFirstBlo), int64(sb.SFirstBlo+sb.SBlockSize)); err != nil {
		return err
	}

	if err := sb.UpdateBitmapBlock(path); err != nil {
		return err
	}

	if level > 0 {
		return sb.CreatePointerBlock(path, level-1)
	}

	return nil
}

// CreateNewInode creates a new inode in the filesystem (File/Folder)
func (sb *SuperBlock) CreateNewInode(path string, filePath []string, indexInode int32, isFile, root bool) error {
	inode := &Inode{}
	inodePath := int64(sb.SInodeStart + indexInode*sb.SInodeSize)

	if err := inode.ReadInode(path, inodePath); err != nil {
		return err
	}

	if inode.IType == '1' {
		return fmt.Errorf("invalid inode type: file")
	}

	if len(filePath) == 1 {
		return sb.CreatePath(path, filePath[0], inode, isFile, indexInode)
	}

	newIndexInode := sb.findInodeInBlock(path, filePath[0], inode)

	if newIndexInode != -1 {
		return sb.CreateNewInode(path, filePath[1:], newIndexInode, isFile, root)
	}

	if root {
		if err := sb.CreatePath(path, filePath[0], inode, false, indexInode); err != nil {
			return err
		}
		newIndexInode := sb.findInodeInBlock(path, filePath[0], inode)
		return sb.CreateNewInode(path, filePath[1:], newIndexInode, isFile, root)
	}

	return nil
}

// CreatePath creates a new path in the filesystem
func (sb *SuperBlock) CreatePath(path, name string, inode *Inode, isFile bool, indexInode int32) error {
	for i, blockIndex := range inode.IBlock[:12] {
		if blockIndex == -1 {
			inode.IBlock[i] = sb.SBlocksCount
			inode.IMTime = float32(time.Now().Unix())

			if err := sb.CreateBlockAndWriteInode(path, name, inode, indexInode); err != nil {
				return err
			}
		} else {
			if condition, err := sb.addContentToFolderBlock(path, name, blockIndex); err != nil {
				return err
			} else if !condition {
				continue
			}
		}

		return sb.CreateInode(path, isFile)
	}

	for i, blockIndex := range inode.IBlock[12:] {
		if blockIndex == -1 {
			inode.IBlock[i+12] = sb.SBlocksCount
			inode.IMTime = float32(time.Now().Unix())

			if err := sb.CreatePointerBlock(path, i); err != nil {
				return err
			}

			if err := sb.CreateBlockAndWriteInode(path, name, inode, indexInode); err != nil {
				return err
			}

		} else {
			if condition, err := sb.addContentToPointerBlock(path, name, blockIndex, indexInode, int32(i)); err != nil {
				return err
			} else if !condition {
				continue
			}
		}

		return sb.CreateInode(path, isFile)
	}

	return fmt.Errorf("no free blocks")
}

// CreateBlockAndWriteInode creates a new block and writes the inode in the filesystem
func (sb *SuperBlock) CreateBlockAndWriteInode(path, name string, inode *Inode, indexInode int32) error {
	inodeStart := int64(sb.SInodeStart + indexInode*sb.SInodeSize)
	inodeEnd := int64(sb.SInodeStart + (indexInode+1)*sb.SInodeSize)

	if err := inode.WriteInode(path, inodeStart, inodeEnd); err != nil {
		return err
	}

	if err := sb.CreateFolderBlock(path, name, indexInode); err != nil {
		return err
	}

	return nil
}

// addContentToFolderBlock adds a new content to a folder block
func (sb *SuperBlock) addContentToFolderBlock(path, name string, blockIndex int32) (bool, error) {
	block := &FolderBlock{}

	if err := block.ReadFolderBlock(path, int64(sb.SBlockStart+blockIndex*sb.SBlockSize)); err != nil {
		return false, err
	}

	if block.BContent[3].BInode != -1 {
		return false, nil
	}

	copy(block.BContent[3].BName[:], name)
	block.BContent[3].BInode = sb.SInodesCount

	if err := block.WriteFolderBlock(path, int64(sb.SBlockStart+blockIndex*sb.SBlockSize),
		int64(sb.SBlockStart+(blockIndex+1)*sb.SBlockSize)); err != nil {
		return false, err
	}

	return true, nil
}

// addContentToPointerBlock adds a new content to a pointer block
func (sb *SuperBlock) addContentToPointerBlock(path, name string, blockIndex, indexInode, level int32) (bool, error) {
	block := &PointerBlock{}

	if err := block.ReadPointerBlock(path, int64(sb.SBlockStart+blockIndex*sb.SBlockSize)); err != nil {
		return false, err
	}

	for i, pointer := range block.PPointers {
		if pointer == -1 {
			block.PPointers[i] = sb.SBlocksCount

			if err := block.WritePointerBlock(path, int64(sb.SBlockStart+blockIndex*sb.SBlockSize),
				int64(sb.SBlockStart+(blockIndex+1)*sb.SBlockSize)); err != nil {
				return false, err
			}
			if level != 0 {
				if err := sb.CreatePointerBlock(path, int(level-1)); err != nil {
					return false, err
				}
			}
			if err := sb.CreateFolderBlock(path, name, sb.SInodesCount); err != nil {
				return false, err
			}
			return true, nil
		}

		if level == 0 {
			if condition, err := sb.addContentToFolderBlock(path, name, pointer); err != nil {
				return false, err
			} else if !condition {
				continue
			} else {
				return true, nil
			}
		} else {
			if condition, err := sb.addContentToPointerBlock(path, name, pointer, indexInode, level-1); err != nil {
				return false, err
			} else if !condition {
				continue
			} else {
				return true, nil
			}
		}
	}

	return false, nil
}

// GetIndexInode returns the index of an inode in a block
func (sb *SuperBlock) GetIndexInode(path, file string, index int32) int32 {
	block := &FolderBlock{}
	blockPath := int64(sb.SBlockStart + index*sb.SBlockSize)

	if err := block.ReadFolderBlock(path, blockPath); err != nil {
		return -1
	}

	for _, entry := range block.BContent[2:] {
		name := strings.TrimRight(string(entry.BName[:]), "\x00")
		if name == file {
			return entry.BInode
		}
	}
	return -1
}

// findInodeInBlock returns the index of an inode in a block
func (sb *SuperBlock) findInodeInBlock(path, part string, inode *Inode) int32 {
	for _, block := range inode.IBlock[:12] {
		if block == -1 {
			return -1
		}

		blockIndex := sb.GetIndexInode(path, part, block)
		if blockIndex != -1 {
			return blockIndex
		}
	}

	for i, block := range inode.IBlock[12:] {
		if block == -1 {
			return -1
		}

		blockIndex := sb.finInodeInPointerBlock(path, part, block, int32(i))
		if blockIndex != -1 {
			return blockIndex
		}

	}

	return -1
}

// finInodeInPointerBlock returns the index of an inode in a pointer block
func (sb *SuperBlock) finInodeInPointerBlock(path, part string, blockIndex, level int32) int32 {
	block := &PointerBlock{}
	blockPath := int64(sb.SBlockStart + blockIndex*sb.SBlockSize)

	if err := block.ReadPointerBlock(path, blockPath); err != nil {
		return -1
	}

	for _, pointer := range block.PPointers {
		if pointer == -1 {
			return -1
		}

		if level == 0 {
			blockIndex := sb.GetIndexInode(path, part, pointer)
			if blockIndex != -1 {
				return blockIndex
			}
		} else {
			blockIndex := sb.finInodeInPointerBlock(path, part, pointer, level-1)
			if blockIndex != -1 {
				return blockIndex
			}
		}
	}

	return -1
}

func (sb *SuperBlock) GetInodeElements(path string, index int32) ([]FolderElement, error) {
	inode := &Inode{}
	inodePath := int64(sb.SInodeStart + index*sb.SInodeSize)

	if err := inode.ReadInode(path, inodePath); err != nil {
		return nil, err
	}

	var allElements []FolderElement

	if index != -1 {
		for _, block := range inode.IBlock[:12] {
			if block == -1 {
				continue
			}
			elements, err := sb.GetFolderElements(path, block)
			if err != nil {
				continue
			}

			allElements = append(allElements, elements...)
		}
	}

	return allElements, nil
}

func (sb *SuperBlock) GetFolderElements(path string, index int32) ([]FolderElement, error) {
	block := &FolderBlock{}
	blockPath := int64(sb.SBlockStart + index*sb.SBlockSize)

	if err := block.ReadFolderBlock(path, blockPath); err != nil {
		return nil, err
	}

	inode := &Inode{}
	var elements []FolderElement

	for _, entry := range block.BContent[2:] {
		if entry.BInode == -1 {
			continue
		}
		inodePath := int64(sb.SInodeStart + entry.BInode*sb.SInodeSize)
		if err := inode.ReadInode(path, inodePath); err != nil {
			continue
		}

		element := FolderElement{
			Name: strings.TrimRight(string(entry.BName[:]), "\x00"),
			ID:   entry.BInode,
			Type: func() string {
				if inode.IType == '1' {
					return "file"
				}
				return "folder"
			}(),
		}

		elements = append(elements, element)
	}

	return elements, nil
}

func (sb *SuperBlock) GetInodeReference(path string, index int32, filePath []string) int32 {
	if len(filePath) == 0 {
		return index
	}

	inode := &Inode{}
	inodePath := int64(sb.SInodeStart + index*sb.SInodeSize)

	if err := inode.ReadInode(path, inodePath); err != nil {
		return -1
	}

	if len(filePath) == 1 {
		for _, block := range inode.IBlock[:12] {
			if block == -1 {
				break
			}
			newIndexInode := sb.GetIndexInode(path, filePath[0], block)
			if newIndexInode != -1 {
				return newIndexInode
			}
		}
	}

	newIndexInode := sb.findInodeInBlock(path, filePath[0], inode)

	if newIndexInode != -1 {
		return sb.GetInodeReference(path, newIndexInode, filePath[1:])
	}

	return -1
}
