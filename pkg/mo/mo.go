package mo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"

	"github.com/execjosh/go-gettext/pkg/domain"
)

const (
	magic          = uint32(0x950412de)
	magicBigEndian = uint32(0xde120495)
)

type header struct {
	Magic       uint32
	Revision    uint32
	NumStrings  uint32
	OrigTblOff  uint32
	TranTblOff  uint32
	HashTblSize uint32
	HashTblOff  uint32
}

type pos struct {
	Length uint32 // length of string, excluding trailing \0
	Offset uint32 // offset of string in file
}

func Load(filename string) (*domain.Domain, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	f := bytes.NewReader(data)

	var byteOrder binary.ByteOrder = binary.LittleEndian
	var h header
	if err := binary.Read(f, binary.LittleEndian, &h); err != nil {
		return nil, err
	}

	if magic != h.Magic {
		if magicBigEndian != h.Magic {
			return nil, fmt.Errorf("not a GNU MO file!?")
		}
		byteOrder = binary.BigEndian
	}

	f.Seek(int64(h.OrigTblOff), 0)

	// load TOC of original strings
	origTOC := make([]pos, h.NumStrings)
	if err := binary.Read(f, byteOrder, &origTOC); err != nil {
		return nil, err
	}

	// load TOC of translated strings
	tranTOC := make([]pos, h.NumStrings)
	if err := binary.Read(f, byteOrder, &tranTOC); err != nil {
		return nil, err
	}

	// load strings into domain
	domain := domain.Domain{}
	for i := uint32(0); i < h.NumStrings; i++ {
		domain[stringAt(data, origTOC[i])] = stringAt(data, tranTOC[i])
	}

	return &domain, nil
}

func stringAt(data []byte, p pos) string {
	beg := p.Offset
	end := p.Offset + p.Length
	return string(data[beg:end])
}
