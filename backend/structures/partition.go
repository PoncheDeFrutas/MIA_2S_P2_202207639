package structures

import (
	"encoding/binary"
	"fmt"
	"math"
	"strings"
)

type Partition struct {
	PartStatus      byte
	PartType        byte
	PartFit         byte
	PartStart       int32
	PartSize        int32
	PartName        [16]byte
	PartCorrelative int32
	PartId          [4]byte
	// Total size of the Partition is 35 bytes
}

func (p *Partition) DefaultValue() {
	p.PartStatus = '9'
	p.PartType = 'P'
	p.PartFit = 'W'
	p.PartStart = -1
	p.PartSize = -1
	copy(p.PartName[:], "$")
	p.PartCorrelative = -1
	copy(p.PartId[:], "$$$$")
}

func (p *Partition) SetPartition(partType string, fit string, start int32, size int32, name string) {
	p.PartStatus = '0'
	p.PartType = partType[0]
	p.PartFit = fit[0]
	p.PartStart = start
	p.PartSize = size
	copy(p.PartName[:], name)
}

func (p *Partition) IsEmpty() bool {
	return p.PartStart == -1 && p.PartSize == -1
}

func (p *Partition) MountPartition(correlative int, id string) error {
	p.PartCorrelative = int32(correlative) + 1
	copy(p.PartId[:], id)
	p.PartStatus = '1'
	return nil
}

func (p *Partition) UnmountPartition() {
	p.PartCorrelative = -1
	copy(p.PartId[:], "$$$$")
}

func (p *Partition) IsMounted() bool {
	return p.PartCorrelative != -1
}

func (p *Partition) CalculateN() int32 {
	numerator := int(p.PartSize) - binary.Size(SuperBlock{})
	denominator := 4 + binary.Size(Inode{}) + 3*binary.Size(FileBlock{})
	return int32(math.Floor(float64(numerator) / float64(denominator)))
}

func (p *Partition) Print() {
	fmt.Println("PartStatus: ", string(p.PartStatus))
	fmt.Println("PartType: ", string(p.PartType))
	fmt.Println("PartFit: ", string(p.PartFit))
	fmt.Println("PartStart: ", p.PartStart)
	fmt.Println("PartSize: ", p.PartSize)
	fmt.Println("PartName: ", string(p.PartName[:]))
	fmt.Println("PartCorrelative: ", p.PartCorrelative)
	fmt.Println("PartId: ", string(p.PartId[:]))
}

func (p *Partition) GetStringBuilder(i int) string {
	var sb strings.Builder

	partName := strings.TrimRight(string(p.PartName[:]), "\x00")
	id := strings.TrimRight(string(p.PartId[:]), "\x00")

	sb.WriteString(fmt.Sprintf("<TR><TD COLSPAN=\"2\" BGCOLOR=\"%s\"><B>Partición %d</B></TD></TR>\n", "#AAAAAA", i+1))
	sb.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">Status</TD><TD WIDTH=\"200\" BGCOLOR=\"%s\">%c</TD></TR>\n", "#DDDDDD", "#DDDDDD", p.PartStatus))
	sb.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">Type</TD><TD WIDTH=\"200\" BGCOLOR=\"%s\">%c</TD></TR>\n", "#FFFFFF", "#FFFFFF", p.PartType))
	sb.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">Fit</TD><TD WIDTH=\"200\" BGCOLOR=\"%s\">%c</TD></TR>\n", "#DDDDDD", "#DDDDDD", p.PartFit))
	sb.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">Start</TD><TD WIDTH=\"200\" BGCOLOR=\"%s\">%d</TD></TR>\n", "#FFFFFF", "#FFFFFF", p.PartStart))
	sb.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">Size</TD><TD WIDTH=\"200\" BGCOLOR=\"%s\">%d</TD></TR>\n", "#DDDDDD", "#DDDDDD", p.PartSize))
	sb.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">Name</TD><TD WIDTH=\"200\" BGCOLOR=\"%s\">%s</TD></TR>\n", "#FFFFFF", "#FFFFFF", partName))
	sb.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">ID</TD><TD WIDTH=\"200\" BGCOLOR=\"%s\">%s</TD></TR>\n", "#DDDDDD", "#DDDDDD", id))
	sb.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">Correlative</TD><TD WIDTH=\"200\" BGCOLOR=\"%s\">%d</TD></TR>\n", "#FFFFFF", "#FFFFFF", p.PartCorrelative))

	return sb.String()
}
