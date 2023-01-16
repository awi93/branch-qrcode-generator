package main

import (
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/fogleman/gg"
	qrcode "github.com/skip2/go-qrcode"
)

func main() {
	file, err := os.Open("csv/example-data.csv")
	if err != nil {
		panic(err)
	}
	r := csv.NewReader(file)

	_, err = r.Read()
	if err != nil {
		log.Fatal(err)
	}
	for {
		row, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		var data = map[string]interface{}{
			"id":            row[0],
			"merchant_id":   row[1],
			"branch_ref_id": row[2],
			"branch_type":   row[3],
			"label":         row[4],
		}
		var qrcodeData = map[string]interface{}{
			"type": "EVENT_POINT_REDEEM",
			"data": data,
		}
		jsonString, err := json.Marshal(qrcodeData)
		if err != nil {
			log.Fatal(err)
		}
		base64 := base64.StdEncoding.EncodeToString(jsonString)
		err = qrcode.WriteFile(base64, qrcode.Medium, 1000, fmt.Sprintf("qrcodes/%s.png", data["id"].(string)))
		if err != nil {
			log.Fatal(err)
		}
		bgImg, err := gg.LoadPNG("templates/template.png")
		if err != nil {
			log.Fatal(err)
		}
		pngImg, err := gg.LoadPNG(fmt.Sprintf("qrcodes/%s.png", data["id"].(string)))
		if err != nil {
			log.Fatal(err)
		}
		dc := gg.NewContext(1410, 2000)
		dc.DrawImage(bgImg, 0, 0)
		dc.DrawImage(pngImg, 205, 500)
		dc.DrawRoundedRectangle(205, 340, 1000, 150, 20)
		dc.SetRGB(1, 1, 1)
		dc.FillPreserve()
		dc.SetRGB(0, 0, 0)
		err = dc.LoadFontFace("fonts/SourceSansPro-Black.ttf", 60)
		if err != nil {
			log.Fatal(err)
		}
		dc.DrawStringAnchored(data["label"].(string), 705, 415, 0.5, 0.5)
		dc.SavePNG(fmt.Sprintf("qrcode-banners/%s.png", data["label"].(string)))
	}
}
