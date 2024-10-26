package structures

import (
	"backend/utils"
	"encoding/binary"
	"fmt"
	"strings"
	"time"
)

type SuperBlock struct {
	SFilesystemType int32
	SInodesCount    int32
	SBlocksCount    int32
	SFreeInodeCount int32
	SFreeBlockCount int32
	SMTime          float32
	SUmTime         float32
	SMntCount       int32
	SMagic          int32
	SInodeSize      int32
	SBlockSize      int32
	SFirstIno       int32
	SFirstBlo       int32
	SBMBlockStart   int32
	SBMInodeStart   int32
	SInodeStart     int32
	SBlockStart     int32
	// Total size of the SuperBlock is 68 bytes
}

func (sb *SuperBlock) CreateSuperBlock(partitionStart int32, n int32) {
	//Bitmaps
	bmInodeStart := partitionStart + int32(binary.Size(sb))
	bmBlockStart := bmInodeStart + n

	//Inodes
	inodeStart := bmBlockStart + (3 * n)

	//Blocks
	blockStart := inodeStart + (int32(binary.Size(Inode{})) * n)

	sb.SFilesystemType = 2
	sb.SInodesCount = 0
	sb.SBlocksCount = 0
	sb.SFreeInodeCount = n
	sb.SFreeBlockCount = n * 3
	sb.SMTime = float32(time.Now().Unix())
	sb.SUmTime = float32(time.Now().Unix())
	sb.SMntCount = 1
	sb.SMagic = 0xEF53
	sb.SInodeSize = int32(binary.Size(Inode{}))
	sb.SBlockSize = int32(binary.Size(FileBlock{}))
	sb.SFirstIno = inodeStart
	sb.SFirstBlo = blockStart
	sb.SBMInodeStart = bmInodeStart
	sb.SBMBlockStart = bmBlockStart
	sb.SInodeStart = inodeStart
	sb.SBlockStart = blockStart
}

func (sb *SuperBlock) WriteSuperBlock(path string, offset int64, maxSize int64) error {
	if err := utils.WriteToFile(path, offset, maxSize, sb); err != nil {
		return err
	}
	return nil
}

func (sb *SuperBlock) ReadSuperBlock(path string, offset int64) error {
	if err := utils.ReadFromFile(path, offset, sb); err != nil {
		return err
	}
	return nil
}

func (sb *SuperBlock) CreateUserFile(path string) error {
	if err := sb.createRootInodeAndBlock(path); err != nil {
		return err
	}

	if err := sb.createUsersFile(path); err != nil {
		return err
	}

	return nil
}

func (sb *SuperBlock) createRootInodeAndBlock(path string) error {
	rootInode := &Inode{}
	rootInode.DefaultValue(sb.SBlocksCount)

	if err := rootInode.WriteInode(path, int64(sb.SFirstIno), int64(sb.SFirstIno+sb.SInodeSize)); err != nil {
		return err
	}

	if err := sb.UpdateBitmapInode(path); err != nil {
		return err
	}

	rootBlock := &FolderBlock{}
	rootBlock.DefaultValue()

	if err := rootBlock.WriteFolderBlock(path, int64(sb.SFirstBlo), int64(sb.SFirstBlo+sb.SBlockSize)); err != nil {
		return err
	}

	if err := sb.UpdateBitmapBlock(path); err != nil {
		return err
	}

	fmt.Println("rootInode")
	rootInode.Print()

	fmt.Println("rootBlock")
	rootBlock.Print()

	return nil
}

func (sb *SuperBlock) createUsersFile(path string) error {
	usersText := "1,G,root\n1,U,root,root,123\n"

	rootInode := &Inode{}
	if err := rootInode.ReadInode(path, int64(sb.SInodeStart+0)); err != nil {
		return err
	}

	rootInode.IAtime = float32(time.Now().Unix())

	if err := rootInode.WriteInode(path, int64(sb.SInodeStart+0), int64(sb.SInodeStart+sb.SInodeSize)); err != nil {
		return err
	}

	rootBlock := &FolderBlock{}
	if err := rootBlock.ReadFolderBlock(path, int64(sb.SBlockStart+0)); err != nil {
		return err
	}

	rootBlock.BContent[2] = FolderContent{BName: [12]byte{'u', 's', 'e', 'r', 's', '.', 't', 'x', 't'}, BInode: sb.SInodesCount}

	if err := rootBlock.WriteFolderBlock(path, int64(sb.SBlockStart+0), int64(sb.SBlockStart+sb.SBlockSize)); err != nil {
		return err
	}

	usersInode := &Inode{}
	usersInode.DefaultValue(sb.SBlocksCount)
	usersInode.ISize = int32(len(usersText))
	usersInode.IType = '1'

	if err := usersInode.WriteInode(path, int64(sb.SFirstIno), int64(sb.SFirstIno+sb.SInodeSize)); err != nil {
		return err
	}

	if err := sb.UpdateBitmapInode(path); err != nil {
		return err
	}

	usersBlock := &FileBlock{}
	copy(usersBlock.BContent[:], usersText)

	if err := usersBlock.WriteFileBlock(path, int64(sb.SFirstBlo), int64(sb.SFirstBlo+sb.SBlockSize)); err != nil {
		return err
	}

	if err := sb.UpdateBitmapBlock(path); err != nil {
		return err
	}

	return nil
}

func (sb *SuperBlock) Print() {
	fmt.Printf("SFilesystemType: %d\n", sb.SFilesystemType)
	fmt.Printf("SInodesCount: %d\n", sb.SInodesCount)
	fmt.Printf("SBlocksCount: %d\n", sb.SBlocksCount)
	fmt.Printf("SFreeInodeCount: %d\n", sb.SFreeInodeCount)
	fmt.Printf("SFreeBlockCount: %d\n", sb.SFreeBlockCount)
	fmt.Printf("SMTime: %s\n", time.Unix(int64(sb.SMTime), 0))
	fmt.Printf("SUMTime: %s\n", time.Unix(int64(sb.SUmTime), 0))
	fmt.Printf("SMntCount: %d\n", sb.SMntCount)
	fmt.Printf("SMagic: %d\n", sb.SMagic)
	fmt.Printf("SInodeSize: %d\n", sb.SInodeSize)
	fmt.Printf("SBlockSize: %d\n", sb.SBlockSize)
	fmt.Printf("SFirstIno: %d\n", sb.SFirstIno)
	fmt.Printf("SFirstBlo: %d\n", sb.SFirstBlo)
	fmt.Printf("SBMBlockStart: %d\n", sb.SBMBlockStart)
	fmt.Printf("SBMInodeStart: %d\n", sb.SBMInodeStart)
	fmt.Printf("SInodeStart: %d\n", sb.SInodeStart)
	fmt.Printf("SBlockStart: %d\n", sb.SBlockStart)
}

func (sb *SuperBlock) GetStringBuilder() string {
	var stringB strings.Builder

	// Main title row
	stringB.WriteString(fmt.Sprintf("\t<TR><TD COLSPAN=\"2\" BGCOLOR=\"%s\"><B>SuperBlock</B></TD></TR>\n", "#333333"))

	// SuperBlock rows
	stringB.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">File System Type</TD><TD WIDTH=\"250\" BGCOLOR=\"%s\">%d</TD></TR>\n", "#DDDDDD", "#DDDDDD", sb.SFilesystemType))
	stringB.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">Inodes Count</TD><TD WIDTH=\"250\" BGCOLOR=\"%s\">%d</TD></TR>\n", "#FFFFFF", "#FFFFFF", sb.SInodesCount))
	stringB.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">Blocks Count</TD><TD WIDTH=\"250\" BGCOLOR=\"%s\">%d</TD></TR>\n", "#DDDDDD", "#DDDDDD", sb.SBlocksCount))
	stringB.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">Free Inode Count</TD><TD WIDTH=\"250\" BGCOLOR=\"%s\">%d</TD></TR>\n", "#FFFFFF", "#FFFFFF", sb.SFreeInodeCount))
	stringB.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">Free Block Count</TD><TD WIDTH=\"250\" BGCOLOR=\"%s\">%d</TD></TR>\n", "#DDDDDD", "#DDDDDD", sb.SFreeBlockCount))
	stringB.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">MTime</TD><TD WIDTH=\"250\" BGCOLOR=\"%s\">%s</TD></TR>\n", "#FFFFFF", "#FFFFFF", time.Unix(int64(sb.SMTime), 0).Format("02-Jan-2006 03:04 PM")))
	stringB.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">UMTime</TD><TD WIDTH=\"250\" BGCOLOR=\"%s\">%s</TD></TR>\n", "#DDDDDD", "#DDDDDD", time.Unix(int64(sb.SUmTime), 0).Format("02-Jan-2006 03:04 PM")))
	stringB.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">MntCount</TD><TD WIDTH=\"250\" BGCOLOR=\"%s\">%d</TD></TR>\n", "#FFFFFF", "#FFFFFF", sb.SMntCount))
	stringB.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">Magic</TD><TD WIDTH=\"250\" BGCOLOR=\"%s\">%d</TD></TR>\n", "#DDDDDD", "#DDDDDD", sb.SMagic))
	stringB.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">Inode Size</TD><TD WIDTH=\"250\" BGCOLOR=\"%s\">%d</TD></TR>\n", "#FFFFFF", "#FFFFFF", sb.SInodeSize))
	stringB.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">Block Size</TD><TD WIDTH=\"250\" BGCOLOR=\"%s\">%d</TD></TR>\n", "#DDDDDD", "#DDDDDD", sb.SBlockSize))
	stringB.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">First Ino</TD><TD WIDTH=\"250\" BGCOLOR=\"%s\">%d</TD></TR>\n", "#FFFFFF", "#FFFFFF", sb.SFirstIno))
	stringB.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">First Blo</TD><TD WIDTH=\"250\" BGCOLOR=\"%s\">%d</TD></TR>\n", "#DDDDDD", "#DDDDDD", sb.SFirstBlo))
	stringB.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">BMBlock Start</TD><TD WIDTH=\"250\" BGCOLOR=\"%s\">%d</TD></TR>\n", "#FFFFFF", "#FFFFFF", sb.SBMBlockStart))
	stringB.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">BMInode Start</TD><TD WIDTH=\"250\" BGCOLOR=\"%s\">%d</TD></TR>\n", "#DDDDDD", "#DDDDDD", sb.SBMInodeStart))
	stringB.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">Inode Start</TD><TD WIDTH=\"250\" BGCOLOR=\"%s\">%d</TD></TR>\n", "#FFFFFF", "#FFFFFF", sb.SInodeStart))
	stringB.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">Block Start</TD><TD WIDTH=\"250\" BGCOLOR=\"%s\">%d</TD></TR>\n", "#DDDDDD", "#DDDDDD", sb.SBlockStart))

	return stringB.String()
}
