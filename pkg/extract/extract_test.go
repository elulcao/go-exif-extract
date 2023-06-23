package extract

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"
	"reflect"
	"sort"
	"testing"
)

// TestRemove tests the removeFromSlice function
func TestRemove(t *testing.T) {
	files := []string{"file1", "file2", "file3"}

	removeFromSlice(&files, "file2")

	expected := []string{"file1", "file3"}

	if !reflect.DeepEqual(files, expected) {
		t.Errorf("Expected %v, but got %v", expected, files)
	}
}

// TestCheckFileIsImage tests the checkFileIsImage function
func TestCheckFileIsImage(t *testing.T) {
	tmpfile := t.TempDir()
	defer os.Remove(tmpfile) // clean up

	// create a dummy file with image exif data
	f, err := os.Create(tmpfile + "/test.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	defer os.Remove(f.Name()) // clean up

	if _, err := f.Write([]byte("this is not an image")); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	// Test the function with the temporary file
	isImage, err := checkFileIsImage(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	if isImage {
		t.Errorf("Expected false, but got true")
	}
}

// TestReadDirAndSubdirs tests the readDirAndSubdirs function
func TestReadDirAndSubdirs(t *testing.T) {
	tmpDir := t.TempDir() + "/test"
	defer os.Remove(tmpDir)

	td := tmpDir
	images := make([]string, 0)
	expected := make([]string, 0)

	for i := 0; i < 30; i++ {
		err := os.MkdirAll(td, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}

		newFile := fmt.Sprintf("%s/test_%d.jpg", td, i)
		file, err := os.Create(newFile)
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()

		width, height := 100, 100
		background := color.RGBA{0, 0xFF, 0, 0xCC}
		image := createImage(width, height, background)

		err = jpeg.Encode(file, image, &jpeg.Options{Quality: 80})
		if err != nil {
			t.Fatal(err)
		}

		expected = append(expected, newFile)
		td = fmt.Sprintf("%s/%d", td, i+1)
	}

	images, err := readDirAndSubdirs(tmpDir, true, &images)
	if err != nil {
		t.Fatal(err)
	}

	sort.Slice(images, func(i, j int) bool {
		return images[i] < images[j]
	})
	sort.Slice(expected, func(i, j int) bool {
		return expected[i] < expected[j]
	})

	if !reflect.DeepEqual(images, expected) {
		t.Errorf("Expected %v, but got %v", expected, images)
	}
}

// createImage creates a new image with the given width, height and background color. Used for testing.
func createImage(width int, height int, background color.RGBA) *image.RGBA {
	rect := image.Rect(0, 0, width, height)
	img := image.NewRGBA(rect)
	draw.Draw(img, img.Bounds(), &image.Uniform{background}, image.Point{}, draw.Src)

	return img
}
