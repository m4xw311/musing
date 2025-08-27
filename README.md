# Musings - A Static Blog Publisher

A simple tool to publish markdown-based static blogs and sync them to platforms like Substack and Medium.

⚠️ This project is currently in development and not ready for use. Currently using the development of this project to test and improve [compell](https://github.com/m4xw311/compell) coding assistant.

## Features

1. Publish blog posts written in markdown to a static website
2. Sync blog posts to external platforms (Substack, Medium)

## Project Structure

```
├── cmd/musings/       # Main application entry point
├── internal/
│   ├── blog/          # Blog management functionality
│   └── site/          # Static site generation
├── posts/             # Markdown blog posts
└── public/            # Generated static website
```

## Usage

```
musings publish  # Generate static website from markdown posts
musings sync     # Sync posts to external platforms
```

## Documentation

- [Implementation Details](IMPLEMENTATION.md) - Detailed information about the current implementation
- [Go Documentation](docs/go-doc-usage.md) - How to use Go documentation with this project
- [Markdown Extensions](docs/markdown-extensions.md) - Information about supported markdown features

## Requirements

See [REQUIREMENTS.md](REQUIREMENTS.md) for detailed requirements.

## Implementation Status

See [IMPLEMENTATION.md](IMPLEMENTATION.md) for current implementation status and next steps.