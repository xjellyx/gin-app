package scontext

import (
	"context"
	"golang.org/x/text/language"
	"strings"
)

type languageCtxTag struct {
}

// SetLanguage set language to context
func SetLanguage(ctx context.Context, lang string) context.Context {
	return context.WithValue(ctx, languageCtxTag{}, lang)
}

// GetLanguage get language by context
func GetLanguage(ctx context.Context) string {
	val, ok := ctx.Value(languageCtxTag{}).(string)
	if ok {
		return strings.ToLower(val)
	}
	return language.Chinese.String()
}

type userUuidCtxTag struct{}

// SetUserUuid set user uuid to context
func SetUserUuid(ctx context.Context, userUuid string) context.Context {
	return context.WithValue(ctx, userUuidCtxTag{}, userUuid)
}

// GetUserUuid get user uuid by context
func GetUserUuid(ctx context.Context) string {
	if val, ok := ctx.Value(userUuidCtxTag{}).(string); ok {
		return val
	}
	return ""
}

type usernameCtxTfg struct{}

// SetUsername 	set username to context
func SetUsername(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, usernameCtxTfg{}, username)
}

// GetUsername get username by context
func GetUsername(ctx context.Context) string {
	if val, ok := ctx.Value(usernameCtxTfg{}).(string); ok {
		return val
	}
	return ""
}

type rolesCtxTag struct{}

// SetRoles set roles to context
func SetRoles(ctx context.Context, roles []string) context.Context {
	return context.WithValue(ctx, rolesCtxTag{}, roles)
}

// GetRoles get roles by context
func GetRoles(ctx context.Context) []string {
	if val, ok := ctx.Value(rolesCtxTag{}).([]string); ok {
		return val
	}
	return nil
}
