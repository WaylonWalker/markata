---
title: "Deploy to AWS S3"
description: "Step-by-step guide to deploy markata-go sites on AWS S3 with CloudFront"
date: 2026-01-24
published: true
tags:
  - documentation
  - deployment
  - aws-s3
---


Amazon S3 provides highly durable object storage that's ideal for hosting static websites. Combined with CloudFront CDN, it offers enterprise-grade performance and security.

## Prerequisites

- An AWS account
- AWS CLI installed and configured
- Your markata-go site ready to build

## Cost

| Service | Free Tier | After Free Tier |
|---------|-----------|-----------------|
| S3 Storage | 5 GB (12 months) | ~$0.023/GB/month |
| S3 Requests | 20,000 GET (12 months) | ~$0.0004/1000 requests |
| CloudFront | 1 TB transfer (12 months) | ~$0.085/GB |
| Route 53 | N/A | $0.50/hosted zone/month |

A typical blog costs $1-5/month after the free tier expires.

## Architecture Options

| Setup | Complexity | HTTPS | Custom Domain | CDN |
|-------|------------|-------|---------------|-----|
| S3 only | Low | No | Limited | No |
| S3 + CloudFront | Medium | Yes | Yes | Yes |
| S3 + CloudFront + Route 53 | High | Yes | Yes (apex) | Yes |

This guide covers the recommended S3 + CloudFront setup.

## Step 1: Create S3 Bucket

### Using AWS Console

1. Go to [S3 Console](https://s3.console.aws.amazon.com/)
2. Click **Create bucket**
3. Enter bucket name (e.g., `my-site-bucket`)
4. Choose a region close to your audience
5. Uncheck **Block all public access** (we'll use CloudFront for access)
6. Click **Create bucket**

### Using AWS CLI

```bash
# Create bucket
aws s3 mb s3://my-site-bucket --region us-east-1

# Enable static website hosting
aws s3 website s3://my-site-bucket --index-document index.html --error-document 404.html
```

## Step 2: Configure Bucket Policy

Create a bucket policy to allow CloudFront access:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "AllowCloudFrontAccess",
      "Effect": "Allow",
      "Principal": {
        "Service": "cloudfront.amazonaws.com"
      },
      "Action": "s3:GetObject",
      "Resource": "arn:aws:s3:::my-site-bucket/*",
      "Condition": {
        "StringEquals": {
          "AWS:SourceArn": "arn:aws:cloudfront::ACCOUNT_ID:distribution/DISTRIBUTION_ID"
        }
      }
    }
  ]
}
```

You'll update this policy after creating the CloudFront distribution.

## Step 3: Create CloudFront Distribution

### Using AWS Console

1. Go to [CloudFront Console](https://console.aws.amazon.com/cloudfront/)
2. Click **Create distribution**
3. Configure origin:
   - **Origin domain**: Select your S3 bucket
   - **Origin access**: Origin access control settings (recommended)
   - Click **Create control setting** > Create
4. Configure default cache behavior:
   - **Viewer protocol policy**: Redirect HTTP to HTTPS
   - **Allowed HTTP methods**: GET, HEAD
   - **Cache policy**: CachingOptimized
5. Configure settings:
   - **Price class**: Use all edge locations (or choose based on audience)
   - **Default root object**: `index.html`
6. Click **Create distribution**

### Update S3 Bucket Policy

After creating the distribution, update the bucket policy with the correct ARN:

```bash
# Get distribution ARN from CloudFront console
# Update bucket policy with the distribution ARN
```

## Step 4: Build and Deploy

### Build Your Site

```bash
# Set the CloudFront URL (or custom domain)
export MARKATA_GO_URL=https://d1234abcd.cloudfront.net

# Build
markata-go build --clean
```

### Deploy to S3

```bash
# Sync all files
aws s3 sync public/ s3://my-site-bucket --delete

# With cache headers for static assets
aws s3 sync public/static/ s3://my-site-bucket/static/ \
  --cache-control "public, max-age=31536000, immutable"

# HTML files with no-cache
aws s3 sync public/ s3://my-site-bucket \
  --exclude "static/*" \
  --cache-control "public, max-age=0, must-revalidate"
```

### Invalidate CloudFront Cache

After deploying, invalidate the cache to see changes immediately:

```bash
aws cloudfront create-invalidation \
  --distribution-id YOUR_DISTRIBUTION_ID \
  --paths "/*"
```

## Step 5: Custom Domain Setup

### Option A: Using Route 53 (Recommended)

1. **Create Hosted Zone**
   ```bash
   aws route53 create-hosted-zone --name example.com --caller-reference $(date +%s)
   ```

2. **Update Domain Nameservers**
   
   Get the nameservers from Route 53 and update them at your registrar.

3. **Request SSL Certificate**
   
   In AWS Certificate Manager (ACM) - **must be in us-east-1 for CloudFront**:
   ```bash
   aws acm request-certificate \
     --domain-name example.com \
     --subject-alternative-names "*.example.com" \
     --validation-method DNS \
     --region us-east-1
   ```

4. **Add Certificate to CloudFront**
   
   Update your distribution:
   - **Alternate domain names**: `example.com`, `www.example.com`
   - **Custom SSL certificate**: Select your ACM certificate

5. **Create DNS Records**
   
   In Route 53, create A records as aliases to your CloudFront distribution.

### Option B: External DNS

1. Request ACM certificate (as above)
2. Validate via DNS by adding CNAME records at your registrar
3. Add alternate domain names to CloudFront
4. Create CNAME record pointing to your CloudFront domain:
   ```
   www.example.com -> d1234abcd.cloudfront.net
   ```

Note: Apex domains (example.com without www) require Route 53 or a DNS provider that supports ALIAS/ANAME records.

## Automation with GitHub Actions

Create `.github/workflows/deploy.yml`:

```yaml
name: Deploy to AWS

on:
  push:
    branches: [main]

env:
  AWS_REGION: us-east-1
  S3_BUCKET: my-site-bucket
  CLOUDFRONT_DISTRIBUTION_ID: E1234567890ABC

jobs:
  deploy:
    runs-on: ubuntu-latest
    permissions:
      id-token: write
      contents: read
    
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'

      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v4
        with:
          role-to-assume: arn:aws:iam::ACCOUNT_ID:role/GitHubActionsRole
          aws-region: ${{ env.AWS_REGION }}

      - name: Build
        run: |
          go install github.com/WaylonWalker/markata-go/cmd/markata-go@latest
          markata-go build --clean
        env:
          MARKATA_GO_URL: https://example.com

      - name: Deploy to S3
        run: |
          # Static assets with long cache
          aws s3 sync public/static/ s3://${{ env.S3_BUCKET }}/static/ \
            --cache-control "public, max-age=31536000, immutable"
          
          # Everything else
          aws s3 sync public/ s3://${{ env.S3_BUCKET }} \
            --exclude "static/*" \
            --cache-control "public, max-age=0, must-revalidate" \
            --delete

      - name: Invalidate CloudFront
        run: |
          aws cloudfront create-invalidation \
            --distribution-id ${{ env.CLOUDFRONT_DISTRIBUTION_ID }} \
            --paths "/*"
```

### IAM Role for GitHub Actions

Create an IAM role with OIDC trust for GitHub:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Federated": "arn:aws:iam::ACCOUNT_ID:oidc-provider/token.actions.githubusercontent.com"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringEquals": {
          "token.actions.githubusercontent.com:aud": "sts.amazonaws.com"
        },
        "StringLike": {
          "token.actions.githubusercontent.com:sub": "repo:YOUR_ORG/YOUR_REPO:*"
        }
      }
    }
  ]
}
```

Attach a policy allowing S3 and CloudFront access:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:PutObject",
        "s3:GetObject",
        "s3:DeleteObject",
        "s3:ListBucket"
      ],
      "Resource": [
        "arn:aws:s3:::my-site-bucket",
        "arn:aws:s3:::my-site-bucket/*"
      ]
    },
    {
      "Effect": "Allow",
      "Action": "cloudfront:CreateInvalidation",
      "Resource": "arn:aws:cloudfront::ACCOUNT_ID:distribution/DISTRIBUTION_ID"
    }
  ]
}
```

## Troubleshooting

### Access Denied Errors

1. Check bucket policy allows CloudFront access
2. Verify Origin Access Control is configured
3. Ensure the distribution ARN in bucket policy is correct

### 404 on Subpages

Configure CloudFront error pages:

1. Go to your distribution > **Error pages**
2. Create custom error response:
   - HTTP error code: 403
   - Response page path: `/index.html`
   - HTTP response code: 200

Or use a CloudFront Function for clean URLs (see below).

### CloudFront Function for Clean URLs

Create a function to handle clean URLs:

```javascript
function handler(event) {
  var request = event.request;
  var uri = request.uri;
  
  // Check whether the URI is missing a file extension
  if (!uri.includes('.')) {
    // If URI doesn't end with /, add it
    if (!uri.endsWith('/')) {
      uri += '/';
    }
    // Append index.html
    request.uri = uri + 'index.html';
  }
  
  return request;
}
```

Attach to your distribution's viewer request.

### Changes Not Appearing

1. Wait for CloudFront cache to expire, or
2. Create a cache invalidation:
   ```bash
   aws cloudfront create-invalidation \
     --distribution-id YOUR_DISTRIBUTION_ID \
     --paths "/*"
   ```

### SSL Certificate Issues

- ACM certificate must be in **us-east-1** region
- Certificate must be validated (check ACM console)
- Alternate domain names must match certificate

## Cost Optimization

### Reduce CloudFront Costs

- Use **Price Class 100** (North America and Europe only) for regional sites
- Enable **Compress objects automatically**
- Set appropriate cache TTLs

### Reduce S3 Costs

- Enable **S3 Intelligent-Tiering** for rarely accessed sites
- Use **S3 Standard** for frequently accessed content
- Clean up old versions if versioning is enabled

### Monitor Costs

Set up AWS Budgets to alert you:

```bash
aws budgets create-budget \
  --account-id YOUR_ACCOUNT_ID \
  --budget file://budget.json \
  --notifications-with-subscribers file://notifications.json
```

## Next Steps

- [AWS S3 Docs](https://docs.aws.amazon.com/s3/) - Official S3 documentation
- [CloudFront Docs](https://docs.aws.amazon.com/cloudfront/) - CDN configuration
- [Configuration Guide](../configuration/) - Customize your markata-go site
