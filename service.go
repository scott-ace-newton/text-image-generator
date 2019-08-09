package main

import (
	"flag"
	"fmt"
	"github.com/golang/freetype"
	"github.com/sirupsen/logrus"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"math/rand"
	"os"
)

var (
	dpi      = flag.Float64("dpi", 144, "screen resolution in Dots Per Inch")
	fontfile = flag.String("fontfile", "./fonts/ComicSansMS3.ttf", "filename of the ttf font")
	hinting  = flag.String("hinting", "none", "none | full")
	size     = flag.Float64("size", 16, "font size in points")
	spacing  = flag.Float64("spacing", 1.5, "line spacing (e.g. 2 means double spaced)")

	images = []string{
		"aaron-burden-c333d6YEhi0-unsplash.jpg",
		"adam-bixby-Ix78f0AuCBI-unsplash.jpg",
		"bekir-donmez-eofm5R5f9Kw-unsplash.jpg",
		"ben-white-7BiMECHFgFY-unsplash.jpg",
		"ben-white-TlBF3ZUVTvE-unsplash.jpg",
		"ben-white-vtCBruWoNqo-unsplash.jpg",
		"clemens-posch-we1tBosANpU-unsplash.jpg",
		"david-kovalenko-YVBXuL6Al2w-unsplash.jpg",
		"dex-ezekiel-GBJ7uFCqsh4-unsplash.jpg",
		"ellen-auer-BaM_1KlFZkc-unsplash.jpg",
		"iulia-mihailov-Jn3_uxVmyuA-unsplash.jpg",
		"james-pond-HUiSySuofY0-unsplash.jpg",
		"jasper-van-der-meij-sRQ0MJsWXvE-unsplash.jpg",
		"joshua-earle-Dn3ATeXQEQ4-unsplash.jpg",
		"joshua-sortino-f3uWi9G-lus-unsplash.jpg",
		"joshua-sortino-GPtnV_XdQQU-unsplash.jpg",
		"joshua-sortino-IlvY3z4KVCI-unsplash.jpg",
		"joshua-sortino-LqKhnDzSF-8-unsplash.jpg",
		"joshua-sortino-lRA_WTczjgw-unsplash.jpg",
		"joshua-sortino-m5P0c6ABWDs-unsplash.jpg",
		"joshua-sortino-m-XLNFYdiVw-unsplash.jpg",
		"joshua-sortino-xZqr8WtYEJ0-unsplash.jpg",
		"jr-korpa-XEJDsC5Tzec-unsplash.jpg",
		"kash-goudarzi-uGFGAwTN_3o-unsplash.jpg",
		"kevin-kristhian-hX0AqODUU20-unsplash.jpg",
		"kid-circus-7vSlK_9gHWA-unsplash.jpg",
		"kristopher-roller-o2LlAeqJzVo-unsplash.jpg",
		"kyle-glenn-kvIAk3J_A1c-unsplash.jpg",
		"li-yang-_vPCiuXL2HE-unsplash.jpg",
		"mark-ivan-_y35BqcsauI-unsplash.jpg",
		"megan-johnston-t1NiXOf5fTI-unsplash.jpg",
		"natalie-grainger-Mw1efRU1qcU-unsplash.jpg",
		"paul-green-fhOGkxwQz0s-unsplash.jpg",
		"paula-brustur-ngzjG6ZhoDw-unsplash.jpg",
		"rik-buiting-Zb-nqiQsLe4-unsplash.jpg",
		"samuel-zeller-rk_Zz3b7G2Y-unsplash.jpg",
		"taylor-leopold-Rr9zn33OMbk-unsplash.jpg",
		"tommy-lisbin-gDzhss2CznA-unsplash.jpg",
		"tommy-lisbin-opt65nQcMZc-unsplash.jpg",
		"tommy-lisbin-wnq68O-6UNs-unsplash.jpg",
		"weroad-3QgXtnXakBw-unsplash.jpg",
		"weroad-WVtPZxbKZYs-unsplash.jpg",
		"wolf-schram-3RJG1Ecx7os-unsplash.jpg",
	}
)

type Service interface {
	CreateImage(text string) (*image.RGBA, error)
}

type service struct {
	le *logrus.Entry
}

func NewService(le *logrus.Entry) Service {
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	return &service{
		le: le,
	}
}

func (s *service) CreateImage(text string) (*image.RGBA, error) {
	log := s.le
	// Read the font data.
	fontBytes, err := ioutil.ReadFile(*fontfile)
	if err != nil {
		log.WithError(err).Error("could not read font file")
		return nil, err
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.WithError(err).Error("could not parse font file")
		return nil, err
	}

	imgFile, err := os.Open(fmt.Sprintf("./images/stock/%s", images[rand.Intn(len(images))]))
	if err != nil {
		log.WithError(err).Error("could not read image file")
		return nil, err
	}
	i, _, err := image.Decode(imgFile)
	if err != nil {
		log.WithError(err).Error("could not decode image")
		return nil, err
	}

	// Initialize the context.
	fg := image.Black
	ruler := color.RGBA{0xdd, 0xdd, 0xdd, 0xff}

	rgba := image.NewRGBA(i.Bounds())
	draw.Draw(rgba, rgba.Bounds(), i, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(*dpi)
	c.SetFont(f)
	c.SetFontSize(*size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
	switch *hinting {
	default:
		c.SetHinting(font.HintingNone)
	case "full":
		c.SetHinting(font.HintingFull)
	}

	// Draw the guidelines.
	for i := 0; i < 200; i++ {
		rgba.Set(10, 10+i, ruler)
		rgba.Set(10+i, 10, ruler)
	}

	// Draw the text.
	pt := freetype.Pt(10, 10+int(c.PointToFixed(*size)>>6))

	_, err = c.DrawString(text, pt)
	if err != nil {
		log.WithError(err).Error("could not draw string")
		return nil, err
	}
	pt.Y += c.PointToFixed(*size * *spacing)
	return rgba, nil
}
