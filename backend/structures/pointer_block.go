package structures

import (
	"backend/utils"
	"fmt"
	"strings"
)

type PointerBlock struct {
	PPointers [16]int32
	// Total size of the PointerBlock is 64 bytes
}

func (p *PointerBlock) DefaultValue() {
	for i := range p.PPointers {
		p.PPointers[i] = -1
	}
}

func (p *PointerBlock) FindFreePointer() int32 {
	for i, pointer := range p.PPointers {
		if pointer == -1 {
			return int32(i)
		}
	}
	return -1
}

func (p *PointerBlock) WritePointerBlock(path string, offset int64, maxSize int64) error {
	if err := utils.WriteToFile(path, offset, maxSize, p); err != nil {
		return err
	}
	return nil
}

func (p *PointerBlock) ReadPointerBlock(path string, offset int64) error {
	if err := utils.ReadFromFile(path, offset, p); err != nil {
		return err
	}
	return nil
}

func (p *PointerBlock) GetStringBuilder(nodeName string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("    %s [label=<\n", nodeName))
	sb.WriteString(fmt.Sprintf("    <TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\">\n"))
	sb.WriteString(fmt.Sprintf("\t<TR><TD COLSPAN=\"2\" BGCOLOR=\"%s\"><B>%s</B></TD></TR>\n", "#333333", nodeName))

	for i, pointer := range p.PPointers {
		sb.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">Pointer %d</TD><TD WIDTH=\"250\" BGCOLOR=\"%s\">%d</TD></TR>\n", "#DDDDDD", i, "#DDDDDD", pointer))
	}

	sb.WriteString("    </TABLE>>]\n")

	return sb.String()

}
