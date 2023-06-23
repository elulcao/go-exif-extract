package cmd

import (
	"github.com/elulcao/go-exif-extract/pkg/exif"
	"github.com/elulcao/go-exif-extract/pkg/extract"
	"github.com/elulcao/go-exif-extract/pkg/writer"

	"github.com/spf13/cobra"
)

// extractCmd represents the extract subcommand
var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extracts the exif data from images in a directory",
	Long: `Extracts the exif data from images in a directory and writes it to a csv file. The info
exteracted includes file path, GPS position (latitude and longitude).
The CLI includes an option to generate an HTML with same info.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_ = args
		eData := exif.ExifData{}

		err = extract.ExifData(&eData, cmd)
		if err != nil {
			return err
		}

		err = writer.Write(&eData, cmd)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	extractCmd.Flags().BoolP("subdir", "", false, "includes the subdirectories in the search")
	extractCmd.Flags().BoolP("html", "", false, "generates an HTML with the info extracted from the images")
	extractCmd.Flags().StringP("dir", "", ".", "path to the directory where the images are located")
	extractCmd.Flags().StringP("output", "", "output.csv", "path to the output file")
}
