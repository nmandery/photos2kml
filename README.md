# photos2kml

Convert the geographical information stored in the EXIF data of JPEG images to generate a [KML](https://developers.google.com/kml/) file.

Features:

* Generates a placemark for each each of the photos which contains geodata as well as a timestamp.
* A path will be added to the KML which connects all files in the order they were recorded. This is particularly useful to create routes of trips using the photos taken during this time.
* Optionally the names of places are added by using the reverse geocoding of [http://nominatim.openstreetmap.org](http://nominatim.openstreetmap.org)

Files without a recording date or without coordinates will be ignored.


## Installation

You need to have [go](https://golang.org/) installed. Then you can install this tool by executing


    go get github.com/nmandery/photos2kml



## Usage

The input file list will be read from standard input (stdin), the KML will be written to stdout.

Example using a prepared list of files:

    photos2kml <my_photo_list.txt >my_photos.kml


Example using `find`:

    find . -name '*.jpg' | photos2kml >my_photos.kml


There are addtional switches the influence the contents of the KML file. See

    photos2kml -h


for a list of options.
