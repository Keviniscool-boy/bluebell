package controller

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	zhtranslations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"strings"
)

// Trans 全局翻译器
var Trans ut.Translator

// InitTrans 初始化翻译器并为 validator 注册翻译。locale 支持 "zh" 或 "en"，默认使用中文。
func InitTrans(locale string) (err error) {
	// 验证 gin 使用的 validator 引擎类型
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return fmt.Errorf("binding.Validator.Engine() not *validator.Validate")
	}

	// 使用 json tag 作为字段名
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// 本地化支持
	zhT := zh.New()
	enT := en.New()
	uni := ut.New(zhT, enT)

	var trans ut.Translator
	trans, ok = uni.GetTranslator(locale)
	if !ok {
		// fallback to zh
		trans, _ = uni.GetTranslator("zh")
	}

	switch locale {
	case "en":
		if err = entranslations.RegisterDefaultTranslations(v, trans); err != nil {
			return err
		}
	default:
		if err = zhtranslations.RegisterDefaultTranslations(v, trans); err != nil {
			return err
		}
	}

	Trans = trans
	return nil
}

// Translate 将 validator 返回的错误翻译为 map，key 为去掉结构体前缀的字段名。
func Translate(err error) map[string]string {
	out := make(map[string]string)
	if err == nil {
		return out
	}

	if errs, ok := err.(validator.ValidationErrors); ok {
		translated := errs.Translate(Trans)
		for k, v := range translated {
			// k 形如 StructName.Field，去掉 StructName.
			parts := strings.SplitN(k, ".", 2)
			if len(parts) == 2 {
				out[parts[1]] = v
			} else {
				out[k] = v
			}
		}
		return out
	}

	out["error"] = err.Error()
	return out
}
