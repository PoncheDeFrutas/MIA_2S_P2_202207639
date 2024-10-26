package structures

import (
	common "backend/utils"
	"encoding/binary"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type MBR struct {
	MbrSize          int32
	MbrCreationDate  float32
	MbrDiskSignature int32
	MbrDiskFit       byte
	MbrPartition     [4]Partition
	// Total size of the MBR is 153 bytes
}

func (m *MBR) CreateMBR(size int, fit string) error {
	m.MbrSize = int32(size)
	m.MbrCreationDate = float32(time.Now().Unix())
	m.MbrDiskSignature = rand.Int31()
	m.MbrDiskFit = fit[0]

	for i := range m.MbrPartition {
		m.MbrPartition[i].DefaultValue()
	}

	return nil
}

func (m *MBR) WriteMBR(path string) error {
	if err := common.WriteToFile(path, int64(0), int64(binary.Size(m)), m); err != nil {
		return err
	}
	return nil
}

func (m *MBR) ReadMBR(path string) error {
	if err := common.ReadFromFile(path, int64(0), m); err != nil {
		return err
	}
	return nil
}

func (m *MBR) FindFreePartition() int {
	for i, partition := range m.MbrPartition {
		if partition.PartStart == -1 {
			return i
		}
	}
	return -1
}

func (m *MBR) FreeNamePartition(name string) bool {
	for _, partition := range m.MbrPartition {
		if strings.TrimRight(string(partition.PartName[:]), "\x00") == name {
			return false
		}
	}
	return true
}

func (m *MBR) GetPartitionByName(name string) (*Partition, int) {
	for i, partition := range m.MbrPartition {
		if strings.TrimRight(string(partition.PartName[:]), "\x00") == name {
			return &m.MbrPartition[i], i
		}
	}
	return nil, -1
}

func (m *MBR) GetPartitionByID(id string) (*Partition, error) {
	for i, partition := range m.MbrPartition {
		if strings.TrimRight(string(partition.PartId[:]), "\x00") == id {
			return &m.MbrPartition[i], nil
		}
	}
	return nil, fmt.Errorf("partition not found")
}

func (m *MBR) ExtendPartitionExist() bool {
	for _, partition := range m.MbrPartition {
		if partition.PartType == 'E' {
			return true
		}
	}
	return false
}

func (m *MBR) GetExtendedPartition() *Partition {
	for i, partition := range m.MbrPartition {
		if partition.PartType == 'E' {
			return &m.MbrPartition[i]
		}
	}
	return nil
}

func (m *MBR) Print() {
	fmt.Println("/*********************** MBR ***********************/")
	fmt.Printf("Size: %d\n", m.MbrSize)
	fmt.Printf("Creation Date: %s\n", time.Unix(int64(m.MbrCreationDate), 0))
	fmt.Printf("Disk Signature: %d\n", m.MbrDiskSignature)
	fmt.Printf("Disk Fit: %c\n", m.MbrDiskFit)
	fmt.Println("/******************** Partitions ********************/")
	for i, partition := range m.MbrPartition {
		fmt.Println("--------------------------------------------------")
		fmt.Printf("Partition %d\n", i)
		partition.Print()
	}
	fmt.Println("--------------------------------------------------")
}

func (m *MBR) GetStringBuilder() string {
	var sb strings.Builder
	// Main title row
	sb.WriteString(fmt.Sprintf("\t<TR><TD COLSPAN=\"2\" BGCOLOR=\"%s\"><B>MBR</B></TD></TR>\n", "#333333"))

	// MBR data rows with uniform cell colors
	sb.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">MBR Tamaño</TD><TD WIDTH=\"200\" BGCOLOR=\"%s\">%d</TD></TR>\n", "#DDDDDD", "#DDDDDD", m.MbrSize))
	sb.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">MBR Fecha Creación</TD><TD WIDTH=\"200\" BGCOLOR=\"%s\">%s</TD></TR>\n", "#DDDDDD", "#DDDDDD", time.Unix(int64(m.MbrCreationDate), 0).Format("02-Jan-2006 03:04 PM")))
	sb.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">MBR Disk Signature</TD><TD WIDTH=\"200\" BGCOLOR=\"%s\">%d</TD></TR>\n", "#DDDDDD", "#DDDDDD", m.MbrDiskSignature))
	sb.WriteString(fmt.Sprintf("<TR><TD WIDTH=\"150\" BGCOLOR=\"%s\">MBR Disk Fit</TD><TD WIDTH=\"200\" BGCOLOR=\"%s\">%c</TD></TR>\n", "#DDDDDD", "#DDDDDD", m.MbrDiskFit))

	return sb.String()
}
