package xlsx_tables

type StringItem struct {
	T string `xml:"t"`
}

type Row struct {
	C []struct {
		T string `xml:"t,attr"`
		V int    `xml:"v"`
	} `xml:"c"`
}
