package main

import (
	"encoding/json"
	"image"

	//"errors"
	"fmt"
	"image/jpeg"
	"os"
	"path/filepath"

	flag "github.com/spf13/pflag"
	"go.uber.org/zap"
	//"go.uber.org/zap/zapcore"
)

// Checks whether or not required arguments are set and returns true or false
// accordingly.
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

// Wraps CheckForEmptyArg causing the app to exit and usage help to print
// when required args are missing.
func ValidateRequiredArgs(logger *zap.Logger) {
	missingArgs := CheckForEmptyArg()
	if missingArgs {
		fmt.Fprintf(os.Stderr, "Usage of maps\n")
		flag.PrintDefaults()
		logger.Fatal("One or more missing args not set.")
	}
}

// Checks the built version of the app variable version and prints when asked.
func CheckVersion(vflag bool) {
	if vflag {
		fmt.Printf("mapscli version: %s\n", version)
		os.Exit(0)
	}
}

// Loads the config.json file  and returns the config defined by APIConfig type.
func LoadConfig(config string, logger *zap.Logger) (cfg APIConfig) {
	readConfig, err := os.ReadFile(filepath.Clean(config))
	if err != nil {
		logger.Fatal("could not read configuration file.")
	}
	if err := json.Unmarshal(readConfig, &cfg); err != nil {
		logger.Fatal("Could not unmarshal json data.")
	}
	return cfg
}

// Writes the mapResp (map response from the maps api) to an image file.
func WriteImage(dir string, filename string, mapResp image.Image, logger *zap.Logger) {
	filename = fmt.Sprintf("%s.jpg", filename)
	fullPath := filepath.Join(dir, filename)
	out, err := os.Create(fullPath)
	if err != nil {
		logger.Fatal("Could not create file.")
	}
	logger.Sugar().Desugar().Sugar().Infof("creating image file: %s", "image.jpg")

	defer out.Close() //#nosec G307 manually closing after writing the file
	if err = jpeg.Encode(out, mapResp, nil); err != nil {
		logger.Sugar().Fatalf("failed to encode %v", err)
	}

	if err := out.Close(); err != nil {
		logger.Sugar().Fatal("failed to close file")
	}

	logger.Sugar().Desugar().Sugar().Infof("new file '%s' has been generated", filename)
}

// Creates the output_dir in the config file if the directory doesn't exist
func CreateDirIfNotExist(dir string, logger *zap.Logger) {
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Sugar().Infof("Directory %s does not exist. Creating...", dir)
			err := os.MkdirAll(dir, 0755)
			if err != nil && !os.IsExist(err) {
				logger.Sugar().Fatalf("error: %s \n creating %s failed, exiting", err, dir)
			}
		}
	} else {
		logger.Sugar().Infof("%s already exists..", dir)
	}
}
