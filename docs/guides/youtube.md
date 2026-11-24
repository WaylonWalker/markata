---
title: "Embedding YouTube Videos"
description: "How to embed YouTube videos in your posts with a simple one-line URL"
date: 2024-01-15
published: true
tags:
  - documentation
  - guides
  - youtube
  - media
---


markata-go automatically converts YouTube URLs into responsive embedded videos. Just paste a YouTube link on its own line, and it becomes a playable video embed.

Note: This guide covers the YouTube plugin (standalone URLs). The embeds plugin's oEmbed-based YouTube renders use Lite YouTube by default.

## Quick Start

Simply paste a YouTube URL on its own line in your markdown:

```markdown
Check out this video:

https://www.youtube.com/watch?v=dQw4w9WgXcQ

More content below...
```

That's it! The URL will be converted into a responsive video player.

## Supported URL Formats

All common YouTube URL formats work:

| Format | Example |
|--------|---------|
| Standard | `https://www.youtube.com/watch?v=VIDEO_ID` |
| Short | `https://youtu.be/VIDEO_ID` |
| Mobile | `https://m.youtube.com/watch?v=VIDEO_ID` |

## Timestamps

Start videos at a specific time by adding a timestamp:

```markdown
https://youtu.be/dQw4w9WgXcQ?t=1m30s
```

Supported timestamp formats:
- Seconds: `?t=90`
- Minutes and seconds: `?t=1m30s`
- Hours, minutes, seconds: `?t=1h2m3s`
- Start parameter: `?start=90`

## Configuration

Configure the YouTube plugin in your `markata-go.toml`:

```toml
[markata-go.youtube]
enabled = true              # Enable/disable embeds (default: true)
privacy_enhanced = true     # Use privacy-enhanced mode (default: true)
container_class = "youtube-embed"  # CSS class for styling
lazy_load = true            # Lazy load videos (default: true)
```

### Privacy-Enhanced Mode

By default, videos use YouTube's privacy-enhanced mode (`youtube-nocookie.com`):

- **No tracking cookies** until the user clicks play
- **GDPR compliant** - helps meet privacy regulations
- **Same video experience** - playback works identically

To use standard YouTube embeds:

```toml
[markata-go.youtube]
privacy_enhanced = false
```

### Lazy Loading

Lazy loading (enabled by default) improves page performance by only loading video iframes when they're about to enter the viewport.

## Styling

Videos are automatically responsive with a 16:9 aspect ratio. The default theme includes styling, but you can customize:

```css
/* Custom video container styling */
.youtube-embed {
  margin: 2rem 0;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}
```

## What Doesn't Get Embedded

The plugin only converts URLs that are **on their own line**. These patterns are NOT converted:

### Inline URLs

```markdown
Check out https://youtu.be/dQw4w9WgXcQ for more info.
```

This stays as a regular link because it's inline with other text.

### Code Blocks

```markdown
To embed a video, use: `https://youtu.be/VIDEO_ID`
```

URLs in code are preserved for documentation purposes.

### Invalid Video IDs

URLs with invalid video IDs (not exactly 11 characters) are left unchanged.

## Output HTML

The plugin generates this HTML structure:

```html
<div class="youtube-embed">
  <iframe
    src="https://www.youtube-nocookie.com/embed/VIDEO_ID"
    title="YouTube video player"
    frameborder="0"
    allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
    allowfullscreen
    loading="lazy">
  </iframe>
</div>
```

## Tips

1. **One URL per line** - Each video should be on its own line with blank lines before and after
2. **Use short URLs** - `youtu.be` links are cleaner in your markdown
3. **Add context** - Include a brief description before your video embed
4. **Test timestamps** - Preview your post to ensure timestamps work correctly

## Troubleshooting

### Video Not Embedding

- Ensure the URL is on its own line (not inline with text)
- Check that the video ID is valid (11 characters)
- Verify the plugin is enabled in your config

### Video Not Playing

- Check if the video is available in your region
- Some videos may have embedding disabled by the uploader
- Try with `privacy_enhanced = false` if issues persist

## Related

- [[markdown|Markdown Features]] - Other markdown enhancements
- [[md-video|Video Embeds]] - Embed local video files
- [[configuration-guide|Configuration]] - Global settings
