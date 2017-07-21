package main

import (
	"bufio"
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"io"
	"os"
	"path"
	"sort"
	"strings"
	"time"
)

var (
	useAbsoluteFilenames = false
)

type Photo struct {
	Filename  string
	Timestamp time.Time
	Lon       float64
	Lat       float64
}

type Photos []*Photo

func (ps Photos) Len() int {
	return len(ps)
}

func (ps Photos) Less(i, j int) bool {
	return ps[i].Timestamp.Before(ps[j].Timestamp)
}

func (ps Photos) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

func Tell(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(os.Stderr, format+"\n", a...)
}

func PlacemarkName(filename string) string {
	if useAbsoluteFilenames {
		return filename
	}
	return path.Base(filename)
}

func ReadPhotosFromList(reader *bufio.Reader) (photos Photos, err error) {
	for {
		filename, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				err = nil // not realy an error in this context
			}
			return photos, err

		}
		filename = strings.TrimRight(filename, "\n")
		if filename == "" {
			continue // skip empty lines
		}

		photo := &Photo{
			Filename: PlacemarkName(filename),
		}

		photoFile, err := os.Open(filename)
		if err != nil {
			return photos, err
		}

		exifdata, err := exif.Decode(photoFile)
		if err != nil {
			return photos, err
		}

		timestamp, terr := exifdata.DateTime()
		if terr != nil {
			Tell("The photo %s has no timestamp -> will be skipped", filename)
			continue
		}
		photo.Timestamp = timestamp

		photo.Lat, photo.Lon, terr = exifdata.LatLong()
		if terr != nil {
			Tell("The photo %s has no location -> will be skipped", filename)
			continue
		}

		photos = append(photos, photo)
	}
	return
}

func WriteKML(w io.Writer, photos Photos) {
	// https://developers.google.com/kml/documentation/kml_tut#paths
	fmt.Fprint(w, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<kml xmlns=\"http://www.opengis.net/kml/2.2\" xmlns:gx=\"http://www.google.com/kml/ext/2.2\">")
	fmt.Fprint(w, "<Document>")

	// single photos
	for _, photo := range photos {
		fmt.Fprintf(w, "<Placemark><name>")
		xml.EscapeText(w, []byte(photo.Filename))
		fmt.Fprintf(w, "</name><TimeStamp><when>%s</when></TimeStamp><Point><coordinates>%f,%f</coordinates></Point></Placemark>", photo.Timestamp.Format(time.RFC3339), photo.Lon, photo.Lat)
	}

	// path consisting of all photos
	fmt.Fprint(w, "<Placemark><name>Path</name><LineString><coordinates>")
	for _, photo := range photos {
		fmt.Fprintf(w, "%f,%f ", photo.Lon, photo.Lat)
	}
	fmt.Fprint(w, "</coordinates></LineString></Placemark>")

	fmt.Fprint(w, "</Document></kml>")
}

func init() {
	flag.BoolVar(&useAbsoluteFilenames, "a", false, "Use the absolute filenames of the photos for the name of the placemarks. Default is using just the basename.")
}

func Usage() {
	fmt.Fprintf(os.Stderr, "Reads a list of filenames pf photos from stdin, generates a KML document with the\nlocations of the photos and a path connecting the photos in chronological order\nand writes the KML to stdout.\n\n")
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])

	flag.PrintDefaults()
	os.Exit(0)
}

func main() {
	flag.Usage = Usage
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)
	Tell("Reading list of photos from stdin ...")

	// register camera makenote data parsing
	exif.RegisterParsers(mknote.All...)

	photos, err := ReadPhotosFromList(reader)
	if err != nil {
		Tell("An error occured: %v", err)
		os.Exit(1)
	}

	Tell("Collected %d photos", len(photos))
	sort.Sort(photos) // sort photos by their timestamp

	WriteKML(os.Stdout, photos)
}
