package selferror

import (
	"encoding/json"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle

func InitI18n() error {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	if _, err := bundle.LoadMessageFile("pkg/selferror/en.json"); err != nil {
		return err
	}
	if _, err := bundle.LoadMessageFile("pkg/selferror/zh.json"); err != nil {
		return err
	}
	return nil
}

func LocalizeError(lang, key string) (string, error) {
	localize := i18n.NewLocalizer(bundle, lang)
	msg, err := localize.Localize(&i18n.LocalizeConfig{
		MessageID: strings.ToLower(key),
	})
	if err != nil {
		return "", err
	}
	return msg, nil
}
