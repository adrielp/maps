package main

import (
	"context"
	"fmt"
	"os"

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

	if err != nil {
		fmt.Printf("Failed to build logger. \nError: \n%s", err)
		os.Exit(1)
	}

	// Replaces the global logger with Zap. "Discouraged" but significantly
	// cleaner than passing the logger around through files and functions
	zap.ReplaceGlobals(logger)

	ValidateRequiredArgs(logger)

	config := LoadConfig(ConfigFile, logger)

	// Create output directory
	CreateDirIfNotExist(config.OutputDirectory, logger)

	client, err := maps.NewClient(maps.WithAPIKey(config.APIKey))
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
	WriteImage(config.OutputDirectory, "image1", mapResp, logger)
}
