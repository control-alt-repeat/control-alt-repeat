package reconcile

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"path"
	"strconv"
	"strings"
)

// readXLSXRows returns the first worksheet of an .xlsx file as rows of strings, so
// Excel exports can be used directly without converting them to CSV first.
func readXLSXRows(filename string) ([][]string, error) {
	zr, err := zip.OpenReader(filename)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", filename, err)
	}
	defer zr.Close()

	files := map[string]*zip.File{}
	for _, f := range zr.File {
		files[f.Name] = f
	}

	shared, err := readSharedStrings(files["xl/sharedStrings.xml"])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", filename, err)
	}

	sheetPath, err := firstSheetPath(files)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", filename, err)
	}
	sheet, ok := files[sheetPath]
	if !ok {
		return nil, fmt.Errorf("%s: worksheet %s not found", filename, sheetPath)
	}

	var ws struct {
		Rows []struct {
			Cells []struct {
				Ref    string `xml:"r,attr"`
				Type   string `xml:"t,attr"`
				Value  string `xml:"v"`
				Inline string `xml:"is>t"`
			} `xml:"c"`
		} `xml:"sheetData>row"`
	}
	if err := decodeXML(sheet, &ws); err != nil {
		return nil, fmt.Errorf("%s: %w", filename, err)
	}

	rows := make([][]string, 0, len(ws.Rows))
	for _, r := range ws.Rows {
		var row []string
		for i, c := range r.Cells {
			col := i
			if c.Ref != "" {
				col = columnIndex(c.Ref)
			}
			for len(row) <= col {
				row = append(row, "")
			}
			switch c.Type {
			case "s":
				idx, err := strconv.Atoi(c.Value)
				if err != nil || idx >= len(shared) {
					return nil, fmt.Errorf("%s: bad shared string index %q in %s", filename, c.Value, c.Ref)
				}
				row[col] = shared[idx]
			case "inlineStr":
				row[col] = c.Inline
			default:
				row[col] = c.Value
			}
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func readSharedStrings(f *zip.File) ([]string, error) {
	if f == nil {
		return nil, nil
	}
	var sst struct {
		Items []struct {
			Text string   `xml:"t"`
			Runs []string `xml:"r>t"`
		} `xml:"si"`
	}
	if err := decodeXML(f, &sst); err != nil {
		return nil, err
	}
	out := make([]string, len(sst.Items))
	for i, si := range sst.Items {
		out[i] = si.Text + strings.Join(si.Runs, "")
	}
	return out, nil
}

// firstSheetPath follows workbook.xml -> its relationships to find the first sheet,
// since sheets are not always stored as sheet1.xml.
func firstSheetPath(files map[string]*zip.File) (string, error) {
	var wb struct {
		Sheets []struct {
			RelID string `xml:"http://schemas.openxmlformats.org/officeDocument/2006/relationships id,attr"`
		} `xml:"sheets>sheet"`
	}
	// Fall back to the conventional location when the workbook metadata is unusual.
	if err := decodeXML(files["xl/workbook.xml"], &wb); err != nil || len(wb.Sheets) == 0 {
		return "xl/worksheets/sheet1.xml", nil
	}
	var rels struct {
		Rels []struct {
			ID     string `xml:"Id,attr"`
			Target string `xml:"Target,attr"`
		} `xml:"Relationship"`
	}
	if err := decodeXML(files["xl/_rels/workbook.xml.rels"], &rels); err != nil {
		return "xl/worksheets/sheet1.xml", nil
	}
	for _, r := range rels.Rels {
		if r.ID == wb.Sheets[0].RelID {
			if strings.HasPrefix(r.Target, "/") {
				return strings.TrimPrefix(r.Target, "/"), nil
			}
			return path.Join("xl", r.Target), nil
		}
	}
	return "", fmt.Errorf("first worksheet not found in workbook relationships")
}

func decodeXML(f *zip.File, v any) error {
	if f == nil {
		return fmt.Errorf("missing part")
	}
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()
	return xml.NewDecoder(io.LimitReader(rc, 256<<20)).Decode(v)
}

// columnIndex turns a cell reference like "AB12" into a zero-based column index.
func columnIndex(ref string) int {
	n := 0
	for _, ch := range ref {
		if ch < 'A' || ch > 'Z' {
			break
		}
		n = n*26 + int(ch-'A'+1)
	}
	return n - 1
}
