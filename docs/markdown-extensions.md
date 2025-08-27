# Markdown Extensions Documentation

This document explains the custom markdown extensions supported by the Musings blog engine.

## Overview

The Musings blog engine now supports enhanced markdown parsing through the `github.com/gomarkdown/markdown` library with extended features enabled. These extensions provide additional formatting options beyond standard markdown. Code syntax highlighting is provided by Prism.js, and math expressions are rendered using MathJax.

## Supported Extensions

### 1. Tables

Create tables using the standard markdown table syntax:

```markdown
| Header 1 | Header 2 | Header 3 |
|----------|----------|----------|
| Cell 1   | Cell 2   | Cell 3   |
| Cell 4   | Cell 5   | Cell 6   |
```

### 2. Footnotes

Add footnotes to your posts using the following syntax:

```markdown
This is a text with a footnote[^1].

[^1]: This is the footnote content.
```

### 3. Fenced Code Blocks

Enhanced code block support with syntax highlighting:

``````markdown
```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, world!")
}
```
``````

**Note**: Syntax highlighting is provided by Prism.js and supports many programming languages automatically.

### 4. Auto-generated Heading IDs

Headers automatically get ID attributes based on their content, allowing for anchor links:

```markdown
# My Section

This creates a heading with an ID attribute that can be linked to directly.
```

### 5. Strikethrough

Add strikethrough formatting to text:

```markdown
This text has ~~strikethrough~~ formatting.
```

### 6. Definition Lists

Create definition lists with terms and definitions:

```markdown
Term 1
:   Definition 1

Term 2
:   Definition 2
    Second paragraph of definition 2
```

### 7. Math/LaTeX Rendering

Render mathematical expressions using LaTeX syntax with MathJax:

```markdown
Inline math: $E = mc^2$

Block math:

$$
\int_{-\infty}^{\infty} e^{-x^2} dx = \sqrt{\pi}
$$
```

**Note**: Math expressions are rendered visually using MathJax JavaScript library.

### 8. Backslash Line Breaks

Force line breaks using backslash:

```markdown
This line has a backslash line break\
This should be on a new line.
```

### 9. Smart Fractions

Automatic rendering of common fractions:

```markdown
1/2 1/4 3/4
```

## Usage

The enhanced markdown parsing is automatically applied to all blog posts when they are processed. No additional configuration is required. Syntax highlighting and math expressions will be rendered visually when the blog is viewed in a web browser.

## Examples

See `posts/test-extensions.md` for a complete example of all supported extensions.

## Notes

- All extensions maintain backward compatibility with standard markdown
- Syntax highlighting and math rendering require JavaScript to be enabled in the browser
- Prism.js automatically detects programming languages from the code block language identifier
- MathJax and Prism.js are loaded asynchronously to avoid blocking page rendering
- Some features like MathJax require specific CSS styling for optimal rendering
- The enhanced parser is more forgiving of spacing issues in markdown syntax
