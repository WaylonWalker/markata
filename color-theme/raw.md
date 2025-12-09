---
title: Theming your markata site
description: Guide to get going with Markata

---

## Default colors

The default markata page template support the following colors to be configured
in the `markata.toml` file. There are two sets of similarly named colors, one
for light theme and one for dark.  The default page template will set these
colors based on the users `prefers-color-scheme`.

```toml
[markata]
color_bg = '#1f2022'
color_bg_code = '#1f2022'
color_text = '#eefbfe'
color_link = '#47cbff' 
color_accent = '#e1bd00c9'
overlay_brightness = '.85'

color_bg_light = '#eefbfe'
color_bg_code_light = '#eefbfe'
color_text_light = '#1f2022'
color_link_light = '#47cbff' 
color_accent_light = '#ffeb00'
overlay_brightness_light = '.95'
```

## Pink and Purple

```toml
[markata]
color_bg = 'deeppink'
color_bg_code = 'rebeccapurple'
color_text = 'white'
color_link = 'aqua' 
color_accent = 'peachpuff'
overlay_brightness = '1.2'
```

## Changing your favicon

Your favicon should be kept in your `assets` directory, which is `./static` by
default.  You can name the icon what you want, `icon.png` is the default, but
you will have to change your icon config in your `markata.toml` to use a
different value.

``` toml
[markata]
assets_dir = "static"
icon = "icon.png"
```
