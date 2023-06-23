package writer

import (
	"bytes"
	"embed"
	"encoding/csv"
	"fmt"
	"github.com/elulcao/go-exif-extract/pkg/exif"
	"html/template"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
)

//go:embed templates/*
var assets embed.FS

// Write writes the exif data to a csv file
func Write(ed *exif.ExifData, cmd *cobra.Command) error {
	csvOutput, _ := cmd.Flags().GetString("output")
	htmlOutput, _ := cmd.Flags().GetBool("html")

	data := [][]string{
		{"File", "Latitude", "Longitude"},
	}

	for _, v := range ed.Exif {
		latitude := strconv.FormatFloat(float64(v.Latitude), 'f', -1, 64)
		longitude := strconv.FormatFloat(float64(v.Longitude), 'f', -1, 64)
		data = append(data, []string{v.FilePath, latitude, longitude})
	}

	file, err := os.Create(csvOutput)
	if err != nil {
		return fmt.Errorf("unable to create csv file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.WriteAll(data)
	if err != nil {
		return fmt.Errorf("unable to write csv file: %w", err)
	}

	if htmlOutput {
		tmpl, err := template.ParseFS(assets, "templates/template.html")
		if err != nil {
			return fmt.Errorf("unable to parse template: %w", err)
		}

		tpl := bytes.Buffer{}
		err = tmpl.Execute(&tpl, ed)
		if err != nil {
			return fmt.Errorf("unable to execute template: %w", err)
		}

		// get path from csvOutput
		path := filepath.Dir(csvOutput)
		htmlFile := path + "/output.html"

		file, err := os.Create(htmlFile)
		if err != nil {
			return fmt.Errorf("unable to create html file: %w", err)
		}

		_, err = file.Write(tpl.Bytes())
		if err != nil {
			return fmt.Errorf("unable to write html file: %w", err)
		}
	}

	return nil
}
