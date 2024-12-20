package logger

import (
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"os"
)

type LoggerConfig struct {
	Level  string            `json:"level" mapstructure:"level"`
	Format string            `json:"format" mapstructure:"format"` // TEXT, JSON, CONSOLE
	Output string            `json:"output" mapstructure:"output"`
	Layers map[string]string `json:"layers" mapstructure:"layers"`
	Models []string          `json:"models" mapstructure:"models"`
}

func NewLoggerConfig(fields any) Logger {

	b, err := json.Marshal(fields)
	if err != nil {
		log.Fatalln(err)
	}

	var cfg LoggerConfig
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		log.Fatalln(err)
	}

	var options slog.HandlerOptions
	options.Level = setLevel(cfg.Level)
	output := setOutput(cfg.Output)
	handler := setHandler(cfg.Format, output, options)

	l := slog.New(handler)
	slog.SetDefault(l)

	return &SLogger{
		Log: l,
	}
}

func setLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}

func setOutput(output string) io.Writer {
	switch output {
	case "stdout":
		return os.Stdout
	case "stderr":
		return os.Stderr
	default:
		return os.Stdout
	}
}

func setHandler(format string, output io.Writer, options slog.HandlerOptions) slog.Handler {
	switch format {
	case "text":
		return slog.NewTextHandler(output, &options)
	default:
		return slog.NewJSONHandler(output, &options)
	}
}
func (l *LoggerConfig) GetLevel() string {
	return l.Level
}
