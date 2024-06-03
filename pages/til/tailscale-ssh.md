---
date: 2024-06-03 16:40:18
templateKey: til
title: tailscale ssh
published: true
tags:
  - linux

---

Tailscale allows you to ssh into all of your tailscale machines, it busts
through firewalls and accross networks without complex setup.  If you have used
tailscale before this is an obvious no brainer.  What is not obvious is that
you can configure tailscale to allow ssh connections from devices within your
tailnet without even a ssh daemon process running right through the tailscale
daemon.

``` bash
tailscale status
tailscale set --ssh
```

I picked this up from the tailscale youtube channel.

[[ https://www.youtube.com/watch?v=08clF9srJ2k&t=35s ]]{.youtube-embed}
