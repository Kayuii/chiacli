package fix

import (
	"fmt"
	"path/filepath"
	"regexp"

	gfind "github.com/kayuii/chiacli/gfind"
)

type Fix struct {
	DirName string
	LogPath string
	PlotID  string
	LogFile string
}

func New() *Fix {
	return &Fix{}
}

type Config struct {
	FarmerKey string `yaml:"FarmerKey"`
	PoolKey   string `yaml:"PoolKey"`
	LocalSk   string `yaml:"LocalSk"`
	Memo      string `yaml:"Memo"`
	FilePath  string `yaml:"FilePath"`
	Pattern   string `yaml:"Pattern"`
}

func (f *Fix) Check(config *Config) error {

	absDir, err := filepath.Abs(config.FilePath)
	if err != nil {
		return err
	}
	dir := filepath.Clean(absDir)

	pattern, err := regexp.Compile(config.Pattern)
	if err != nil {
		return err
	}

	finder := gfind.NewFinder(pattern)

	matches, _ := finder.Find(dir)
	for _, match := range matches {
		fmt.Println(match)
	}

	return nil
}
