package structures

import (
	"backend/utils"
	"fmt"
	"strings"
)

type FolderBlock struct {
	BContent [4]FolderContent
	// Total size of the FolderBlock is 64 bytes
}

type FolderContent struct {
	BName  [12]byte
	BInode int32
	// Total size of the FolderContent is 16 bytes
}

func (f *FolderBlock) DefaultValue() {
	f.BContent = [4]FolderContent{
		{BName: [12]byte{'.'}, BInode: 0},
		{BName: [12]byte{'.', '.'}, BInode: 0},
		{BName: [12]byte{'-'}, BInode: -1},
		{BName: [12]byte{'-'}, BInode: -1},
	}
}

func (f *FolderBlock) WriteFolderBlock(path string, offset int64, maxSize int64) error {
	if err := utils.WriteToFile(path, offset, maxSize, f); err != nil {
		return err
	}
	return nil
}

func (f *FolderBlock) ReadFolderBlock(path string, offset int64) error {
	if err := utils.ReadFromFile(path, offset, f); err != nil {
		return err
	}
	return nil
}

func (f *FolderBlock) Print() {
	for i, content := range f.BContent {
		content.Print(i)
	}
}

func (f *FolderContent) Print(i int) {
	println("FolderContent", i)
	println("BName", string(f.BName[:]))
	println("BInode", f.BInode)
	println()
}

func (f *FolderBlock) GetStringBuilder(nodeName string) string {
	var sb strings.Builder

	if f.BContent[2].BInode != -1 {
		sb.WriteString(fmt.Sprintf("%s -> Inodo_%d\n", nodeName, f.BContent[2].BInode))
		//nodeName -> Inodo_(BINODE[2])
	}
	if f.BContent[3].BInode != -1 {
		sb.WriteString(fmt.Sprintf("%s -> Inodo_%d\n", nodeName, f.BContent[3].BInode))
	}
	sb.WriteString(fmt.Sprintf("    %s [label=<\n", nodeName))
	sb.WriteString(fmt.Sprintf("    <TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\">\n"))

	sb.WriteString(fmt.Sprintf("\t<TR><TD COLSPAN=\"2\" BGCOLOR=\"%s\"><B>%s</B></TD></TR>\n", "#333333", nodeName))
	sb.WriteString(fmt.Sprintf("\t<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">BName</TD><TD WIDTH=\"150\" BGCOLOR=\"%s\">BINode</TD></TR>\n", "#DDDDDD", "#DDDDDD"))
	for _, Content := range f.BContent {
		sb.WriteString(Content.GetStringBuilder())
	}

	sb.WriteString("    </TABLE>>]\n")

	return sb.String()
}

func (f *FolderContent) GetStringBuilder() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("\t<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">%s</TD><TD WIDTH=\"150\" BGCOLOR=\"%s\">%d</TD></TR>\n", "#DDDDDD", strings.TrimRight(string(f.BName[:]), "\x00"), "#DDDDDD", f.BInode))

	return sb.String()
}
