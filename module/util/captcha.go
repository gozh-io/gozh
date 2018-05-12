package util

import (
	"sync"

	"github.com/mojocn/base64Captcha"
)

type CaptchaConfig struct {
	Id              string
	CaptchaType     string
	VerifyValue     string
	ConfigAudio     base64Captcha.ConfigAudio
	ConfigCharacter base64Captcha.ConfigCharacter
	ConfigDigit     base64Captcha.ConfigDigit
}

var (
	captchaConfig     *CaptchaConfig
	captchaConfigOnce sync.Once
)

// 获取base64验证码基本配置
func GetCaptchaConfig() *CaptchaConfig {
	captchaConfigOnce.Do(func() {
		captchaConfig = &CaptchaConfig{
			Id:          "",
			CaptchaType: "character",
			VerifyValue: "",
			ConfigAudio: base64Captcha.ConfigAudio{},
			ConfigCharacter: base64Captcha.ConfigCharacter{
				Height:             60,
				Width:              240,
				Mode:               2,
				IsUseSimpleFont:    false,
				ComplexOfNoiseText: 0,
				ComplexOfNoiseDot:  0,
				IsShowHollowLine:   false,
				IsShowNoiseDot:     false,
				IsShowNoiseText:    false,
				IsShowSlimeLine:    false,
				IsShowSineLine:     false,
				CaptchaLen:         0,
			},
			ConfigDigit: base64Captcha.ConfigDigit{},
		}
	})
	return captchaConfig
}
