package main

import (
	"context"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/kr/pretty"
	flag "github.com/spf13/pflag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"googlemaps.github.io/maps"
)

// Set default constants for flag usage messages.
const (
	configFileUsage   = "Path to the configuration file."
	inputFolderUsage  = "Path to folder containing exported KML files."
	outputFolderUsage = "Path to folder to place generated images."
	versionUsageUsage = "Prints out the version of schemacheck"
	disableColorUsage = "Disables colored logging output."
	versionUsage      = "Outputs the current version of the maps CLI utility"
	debugUsage        = "Enables debug logging levels."
)

// Core variables for flag pointers and info, warning, and error loggers.
var (
	// Core flag variables
	ConfigFile   string
	InputFolder  string
	OutputFolder string
	DisableColor bool
	VersionFlag  bool
	DebugLogging bool

	// version is set through ldflags by GoReleaser upon build, taking in the most recent tag
	// and appending -snapshot in the event that --snapshot is set in GoReleaser.
	version string
)

// Initialize the flags from the command line and their shorthand counterparts.
func init() {
	flag.StringVarP(&ConfigFile, "config-file", "f", "", configFileUsage)
	flag.StringVarP(&InputFolder, "input-folder", "i", "", inputFolderUsage)
	flag.StringVarP(&OutputFolder, "output-folder", "o", "", outputFolderUsage)
	flag.BoolVar(&DisableColor, "disable-color", false, disableColorUsage)
	flag.BoolVarP(&VersionFlag, "version", "v", false, versionUsage)
	flag.BoolVar(&DebugLogging, "debug", true, debugUsage)
}

// Check whether or not required arguments are set
func CheckForEmptyArg() bool {
	configFileArgEmpty := true
	flag.VisitAll(func(f *flag.Flag) {
		if f.Name == "config-file" {
			if f.Changed {
				configFileArgEmpty = false
			}
		}
	})
	return configFileArgEmpty
}

func main() {
	// Parse the flags set in the init() function
	flag.Parse()

	// Configures logging based on argument logger choice
	var logCfg zap.Config
	if DebugLogging {
		logCfg = zap.NewDevelopmentConfig()
	} else {
		logCfg = zap.NewProductionConfig()
	}
	logCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := logCfg.Build()
	//check(logger, err)
	if err != nil {
		fmt.Printf("Failed to build logger. \nError: \n%s", err)
		os.Exit(1)
	}
	// Replaces the global logger with Zap. "Discouraged" but significantly
	// cleaner than passing the logger around through files and functions
	zap.ReplaceGlobals(logger)

	// If version flag is set, output version of app and exit
	if VersionFlag {
		fmt.Printf("maps version: %s\n", version)
		os.Exit(0)
	}

	// Check to ensure required flags aren't empty
	missingArgs := CheckForEmptyArg()
	if missingArgs {
		fmt.Fprintf(os.Stderr, "Usage of maps\n")
		flag.PrintDefaults()
		logger.Fatal("One or more missing args not set.")
	}

	// Read in the json configuration file containing essential config
	// to make API calls to Maps API
	jsonConfig, err := os.ReadFile(filepath.Clean(ConfigFile))
	if err != nil {
		logger.Fatal("Could not read json configuration file.")
	}

	// Unmarshal the configuration file and automatically
	// pulls API key based on APIConfig type
	var cfg APIConfig
	if err := json.Unmarshal(jsonConfig, &cfg); err != nil {
		logger.Fatal("Could not unmarshal json data.")
	}

	// Build the URL to make the request
	elems := []string{
		"center=Berkley,CA&zoom=14&size=400x400",
		"&key=",
		cfg.APIKey,
	}
	urlKey := strings.Join(elems, "")

	baseURL := &url.URL{
		Scheme:   "https",
		Host:     "maps.googleapis.com",
		Path:     "maps/api/staticmap",
		RawQuery: urlKey,
	}

	if err != nil {
		logger.Fatal("Failed to create base url")
	}

	resp, err := http.Get(baseURL.String())
	if err != nil {
		logger.Fatal("Could not make request")
	}

	logger.Sugar().Infof("Response is: %s", resp)

	// This code is similar to https://github.com/googlemaps/google-maps-services-go/blob/v1.3.2/examples/staticmap/cmdline/main.go
	// and https://github.com/googlemaps/google-maps-services-go/blob/v1.3.2/staticmap.go
	// Apparently they implemented their own CLI tool to do this. But it's very much 1:1
	// In terms of generating images. I want to extend this so that I can read in the KMls generated
	// from people creating maps, parsing the data of each one, then generating out the respective image.
	client, err := maps.NewClient(maps.WithAPIKey(cfg.APIKey))
	if err != nil {
		logger.Sugar().Fatalf("fatal err: %s", err)
	}

	//"center=Berkley,CA&zoom=14&size=400x400",
	r := &maps.StaticMapRequest{
		Center: "Berkeley",
		Zoom:   14,
		Size:   "400x400",
		//Scale:    *scale,
		//Format:   maps.Format(*format),
		//Language: *language,
		//Region:   *region,
		//MapType:  maps.MapType(*maptype),
	}

	mapResp, err := client.StaticMap(context.Background(), r)
	if err != nil {
		logger.Fatal("Could not make request")
	}
	pretty.Println(mapResp)

	out, err := os.Create("image.jpg")
	if err != nil {
		logger.Fatal("Could not create file.")
	}
	defer out.Close()
	if err = jpeg.Encode(out, mapResp, nil); err != nil {
		logger.Sugar().Fatalf("failed to encode %v", err)
	}
}
