---
title: "Calendar Feed Views"
description: "Use the built-in calendar view available on every Markata-Go HTML feed."
date: 2026-09-28
published: true
tags:
  - documentation
  - feeds
  - templates
---


Every normal Markata-Go HTML feed includes a **Calendar** view by default. You do not need a second feed, a template override, or special calendar configuration.

The interaction is inspired by [Jim Nielsen's calendar archive](https://blog.jim-nielsen.com/2026/blog-calendar-view/): the normal archive and its calendar are peer views of the same collection, and dates with posts are visibly marked. Select a marked date to reveal the title or titles published that day.

## Open Calendar

Open any feed normally and choose **Calendar** from its view control. For a feed at `/blog/`, the same state can be linked directly as:

```text
/blog/?view=calendar
```

For the root feed, use:

```text
/?view=calendar
```

The **Simple** view also links to Calendar. When feed sidebars are enabled on posts, the sidebar includes a `calendar` link for the currently selected feed.

## What the view shows

Calendar uses the feed's complete `feed.posts` collection, so it represents the full feed even when the normal HTML presentation is paginated.

For each represented year it renders all twelve months. Every day appears in the correct weekday column. Dates containing posts get an activity marker and an expandable list of titles. If multiple posts share a date, they all appear beneath that day.

Posts without a valid date cannot be placed in a month grid. They remain available through the normal feed presentation.

## Three human-facing views

A default feed now has three presentations:

- **Posts** — the feed's normal card/list/template presentation
- **Simple** — the existing compact HTML list at `/simple/`, when enabled
- **Calendar** — the publishing-rhythm view, opened in place and deep-linkable with `?view=calendar`

RSS, Atom, JSON, Markdown, and text remain feed export/subscription formats rather than presentation modes.

Built-in alternate HTML layouts such as `feed-photo-grid.html` keep their existing primary presentation and gain Calendar as a peer view too.

## Accessibility and progressive enhancement

The normal feed remains the server-rendered default. Calendar controls stay hidden until JavaScript successfully initializes, so visitors without JavaScript keep the existing feed experience.

The page embeds a hidden dated source built from the complete feed collection. Calendar day disclosures use native `<details>` and `<summary>` controls, so marked dates are keyboard-operable without a custom keybinding layer. Weekday labels and day summaries expose useful accessible names while decorative activity marks stay hidden from assistive technology.

## Pagination

You do not need to disable pagination. The visible Posts view may render `page.posts`, while Calendar intentionally reads the complete `feed.posts` collection.

## Customizing Calendar

Calendar styling and behavior live in the default theme assets:

- `css/calendar-feed.css`
- `js/calendar-feed.js`

The CSS uses theme variables and `currentColor`-derived borders/backgrounds so it follows custom palettes without calendar-specific color configuration.

The previously shipped `calendar-feed.html` remains available for sites that intentionally want a dedicated calendar-first custom template. It is no longer required to enable Calendar on ordinary feeds.
