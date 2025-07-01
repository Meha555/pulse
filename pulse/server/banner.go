package server

import (
	"os"
	"pulse/utils"
)

type Banner struct {
	banner string
}

func (b *Banner) Show() {
	println(b.banner)
	println("Version:", utils.GetVersion())
}

func NewBanner() *Banner {
	bannerFile := os.Getenv("PULSE_BANNER_FILE")
	content, err := os.ReadFile(bannerFile)
	if err != nil {
		return nil
	}
	return &Banner{banner: string(content)}
}
