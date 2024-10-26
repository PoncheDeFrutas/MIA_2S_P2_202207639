package structures

import (
	"backend/utils"
	"fmt"
	"strings"
)

type FileBlock struct {
	BContent [64]byte
	// Total size of the FileBlock is 64 bytes
}

func (f *FileBlock) WriteFileBlock(path string, offset int64, maxSize int64) error {
	if err := utils.WriteToFile(path, offset, maxSize, f); err != nil {
		return err
	}
	return nil
}

func (f *FileBlock) ReadFileBlock(path string, offset int64) error {
	if err := utils.ReadFromFile(path, offset, f); err != nil {
		return err
	}
	return nil
}

func (f *FileBlock) Print() {
	fmt.Printf("--*-- FileBlock --*--\n")
	fmt.Printf("BContent: %s\n", f.BContent)
}

func (f *FileBlock) GetStringBuilder(nodeName string) string {
	var sb strings.Builder

	content := strings.ReplaceAll(string(f.BContent[:]), "\x00", "")
	content2 := strings.ReplaceAll(content, "\n", "<br/>")

	sb.WriteString(fmt.Sprintf("    %s [label=<\n", nodeName))
	sb.WriteString(fmt.Sprintf("    <TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\">\n"))

	sb.WriteString(fmt.Sprintf("\t<TR><TD COLSPAN=\"2\" BGCOLOR=\"%s\"><B>%s</B></TD></TR>\n", "#333333", nodeName))
	sb.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">Content</TD><TD WIDTH=\"250\" BGCOLOR=\"%s\">%s</TD></TR>\n", "#DDDDDD", "#DDDDDD", content2))

	sb.WriteString("    </TABLE>>]\n")

	return sb.String()
}
