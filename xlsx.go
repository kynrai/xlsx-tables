package xlsx_tables

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
)

type Reader struct {
	Worksheet     string
	sharedStrings []string
	zr            *zip.ReadCloser
	dec           *xml.Decoder
	numLine       int
}

func NewReader(zr *zip.ReadCloser) *Reader {
	return &Reader{
		zr:        zr,
		Worksheet: "sheet1",
	}
}

func (r *Reader) Read() ([]string, error) {
	if r.sharedStrings == nil {
		if err := r.loadStrings(); err != nil {
			return nil, err
		}
	}
	if r.dec == nil {
		err := r.openWorksheet()
		if err != nil {
			return nil, err
		}
	}

	// find the first row from wherever we are in the buffer
	for {
		token, err := r.dec.Token()
		if err != nil {
			return nil, err
		}
		switch se := token.(type) {
		case xml.StartElement:
			ele := se.Name.Local
			if ele != "row" {
				continue
			}
			var row Row
			err = r.dec.DecodeElement(&row, &se)
			if err != nil {
				return nil, err
			}
			rec := []string{}
			for _, v := range row.C {
				if v.T == "n" {
					rec = append(rec, v.V)
					continue
				}
				i, err := strconv.Atoi(v.V)
				if err != nil {
					return nil, err
				}
				rec = append(rec, r.sharedStrings[i])
			}
			return rec, nil
		default:
		}
	}

	return nil, nil
}

func (r *Reader) loadStrings() (err error) {
	var ssrc io.ReadCloser
	defer func() {
		if ssrc != nil {
			ssrc.Close()
		}
	}()

	for _, v := range r.zr.File {
		if v.Name != "xl/sharedStrings.xml" {
			continue
		}
		ssrc, err = v.Open()
		if err != nil {
			return err
		}
		break
	}

	dec := xml.NewDecoder(ssrc)
	for {
		token, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		// Inspect the type of the token just read.
		switch se := token.(type) {
		case xml.StartElement:
			ele := se.Name.Local
			if ele != "si" {
				continue
			}
			var si StringItem
			if err := dec.DecodeElement(&si, &se); err != nil {
				return err
			}
			r.sharedStrings = append(r.sharedStrings, si.T)
		default:
		}
	}
	return nil
}

func (r *Reader) openWorksheet() error {
	for _, v := range r.zr.File {
		if v.Name != fmt.Sprintf("xl/worksheets/%s.xml", r.Worksheet) {
			continue
		}
		rc, err := v.Open()
		if err != nil {
			return err
		}
		r.dec = xml.NewDecoder(rc)
		break
	}
	if r.dec == nil {
		return fmt.Errorf("no worksheet with the name %q found", r.Worksheet)
	}
	return nil
}
