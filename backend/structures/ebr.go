package structures

import (
	"backend/utils"
	"fmt"
	"strings"
)

type EBR struct {
	PartMount byte
	PartFit   byte
	PartStart int32
	PartSize  int32
	PartNext  int32
	PartName  [16]byte
	// Total size of the EBR is 30 bytes
}

func (e *EBR) DefaultValue() {
	e.PartMount = '9'
	e.PartFit = 'W'
	e.PartStart = -1
	e.PartSize = -1
	e.PartNext = -1
	copy(e.PartName[:], "EBR-LOGIC")
}

func (e *EBR) SetEBR(fit string, start int32, size int32, next int32, name string) {
	e.PartMount = '0'
	e.PartFit = fit[0]
	e.PartStart = start
	e.PartSize = size
	e.PartNext = next
	copy(e.PartName[:], name)
}

func (e *EBR) WriteEBR(path string, offset int64, maxSize int64) error {
	if err := utils.WriteToFile(path, offset, maxSize, e); err != nil {
		return err
	}
	return nil
}

func (e *EBR) ReadEBR(path string, offset int64) error {
	if err := utils.ReadFromFile(path, offset, e); err != nil {
		return err
	}
	return nil
}

func (e *EBR) Print() {
	fmt.Println("PartMount: ", string(e.PartMount))
	fmt.Println("PartFit: ", string(e.PartFit))
	fmt.Println("PartStart: ", e.PartStart)
	fmt.Println("PartSize: ", e.PartSize)
	fmt.Println("PartNext: ", e.PartNext)
	fmt.Println("PartName: ", string(e.PartName[:]))
}

func (e *EBR) GetStringBuilder() string {
	var sb strings.Builder
	eName := strings.TrimRight(string(e.PartName[:]), "\x00")
	sb.WriteString(fmt.Sprintf("<TR><TD COLSPAN=\"2\" BGCOLOR=\"%s\"><B>Partición Lógica (EBR)</B></TD></TR>\n", "#AAAAAA"))
	sb.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"100\" BGCOLOR=\"%s\">Status</TD><TD WIDTH=\"200\" BGCOLOR=\"%s\">%c</TD></TR>\n", "#DDDDDD", "#DDDDDD", e.PartMount))
	sb.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"100\" BGCOLOR=\"%s\">Fit</TD><TD WIDTH=\"200\" BGCOLOR=\"%s\">%c</TD></TR>\n", "#FFFFFF", "#FFFFFF", e.PartFit))
	sb.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"100\" BGCOLOR=\"%s\">Start</TD><TD WIDTH=\"200\" BGCOLOR=\"%s\">%d</TD></TR>\n", "#DDDDDD", "#DDDDDD", e.PartStart))
	sb.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"100\" BGCOLOR=\"%s\">Size</TD><TD WIDTH=\"200\" BGCOLOR=\"%s\">%d</TD></TR>\n", "#FFFFFF", "#FFFFFF", e.PartSize))
	sb.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"100\" BGCOLOR=\"%s\">Next</TD><TD WIDTH=\"200\" BGCOLOR=\"%s\">%d</TD></TR>\n", "#DDDDDD", "#DDDDDD", e.PartNext))
	sb.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"100\" BGCOLOR=\"%s\">Name</TD><TD WIDTH=\"200\" BGCOLOR=\"%s\">%s</TD></TR>\n", "#FFFFFF", "#FFFFFF", eName))

	return sb.String()
}
