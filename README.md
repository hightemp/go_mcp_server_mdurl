# go_mcp_server_mdurl

A simple MCP (Model Context Protocol) server that provides tools for converting web content to Markdown format.

Server url: `http://127.0.0.1:8888/sse`

## Tools

1. `markdown_content_of_url` - Extracts the main article content from a URL and converts it to Markdown
2. `markdown_all_html_of_url` - Converts the entire HTML content from a URL to Markdown

## Usage

The server can operate in two modes: stdio and sse (Server-Sent Events). By default, it uses the sse mode.

```bash
go_mcp_server_mdurl -t sse -h 0.0.0.0 -p 8888
# or
go_mcp_server_mdurl -t stdio
```