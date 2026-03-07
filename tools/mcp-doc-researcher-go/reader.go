package main

import (
	"fmt"
	"io"
	"strings"

	"code.sajari.com/docconv"
	"github.com/ledongthuc/pdf"
)

// ReadLocalDoc Extracts text from a local .doc, .docx, or .pdf file.
func ReadLocalDoc(filePath string) (string, error) {
	if strings.HasSuffix(strings.ToLower(filePath), ".pdf") {
		return readPdf(filePath)
	}

	res, err := docconv.ConvertPath(filePath)
	if err != nil {
		return "", fmt.Errorf("error reading %s: %v", filePath, err)
	}
	return strings.TrimSpace(res.Body), nil
}

func readPdf(path string) (string, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return "", fmt.Errorf("error opening pdf %s: %v", path, err)
	}
	defer f.Close()

	b, err := r.GetPlainText()
	if err != nil {
		return "", fmt.Errorf("error getting text from pdf %s: %v", path, err)
	}

	bytes, err := io.ReadAll(b)
	if err != nil {
		return "", fmt.Errorf("error reading text from pdf %s: %v", path, err)
	}
	return strings.TrimSpace(string(bytes)), nil
}
