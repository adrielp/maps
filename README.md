# Maps
Repsitory hosting files to create territory maps & download the images by
using the Maps Static API from the [Google Maps Platform](https://developers.google.com/maps/documentation).


## TODO
* make requests to static map via google maps library
* ensure maps library outputs image file
* read in kml file & parse data
* write unit tests
* add additional config (like default sizing of images & zoom) into `config.json`
    * add size and zoom level for example

## Requirements
You need an [API Key](https://developers.google.com/maps/documentation/maps-static/get-api-key)
and a [Digital Signature](https://developers.google.com/maps/documentation/maps-static/digital-signature)

Take your API Key and create a directory in your home named `~/.maps`. Then place your API key in a json file named `config.json`.
The full path should resemble `~/.maps/config.json`.

The JSON file should resemble:

```json
{
    "api_key": "myapikey1234"
}
```

## Install
There are a couple different methods to install `maps`.

### Preferred methods
* Via `go` (recommended): `go install github.com/adrielp/maps`
* Via `brew`: `brew install adrielp/tap/maps` (Mac / Linux)


### Mac/Linux during local development
* Clone down this repository and run `make install`

### Windows
There's a binary for that, but it's not directly supported or tested because `#windows`

## Getting Started
### Prereqs
* Have [make](https://www.gnu.org/software/make/) installed
* Have [GoReleaser](https://goreleaser.com/) installed

### Instructions
* Clone down this repository
* Run commands in the [Makefile](./Makefile) like `make build`
