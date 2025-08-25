# Musings - A Static Blog Publisher

A simple tool to publish markdown-based static blogs and sync them to platforms like Substack and Medium.

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

## Requirements

See [REQUIREMENTS.md](REQUIREMENTS.md) for detailed requirements.

## Implementation Status

See [IMPLEMENTATION.md](IMPLEMENTATION.md) for current implementation status and next steps.