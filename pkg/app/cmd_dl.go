package app

import (
	"fmt"

	"gitar/pkg/config"
)

func DownloadArchive(url string) error {
	println(url)
	cfg, err := config.LoadConfig()
	fmt.Printf("%+v", cfg)
	return err
}
