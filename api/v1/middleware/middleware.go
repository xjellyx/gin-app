package middleware

import (
	"gin-app/internal/domain"
	"gin-app/pkg/scontext"
	"gin-app/pkg/serror"
	"github.com/gin-gonic/gin/binding"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/locales/zh_Hant_TW"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translation "github.com/go-playground/validator/v10/translations/en"
	zh_translation "github.com/go-playground/validator/v10/translations/zh"
	"github.com/go-playground/validator/v10/translations/zh_tw"
)

// HandlerError 错误统一处理
func HandlerError() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Errors != nil {
			resp := domain.Response{
				Code: "FAIL",
				Data: nil,
			}
			last := c.Errors.Last()
			// 判断错误是否是自定义错误
			selfError, ok := last.Err.(serror.SelfError)
			if ok {
				resp.Code = selfError.Code()
				resp.Msg = selfError.Error()
			}
			validationErrors, ok := last.Err.(validator.ValidationErrors)
			// 翻译验证错误消息
			translatedErrors := make(serror.TranslateErr, 0, len(validationErrors))
			if ok {
				for _, e := range validationErrors {
					translatedErrors = append(translatedErrors, e.Translate(translate(scontext.GetLanguage(c.Request.Context()))))
				}
				resp.Msg = translatedErrors.Error()
			}
			c.JSON(200, resp)
			return
		}
	}
}

func translate(language string) ut.Translator {
	var trans ut.Translator
	var validate = binding.Validator.Engine().(*validator.Validate)
	switch language {
	case "en":
		defaultEn := en.New()
		uni := ut.New(defaultEn, defaultEn)
		trans, _ = uni.GetTranslator(language)
		_ = en_translation.RegisterDefaultTranslations(validate, trans)
	case "zh-tw":
		defaultZhTw := zh_Hant_TW.New()
		uni := ut.New(defaultZhTw, defaultZhTw)
		trans, _ = uni.GetTranslator(language)

		_ = zh_tw.RegisterDefaultTranslations(validate, trans)
	default:
		defaultZh := zh.New()
		uni := ut.New(defaultZh, defaultZh)
		trans, _ = uni.GetTranslator(language)
		_ = zh_translation.RegisterDefaultTranslations(validate, trans)
	}
	return trans
}

// HandlerHeadersCtx 保存需要上下文传递的头部数据
func HandlerHeadersCtx() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.GetHeader("X-Language")
		ctx := scontext.SetLanguage(c.Request.Context(), lang)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
