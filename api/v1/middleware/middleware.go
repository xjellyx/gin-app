package middleware

import (
	"gin-app/internal/bootstrap"
	"gin-app/internal/domain"
	"gin-app/internal/infra/cache"
	"gin-app/internal/usecase"
	"gin-app/pkg/scontext"
	"gin-app/pkg/serror"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/locales/zh_Hant_TW"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translation "github.com/go-playground/validator/v10/translations/en"
	zh_translation "github.com/go-playground/validator/v10/translations/zh"
	"github.com/go-playground/validator/v10/translations/zh_tw"
	"go.uber.org/zap"
)

// HandlerError 错误统一处理
func HandlerError(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Errors != nil {
			resp := domain.Response{
				Code: "FAIL",
				Data: nil,
			}
			last := c.Errors.Last()
			switch last.Err.(type) {
			case serror.SelfError: // 判断错误是否是自定义错误
				selfError := last.Err.(serror.SelfError)
				resp.Code = selfError.Code()
				resp.Msg = selfError.Error()
			case validator.ValidationErrors: // 判断错误是否是validator错误
				validationErrors := last.Err.(validator.ValidationErrors)
				// 翻译验证错误消息
				translatedErrors := make(serror.TranslateErr, 0, len(validationErrors))
				for _, e := range validationErrors {
					translatedErrors = append(translatedErrors, e.Translate(translate(scontext.GetLanguage(c.Request.Context()))))
				}
				resp.Msg = translatedErrors.Error()
			default:
				log.Error("HandlerError | "+c.Request.RequestURI, zap.Error(last.Err))
				resp.Code = string(serror.ErrCodeInternalServerError)
				resp.Msg = serror.Error(serror.ErrCodeInternalServerError, scontext.GetLanguage(c.Request.Context())).Error()
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
		if lang == "" {
			c.Next()
			return
		}
		ctx := scontext.SetLanguage(c.Request.Context(), lang)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func HandlerAuth(ca cache.Cache) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			token = ctx.Query("token")
		}
		msg := serror.Error(serror.ErrUnauthorized, scontext.GetLanguage(ctx.Request.Context())).Error()
		if token == "" {
			ctx.AbortWithStatusJSON(401, gin.H{
				"code": serror.ErrUnauthorized,
				"msg":  msg,
			})
			return
		}
		cla := usecase.Claims{}
		if !usecase.ParseToken(token, &cla, bootstrap.GetConfig().JWTSigningKey) {
			ctx.AbortWithStatusJSON(401, gin.H{
				"code": serror.ErrUnauthorized,
				"msg":  msg,
			})
			return
		}
		get, err := ca.Get(ctx.Request.Context(), cla.CacheKey)
		if err != nil {
			return
		}
		if get != token {
			ctx.AbortWithStatusJSON(401, gin.H{
				"code": serror.ErrUnauthorized,
				"msg":  msg,
			})
			return
		}
		ctx.Request = ctx.Request.WithContext(scontext.SetUserUuid(ctx.Request.Context(), cla.UserUuid))
		ctx.Next()
	}
}
