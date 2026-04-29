package markdown

import (
    "strings"
    "testing"
)

func TestToHTML(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {
            name:     "Convert plain text",
            input:    "This is plain text",
            expected: "<p>This is plain text</p>\n",
        },
        {
            name:     "Convert bold text",
            input:    "This is **bold** text",
            expected: "<p>This is <strong>bold</strong> text</p>\n",
        },
        {
            name:     "Convert italic text",
            input:    "This is *italic* text",
            expected: "<p>This is <em>italic</em> text</p>\n",
        },
        {
            name:     "Convert bullet list",
            input:    "* Item 1\n* Item 2\n* Item 3",
            expected: "<ul>\n<li>Item 1</li>\n<li>Item 2</li>\n<li>Item 3</li>\n</ul>\n",
        },
        {
            name:     "Convert link",
            input:    "[Link text](https://example.com)",
            expected: "<p><a href=\"https://example.com\">Link text</a></p>\n",
        },
        {
            name:     "Convert image",
            input:    "![Alt text](https://example.com/image.jpg)",
            expected: "<p><img src=\"https://example.com/image.jpg\" alt=\"Alt text\" /></p>\n",
        },
    }

    markdown := &Markdown{}
    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            result := markdown.ToHTML(test.input)
            if result != test.expected {
                t.Errorf("Expected %q, got %q", test.expected, result)
            }
        })
    }
}

func TestFromHTML(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {
            name:     "Convert paragraph",
            input:    "<p>This is a paragraph</p>",
            expected: "This is a paragraph",
        },
        {
            name:     "Convert multiple paragraphs",
            input:    "<p>Paragraph 1</p><p>Paragraph 2</p>",
            expected: "Paragraph 1\n\nParagraph 2",
        },
        {
            name:     "Convert bold",
            input:    "<p>This is <strong>bold</strong> text</p>",
            expected: "This is **bold** text",
        },
        {
            name:     "Convert italic",
            input:    "<p>This is <em>italic</em> text</p>",
            expected: "This is *italic* text",
        },
        {
            name:     "Convert list",
            input:    "<ul><li>Item 1</li><li>Item 2</li><li>Item 3</li></ul>",
            expected: "* Item 1\n* Item 2\n* Item 3",
        },
        {
            name:     "Convert link",
            input:    "<p><a href=\"https://example.com\">Link text</a></p>",
            expected: "[Link text](https://example.com)",
        },
        {
            name:     "Convert image",
            input:    "<p><img src=\"https://example.com/image.jpg\" alt=\"Alt text\"></p>",
            expected: "![Alt text](https://example.com/image.jpg)",
        },
    }

    markdown := &Markdown{}
    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            result := markdown.FromHTML(test.input)
            // Trim to handle any extra whitespace differences
            if strings.TrimSpace(result) != strings.TrimSpace(test.expected) {
                t.Errorf("Expected %q, got %q", test.expected, result)
            }
        })
    }
}