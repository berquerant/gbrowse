package mcp_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/berquerant/gbrowse/mcp"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMCPServer(t *testing.T) {
	ctx := context.Background()

	server := mcp.NewMCPServer(mcp.ServerConfig{
		GitCommand: "git",
	})

	serverTransport, clientTransport := mcpsdk.NewInMemoryTransports()

	go func() {
		_ = server.Run(ctx, serverTransport)
	}()

	client := mcpsdk.NewClient(&mcpsdk.Implementation{
		Name: "test-client",
	}, nil)

	clientSession, err := client.Connect(ctx, clientTransport, nil)
	require.NoError(t, err)
	defer clientSession.Close()

	parseResult := func(res *mcpsdk.CallToolResult) (mcp.URLToolResult, error) {
		var toolResult mcp.URLToolResult
		if res.StructuredContent != nil {
			b, err := json.Marshal(res.StructuredContent)
			if err != nil {
				return toolResult, err
			}
			if err := json.Unmarshal(b, &toolResult); err != nil {
				return toolResult, err
			}
			return toolResult, nil
		}
		if len(res.Content) > 0 {
			if textContent, ok := res.Content[0].(*mcpsdk.TextContent); ok {
				if err := json.Unmarshal([]byte(textContent.Text), &toolResult); err != nil {
					return toolResult, err
				}
				return toolResult, nil
			}
		}
		return toolResult, nil
	}

	t.Run("list tools", func(t *testing.T) {
		res, err := clientSession.ListTools(ctx, nil)
		require.NoError(t, err)
		assert.Len(t, res.Tools, 1)
		assert.Equal(t, "gbrowse_url", res.Tools[0].Name)
	})

	t.Run("call tool mutually exclusive params returns error", func(t *testing.T) {
		res, err := clientSession.CallTool(ctx, &mcpsdk.CallToolParams{
			Name: "gbrowse_url",
			Arguments: map[string]any{
				"commit":  true,
				"compare": "main",
			},
		})
		require.NoError(t, err)
		assert.True(t, res.IsError)
	})

	t.Run("call tool successfully in git repo", func(t *testing.T) {
		res, err := clientSession.CallTool(ctx, &mcpsdk.CallToolParams{
			Name: "gbrowse_url",
			Arguments: map[string]any{
				"target": "README.md:10",
			},
		})
		require.NoError(t, err)
		assert.False(t, res.IsError)

		toolResult, err := parseResult(res)
		require.NoError(t, err)
		assert.Contains(t, toolResult.URL, "github.com")
		assert.Contains(t, toolResult.URL, "/blob/")
		assert.Contains(t, toolResult.URL, "#L10")
	})
}
