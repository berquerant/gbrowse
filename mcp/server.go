package mcp

import (
	"context"
	"fmt"

	"github.com/berquerant/gbrowse/ctxlog"
	"github.com/berquerant/gbrowse/git"
	"github.com/berquerant/gbrowse/parse"
	"github.com/berquerant/gbrowse/urlx"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// URLToolParams defines the input parameters for the gbrowse_url tool.
//
// Target specifies a file path (optionally with :LINUM) or directory path within
// the repository. If Commit is true, Target is interpreted as the commit hash
// to view (or current commit if empty). If Compare is set, Target is used as
// the target ref to compare against (or current commit if empty).
type URLToolParams struct {
	Target  string `json:"target,omitempty" jsonschema:"target file path or FILE:LINUM or directory (default: repo root)"`
	Commit  bool   `json:"commit,omitempty" jsonschema:"open commit page (use target as commit hash or current commit if empty)"`
	Compare string `json:"compare,omitempty" jsonschema:"open compare page between the specified ref and target (or current commit if target is empty)"`
}

// URLToolResult defines the output of the gbrowse_url tool containing the generated repository URL.
type URLToolResult struct {
	URL string `json:"url" jsonschema:"the generated repository URL"`
}

// ServerConfig configures the MCP server, specifying the git command executable
// and optional logger for debug/error tracing.
type ServerConfig struct {
	GitCommand string
	Logger     ctxlog.Logger
}

// NewMCPServer creates and returns a configured MCP Server exposing the gbrowse_url tool.
func NewMCPServer(cfg ServerConfig) *mcp.Server {
	s := mcp.NewServer(
		&mcp.Implementation{
			Name:    "gbrowse",
			Version: "dev",
		},
		nil,
	)

	tool := &mcp.Tool{
		Name:        "gbrowse_url",
		Description: "Generate a repository URL for a file/directory path, commit, or comparison in the current git repository.",
	}

	mcp.AddTool(s, tool, func(ctx context.Context, req *mcp.CallToolRequest, params URLToolParams) (*mcp.CallToolResult, URLToolResult, error) {
		res, err := runURLTool(ctx, cfg, params)
		if err != nil {
			return nil, URLToolResult{}, err
		}
		return nil, res, nil
	})

	return s
}

// Run starts the MCP server over standard input and output (stdio).
// It blocks until the transport is closed or the context is cancelled.
func Run(ctx context.Context, cfg ServerConfig) error {
	s := NewMCPServer(cfg)
	t := &mcp.StdioTransport{}
	return s.Run(ctx, t)
}

func runURLTool(ctx context.Context, cfg ServerConfig, params URLToolParams) (URLToolResult, error) {
	if cfg.Logger != nil {
		ctx = ctxlog.With(ctx, cfg.Logger)
	}

	if params.Commit && params.Compare != "" {
		return URLToolResult{}, fmt.Errorf("commit and compare cannot be used together")
	}

	gitCmdName := cfg.GitCommand
	if gitCmdName == "" {
		gitCmdName = "git"
	}
	gitCommand := git.New(git.WithGitCommand(gitCmdName))

	var (
		targetURL string
		err       error
	)

	switch {
	case params.Commit:
		targetURL, err = urlx.BuildCommit(ctx, gitCommand, params.Target)
	case params.Compare != "":
		targetURL, err = urlx.BuildCompare(ctx, gitCommand, params.Compare, params.Target)
	default:
		var target *parse.Target
		target, err = parse.ReadTarget(params.Target)
		if err != nil {
			return URLToolResult{}, fmt.Errorf("parse target: %w", err)
		}
		targetURL, err = urlx.Build(ctx, gitCommand, target)
	}

	if err != nil {
		return URLToolResult{}, fmt.Errorf("build url: %w", err)
	}

	return URLToolResult{
		URL: targetURL,
	}, nil
}
