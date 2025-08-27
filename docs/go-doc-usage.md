# Musings Documentation

This document explains how to use the documentation features of the Musings project.

## Go Documentation

The Musings project now includes proper Go documentation that works with the `go doc` tool. You can explore the documentation using the following commands:

### Package Documentation

To view documentation for a specific package:

```bash
go doc ./internal/blog        # Blog package documentation
go doc ./internal/site        # Site package documentation
go doc ./internal/template    # Template package documentation
```

### Type Documentation

To view documentation for specific types:

```bash
go doc ./internal/blog.Post           # Post struct documentation
go doc ./internal/blog.Blog           # Blog struct documentation
go doc ./internal/site.StaticSiteGenerator  # StaticSiteGenerator struct documentation
```

### Function Documentation

To view documentation for specific functions:

```bash
go doc ./internal/blog.ConvertMarkdownToHTML  # ConvertMarkdownToHTML function
go doc ./internal/blog.NewBlog                # NewBlog function
go doc ./internal/site.NewStaticSiteGenerator # NewStaticSiteGenerator function
```

## Markdown Extensions Documentation

For information about the supported markdown extensions, see [docs/markdown-extensions.md](docs/markdown-extensions.md).

## Implementation Details

For detailed implementation information, see [IMPLEMENTATION.md](IMPLEMENTATION.md).