package writer

import (
	"github.com/elulcao/go-exif-extract/pkg/exif"
	"os"
	"testing"

	"github.com/spf13/cobra"
)

// TestWrite tests the Write function
func TestWrite(t *testing.T) {
	eData := exif.ExifData{
		Exif: map[string]exif.Exif{
			"test1.jpg": {
				FilePath:  "test1.jpg",
				Latitude:  1.1,
				Longitude: 1.2,
			},
			"test2.jpg": {
				FilePath:  "test2.jpg",
				Latitude:  2.1,
				Longitude: 2.2,
			},
		},
	}

	tmpDir := t.TempDir() + "/test.csv"
	defer os.Remove(tmpDir)
	cmd := &cobra.Command{}
	cmd.Flags().StringP("output", "", tmpDir, "path to the output file")

	err := Write(&eData, cmd)
	if err != nil {
		t.Errorf("Error writing csv file: %v", err)
	}
}
