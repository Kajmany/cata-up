package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"

	_ "embed"
)

// TODO get better errors here (especially ones I know how to work them into bubbletea TUI/log)

// baseConfig is the 'stock' TOML config. It should have sensible defaults and annotations for adding to it.
//
//go:embed config.toml
var baseConfig string

// TODO this should be fetched from package somehow
const appName = "cata-up"

type Source struct {
	Name string `toml:"name"`
	URI  string `toml:"URI"` // TODO make sure this is decoded as a proper URI
}

type Config struct {
	Sources []Source `toml:"sources"`
}

func GetConfig() (*Config, error) {
	cfgPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	// No cfg where it should be? Write our embedded cfg file to the path!
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		err = writeDefaultConfig(cfgPath)
		if err != nil {
			return nil, fmt.Errorf("no config found, but unable to write the default configuration file: %w", err)
		}
	}

	// Now, there should be *something* to parse.
	var config Config
	_, err = toml.DecodeFile(cfgPath, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func writeDefaultConfig(path string) error {
	// Windows: Reason why this is a function - might have to create a directory.
	// Linux: Not my problem if the directory doesn't already exist
	if runtime.GOOS == "windows" {
		dir := filepath.Dir(path)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err = os.Mkdir(dir, 0755)
			if err != nil {
				return fmt.Errorf("unable to create an APPDATA folder for config: %w", err)
			}
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(baseConfig)
	if err != nil {
		return err
	}
	return nil
}

func getConfigPath() (string, error) {
	// Windows: We want our own folder inside AppData/Roaming with a generic name
	// Linux: Just put the file inside the XDG_CONFIG_HOME with an app-related name

	// This could be a compile-time constant but that's two more files
	// And the rest of this is already done at runtime, so whatever
	var fileName string
	var configDir string

	switch runtime.GOOS {
	case "windows":
		fileName = "config.toml"
		appData := os.Getenv("APPDATA")
		if appData == "" {
			return "", fmt.Errorf("APPDATA environment variable not set")
		}
		configDir = filepath.Join(appData, appName)

	case "linux":
		fileName = appName + ".toml"
		xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
		if xdgConfigHome == "" {
			home := os.Getenv("HOME")
			if home == "" {
				return "", fmt.Errorf("HOME environment variable not set")
			}
			xdgConfigHome = filepath.Join(home, ".config")
		}
		configDir = xdgConfigHome

	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	configPath := filepath.Join(configDir, fileName)
	return configPath, nil
}
