package exif

type ExifData struct {
	Exif map[string]Exif
}

type Exif struct {
	FilePath  string
	Latitude  float64
	Longitude float64
}
