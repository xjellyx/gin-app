package middleware

import (
	"bytes"
	"errors"
	"io"
	"log/slog"
	"time"

	"gin-app/internal/bootstrap"
	"gin-app/internal/domain"
	"gin-app/internal/usecase"
	"gin-app/pkg/scontext"
	"gin-app/pkg/serror"

	"github.com/gin-contrib/cors"
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
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
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
			var selfError serror.SelfError
			var validationErrors validator.ValidationErrors
			switch {
			case errors.As(last.Err, &selfError): // 判断错误是否是自定义错误
				resp.Code = selfError.Code()
				resp.Msg = selfError.Error()
			case errors.As(last.Err, &validationErrors): // 判断错误是否是validator错误
				// 翻译验证错误消息
				translatedErrors := make(serror.TranslateErr, 0, len(validationErrors))
				for _, e := range validationErrors {
					translatedErrors = append(translatedErrors, e.Translate(translate(scontext.GetLanguage(c.Request.Context()))))
				}
				resp.Msg = translatedErrors.Error()
			default:
				slog.Error("HandlerError", "uri", c.Request.RequestURI, "err", c.Errors)
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

func HandlerAuth() gin.HandlerFunc {
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

		ctx.Request = ctx.Request.WithContext(scontext.SetUserUuid(ctx.Request.Context(), cla.UserUuid))
		ctx.Request = ctx.Request.WithContext(scontext.SetUsername(ctx.Request.Context(), cla.Username))
		ctx.Next()
	}
}

func LimitRequestRate(limit *limiter.Limiter) gin.HandlerFunc {
	return mgin.NewMiddleware(limit)
}

func CorsHandler() gin.HandlerFunc {
	corsCfg := cors.DefaultConfig()
	corsCfg.AllowAllOrigins = true
	return cors.New(corsCfg)
}

type formatter struct {
	Address    string        `json:"address"`
	Time       string        `json:"time"`
	Protoc     string        `json:"protoc"`
	Method     string        `json:"method"`
	Path       string        `json:"path"`
	Code       int           `json:"code"`
	Body       any           `json:"body"`
	BodySize   int           `json:"bodySize"`
	Latency    time.Duration `json:"latency"`
	ErrMessage string        `json:"errMessage"`
}
type LogFormatter func(params gin.LogFormatterParams) formatter

type bodyLogWriter struct {
	gin.ResponseWriter
	bodyBuf *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.bodyBuf.Write(b)
	return w.ResponseWriter.Write(b)
}
func format(f LogFormatter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 读取body,复用body信息
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
		}
		// 把刚刚读出来的再写进去
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		// 复用response的信息
		var blw bodyLogWriter
		blw = bodyLogWriter{bodyBuf: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		// Process request
		c.Next()

		// Log only when path is not being skipped
		param := gin.LogFormatterParams{
			Request: c.Request,
			Keys:    c.Keys,
		}
		// Stop timer
		param.TimeStamp = time.Now()
		param.Latency = param.TimeStamp.Sub(start)
		param.ClientIP = c.ClientIP()
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
		param.BodySize = blw.Size()
		if raw != "" {
			path = path + "?" + raw
		}
		param.Path = path
		slog.Info("HTTP Request ", "HTTP", f(param))
	}
}
func Format() gin.HandlerFunc {
	return format(func(params gin.LogFormatterParams) formatter {
		var res = formatter{
			Address:    params.ClientIP,
			Time:       params.TimeStamp.Format("2006-01-02 15:04:05"),
			Protoc:     params.Request.Proto,
			Method:     params.Request.Method,
			Path:       params.Path,
			Code:       params.StatusCode,
			BodySize:   params.BodySize,
			Latency:    params.Latency,
			ErrMessage: params.ErrorMessage,
		}
		return res
	})
}
