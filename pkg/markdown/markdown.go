package markdown

import (
    "regexp"
    "strings"

    "github.com/gomarkdown/markdown"
    "github.com/gomarkdown/markdown/html"
    "github.com/gomarkdown/markdown/parser"
)

// Markdown implements the MarkdownProcessor interface
type Markdown struct{}

// ToHTML converts markdown text to HTML
func (m *Markdown) ToHTML(input string) string {
    // Create markdown parser with limited extensions
    // Only enable bullet lists, images, links, emphasis (strong and italics)
    extensions := parser.NoIntraEmphasis |
        parser.FencedCode |
        parser.Tables |
        parser.Autolink |
        parser.Strikethrough |
        parser.SpaceHeadings

    p := parser.NewWithExtensions(extensions)

    // Parse the markdown text 
    nodes := p.Parse([]byte(input))

    // Create HTML renderer with specific rendering flags
    htmlFlags := html.CommonFlags
    // Filter out heading tags by using a custom renderer if needed

    opts := html.RendererOptions{Flags: htmlFlags}
    renderer := html.NewRenderer(opts)

    return string(markdown.Render(nodes, renderer))
}

// FromHTML converts HTML text to Markdown (best effort for migration)
func (m *Markdown) FromHTML(input string) string {
    // Basic HTML to Markdown conversion for the limited set of tags we support
    // Replace paragraph tags
    output := strings.ReplaceAll(input, "<p>", "")
    output = strings.ReplaceAll(output, "</p>", "\n\n")
    
    // Replace strong/bold tags
    output = strings.ReplaceAll(output, "<strong>", "**")
    output = strings.ReplaceAll(output, "</strong>", "**")
    output = strings.ReplaceAll(output, "<b>", "**")
    output = strings.ReplaceAll(output, "</b>", "**")
    
    // Replace italic tags
    output = strings.ReplaceAll(output, "<em>", "*")
    output = strings.ReplaceAll(output, "</em>", "*")
    output = strings.ReplaceAll(output, "<i>", "*")
    output = strings.ReplaceAll(output, "</i>", "*")
    
    // Replace unordered list
    output = strings.ReplaceAll(output, "<ul>", "")
    output = strings.ReplaceAll(output, "</ul>", "\n")
    output = strings.ReplaceAll(output, "<li>", "* ")
    output = strings.ReplaceAll(output, "</li>", "\n")
    
    // Replace links
    linkRegex := regexp.MustCompile(`<a href="([^"]+)"[^>]*>([^<]+)</a>`)
    output = linkRegex.ReplaceAllString(output, "[$2]($1)")
    
    // Replace images
    imgRegex := regexp.MustCompile(`<img src="([^"]+)"[^>]*alt="([^"]*)"[^>]*>`)
    output = imgRegex.ReplaceAllString(output, "![$2]($1)")
    
    // Clean up extra whitespace
    output = strings.TrimSpace(output)
    
    return output
}