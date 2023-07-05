package orm

import (
	"go.uber.org/zap"
)

type options struct {
	logger            *zap.Logger
	opentracingPlugin *OpentracingPlugin
	autoMigrate       bool
	autoMigrateDst    []any
}

type Option interface {
	apply(*options)
}

type loggerOption struct {
	log *zap.Logger
}

func (l loggerOption) apply(opts *options) {
	opts.logger = l.log
}

func WithLogger(log *zap.Logger) Option {
	return loggerOption{log: log}
}

type autoMigrateOption bool

func (a autoMigrateOption) apply(opts *options) {
	opts.autoMigrate = bool(a)
}

func WithAutoMigrate(a bool) Option {
	return autoMigrateOption(a)
}

type opentracingPluginOption struct {
	opentracingPlugin *OpentracingPlugin
}

func (o opentracingPluginOption) apply(opts *options) {
	opts.opentracingPlugin = o.opentracingPlugin
}

func WithOpentracingPlugin(op *OpentracingPlugin) Option {
	return &opentracingPluginOption{opentracingPlugin: op}
}

type autoMigrateDstOption []any

func (a autoMigrateDstOption) apply(opts *options) {
	opts.autoMigrateDst = a
}

func WithAutoMigrateDst(models []any) Option {
	return autoMigrateDstOption(models)
}
