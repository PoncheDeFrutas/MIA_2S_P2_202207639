package structures

import (
	"backend/utils"
	"fmt"
	"strings"
	"time"
)

type Inode struct {
	IuId   int32
	IGid   int32
	ISize  int32
	IAtime float32
	ICTime float32
	IMTime float32
	IBlock [15]int32
	IType  byte
	IPerm  [3]byte
}

func (i *Inode) DefaultValue(blockCount int32) {
	i.IuId = 1
	i.IGid = 1
	i.ISize = 0
	i.IAtime = float32(time.Now().Unix())
	i.ICTime = float32(time.Now().Unix())
	i.IMTime = float32(time.Now().Unix())
	i.IBlock[0] = blockCount
	for j := 1; j < 15; j++ {
		i.IBlock[j] = -1
	}
	i.IType = '0'
	for j := 0; j < 3; j++ {
		i.IPerm[j] = '7'
	}
}

func (i *Inode) WriteInode(path string, offset int64, maxSize int64) error {
	if err := utils.WriteToFile(path, offset, maxSize, i); err != nil {
		return err
	}
	return nil
}

func (i *Inode) ReadInode(path string, offset int64) error {
	if err := utils.ReadFromFile(path, offset, i); err != nil {
		return err
	}
	return nil
}

func (i *Inode) Print() {
	fmt.Printf("IuId: %d\n", i.IuId)
	fmt.Printf("IGid: %d\n", i.IGid)
	fmt.Printf("ISize: %d\n", i.ISize)
	fmt.Printf("IAtime: %s\n", time.Unix(int64(i.IAtime), 0))
	fmt.Printf("ICTime: %s\n", time.Unix(int64(i.ICTime), 0))
	fmt.Printf("IMTime: %s\n", time.Unix(int64(i.IMTime), 0))
	fmt.Printf("IBlock: %v\n", i.IBlock)
	fmt.Printf("IType: %c\n", i.IType)
	fmt.Printf("IPerm: %s\n", string(i.IPerm[:]))
}

func (i *Inode) GetStringBuilder(nodeName string) string {
	var sb strings.Builder

	Perm := strings.TrimRight(string(i.IPerm[:]), "\x00")

	sb.WriteString(fmt.Sprintf("    %s [label=<\n", nodeName))
	sb.WriteString(fmt.Sprintf("    <TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\">\n"))

	sb.WriteString(fmt.Sprintf("\t<TR><TD COLSPAN=\"2\" BGCOLOR=\"%s\"><B>%s</B></TD></TR>\n", "#333333", nodeName))

	sb.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">IuId</TD><TD WIDTH=\"250\" BGCOLOR=\"%s\">%d</TD></TR>\n", "#DDDDDD", "#DDDDDD", i.IuId))
	sb.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">IGid</TD><TD WIDTH=\"250\" BGCOLOR=\"%s\">%d</TD></TR>\n", "#FFFFFF", "#FFFFFF", i.IGid))
	sb.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">ISize</TD><TD WIDTH=\"250\" BGCOLOR=\"%s\">%d</TD></TR>\n", "#DDDDDD", "#DDDDDD", i.ISize))
	sb.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">IAtime</TD><TD WIDTH=\"250\" BGCOLOR=\"%s\">%s</TD></TR>\n", "#FFFFFF", "#FFFFFF", time.Unix(int64(i.IAtime), 0).Format("02-Jan-2006 03:04 PM")))
	sb.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">ICTime</TD><TD WIDTH=\"250\" BGCOLOR=\"%s\">%s</TD></TR>\n", "#DDDDDD", "#DDDDDD", time.Unix(int64(i.ICTime), 0).Format("02-Jan-2006 03:04 PM")))
	sb.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">IMTime</TD><TD WIDTH=\"250\" BGCOLOR=\"%s\">%s</TD></TR>\n", "#FFFFFF", "#FFFFFF", time.Unix(int64(i.IMTime), 0).Format("02-Jan-2006 03:04 PM")))
	sb.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">IType</TD><TD WIDTH=\"250\" BGCOLOR=\"%s\">%c</TD></TR>\n", "#DDDDDD", "#DDDDDD", i.IType))
	sb.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">IPerm</TD><TD WIDTH=\"250\" BGCOLOR=\"%s\">%s</TD></TR>\n", "#FFFFFF", "#FFFFFF", Perm))

	sb.WriteString(fmt.Sprintf("\t<TR><TD COLSPAN=\"2\" BGCOLOR=\"%s\"><B>Direct Blocks</B></TD></TR>\n", "#333333"))
	for j := 0; j < 12; j++ {
		bgColor := "#DDDDDD"
		if j%2 == 0 {
			bgColor = "#FFFFFF"
		}
		sb.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">IBlock %d</TD><TD WIDTH=\"250\" BGCOLOR=\"%s\">%d</TD></TR>\n", bgColor, j, bgColor, i.IBlock[j]))
	}

	sb.WriteString(fmt.Sprintf("\t<TR><TD COLSPAN=\"2\" BGCOLOR=\"%s\"><B>Indirect Blocks</B></TD></TR>\n", "#333333"))

	for j := 12; j < 15; j++ {
		bgColor := "#DDDDDD"
		if j%2 == 0 {
			bgColor = "#FFFFFF"
		}
		sb.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">IBlock %d</TD><TD WIDTH=\"250\" BGCOLOR=\"%s\">%d</TD></TR>\n", bgColor, j, bgColor, i.IBlock[j]))
	}

	sb.WriteString("    </TABLE>>];\n")

	return sb.String()
}
