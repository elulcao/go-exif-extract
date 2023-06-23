package extract

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	pkgExif "github.com/elulcao/go-exif-extract/pkg/exif"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/spf13/cobra"
)

// ExifData extracts the exif data from images in a directory and writes it to a csv file. The info
// exteracted includes file path, GPS position (latitude and longitude). If html flag is true, it will
// generate an HTML with same info.
func ExifData(ed *pkgExif.ExifData, cmd *cobra.Command) error {
	dir, _ := cmd.Flags().GetString("dir")
	subdir, _ := cmd.Flags().GetBool("subdir")
	files := make([]string, 0)

	files, err := readDirAndSubdirs(dir, subdir, &files)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	ed.Exif = make(map[string]pkgExif.Exif, 0)

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", file, err)
		}
		defer f.Close()

		x, err := exif.Decode(f)
		if err != nil {
			fmt.Printf("Failed to decode EXIF data %s: %v\n", file, err)
			ed.Exif[file] = pkgExif.Exif{FilePath: file}
			continue
		}

		lat, lon, err := x.LatLong()
		if err != nil {
			fmt.Printf("Failed to extract latitude and longitude %s:: %v\n", file, err)
			ed.Exif[file] = pkgExif.Exif{FilePath: file}
			continue
		}

		ed.Exif[file] = pkgExif.Exif{
			FilePath:  file,
			Latitude:  lat,
			Longitude: lon,
		}
	}

	return nil
}

// readDirAndSubdirs reads the directory and returns a slice of strings with the paths of the files
// in the directory and subdirectories.
func readDirAndSubdirs(path string, subdir bool, images *[]string) (files []string, err error) {
	dir, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open directory %s: %w", path, err)
	}
	defer dir.Close()

	fileInfo, err := dir.Readdir(-1)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", path, err)
	}

	for _, file := range fileInfo {
		if file.IsDir() {
			if subdir {
				subfiles, err := readDirAndSubdirs(filepath.Join(path, file.Name()), subdir, images)
				if err != nil {
					return nil, fmt.Errorf("failed to read subdirectory %s: %w", file.Name(), err)
				}

				files = append(files, subfiles...)
			}
		} else {
			files = append(files, filepath.Join(path, file.Name()))
		}
	}

	// Filter the files to only include jpg images
	for _, file := range files {
		isImage, err := checkFileIsImage(file)
		if err != nil {
			return nil, fmt.Errorf("failed to check if file is image %s: %w", file, err)
		}

		if !isImage {
			removeFromSlice(&files, file)
		}
	}

	*images = append(*images, files...)
	return files, nil
}

// checkFileIsImage checks if the file is an image
func checkFileIsImage(file string) (isImage bool, err error) {
	f, err := os.Open(file)
	if err != nil {
		return false, fmt.Errorf("failed to open file %s: %w", file, err)
	}
	defer f.Close()

	buffer := make([]byte, 512)
	_, err = f.Read(buffer)
	if err != nil {
		return false, fmt.Errorf("failed to read file %s: %w", file, err)
	}

	contentType := http.DetectContentType(buffer)

	switch contentType {
	case "image/jpeg", "image/jpg", "image/png", "image/gif":
		return true, nil
	default:
		return false, nil
	}
}

// removeFromSlice removes an element from the slice
func removeFromSlice(files *[]string, file string) {
	for i, f := range *files {
		if f == file {
			*files = append((*files)[:i], (*files)[i+1:]...)
		}
	}
}
