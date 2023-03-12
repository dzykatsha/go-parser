package main

import (
	"flag"
	"github.com/rs/zerolog/log"
	"go-parser/parser"
	djvu_parser "go-parser/parser/djvu"
	doc_parser "go-parser/parser/doc"
	docx_parser "go-parser/parser/docx"
	html_parser "go-parser/parser/html"
	pdf_parser "go-parser/parser/pdf"
	"regexp"
)

func main() {
	var (
		inputFile  = flag.String("i", "", "path to input file")
		outputFile = flag.String("o", "output.txt", "path to output file")
		plain      = flag.Bool("plain", false, "get plain text from PDF file (use only one of plain, styled or grouped)")
		styled     = flag.Bool("styled", false, "get all text with styles from PDF file (use only one of plain, styled or grouped)")
		grouped    = flag.Bool("grouped", false, "get text grouped by rows from PDF file (use only one of plain, styled or grouped)")
	)
	flag.Parse()

	html, _ := regexp.MatchString(`.*\.html$`, *inputFile)
	doc, _ := regexp.MatchString(`.*\.doc$`, *inputFile)
	docx, _ := regexp.MatchString(`.*\.docx$`, *inputFile)
	pdf, _ := regexp.MatchString(`.*\.pdf$`, *inputFile)
	djvu, _ := regexp.MatchString(`.*\.djvu$`, *inputFile)

	if *inputFile == "" {
		log.Fatal().Msg("no input file")
	}

	switch {
	case html:
		text, err := html_parser.ReadHtmlFile(*inputFile)
		if err != nil {
			log.Fatal().Msgf("error from html parser: %s", err.Error())
		}
		content := html_parser.ParseHtml(text)
		err = parser.WriteContent(content, *outputFile)
		if err != nil {
			log.Fatal().Msgf("error from html parser: %s", err.Error())
		}
	case doc:
		err := doc_parser.ParseAndWriteDoc(*inputFile, *outputFile)
		if err != nil {
			log.Fatal().Msgf("error from doc parser: %s", err.Error())
		}
	case docx:
		reader, err := docx_parser.ReadDocxFile(*inputFile)
		if err != nil {
			log.Fatal().Msgf("error from docx parser: %s", err.Error())
		}
		content, metadata, err := docx_parser.ParseDocx(reader)
		if err != nil {
			log.Fatal().Msgf("error from docx parser: %s", err.Error())
		}
		log.Info().Msgf("Metadata: %v", metadata)
		err = parser.WriteContent(content, *outputFile)
		if err != nil {
			log.Fatal().Msgf("error from html parser: %s", err.Error())
		}
	case pdf:
		err := pdf_parser.ValidatePdfArgs(*plain, *styled, *grouped)
		if err != nil {
			log.Fatal().Msgf("error from pdf parser: %s", err.Error())
		}

		var content string
		switch {
		case *plain:
			content, err = pdf_parser.ReadPlainPdf(*inputFile)
			if err != nil {
				log.Fatal().Msgf("error from pdf parser: %s", err.Error())
			}
		case *styled:
			content, err = pdf_parser.ReadStyledPdf(*inputFile)
			if err != nil {
				log.Fatal().Msgf("error from pdf parser: %s", err.Error())
			}
		case *grouped:
			content, err = pdf_parser.ReadGroupedPdf(*inputFile)
			if err != nil {
				log.Fatal().Msgf("error from pdf parser: %s", err.Error())
			}
		default:
			log.Fatal().Msg("Unknown type of pdf parameter")
		}

		err = parser.WriteContent(content, *outputFile)
		if err != nil {
			log.Fatal().Msgf("error from pdf parser: %s", err.Error())
		}
	case djvu:
		content, err := djvu_parser.ParseDjvu(*inputFile)
		if err != nil {
			log.Fatal().Msgf("error from djvu parser: %s", err.Error())
		}
		err = parser.WriteContent(content, *outputFile)
		if err != nil {
			log.Fatal().Msgf("error from html parser: %s", err.Error())
		}
	default:
		log.Fatal().Msg("Unknown type of file")
	}
}
