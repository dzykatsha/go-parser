package docx

import (
	"archive/zip"
	"fmt"
	"time"
)

func ReadDocxFile(path string) (*zip.ReadCloser, error) {
	reader, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	return reader, nil
}

func ParseDocx(reader *zip.ReadCloser) (string,  map[string]string, error) {
	zipFiles := mapZipFiles(reader.File)

	contentTypeDefinition, err := getContentTypeDefinition(zipFiles["[Content_Types].xml"])
	if err != nil {
		return "", nil, err
	}

	metadata := make(map[string]string)
	var textHeader, textBody, textFooter string

	for _, override := range contentTypeDefinition.Overrides {
		file := zipFiles[override.PartName]
		switch {
		case override.ContentType == "application/vnd.openxmlformats-package.core-properties+xml":
			readCloser, err := file.Open()
			if err != nil {
				return "", nil, fmt.Errorf("error opening '%v': %v", file.Name, err)
			}
			defer readCloser.Close()

			metadata, err = xmlToMap(readCloser)
			if err != nil {
				return "", nil, fmt.Errorf("error parsing '%v': %v", file.Name, err)
			}

			if modified, ok := metadata["modified"]; ok {
				if timeModificaion, err := time.Parse(time.RFC3339, modified); err == nil {
					metadata["ModifiedDate"] = fmt.Sprintf("%d", timeModificaion.Unix())
				}
			}
			if created, ok := metadata["created"]; ok {
				if timeCreation, err := time.Parse(time.RFC3339, created); err == nil {
					metadata["CreatedDate"] = fmt.Sprintf("%d", timeCreation.Unix())
				}
			}
		case override.ContentType == "application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml":
			body, err := parseDocxText(file)
			if err != nil {
				return "", nil, err
			}
			textBody += body + "\n"
		case override.ContentType == "application/vnd.openxmlformats-officedocument.wordprocessingml.footer+xml":
			footer, err := parseDocxText(file)
			if err != nil {
				return "", nil, err
			}
			textFooter += footer + "\n"
		case override.ContentType == "application/vnd.openxmlformats-officedocument.wordprocessingml.header+xml":
			header, err := parseDocxText(file)
			if err != nil {
				return "", nil, err
			}
			textHeader += header + "\n"
		}

	}
	return textHeader + "\n" + textBody + "\n" + textFooter, metadata, nil
}
