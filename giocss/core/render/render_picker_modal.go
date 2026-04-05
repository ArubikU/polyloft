package render

import (
	"image"
	"image/color"
)

type PickerModalModel struct {
	OverlayColor color.NRGBA
	ModalBg      color.NRGBA
	BorderColor  color.NRGBA
	TitleColor   color.NRGBA
	ValueColor   color.NRGBA
	DecColor     color.NRGBA
	IncColor     color.NRGBA
	CloseColor   color.NRGBA

	ModalRadius  int
	ButtonRadius int

	ModalRect image.Rectangle
	TitleRect image.Rectangle
	SepRect   image.Rectangle
	ValueRect image.Rectangle
	DecRect   image.Rectangle
	IncRect   image.Rectangle
	CloseRect image.Rectangle

	TitleText string
	ValueText string
	DecPath   string
	IncPath   string
	ClosePath string
}

func BuildPickerModalModel(screenW, screenH int, pickerType, pickerValue, pickerModalOpen string) PickerModalModel {
	modalW := minInt(400, screenW-40)
	if modalW < 160 {
		modalW = maxInt(1, screenW)
	}
	modalH := 300
	if modalH > screenH-20 {
		modalH = maxInt(120, screenH-20)
	}
	modalX := (screenW - modalW) / 2
	modalY := (screenH - modalH) / 2

	titleY := modalY + 15
	sepY := titleY + 35
	contentY := sepY + 20

	btnW := 50
	btnH := 30
	btnY := contentY + 50
	decX := modalX + 20
	incX := modalX + modalW - btnW - 20

	closeW := 80
	closeH := 30
	closeY := modalY + modalH - 40
	closeX := modalX + (modalW-closeW)/2

	valueText := pickerValue
	if valueText == "" {
		valueText = "(no value)"
	}

	return PickerModalModel{
		OverlayColor: color.NRGBA{R: 0, G: 0, B: 0, A: 200},
		ModalBg:      color.NRGBA{R: 0x1F, G: 0x29, B: 0x37, A: 0xFF},
		BorderColor:  color.NRGBA{R: 0x47, G: 0x55, B: 0x69, A: 0xFF},
		TitleColor:   color.NRGBA{R: 0xE2, G: 0xE8, B: 0xF0, A: 0xFF},
		ValueColor:   color.NRGBA{R: 0x94, G: 0xA3, B: 0xB8, A: 0xFF},
		DecColor:     color.NRGBA{R: 0x47, G: 0x55, B: 0x69, A: 0xFF},
		IncColor:     color.NRGBA{R: 0x22, G: 0xC5, B: 0x5E, A: 0xFF},
		CloseColor:   color.NRGBA{R: 0x47, G: 0x55, B: 0x69, A: 0xFF},

		ModalRadius:  12,
		ButtonRadius: 6,

		ModalRect: image.Rect(modalX, modalY, modalX+modalW, modalY+modalH),
		TitleRect: image.Rect(modalX+10, titleY, modalX+modalW-10, titleY+30),
		SepRect:   image.Rect(modalX+10, sepY, modalX+modalW-10, sepY+1),
		ValueRect: image.Rect(modalX+20, contentY, modalX+modalW-20, contentY+40),
		DecRect:   image.Rect(decX, btnY, decX+btnW, btnY+btnH),
		IncRect:   image.Rect(incX, btnY, incX+btnW, btnY+btnH),
		CloseRect: image.Rect(closeX, closeY, closeX+closeW, closeY+closeH),

		TitleText: "Select " + pickerType,
		ValueText: valueText,
		DecPath:   pickerModalOpen + "__modal_dec",
		IncPath:   pickerModalOpen + "__modal_inc",
		ClosePath: pickerModalOpen + "__modal_close",
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}





