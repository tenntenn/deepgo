package internal

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mark3labs/mcp-go/server"

	"github.com/tenntenn/deepgo"
	"github.com/tenntenn/deepgo/cmd/deepgomcp/internal/document"
	"github.com/tenntenn/deepgo/cmd/deepgomcp/internal/proposal"
	"github.com/tenntenn/deepgo/cmd/deepgomcp/internal/style"
	"github.com/tenntenn/deepgo/cmd/deepgomcp/internal/util"
)

type MCPServer struct {
	mcpServer   *server.MCPServer
	stdioServer *server.StdioServer
	logger      *slog.Logger
}

func New(ctx context.Context) (*MCPServer, error) {
	logger, err := newLogger(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	mcpServer := server.NewMCPServer(
		"DeepGo MCP",
		deepgo.Version,
	)
	stdioServer := server.NewStdioServer(mcpServer)
	slog.SetDefault(logger)
	slog.SetLogLoggerLevel(slog.LevelError)
	stdioServer.SetErrorLogger(log.Default())

	s := &MCPServer{
		mcpServer:   mcpServer,
		stdioServer: stdioServer,
		logger:      logger,
	}

	s.initTools()

	return s, nil
}

func newLogger(ctx context.Context) (*slog.Logger, error) {
	gopath, err := getGOPATH(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get GOPATH: %w", err)
	}

	filename := filepath.Join(gopath, "deepgo", "mcpserver.log")
	if err := os.MkdirAll(filepath.Dir(filename), 0o644); err != nil {
		return nil, fmt.Errorf("failed to create directory %q: %w", filepath.Dir(filename), err)
	}

	logfile, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o600)
	if err != nil {
		return nil, fmt.Errorf("cannot open file %q: %w", filename, err)
	}

	h := slog.NewJSONHandler(logfile, nil)
	return slog.New(h), nil
}

func getGOPATH(ctx context.Context) (string, error) {
	var stdout bytes.Buffer
	cmd := exec.CommandContext(ctx, "go", "env", "GOPATH")
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to run go env GOPATH: %w", err)
	}

	return stdout.String(), nil
}

func (s *MCPServer) initTools() {
	s.mcpServer.AddTools(
		proposal.NewReviewMeetingMinutesTool(),
		document.NewReleaseNoteTool(),
		document.NewLatestGoVersionTool(),
		style.NewModernizeTool(),
		style.NewGoStyleTool(),
		style.NewSkeleton(),
		util.NewCopyTxtarTool(),
	)
}

func (s *MCPServer) Start(ctx context.Context) error {
	s.logger.Info("MCP server listen...")
	if err := s.stdioServer.Listen(ctx, os.Stdin, os.Stdout); err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("failed to listen on server: %w", err)
	}
	return nil
}
