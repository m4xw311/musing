# Musings
A simple tool to publish a markdown based static blog and publish the same to platforms like substack and medium.
## Technology
1. Go
2. Terraform
## Features:
1. Publish blog posts written in markdown in a cloud hosted static website
  1.1 Allow user to write blog posts in markdown format
  1.2 Ability to render the markdown content in a static web site
  1.3 The blog posts should be in a simple directory structure
  1.4 The home page of the website should show the latest post prominently and have a list of few recent posts
  1.5 It is acceptable if some well known static site generators or serving libraries are used
  1.6 Support for bootstrapping the static website in cloud providers - Implement later
    1.6.1 Hosting
     1.6.1.1 AWS S3
     1.6.1.2 Azure, GCP etc - Implement later
    1.6.2 Bootstrap using standard IaC template in terraform
2. Sync the blog posts to cloud platforms like substack and medium
  1.1 Substack - Implement later
  1.2 Medium - Implement later
