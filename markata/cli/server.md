---
date: 2025-12-09
description: None
published: false
slug: markata/cli/server
title: server.py


---

---

None

---

!!! function
    <h2 id="find_port" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">find_port <em class="small">function</em></h2>

    Find a port not in ues starting at given port

???+ source "find_port <em class='small'>source</em>"
    ```python
    def find_port(port: int = 8000) -> int:
        """Find a port not in ues starting at given port"""
        import socket

        with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
            if s.connect_ex(("localhost", port)) == 0:
                return find_port(port=port + 1)
            return port
    ```
!!! function
    <h2 id="serve" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">serve <em class="small">function</em></h2>

    Serve the site locally.

???+ source "serve <em class='small'>source</em>"
    ```python
    def serve():
            """
            Serve the site locally.
            """
            run_server()
    ```