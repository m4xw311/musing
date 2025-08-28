# Musings Implementation Details

This document outlines the detailed functional structure of the Musings application based on the current implementation. It serves as a guide for anyone wanting to contribute to the development of the project.

## Current Implementation

### Core Components

#### 1. Command Line Interface
Using the Cobra framework, the CLI provides two main commands:

- `musings publish` - Generates a static website from markdown blog posts
- `musings sync` - Syncs blog posts to external platforms (currently a placeholder)

The CLI entry point is at `cmd/musings/main.go`, with command implementations in the `cmd/musings/cmd/` directory.

#### 2. Blog Management System
Located in `internal/blog/`, this package handles:

- **Post Struct**: Represents individual blog posts with fields for title, content, creation/update dates, tags, slug, and publication status
- **Blog Struct**: Manages a collection of posts
- **Markdown Parsing**: Converts markdown files to Post structs with:
  - Frontmatter parsing for metadata (CreatedDate, UpdatedDate, Tags, Published status)
  - Automatic title extraction from first H1 heading
  - Slug generation from titles
  - **Enhanced Markdown to HTML conversion** using the `github.com/gomarkdown/markdown` library with extended features:
    - Tables
    - Footnotes
    - Fenced code blocks
    - Auto-generated heading IDs
    - Strikethrough
    - Definition lists
    - Math/LaTeX rendering (with MathJax support)
    - Backslash line breaks
    - Smart fractions
  - Content snippet generation for post previews
  - Automatic addition of missing date fields to markdown files

#### 3. Static Site Generation
Located in `internal/site/`, this package handles:

- **HTML Generation**: Creates index and individual post pages using Go templates
- **Asset Handling**: Copies CSS stylesheets and image assets to the output directory
- **Feed Generation**: Creates both RSS and Atom feeds for the blog content
- **Directory Structure**:
  - `posts/` - Source directory for markdown blog posts
  - `public/` - Output directory for generated static site

#### 4. Template System
Located in `internal/template/`, includes:

- `index.html` - Main page showing latest posts and a list of all posts
- `post.html` - Template for individual blog post pages
- `style.css` - Styling for the generated site with responsive design
- Feed templates implemented in code for RSS and Atom formats
- **MathJax support** for rendering mathematical expressions
- **Prism.js support** for syntax highlighting of code blocks

### AWS Infrastructure (Newly Implemented)

Located in `infrastructure/aws-cdk/`, this package handles:

- **AWS CDK Infrastructure** for deploying the static site to AWS
- **S3 Bucket** for hosting static website files
- **CloudFront CDN** for global content delivery
- **Deployment Pipeline** for CI/CD using AWS CodePipeline
- **Monitoring and Alerting** with CloudWatch alarms
- **HTTPS Support** with ACM certificates
- **Custom Domain Support** with Route53 DNS integration

### Code Structure

```
.
├── cmd/
│   └── musings/
│       ├── main.go           # Entry point
│       └── cmd/              # Command implementations
│           ├── publish.go
│           ├── root.go
│           └── sync.go
├── internal/
│   ├── blog/                 # Blog management
│   │   └── blog.go
│   ├── site/                 # Static site generation
│   │   ├── site.go
│   │   └── feed.go
│   └── template/             # HTML templates and CSS
│       ├── index.html
│       ├── post.html
│       └── style.css
├── infrastructure/
│   └── aws-cdk/              # AWS CDK infrastructure code
│       ├── bin/
│       ├── lib/
│       ├── test/
│       ├── package.json
│       ├── cdk.json
│       ├── tsconfig.json
│       └── README.md
├── posts/                    # Source markdown files
└── public/                   # Generated static site
```

### Key Features Implemented

1. **Automatic Metadata Handling**:
   - Parses frontmatter in markdown files for metadata
   - Automatically adds missing CreatedDate/UpdatedDate fields
   - Extracts titles from H1 headings if not in frontmatter

2. **Complete Static Site Generation**:
   - Generates index page with latest posts and full post listing
   - Creates individual HTML pages for each post
   - Produces RSS and Atom feeds for syndication
   - Copies static assets (CSS, images)

3. **Responsive Design**:
   - Mobile-friendly CSS styling
   - Card-based layout for latest posts
   - Proper typography and spacing

4. **Enhanced Markdown Support**:
   - Extended markdown parsing with support for tables, footnotes, and more
   - Full markdown parsing to HTML
   - Math/LaTeX rendering with MathJax
   - Syntax highlighting with Prism.js
   - Image handling with proper sizing

5. **AWS Infrastructure**:
   - Infrastructure-as-code using AWS CDK
   - S3 hosting with CloudFront CDN
   - CI/CD pipeline with GitHub integration
   - Monitoring and alerting with CloudWatch
   - HTTPS support with ACM certificates

## Getting Started for Contributors

1. **Understanding the Data Flow**:
   - CLI commands initialize Blog and StaticSiteGenerator instances
   - Blog.LoadPosts() reads and parses markdown files
   - StaticSiteGenerator.Generate() orchestrates the site creation

2. **Key Implementation Details**:
   - Frontmatter parsing expects YAML-like format with specific field names
   - Date formats must follow "YYYY-MM-DD HH:MM:SS"
   - Slug generation automatically creates URL-friendly identifiers
   - Markdown parsing now supports extended syntax features including MathJax and Prism.js

3. **Adding New Features**:
   - New CLI commands should be added in `cmd/musings/cmd/`
   - Blog-related features go in `internal/blog/`
   - Site generation enhancements belong in `internal/site/`
   - Templates can be modified in `internal/template/`
   - Infrastructure enhancements belong in `infrastructure/aws-cdk/`

4. **AWS Deployment**:
   - Requires AWS credentials configured via `aws configure`
   - Deploy with `cdk deploy` from the `infrastructure/aws-cdk` directory
   - Supports custom domains with automatic DNS and SSL configuration