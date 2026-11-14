---
title: "Adaptive Loading"
description: "How Markata-Go keeps sites useful and calm on slow or unstable connections."
date: 2026-10-07
published: true
tags:
  - performance
  - themes
---


Markata-Go sites keep article text available immediately and let the browser
load images and video with native HTML behavior. Visitors can choose how much
optional media work the site should do.

## Choose a loading mode

Open **Loading** in the site header and choose one of these modes:

- **Auto** uses recent browser connection hints and a few real resource timings.
- **Save Data** limits optional media work and speculative requests.
- **Full Quality** allows optional work to begin sooner.

Auto is the default. A saved Save Data or Full Quality choice always takes
precedence over Auto. The preference is stored in the current browser only.

## What Auto does

Auto uses `navigator.connection` when the browser exposes it, but the site
works without that API. It also learns from a small number of recent navigation
and resource timings. One slow image or request does not change the mode. Auto
requires two meaningful transfers below 2 Mbps to constrain optional work and
four consecutive transfers above 4 Mbps with low response latency to recover.
Generic 4g/downlink hints initialize Auto but cannot erase measured evidence.
Browser Save Data intent constrains Auto immediately. Each page observes at
most twelve transfers for thirty seconds; later pages can provide fresh evidence.
Offline constrains Auto immediately, and reconnecting needs fresh evidence.
Suggestions require four poor observations for Full Quality or five good ones
for Save Data.

If Auto becomes confident that your manual choice no longer fits the
connection, a small status prompt offers **Use Auto** or **Keep [mode]**.
Keeping or dismissing the choice quiets the prompt for the current session.

## Media and JavaScript

Images use browser-native loading and retain their dimensions while loading.
Derived video posters reserve a 16:9 frame when the author has not supplied
height or inline styling. Supply your own height/style for a different frame.
Video previews show a poster and do not fetch the video until playback is
requested. Article video keeps native controls and remains usable if JavaScript
does not run. The adaptive script does not replace an image that is already
visible or restart a video when the connection changes.

Markdown videos on individual posts (including shot posts) now default to
muted, looping autoplay **when the video becomes visible**, with
`preload = "none"` until then. Feed/preview cards are separate and remain
poster-only until deliberate interaction. Raw HTML retains autoplay intent
only when its author explicitly supplies `autoplay`; that intent is stored as
`data-authored-autoplay="true"` rather than a native autoplay attribute.

To opt out of autoplay for Markdown videos site-wide:

```toml
[markata-go.md_video]
autoplay = false
```

Auto starts marked videos once visible unless connection measurements indicate
constrained bandwidth. Save Data (including the browser's Save Data hint),
reduced motion, and offline mode suppress automatic starts even in Full Quality.
The controller pauses videos **it started** when they leave the viewport or
the tab becomes hidden, then resumes them when visible again unless the reader
manually paused. User-started videos are not hijacked. Raw video without
authored autoplay stays manual even when Markdown autoplay is enabled.

With JavaScript disabled or IntersectionObserver unavailable, videos remain
click-to-play through their native controls. The controller keeps posters and
`preload="none"` before viewport playback and does not automatically load
offscreen video data.

The controller is a small theme asset. Video and embed helpers remain
conditional, so a text-only page does not load a player just in case.


Autoplay starts only when the video enters the viewport and loading policy
allows it. The initial page waits for its load event so resource measurements
can detect constrained bandwidth before autoplay. Switching to Save Data pauses
only playback managed by this controller; manual playback remains untouched.


Keep native controls enabled (or supply an accessible play action) for deferred
videos. A video deliberately authored without controls needs its own interaction
path in Save Data and no-JS mode; autoplay intent alone is not a play control.


The existing full-document View Transition navigator preserves loading-mode
attributes and the persistent control. Each completed route initializes only
its newly inserted authored media, expires stale session evidence, and opens a
fresh bounded timing window. It does not rerun the controller or reset the
visitor's manual choice.
