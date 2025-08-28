# AWS CDK Infrastructure for Musings

This directory contains the AWS Cloud Development Kit (CDK) code to deploy the Musings static blog to AWS infrastructure.

## Architecture Overview

The CDK stack creates the following AWS resources:

1. **S3 Bucket** - Hosts the static website files
2. **CloudFront Distribution** - CDN for global content delivery
3. **Route53 Records** - DNS records for the custom domain
4. **SSL Certificate** - Managed by AWS Certificate Manager
5. **Deployment Pipeline** - CI/CD pipeline for automatic deployments
6. **Monitoring** - CloudWatch alarms for site health and performance

## Prerequisites

1. Install AWS CDK:
   ```bash
   npm install -g aws-cdk
   ```

2. Install project dependencies:
   ```bash
   npm install
   ```

3. Configure AWS credentials:
   ```bash
   aws configure
   ```

## Deployment

To deploy the infrastructure:

1. Synthesize the CloudFormation template:
   ```bash
   cdk synth
   ```

2. Deploy the stack:
   ```bash
   cdk deploy
   ```

## Configuration

The stack can be configured through context values in `cdk.json` or via command line parameters:

- `domainName` - Your custom domain name
- `siteSubDomain` - Subdomain for the site (optional)
- `repoOwner` - GitHub repository owner
- `repoName` - GitHub repository name

Example:
```bash
cdk deploy -c domainName=example.com -c siteSubDomain=blog -c repoOwner=yourname -c repoName=musings
```

## Stacks

- `StaticSiteStack` - Main stack that creates the S3 bucket, CloudFront distribution, and Route53 records

## Current Implementation Status

✅ **8.1 Implement AWS CDK configurations for AWS S3 hosting**  
✅ **8.2 Add CloudFront CDN integration**  
🔲 8.3 Implement Route53 DNS configuration  
🔲 8.4 Add SSL certificate management  
✅ 8.5 Implement deployment pipelines  
✅ 8.6 Add monitoring and logging