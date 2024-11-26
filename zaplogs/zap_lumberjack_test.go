package zaplogs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewLumberjackZapLog(t *testing.T) {
	cfgs := []*LumberjackZapCfg{
		NewLumberjackZapCFG(NewLumberjackConfig("stdout", "debug")),
		NewLumberjackZapCFG(NewLumberjackConfig("stderr", "error")),
	}

	{
		zapLog := NewLumberjackZapLog(cfgs, true, 0)
		zapLog.Info("123", zap.String("k", "v"))
		zapLog.Debug("abc", zap.String("k", "v"))
		zapLog.Error("xyz", zap.String("k", "v")) // will be print twice(both to stdout and stderr output)
		zapLog.Warn("uvw", zap.String("k", "v"))
	}
	{
		zapLog := NewLumberjackZapLog(cfgs, false, 0)
		zapLog.Info("123", zap.String("k", "v"))
		zapLog.Debug("abc", zap.String("k", "v"))
		zapLog.Error("xyz", zap.String("k", "v")) // will be print twice(both to stdout and stderr output)
		zapLog.Warn("uvw", zap.String("k", "v"))
	}
}

func TestNewLumberjackZapLOG(t *testing.T) {
	temp, err := os.MkdirTemp("", "zaplogs_case_simple")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.RemoveAll(temp))
	}()
	t.Log(temp)

	debugPath := filepath.Join(temp, "debug.log")
	errorPath := filepath.Join(temp, "error.log")

	cfgs := []*LumberjackZapCfg{
		NewLumberjackZapCFG(NewLumberjackConfig(debugPath, "debug")),
		NewLumberjackZapCFG(NewLumberjackConfig(errorPath, "error")),
	}
	defer func() {
		for _, cfg := range cfgs {
			require.NoError(t, cfg.Close())
		}
	}()

	zapLog := NewLumberjackZapLOG(cfgs)
	for i := 0; i < 3; i++ {
		zapLog.Info("123", zap.String("k", "v"))
		zapLog.Debug("abc", zap.String("k", "v"))
		zapLog.Error("xyz", zap.String("k", "v"))
		zapLog.Warn("uvw", zap.String("k", "v"))
	}
	require.NoError(t, zapLog.Sync())

	showContent(t, debugPath)
	showContent(t, errorPath)
}

func showContent(t *testing.T, path string) {
	t.Log("path:", path)
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	t.Log("data:", string(data))
}
