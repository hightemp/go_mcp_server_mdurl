package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/go-shiori/go-readability"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	var transport string
	var host string
	var port string
	flag.StringVar(&transport, "t", "sse", "Transport type (stdio or sse)")
	flag.StringVar(&host, "h", "0.0.0.0", "Host of sse server")
	flag.StringVar(&port, "p", "8890", "Port of sse server")
	flag.Parse()

	mcpServer := server.NewMCPServer(
		"go_mcp_server_mdurl",
		"1.0.0",
	)

	tool := mcp.NewTool("markdown_content_of_url",
		mcp.WithDescription("Get markdowned content of article from URL"),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("url of site"),
		),
	)

	mcpServer.AddTool(tool, helloHandler)

	tool2 := mcp.NewTool("markdown_all_html_of_url",
		mcp.WithDescription("Get markdowned content of all html from URL"),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("url of site"),
		),
	)

	mcpServer.AddTool(tool2, helloHandler2)

	if transport == "sse" {
		sseServer := server.NewSSEServer(mcpServer, server.WithBaseURL(fmt.Sprintf("http://localhost:%s", port)))
		log.Printf("SSE server listening on %s:%s URL: http://127.0.0.1:%s/sse", host, port, port)
		if err := sseServer.Start(fmt.Sprintf("%s:%s", host, port)); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	} else {
		if err := server.ServeStdio(mcpServer); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}
}

func helloHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	url, ok := request.Params.Arguments["url"].(string)
	if !ok {
		return nil, errors.New("url must be a string")
	}

	article, err := readability.FromURL(url, 30*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to extract article: %v", err)
	}

	converter := md.NewConverter("", true, nil)
	markdown, err := converter.ConvertString(article.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to convert HTML to markdown: %v", err)
	}

	result := fmt.Sprintf("# %s\n\n%s", article.Title, markdown)
	return mcp.NewToolResultText(result), nil
}

func helloHandler2(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	url, ok := request.Params.Arguments["url"].(string)
	if !ok {
		return nil, errors.New("url must be a string")
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %v", err)
	}
	defer resp.Body.Close()

	html, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read HTML: %v", err)
	}

	converter := md.NewConverter("", true, nil)
	markdown, err := converter.ConvertString(string(html))
	if err != nil {
		return nil, fmt.Errorf("failed to convert HTML to markdown: %v", err)
	}

	markdown = strings.ReplaceAll(markdown, "\n\n\n", "\n\n")

	return mcp.NewToolResultText(markdown), nil
}
