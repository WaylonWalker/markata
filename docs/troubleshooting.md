---
title: "Troubleshooting"
description: "Solutions for common issues when working with markata-go"
date: 2024-01-15
published: true
tags:
  - documentation
  - troubleshooting
---


This guide covers common issues users encounter when working with markata-go and how to resolve them.

## Table of Contents

- [Installation Issues](#installation-issues)
- [Configuration Issues](#configuration-issues)
- [Content Issues](#content-issues)
- [Build Issues](#build-issues)
- [Feed Issues](#feed-issues)
- [Serve Issues](#serve-issues)
- [Deployment Issues](#deployment-issues)
- [Plugin Issues](#plugin-issues)
- [Getting Help](#getting-help)

---

## Installation Issues

### Go Not Installed

**Symptom:**
```
bash: go: command not found
```
Or on Windows:
```
'go' is not recognized as an internal or external command
```

**Cause:** Go is not installed on your system or not in your PATH.

**Solution:**

1. Download Go from [https://go.dev/dl/](https://go.dev/dl/)
2. Follow the installation instructions for your operating system
3. Verify installation:
   ```bash
   go version
   ```

**Example output:**
```
go version go1.22.2 linux/amd64
```

---

### Wrong Go Version

**Symptom:**
```
note: module requires Go 1.22
```
Or build errors mentioning unsupported language features.

**Cause:** markata-go requires Go 1.22 or later, but an older version is installed.

**Solution:**

1. Check your current version:
   ```bash
   go version
   ```

2. If it's older than 1.22, update Go:
   - Download the latest version from [https://go.dev/dl/](https://go.dev/dl/)
   - Remove the old version and install the new one
   - Or use a version manager like [gvm](https://github.com/moovweb/gvm) or [goenv](https://github.com/syndbg/goenv)

3. Verify the update:
   ```bash
   go version
   # Should show go1.22.x or higher
   ```

---

### PATH Issues

**Symptom:**
```
bash: markata-go: command not found
```
After running `go install` successfully.

**Cause:** The Go bin directory is not in your PATH.

**Solution:**

1. Find your Go bin directory:
   ```bash
   go env GOPATH
   # Typically: /home/username/go (Linux/macOS) or C:\Users\username\go (Windows)
   ```

2. Add the bin directory to your PATH:

   **Linux/macOS (bash/zsh):**
   ```bash
   # Add to ~/.bashrc, ~/.zshrc, or ~/.profile
   export PATH=$PATH:$(go env GOPATH)/bin

   # Reload the shell
   source ~/.bashrc  # or ~/.zshrc
   ```

   **Windows (PowerShell):**
   ```powershell
   # Add to your PowerShell profile
   $env:PATH += ";$(go env GOPATH)\bin"
   ```

   **Windows (permanently):**
   - Open System Properties > Environment Variables
   - Edit the `PATH` variable
   - Add `%USERPROFILE%\go\bin`

3. Verify markata-go is accessible:
   ```bash
   markata-go --version
   ```

---

### Installation Fails with Network Error

**Symptom:**
```
go: github.com/WaylonWalker/markata-go@latest: module lookup disabled by GOPROXY=off
```
Or timeout errors during installation.

**Cause:** Network issues or proxy misconfiguration.

**Solution:**

1. Check your proxy settings:
   ```bash
   go env GOPROXY
   ```

2. Reset to default if needed:
   ```bash
   go env -w GOPROXY=https://proxy.golang.org,direct
   ```

3. If behind a corporate proxy, configure it:
   ```bash
   export HTTP_PROXY=http://proxy.example.com:8080
   export HTTPS_PROXY=http://proxy.example.com:8080
   ```

4. Retry installation:
   ```bash
   go install github.com/WaylonWalker/markata-go/cmd/markata-go@latest
   ```

---

## Configuration Issues

### Config File Not Found

**Symptom:** markata-go uses default values when you expected it to load a site configuration.

**Cause:** No configuration file exists in the expected locations, or it has the wrong name. When no configuration file is found, markata-go continues with built-in defaults and environment-variable overrides.

**Solution:**

1. Check that a config file exists with the correct name:
   ```bash
   ls -la markata-go.toml markata-go.yaml markata-go.yml markata-go.json 2>/dev/null
   ```

2. Create a config file if missing:
   ```bash
   markata-go config init
   ```

3. Or specify a custom config path:
   ```bash
   markata-go build --config path/to/my-config.toml
   ```

**Supported local config file names (in priority order):**
- `markata-go.toml` (recommended)
- `markata-go.yaml`
- `markata-go.yml`
- `markata-go.json`

If none of those exists, markata-go also checks `~/.config/markata-go/config.toml` before falling back to built-in defaults.

See the [[configuration-guide|Configuration Guide]] for details on config file locations.

---

### Invalid TOML/YAML Syntax

**Symptom:**
```
Error: failed to parse config: toml: line 15: expected '=' after key
```
Or:
```
Error: yaml: line 10: did not find expected key
```

**Cause:** Syntax error in your configuration file.

**Solution:**

1. Validate your config file:
   ```bash
   markata-go config validate
   ```

2. Common TOML mistakes:

   **Wrong - missing quotes on strings with special characters:**
   ```toml
   title = My Site: The Best!
   ```

   **Correct:**
   ```toml
   title = "My Site: The Best!"
   ```

   **Wrong - using colons instead of equals:**
   ```toml
   title: "My Site"
   ```

   **Correct:**
   ```toml
   title = "My Site"
   ```

3. Common YAML mistakes:

   **Wrong - incorrect indentation:**
   ```yaml
   markata-go:
   title: "My Site"
   ```

   **Correct:**
   ```yaml
   markata-go:
     title: "My Site"
   ```

   **Wrong - tabs instead of spaces:**
   ```yaml
   markata-go:
   	title: "My Site"  # Tab character
   ```

   **Correct:**
   ```yaml
   markata-go:
     title: "My Site"  # Two spaces
   ```

4. Use an online validator:
   - TOML: [toml-lint.com](https://www.toml-lint.com/)
   - YAML: [yamllint.com](http://www.yamllint.com/)

---

### Unknown Configuration Options

**Symptom:**
```
Warning: unknown configuration key: output-dir
```
Or a setting doesn't seem to take effect.

**Cause:** Typo in configuration key or using the wrong format.

**Solution:**

1. Check for typos. Common mistakes:
   - `output-dir` should be `output_dir`
   - `templatesDir` should be `templates_dir`
   - Keys are case-sensitive

2. Verify the correct key names:
   ```bash
   markata-go config show
   ```

3. Reference the [[configuration-guide|Configuration Guide]] for the complete list of options.

**Example of correct configuration:**
```toml
[markata-go]
output_dir = "public"      # Correct: underscore
templates_dir = "templates"
assets_dir = "static"
```

---

### How to Validate Config

**Solution:**

Run the validation command:
```bash
markata-go config validate
```

**Example success output:**
```
Configuration is valid
```

**Example error output:**
```
Configuration errors:
  - url: URL must include a scheme (e.g., https://)
  - concurrency: must be >= 0 (0 means auto-detect)

Configuration warnings:
  - glob.patterns: no glob patterns specified, no files will be processed
```

To see the fully resolved configuration:
```bash
markata-go config show
```

To see a specific value:
```bash
markata-go config get output_dir
markata-go config get glob.patterns
```

---

## Content Issues

### Invalid Frontmatter YAML

**Symptom:**
```
Error: failed to parse frontmatter in posts/my-post.md: yaml: line 5: could not find expected ':'
```

**Cause:** YAML syntax error in the frontmatter block.

**Solution:**

1. Check the frontmatter structure:
   ```markdown
   ---
   title: "My Post"
   date: 2024-01-15
   tags:
     - go
     - tutorial
   ---
   ```

2. Common frontmatter mistakes:

   **Wrong - missing closing delimiter:**
   ```markdown
   ---
   title: "My Post"

   Content starts here...
   ```

   **Correct:**
   ```markdown
   ---
   title: "My Post"
   ---

   Content starts here...
   ```

   **Wrong - invalid date format:**
   ```markdown
   ---
   date: January 15, 2024
   ---
   ```

   **Correct:**
   ```markdown
   ---
   date: 2024-01-15
   ---
   ```

   **Wrong - unquoted special characters:**
   ```markdown
   ---
   title: My Post: A Guide
   ---
   ```

   **Correct:**
   ```markdown
   ---
   title: "My Post: A Guide"
   ---
   ```

3. Validate YAML online at [yamllint.com](http://www.yamllint.com/)

See the [[frontmatter-guide|Frontmatter Guide]] for complete frontmatter documentation.

---

### Posts Not Appearing

**Symptom:** You created a post but it doesn't appear on your site.

**Cause:** Multiple possible causes:
- Post not published
- Filtered out by feed configuration
- File not matched by glob patterns
- File in .gitignore

**Solution:**

1. **Check if the post is published:**
   ```markdown
   ---
   title: "My Post"
   published: true    # Must be true
   draft: false       # Should be false for production
   ---
   ```

2. **Check your feed filter:**
   ```toml
   [[markata-go.feeds]]
   slug = "blog"
   filter = "published == True"  # Make sure your post matches
   ```

3. **Check glob patterns:**
   ```toml
   [markata-go.glob]
   patterns = ["posts/**/*.md"]  # Does your file match this pattern?
   ```

4. **Check if the file is gitignored:**
   ```bash
   git check-ignore posts/my-post.md
   # If it returns the filename, the file is ignored
   ```

5. **Run a verbose build to see what's being processed:**
   ```bash
   markata-go build -v
   ```

---

### Drafts Showing in Production

**Symptom:** Draft posts appear on the live site.

**Cause:** Feed filters don't exclude drafts.

**Solution:**

1. Add draft filtering to your feed configuration:
   ```toml
   [[markata-go.feeds]]
   slug = "blog"
   filter = "published == True and draft == False"
   ```

2. Or filter only by published status (drafts should have `published: false`):
   ```toml
   [[markata-go.feeds]]
   slug = "blog"
   filter = "published == True"
   ```

3. Verify draft status in frontmatter:
   ```markdown
   ---
   title: "My Draft"
   published: false
   draft: true
   ---
   ```

---

### Date Format Problems

**Symptom:**
```
Error: cannot parse "January 15, 2024" as "2006-01-02"
```
Or dates display incorrectly.

**Cause:** Using an unsupported date format in frontmatter.

**Solution:**

1. Use ISO 8601 format (recommended):
   ```yaml
   date: 2024-01-15
   ```

2. Supported date formats:
   ```yaml
   date: 2024-01-15              # YYYY-MM-DD (recommended)
   date: 2024-01-15T10:30:00     # With time
   date: 2024-01-15T10:30:00Z    # With timezone
   date: "2024-01-15"            # Quoted string
   ```

3. **NOT supported:**
   ```yaml
   date: January 15, 2024        # Wrong
   date: 15/01/2024              # Wrong
   date: 01-15-2024              # Wrong (MM-DD-YYYY)
   ```

---

### Encoding Issues (UTF-8)

**Symptom:**
- Special characters display as `?` or garbled text
- Build errors mentioning encoding

**Cause:** File is not saved as UTF-8, or has a BOM (Byte Order Mark).

**Solution:**

1. **Check file encoding:**
   ```bash
   file -i posts/my-post.md
   # Should show: text/plain; charset=utf-8
   ```

2. **Convert to UTF-8 (Linux/macOS):**
   ```bash
   iconv -f ISO-8859-1 -t UTF-8 posts/my-post.md -o posts/my-post-utf8.md
   mv posts/my-post-utf8.md posts/my-post.md
   ```

3. **Remove BOM if present:**
   ```bash
   sed -i '1s/^\xEF\xBB\xBF//' posts/my-post.md
   ```

4. **Configure your editor to use UTF-8:**
   - VS Code: Bottom status bar shows encoding, click to change
   - Vim: `:set encoding=utf-8`
   - Ensure files are saved without BOM

5. **Verify your HTML templates include charset:**
   ```html
   <meta charset="UTF-8">
   ```

---

## Build Issues

### Empty Output Directory

**Symptom:** After building, the output directory is empty or missing files.

**Cause:** No files matched the glob patterns, or all posts are filtered out.

**Solution:**

1. **Check glob patterns match your files:**
   ```toml
   [markata-go.glob]
   patterns = ["posts/**/*.md", "pages/*.md"]
   ```

2. **Verify files exist in the expected location:**
   ```bash
   ls -la posts/
   ```

3. **Check if files are being processed:**
   ```bash
   markata-go build -v
   ```

4. **Check feed filters aren't excluding everything:**
   ```toml
   [[markata-go.feeds]]
   filter = "published == True"  # Do any posts have published: true?
   ```

5. **Do a dry run to see what would be built:**
   ```bash
   markata-go build --dry-run
   ```

---

### Missing Templates

**Symptom:**
```
Error: template "post.html" not found
```
Or:
```
Error: template "custom-layout.html" not found
```

**Cause:** Template file doesn't exist in the templates directory.

**Solution:**

1. **Check templates directory exists:**
   ```bash
   ls -la templates/
   ```

2. **Verify template filename matches what's specified:**
   ```yaml
   # In frontmatter
   template: "post.html"  # Must exist as templates/post.html
   ```

3. **Create missing templates or use defaults:**
   ```bash
   # Create a minimal post template
   mkdir -p templates
   cat > templates/post.html << 'EOF'
   {% extends "base.html" %}
   {% block content %}
   <article>
     <h1>{{ post.Title }}</h1>
     {{ body|safe }}
   </article>
   {% endblock %}
   EOF
   ```

4. **Check templates_dir configuration:**
   ```toml
   [markata-go]
   templates_dir = "templates"  # Must match your directory
   ```

See the [[templates-guide|Templates Guide]] for template setup.

---

### Template Errors

**Symptom:**
```
Error: template error in post.html: unexpected tag "endif"
```
Or:
```
Error: template error: variable "post.title" not found
```

**Cause:** Syntax error in template or incorrect variable access.

**Solution:**

1. **Check template syntax:**

   **Wrong - mismatched tags:**
   ```html
   {% if post.Tags %}
   <ul>...</ul>
   {% end %}  <!-- Wrong -->
   ```

   **Correct:**
   ```html
   {% if post.Tags %}
   <ul>...</ul>
   {% endif %}
   ```

2. **Check variable names are correct (case-sensitive):**

   **Wrong:**
   ```html
   {{ post.title }}  <!-- Wrong case -->
   ```

   **Correct:**
   ```html
   {{ post.Title }}  <!-- Correct -->
   ```

3. **Use safe filter for HTML content:**
   ```html
   {{ body|safe }}
   {{ post.ArticleHTML|safe }}
   ```

4. **Handle missing values:**
   ```html
   {{ post.Description|default_if_none:"No description" }}
   ```

**Common template variables:**
| Variable | Description |
|----------|-------------|
| `post.Title` | Post title |
| `post.Slug` | URL slug |
| `post.Href` | Relative URL path |
| `post.Date` | Publication date |
| `post.Tags` | List of tags |
| `post.ArticleHTML` | Rendered HTML content |
| `body` | Alias for ArticleHTML |
| `config.Title` | Site title |
| `config.URL` | Site URL |

---

### Slow Builds

**Symptom:** Build takes much longer than expected.

**Cause:** Too many files, inefficient configuration, or resource constraints.

**Solution:**

1. **Increase concurrency:**
   ```toml
   [markata-go]
   concurrency = 8  # Adjust based on CPU cores
   ```

2. **Narrow glob patterns:**
   ```toml
   [markata-go.glob]
   # Instead of:
   # patterns = ["**/*.md"]

   # Be more specific:
   patterns = ["posts/*.md", "pages/*.md"]
   ```

3. **Enable gitignore filtering:**
   ```toml
   [markata-go.glob]
   use_gitignore = true
   ```

4. **Profile the build:**
   ```bash
   time markata-go build -v
   ```

5. **Check for large files or images in content directories:**
   ```bash
   find posts/ -type f -size +1M
   ```

---

## Feed Issues

### RSS/Atom Not Generating

**Symptom:** No `rss.xml` or `atom.xml` files in the output.

**Cause:** Feed formats are not enabled in configuration.

**Solution:**

1. **Enable RSS/Atom formats:**
   ```toml
   [[markata-go.feeds]]
   slug = "blog"
   title = "Blog"
   filter = "published == True"

   [markata-go.feeds.formats]
   html = true
   rss = true
   atom = true
   ```

2. **Or set defaults for all feeds:**
   ```toml
   [markata-go.feed_defaults.formats]
   html = true
   rss = true
   atom = true
   ```

3. **Verify the feed is being generated:**
   ```bash
   markata-go build -v
   ls -la output/blog/rss.xml output/blog/atom.xml
   ```

See the [[feeds-guide|Feeds Guide]] for detailed feed configuration.

---

### Empty Feeds

**Symptom:** RSS/Atom files are generated but contain no items.

**Cause:** No posts match the feed filter, or posts lack required fields.

**Solution:**

1. **Check feed filter matches posts:**
   ```toml
   [[markata-go.feeds]]
   slug = "blog"
   filter = "published == True"  # Do posts have published: true?
   ```

2. **Verify posts have required frontmatter:**
   ```yaml
   ---
   title: "My Post"        # Required for feed items
   date: 2024-01-15        # Required for feed ordering
   published: true         # Required to match filter
   ---
   ```

3. **Check syndication settings:**
   ```toml
   [markata-go.feed_defaults.syndication]
   max_items = 20          # Not set to 0
   ```

4. **Debug by removing the filter temporarily:**
   ```toml
   [[markata-go.feeds]]
   slug = "debug"
   filter = ""  # No filter - should include all posts
   ```

---

### Missing Site URL

**Symptom:**
- Feed URLs show as relative paths or empty
- Feed validators report "missing link" errors
- Social sharing doesn't work

**Cause:** The `url` configuration is not set.

**Solution:**

1. **Set the site URL in configuration:**
   ```toml
   [markata-go]
   url = "https://example.com"
   ```

2. **Or set via environment variable:**
   ```bash
   MARKATA_GO_URL=https://example.com markata-go build
   ```

3. **Verify the URL is set:**
   ```bash
   markata-go config get url
   ```

---

### Feed Validation Errors

**Symptom:** Feed validators report errors like:
- "Missing required element"
- "Invalid date format"
- "URL must be absolute"

**Cause:** Missing configuration or invalid content.

**Solution:**

1. **Ensure site URL is set (required for valid feeds):**
   ```toml
   [markata-go]
   url = "https://example.com"
   ```

2. **Ensure all posts have dates:**
   ```yaml
   ---
   date: 2024-01-15
   ---
   ```

3. **Validate generated feeds:**
   ```bash
   # RSS
   xmllint --noout output/blog/rss.xml

   # Atom  
   xmllint --noout output/blog/atom.xml
   ```

4. **Use online validators:**
   - [W3C Feed Validator](https://validator.w3.org/feed/)
   - [Feed Validator](https://www.feedvalidator.org/)

5. **Check feed templates if using custom ones:**
   ```toml
   [markata-go.feed_defaults.templates]
   rss = "rss.xml"
   atom = "atom.xml"
   ```

---

## Serve Issues

### Port Already in Use

**Symptom:**
```
Error: listen tcp :8000: bind: address already in use
```

**Cause:** Another process is using port 8000.

**Solution:**

1. **Use a different port:**
   ```bash
   markata-go serve -p 3000
   markata-go serve -p 8080
   ```

2. **Find what's using the port:**
   ```bash
   # Linux/macOS
   lsof -i :8000

   # Windows
   netstat -ano | findstr :8000
   ```

3. **Kill the process using the port:**
   ```bash
   # Linux/macOS
   kill $(lsof -t -i :8000)

   # Or force kill
   kill -9 $(lsof -t -i :8000)
   ```

---

### Live Reload Not Working

**Symptom:** Changes to files don't trigger a browser refresh.

**Cause:** File watcher not detecting changes, or browser cache.

**Solution:**

1. **Ensure file watching is enabled (default):**
   ```bash
   markata-go serve  # Watching is on by default

   # NOT this:
   markata-go serve --no-watch  # This disables watching
   ```

2. **Check for file system notification limits (Linux):**
   ```bash
   # Check current limit
   cat /proc/sys/fs/inotify/max_user_watches

   # Increase if needed
   echo fs.inotify.max_user_watches=524288 | sudo tee -a /etc/sysctl.conf
   sudo sysctl -p
   ```

3. **Hard refresh your browser:**
   - Chrome/Firefox: `Ctrl+Shift+R` (Windows/Linux) or `Cmd+Shift+R` (Mac)
   - Or open in incognito/private mode

4. **Check the terminal for rebuild messages:**
   - The server should show "Rebuilding..." when files change
   - If not, the file watcher isn't detecting changes

5. **Try verbose mode:**
   ```bash
   markata-go serve -v
   ```

---

### Changes Not Reflecting

**Symptom:** You made changes but the browser shows old content.

**Cause:** Browser caching, or changes to files outside watched directories.

**Solution:**

1. **Hard refresh the browser:**
   - `Ctrl+Shift+R` or `Cmd+Shift+R`

2. **Clear browser cache:**
   - Developer Tools > Application > Clear Storage

3. **Verify the file is in a watched directory:**
   - Files must match glob patterns to be watched
   - Check `patterns` in your config

4. **Rebuild manually if needed:**
   ```bash
   # Stop the server and rebuild
   markata-go build --clean
   markata-go serve
   ```

5. **Check file permissions:**
   ```bash
   ls -la posts/my-post.md
   # Ensure the file is readable
   ```

---

## Deployment Issues

### 404 Errors

**Symptom:** Pages return 404 Not Found after deployment.

**Cause:** Wrong base URL, incorrect paths, or server configuration.

**Solution:**

1. **Check base URL matches deployment:**
   ```toml
   [markata-go]
   url = "https://yourdomain.com"  # Must match exactly
   ```

2. **For GitHub Pages with repo subdirectory:**
   ```toml
   # If deployed to https://username.github.io/repo-name/
   [markata-go]
   url = "https://username.github.io/repo-name"
   ```

3. **Verify files exist in output:**
   ```bash
   ls -la output/
   ls -la output/my-post/index.html
   ```

4. **Check server is configured for clean URLs:**
   - Ensure `/my-post/` serves `/my-post/index.html`
   - Most static hosts handle this automatically

5. **For custom servers (nginx), ensure proper config:**
   ```nginx
   location / {
       try_files $uri $uri/ =404;
   }
   ```

---

### Broken Links

**Symptom:** Links on the site lead to 404 pages.

**Cause:** Internal links are incorrect or posts have moved/been deleted.

**Solution:**

1. **Validate links before deploying:**
   ```bash
   # Install linkinator
   npm install -g linkinator

   # Check for broken links
   linkinator output --recurse
   ```

2. **Check wikilinks are correct:**
   ```markdown
   [[other-post]]           # Must match a slug exactly
   [[other-post|Link Text]] # Custom link text
   ```

3. **Verify slug values:**
   ```yaml
   ---
   slug: "my-post"  # Links should reference this exact slug
   ---
   ```

4. **Add redirects for moved content:**

   **Netlify (_redirects or netlify.toml):**
   ```
   /old-post  /new-post  301
   ```

   **Vercel (vercel.json):**
   ```json
   {
     "redirects": [
       { "source": "/old-post", "destination": "/new-post", "permanent": true }
     ]
   }
   ```

---

### Missing Assets

**Symptom:** CSS, JavaScript, or images don't load.

**Cause:** Assets not copied to output, wrong paths, or base URL issues.

**Solution:**

1. **Check assets_dir is configured:**
   ```toml
   [markata-go]
   assets_dir = "static"
   ```

2. **Verify assets are in the correct directory:**
   ```bash
   ls -la static/
   # Should contain css/, js/, images/, etc.
   ```

3. **Check assets were copied to output:**
   ```bash
   ls -la output/static/
   # or
   ls -la output/css/
   ```

4. **Use correct paths in templates:**
   ```html
   <!-- Absolute path (recommended) -->
   <link rel="stylesheet" href="/css/style.css">

   <!-- NOT relative -->
   <link rel="stylesheet" href="css/style.css">  <!-- May break on subpages -->
   ```

5. **For GitHub Pages subdirectory deployment:**
   ```html
   <!-- If deployed to /repo-name/ -->
   <link rel="stylesheet" href="/repo-name/css/style.css">
   ```

---

### Wrong Base URL

**Symptom:**
- Absolute URLs point to wrong domain
- RSS/Atom feeds have incorrect links
- Social sharing shows wrong URLs

**Cause:** The `url` config doesn't match the deployment URL.

**Solution:**

1. **Set the correct URL for each environment:**

   **Production (config file):**
   ```toml
   [markata-go]
   url = "https://example.com"
   ```

   **Staging (environment variable):**
   ```bash
   MARKATA_GO_URL=https://staging.example.com markata-go build
   ```

2. **For different deployment environments, use CI/CD variables:**

   **GitHub Actions:**
   ```yaml
   - name: Build
     run: markata-go build
     env:
       MARKATA_GO_URL: ${{ vars.SITE_URL }}
   ```

   **Netlify:**
   ```toml
   [context.production.environment]
   MARKATA_GO_URL = "https://example.com"

   [context.deploy-preview.environment]
   MARKATA_GO_URL = ""  # Use relative URLs for previews
   ```

3. **Verify the URL is correct:**
   ```bash
   markata-go config get url
   ```

---

## Plugin Issues

### Plugin Not Loading

**Symptom:**
- Plugin features don't work
- Warning about unknown plugin name

**Cause:** Plugin not registered, typo in name, or disabled.

**Solution:**

1. **Check hooks configuration:**
   ```toml
   [markata-go]
   hooks = ["default"]  # Loads all default plugins
   ```

2. **Check disabled_hooks doesn't include your plugin:**
   ```toml
   [markata-go]
   disabled_hooks = []  # Make sure it's not listed here
   ```

3. **Verify the plugin name against the [[built-in-plugins|Plugin Reference]].**

4. **For custom plugins, ensure they're registered properly:**
   - See the [[plugin-development|Plugin Development Guide]]

---

### Plugin Errors

**Symptom:**
```
Error: plugin "my_plugin" failed: <specific error message>
```

**Cause:** Plugin-specific issue, usually configuration or data related.

**Solution:**

1. **Check plugin configuration:**
   ```toml
   [markata-go]
   # Plugin-specific options go under the main namespace
   my_plugin_option = "value"
   ```

2. **Enable verbose output to see detailed errors:**
   ```bash
   markata-go build -v
   ```

3. **Check posts have required fields for the plugin:**
   - Reading time plugin needs content
   - Feed plugins need dates for sorting

4. **Review the plugin documentation for requirements**

5. **Try disabling the plugin to isolate the issue:**
   ```toml
   [markata-go]
   disabled_hooks = ["my_plugin"]
   ```

---

## Getting Help

If you can't resolve your issue with this guide:

1. **Check existing issues:**
   - [GitHub Issues](https://github.com/WaylonWalker/markata-go/issues)

2. **Search the documentation:**
   - [[getting-started|Getting Started]]
   - [[configuration-guide|Configuration Guide]]
   - [[feeds-guide|Feeds Guide]]
   - [[templates-guide|Templates Guide]]

3. **Create a minimal reproduction:**
   - Create the smallest possible project that shows the issue
   - Include your config file and a sample post

4. **Open a new issue with:**
   - markata-go version (`markata-go --version`)
   - Go version (`go version`)
   - Operating system
   - Full error message
   - Steps to reproduce
   - Relevant configuration

5. **Use verbose mode when reporting issues:**
   ```bash
   markata-go build -v 2>&1 | tee build-output.txt
   ```

---

## Quick Reference: Common Commands

```bash
# Validate configuration
markata-go config validate

# Show resolved configuration
markata-go config show

# Get a specific config value
markata-go config get url

# Build with verbose output
markata-go build -v

# Dry run (preview without writing)
markata-go build --dry-run

# Clean build
markata-go build --clean

# Serve on a different port
markata-go serve -p 3000

# Create a new post
markata-go new "My Post Title"

# Check version
markata-go --version
```

## Agent-assisted troubleshooting

Run `markata triage` to prepare a diagnostic prompt and choose an installed coding agent, clipboard, or stdout. See the [triage guide](/docs/guides/triage/).
