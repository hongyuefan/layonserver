package verifycode

import (
	"github.com/mojocn/base64Captcha"
)

var store = base64Captcha.DefaultMemStore

//  DriverAudio   *base64Captcha.DriverAudio
// 	DriverString  *base64Captcha.DriverString
// 	DriverChinese *base64Captcha.DriverChinese
// 	DriverMath    *base64Captcha.DriverMath
// 	DriverDigit   *base64Captcha.DriverDigit

func VCodeGenerate(length int) (capId, pngBase64 string, err error) {

	// driver = param.DriverAudio
	// driver = param.DriverString.ConvertFonts()
	// driver = param.DriverMath.ConvertFonts()
	// driver = param.DriverChinese.ConvertFonts()
	// driver = param.DriverDigit

	c := base64Captcha.NewCaptcha(&base64Captcha.DriverDigit{Height: 40, Width: 120, Length: length, MaxSkew: 0.7, DotCount: 10}, store)

	capId, pngBase64, err = c.Generate()

	return
}

func VCodeValidate(identifier, verifyValue string) bool {
	return store.Verify(identifier, verifyValue, true)
}
