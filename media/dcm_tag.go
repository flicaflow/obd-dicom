package media

import (
	"bytes"
	"encoding/binary"
	"strconv"
	"strings"
)

// DcmTag DICOM tag structure
type DcmTag struct {
	Name        string
	Description string
	Group       uint16
	Element     uint16
	Length      uint32
	VR          string
	VM          string
	Data        []byte
	BigEndian   bool
}

// GetUShort convert tag.Data to uint16
func (tag *DcmTag) GetUShort() uint16 {
	if tag.Length == 2 {
		if tag.BigEndian {
			return binary.BigEndian.Uint16(tag.Data)
		}
		return binary.LittleEndian.Uint16(tag.Data)
	}
	return 0
}

// GetUInt convert tag.Data to uint32
func (tag *DcmTag) GetUInt() uint32 {
	var val uint32
	if tag.Length == 4 {
		if tag.BigEndian {
			val = binary.BigEndian.Uint32(tag.Data)
		} else {
			val = binary.LittleEndian.Uint32(tag.Data)
		}
	}
	return val
}

// GetString convert tag.Data to string
func (tag *DcmTag) GetString() string {
	n := bytes.IndexByte(tag.Data, 0)
	if n == -1 {
		n = int(tag.Length)
	}
	return strings.TrimSpace(string(tag.Data[:n]))
}

// GetFloat convert tag.Data to float32
func (tag *DcmTag) GetFloat() float32 {
	val := tag.GetString()
	if s, err := strconv.ParseFloat(val, 32); err == nil {
		return float32(s)
	}
	return 0.0
}

// WriteSeq - Create an SQ tag from a DICOM Object
func (tag *DcmTag) WriteSeq(group uint16, element uint16, seq DcmObj) {
	bufdata := &bufData{
		BigEndian: false,
		MS:        NewEmptyMemoryStream(),
	}

	bufdata.BigEndian = seq.IsBigEndian()
	tag.BigEndian = seq.IsBigEndian()
	tag.Group = group
	tag.Element = element
	if tag.Group == 0xFFFE {
		tag.VR = ""
	} else {
		tag.VR = "SQ"
	}
	for i := 0; i < seq.TagCount(); i++ {
		temptag := seq.GetTagAt(i)
		bufdata.WriteTag(temptag, seq.IsExplicitVR())
	}
	tag.Length = uint32(bufdata.GetSize())
	if tag.Length%2 == 1 {
		tag.Length++
		bufdata.MS.Write([]byte{0x00}, 1)
	}
	if tag.Length > 0 {
		bufdata.SetPosition(0)
		data, _ := bufdata.MS.Read(int(tag.Length))
		tag.Data = data
	}
}

// ReadSeq - reads a dicom sequence
func (tag *DcmTag) ReadSeq(ExplicitVR bool) DcmObj {
	seq := NewEmptyDCMObj()
	bufdata := &bufData{
		BigEndian: false,
		MS:        NewEmptyMemoryStream(),
	}

	bufdata.Write(tag.Data, int(tag.Length))
	bufdata.MS.SetPosition(0)

	for bufdata.MS.GetPosition() < bufdata.MS.GetSize() {
		temptag, err := bufdata.ReadTag(ExplicitVR)
		if err != nil {
			continue
		}

		if !ExplicitVR {
			temptag.VR = GetDictionaryVR(tag.Group, tag.Element)
		}
		seq.Add(temptag)
	}
	return seq
}
