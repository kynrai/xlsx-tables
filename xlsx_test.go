package xlsx_tables

import (
	"archive/zip"
	"io"
	"strings"
	"testing"
)

func TestReader(t *testing.T) {
	expected := []string{
		"name,country",
		"tellus s1,uk",
		"tellus s2,us",
		"tellus s3,us",
	}
	f, err := zip.OpenReader("testdata/xlsx_sample.xlsx")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	r := NewReader(f)

	res := []string{}
	for {
		row, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Fatal(err)
		}
		res = append(res, strings.Join(row, ","))
	}

	for k, v := range expected {
		if res[k] != v {
			t.Fatalf("unexpected results, got: %q, want: %q", res[k], v)
		}
	}
}

func TestReaderWorksheet2(t *testing.T) {
	expected := []string{
		"code,product",
		"1,oil",
		"2,grease",
		"3.1,tellus",
	}
	f, err := zip.OpenReader("testdata/xlsx_sample.xlsx")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	r := NewReader(f)
	r.Worksheet = "sheet2"

	res := []string{}
	for {
		row, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Fatal(err)
		}
		res = append(res, strings.Join(row, ","))
	}

	for k, v := range expected {
		if res[k] != v {
			t.Fatalf("unexpected results, got: %q, want: %q", res[k], v)
		}
	}
}

func BenchmarkRead(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		f, err := zip.OpenReader("testdata/xlsx_sample.xlsx")
		if err != nil {
			b.Fatal(err)
		}
		defer f.Close()
		r := NewReader(f)

		for {
			_, err := r.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				b.Fatal(err)
			}
		}
	}
}
