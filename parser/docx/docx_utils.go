package docx

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
)

const maxBytes = 20 << 20

type typeOverride struct {
	XMLName     xml.Name `xml:"Override"`
	ContentType string   `xml:"ContentType,attr"`
	PartName    string   `xml:"PartName,attr"`
}

type contentTypeDefinition struct {
	XMLName   xml.Name       `xml:"Types"`
	Overrides []typeOverride `xml:"Override"`
}

func mapZipFiles(files []*zip.File) map[string]*zip.File {
	filesMap := make(map[string]*zip.File, len(files))
	for _, f := range files {
		filesMap[f.Name] = f
		filesMap["/"+f.Name] = f
	}
	return filesMap
}

func getContentTypeDefinition(file *zip.File) (*contentTypeDefinition, error) {
	f, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	contentDefinition := &contentTypeDefinition{}
	if err := xml.NewDecoder(f).Decode(contentDefinition); err != nil {
		return nil, err
	}
	return contentDefinition, nil
}

func xmlToMap(r io.Reader) (map[string]string, error) {
	m := make(map[string]string)
	dec := xml.NewDecoder(io.LimitReader(r, maxBytes))
	var tagName string
	for {
		t, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		switch v := t.(type) {
		case xml.StartElement:
			tagName = string(v.Name.Local)
		case xml.CharData:
			m[tagName] = string(v)
		}
	}
	return m, nil
}

func docxXMLToText(reader io.Reader) (string, error) {
	return xmlToText(reader, []string{"br", "p", "tab"}, []string{"instrText", "script"}, true)
}

func xmlToText(r io.Reader, breaks []string, skip []string, strict bool) (string, error) {
	var result string

	decoder := xml.NewDecoder(io.LimitReader(r, maxBytes))
	decoder.Strict = strict
	for {
		t, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		switch v := t.(type) {
		case xml.CharData:
			result += string(v)
		case xml.StartElement:
			for _, breakElement := range breaks {
				if v.Name.Local == breakElement {
					result += "\n"
				}
			}
			for _, skipElement := range skip {
				if v.Name.Local == skipElement {
					depth := 1
					for {
						t, err := decoder.Token()
						if err != nil {
							return "", err
						}

						switch t.(type) {
						case xml.StartElement:
							depth++
						case xml.EndElement:
							depth--
						}

						if depth == 0 {
							break
						}
					}
				}
			}
		}
	}
	return result, nil
}

func parseDocxText(file *zip.File) (string, error) {
	reader, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("error opening '%v' from archive: %v", file.Name, err)
	}
	defer reader.Close()

	text, err := docxXMLToText(reader)
	if err != nil {
		return "", fmt.Errorf("error parsing '%v': %v", file.Name, err)
	}
	return text, nil
}
