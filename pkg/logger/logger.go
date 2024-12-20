package logger

import (
	"fmt"
	"log/slog"
	"runtime"
	"strings"
)

type Logger interface {
	Info(msg string, attrs ...Attr)
	Debug(msg string, attrs ...Attr)
	Warn(lm *Message)
	Error(*Message) error
	With(attrs ...slog.Attr)
}

type SLogger struct {
	Log *slog.Logger
}

type Value struct {
	Str string
	Bol bool
	Int int
	Any any
}
type Attr struct {
	Key   string
	Value Value
	// Code  ResponseCode
}

type Message struct {
	Code ResponseCode `json:"code"`
	Body string       `json:"body"`
}

type Content struct {
	Module string `json:"module"`
	Method string `json:"method"`
	Layer  string `json:"layer"`
	Message
}

func (l *SLogger) Info(msg string, attrs ...Attr) {

	var attrsLog []slog.Attr
	for _, attr := range attrs {
		attrsLog = append(attrsLog, slog.String(attr.Key, attr.Value.Str))
	}

	l.Log.Info(msg, attrsLog)
}

func (l *SLogger) Debug(msg string, attrs ...Attr) {
	var attrsLog []slog.Attr
	for _, attr := range attrs {
		attrsLog = append(attrsLog, slog.String(attr.Key, attr.Value.Str))
	}

	l.Log.Debug(msg, attrsLog)
}

func (l *SLogger) Warn(lm *Message) {
	content := setLayer(lm)
	logContent := setMessage("warn", content)
	logContent.Warn(lm.Body)
}

func (l *SLogger) Error(lm *Message) error {
	content := setLayer(lm)
	logContent := setMessage("error", content)
	err := setError(content)
	logContent.Error(lm.Body)

	return err
}

func (l *SLogger) With(attrs ...slog.Attr) {
	l.Log.With(attrs)
}

func setError(content *Content) error {
	return fmt.Errorf("module: %s, method: %s, layer: %s, code: %d, message: %s", content.Module, content.Method, content.Layer, content.Code, content.Body)
}

func setMessage(t string, content *Content) *slog.Logger {

	var group string
	switch strings.ToLower(t) {
	case "info":
		group = "info"
	case "debug":
		group = "debug"
	case "warn":
		group = "warn"
	default:
		group = "error"
	}

	chidren := slog.With(
		slog.Group(group,
			slog.String("module", content.Module),
			slog.String("method", content.Method),
			slog.String("layer", content.Layer),
			slog.Int("code", int(content.Code)),
			slog.String("message", content.Body),
		))
	return chidren
}

func setLayer(lm *Message) *Content {

	var content Content
	pc, file, _, _ := runtime.Caller(2)

	fn := runtime.FuncForPC(pc)
	var funcParts []string = []string{"entity"}

	if fn != nil {

		funcParts = strings.Split(fn.Name(), "/")
		funcPart := funcParts[len(funcParts)-1]
		funcParts = strings.Split(funcPart, ".")
		content.Method = funcParts[len(funcParts)-1]
	}

	fileParts := strings.Split(file, "/")
	filePart := fileParts[len(fileParts)-1]
	filePart = strings.Split(filePart, ".")[0]
	content.GetModule(&filePart)
	content.GetLayer(funcParts[0])

	content.Code = lm.Code
	content.Body = lm.Body

	return &content
}

type ModuleLayer string

const (
	ModuleLayerUser           ModuleLayer = "user"
	ModuleLayerTenant         ModuleLayer = "tenant"
	ModuleLayerPlan           ModuleLayer = "plan"
	ModuleLayerAuthorization  ModuleLayer = "authorization"
	ModuleLayerMiddleware     ModuleLayer = "middleware"
	ModuleLayerCategory       ModuleLayer = "transaction_category"
	ModuleLayerWallet         ModuleLayer = "wallet"
	ModuleLayerTransaction    ModuleLayer = "transaction"
	ModuleLayerInicialization ModuleLayer = "main"
)

func (m *Content) GetModule(module *string) {

	switch strings.ToLower(*module) {
	case "tenant":
		m.Module = string(ModuleLayerTenant)
	case "user":
		m.Module = string(ModuleLayerUser)
	case "transactioncategory":
		m.Module = string(ModuleLayerCategory)
	case "plan":
		m.Module = string(ModuleLayerPlan)
	case "authorization":
		m.Module = string(ModuleLayerAuthorization)
	case "middleware":
		m.Module = string(ModuleLayerMiddleware)
	case "wallet":
		m.Module = string(ModuleLayerWallet)
	case "transaction":
		m.Module = string(ModuleLayerTransaction)
	case "main":
		m.Module = string(ModuleLayerInicialization)
	default:
		m.Module = "unknown"
	}
}

type ApplicationLayer string

const (
	ApplicationLayerRepository     ApplicationLayer = "Data Access Layer"
	ApplicationLayerService        ApplicationLayer = "Business Logic"
	ApplicationLayerHandler        ApplicationLayer = "Controllers"
	ApplicationLayerMiddleware     ApplicationLayer = "Middleware"
	ApplicationLayerEntity         ApplicationLayer = "Model"
	ApplicationLayerInicialization ApplicationLayer = "Main"
)

func (m *Content) GetLayer(layer string) {
	switch strings.ToLower(layer) {
	case "web":
		m.Layer = string(ApplicationLayerHandler)
	case "service":
		m.Layer = string(ApplicationLayerService)
	case "repository":
		m.Layer = string(ApplicationLayerRepository)
	case "middleware":
		m.Layer = string(ApplicationLayerMiddleware)
	case "entity":
		m.Layer = string(ApplicationLayerEntity)
	case "main":
		m.Layer = string(ApplicationLayerInicialization)
	default:
		m.Layer = string(ApplicationLayerRepository)
	}
}
