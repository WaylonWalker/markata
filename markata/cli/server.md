---
content: "None\n\n\n!! function <h2 id='find_port' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>find_port <em class='small'>function</em></h2>\n    Find a port not in ues
  starting at given port\n???+ source \"find_port <em class='small'>source</em>\"\n\n```python\n\n
  \       def find_port(port: int = 8000) -> int:\n            \"\"\"Find a port not
  in ues starting at given port\"\"\"\n            import socket\n\n            with
  socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:\n                if s.connect_ex((\"localhost\",
  port)) == 0:\n                    return find_port(port=port + 1)\n                return
  port\n```\n\n\n!! class <h2 id='Server' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>Server <em class='small'>class</em></h2>\n\n???+ source \"Server <em class='small'>source</em>\"\n\n```python\n\n
  \       class Server:\n            def __init__(\n                self: \"Server\",\n
  \               *,\n                auto_restart: bool = True,\n                directory:
  Union[str, \"Path\"] = None,\n                port: int = 8000,\n            ) ->
  None:\n                if directory is None:\n                    from markata import
  Markata\n\n                    m = Markata()\n                    directory = Path(str(m.config[\"output_dir\"]))\n\n
  \               self.auto_restart = auto_restart\n                self.directory
  = directory\n                self.port = find_port(port=port)\n                self.start_server()\n
  \               atexit.register(self.kill)\n\n            def start_server(self)
  -> None:\n                import subprocess\n\n                self.cmd = [\n                    \"python\",\n
  \                   \"-m\",\n                    \"http.server\",\n                    str(self.port),\n
  \                   \"--directory\",\n                    self.directory,\n                ]\n\n
  \               self.proc = subprocess.Popen(\n                    self.cmd,\n                    stderr=subprocess.PIPE,\n
  \                   stdout=subprocess.PIPE,\n                )\n                self.start_time
  = time.time()\n\n            @property\n            def uptime(self) -> int:\n                return
  round(time.time() - self.start_time)\n\n            @property\n            def title(self)
  -> str:\n                return f\"server ({self.uptime})\"\n\n            def kill(self)
  -> None:\n                self.auto_restart = False\n                self.proc.stdout.close()\n
  \               self.proc.stderr.close()\n                self.proc.kill()\n                self.proc.wait()\n\n
  \           def __rich__(self) -> Panel:\n                if not self.proc.poll():\n
  \                   return Panel(\n                        (\n                            f\"[green]serving
  on port: [gold1]{self.port} \"\n                            f\"[green]using pid:
  [gold1]{self.proc.pid} \"\n                            f\"[green]uptime: [gold1]{self.uptime}
  \"\n                            f\"[green]link: [gold1] http://localhost:{self.port}[/]\"\n
  \                       ),\n                        border_style=\"blue\",\n                        title=self.title,\n
  \                       expand=True,\n                    )\n\n                else:\n
  \                   return Panel(\n                        \"[red]server died\",\n
  \                       title=self.title,\n                        border_style=\"red\",\n
  \                       expand=True,\n                    )\n```\n\n\n!! function
  <h2 id='configure' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>configure
  <em class='small'>function</em></h2>\n\n???+ source \"configure <em class='small'>source</em>\"\n\n```python\n\n
  \       def configure(markata: \"Markata\") -> None:\n            def get_server(self):\n
  \               try:\n                    return self._server\n                except
  AttributeError:\n                    self._server: Server = Server(directory=str(self.config[\"output_dir\"]))\n
  \                   return self._server\n\n            from markata import Markata\n\n
  \           Markata.server = property(get_server)\n```\n\n\n!! function <h2 id='run_server'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>run_server <em class='small'>function</em></h2>\n\n???+
  source \"run_server <em class='small'>source</em>\"\n\n```python\n\n        def
  run_server() -> None:\n            from rich.live import Live\n\n            from
  .cli import run_until_keyboard_interrupt\n\n            with Live(Server(), refresh_per_second=1,
  screen=True):\n                run_until_keyboard_interrupt()\n```\n\n\n!! function
  <h2 id='cli' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>cli <em
  class='small'>function</em></h2>\n\n???+ source \"cli <em class='small'>source</em>\"\n\n```python\n\n
  \       def cli(app: typer.Typer, markata: \"Markata\") -> None:\n            server_app
  = typer.Typer()\n            app.add_typer(server_app)\n\n            @server_app.callback(invoke_without_command=True)\n
  \           def serve():\n                \"\"\"\n                Serve the site
  locally.\n                \"\"\"\n                run_server()\n```\n\n\n!! method
  <h2 id='__init__' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>__init__
  <em class='small'>method</em></h2>\n\n???+ source \"__init__ <em class='small'>source</em>\"\n\n```python\n\n
  \       def __init__(\n                self: \"Server\",\n                *,\n                auto_restart:
  bool = True,\n                directory: Union[str, \"Path\"] = None,\n                port:
  int = 8000,\n            ) -> None:\n                if directory is None:\n                    from
  markata import Markata\n\n                    m = Markata()\n                    directory
  = Path(str(m.config[\"output_dir\"]))\n\n                self.auto_restart = auto_restart\n
  \               self.directory = directory\n                self.port = find_port(port=port)\n
  \               self.start_server()\n                atexit.register(self.kill)\n```\n\n\n!!
  method <h2 id='start_server' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>start_server <em class='small'>method</em></h2>\n\n???+ source \"start_server
  <em class='small'>source</em>\"\n\n```python\n\n        def start_server(self) ->
  None:\n                import subprocess\n\n                self.cmd = [\n                    \"python\",\n
  \                   \"-m\",\n                    \"http.server\",\n                    str(self.port),\n
  \                   \"--directory\",\n                    self.directory,\n                ]\n\n
  \               self.proc = subprocess.Popen(\n                    self.cmd,\n                    stderr=subprocess.PIPE,\n
  \                   stdout=subprocess.PIPE,\n                )\n                self.start_time
  = time.time()\n```\n\n\n!! method <h2 id='uptime' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>uptime <em class='small'>method</em></h2>\n\n???+ source \"uptime <em class='small'>source</em>\"\n\n```python\n\n
  \       def uptime(self) -> int:\n                return round(time.time() - self.start_time)\n```\n\n\n!!
  method <h2 id='title' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>title
  <em class='small'>method</em></h2>\n\n???+ source \"title <em class='small'>source</em>\"\n\n```python\n\n
  \       def title(self) -> str:\n                return f\"server ({self.uptime})\"\n```\n\n\n!!
  method <h2 id='kill' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>kill
  <em class='small'>method</em></h2>\n\n???+ source \"kill <em class='small'>source</em>\"\n\n```python\n\n
  \       def kill(self) -> None:\n                self.auto_restart = False\n                self.proc.stdout.close()\n
  \               self.proc.stderr.close()\n                self.proc.kill()\n                self.proc.wait()\n```\n\n\n!!
  method <h2 id='__rich__' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>__rich__ <em class='small'>method</em></h2>\n\n???+ source \"__rich__ <em
  class='small'>source</em>\"\n\n```python\n\n        def __rich__(self) -> Panel:\n
  \               if not self.proc.poll():\n                    return Panel(\n                        (\n
  \                           f\"[green]serving on port: [gold1]{self.port} \"\n                            f\"[green]using
  pid: [gold1]{self.proc.pid} \"\n                            f\"[green]uptime: [gold1]{self.uptime}
  \"\n                            f\"[green]link: [gold1] http://localhost:{self.port}[/]\"\n
  \                       ),\n                        border_style=\"blue\",\n                        title=self.title,\n
  \                       expand=True,\n                    )\n\n                else:\n
  \                   return Panel(\n                        \"[red]server died\",\n
  \                       title=self.title,\n                        border_style=\"red\",\n
  \                       expand=True,\n                    )\n```\n\n\n!! function
  <h2 id='get_server' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>get_server
  <em class='small'>function</em></h2>\n\n???+ source \"get_server <em class='small'>source</em>\"\n\n```python\n\n
  \       def get_server(self):\n                try:\n                    return
  self._server\n                except AttributeError:\n                    self._server:
  Server = Server(directory=str(self.config[\"output_dir\"]))\n                    return
  self._server\n```\n\n\n!! function <h2 id='serve' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>serve <em class='small'>function</em></h2>\n    Serve the site locally.\n???+
  source \"serve <em class='small'>source</em>\"\n\n```python\n\n        def serve():\n
  \               \"\"\"\n                Serve the site locally.\n                \"\"\"\n
  \               run_server()\n```\n\n"
date: 0001-01-01
description: None ! ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  !
  ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
  source  ! ?
html:
  index: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Server.Py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"None ! ! ???+ source  ! ???+ source  ! ???+ source
    \ ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  !
    ???+ source  ! ???+ source  ! ?\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Server.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"None ! ! ???+ source  ! ???+ source  !
    ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
    source  ! ???+ source  ! ???+ source  ! ?\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    </head>\n
    \   <body>\n<div class='container flex flex-row min-h-screen'>\n    <div>\n    </div>\n
    \   <div class='flex-grow px-8 mx-auto min-h-screen'>\n<header class='flex justify-center
    items-center p-8'>\n\n    <nav class='flex justify-center items-center my-8'>\n
    \       <a\n            href='/'>markata</a>\n        <a\n            href='https://github.com/WaylonWalker/markata'>GitHub</a>\n
    \       <a\n            href='https://markata.dev/docs/'>docs</a>\n        <a\n
    \           href='https://markata.dev/plugins/'>plugins</a>\n    </nav>\n\n    <div>\n
    \       <label id=\"theme-switch\" class=\"theme-switch\" for=\"checkbox-theme\"
    title=\"light/dark mode toggle\">\n            <input type=\"checkbox\" id=\"checkbox-theme\"
    />\n            <div class=\"slider round\"></div>\n        </label>\n    </div>\n</header><article
    class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n        Server.Py\n
    \   </h1>\n</section>    <section class=\"body\">\n        <p>None</p>\n<p>!!
    function <h2 id='find_port' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>find_port <em class='small'>function</em></h2>\nFind a port not in ues
    starting at given port</p>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">find_port <em class='small'>source</em></p>\n</div>\n<pre
    class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button class='copy' title='copy
    code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">find_port</span><span class=\"p\">(</span><span class=\"n\">port</span><span
    class=\"p\">:</span> <span class=\"nb\">int</span> <span class=\"o\">=</span>
    <span class=\"mi\">8000</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"nb\">int</span><span class=\"p\">:</span>\n<span class=\"w\">            </span><span
    class=\"sd\">&quot;&quot;&quot;Find a port not in ues starting at given port&quot;&quot;&quot;</span>\n
    \           <span class=\"kn\">import</span> <span class=\"nn\">socket</span>\n\n
    \           <span class=\"k\">with</span> <span class=\"n\">socket</span><span
    class=\"o\">.</span><span class=\"n\">socket</span><span class=\"p\">(</span><span
    class=\"n\">socket</span><span class=\"o\">.</span><span class=\"n\">AF_INET</span><span
    class=\"p\">,</span> <span class=\"n\">socket</span><span class=\"o\">.</span><span
    class=\"n\">SOCK_STREAM</span><span class=\"p\">)</span> <span class=\"k\">as</span>
    <span class=\"n\">s</span><span class=\"p\">:</span>\n                <span class=\"k\">if</span>
    <span class=\"n\">s</span><span class=\"o\">.</span><span class=\"n\">connect_ex</span><span
    class=\"p\">((</span><span class=\"s2\">&quot;localhost&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">port</span><span class=\"p\">))</span> <span class=\"o\">==</span>
    <span class=\"mi\">0</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">find_port</span><span class=\"p\">(</span><span
    class=\"n\">port</span><span class=\"o\">=</span><span class=\"n\">port</span>
    <span class=\"o\">+</span> <span class=\"mi\">1</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">port</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='Server' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Server
    <em class='small'>class</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">Server <em class='small'>source</em></p>\n</div>\n<pre
    class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button class='copy' title='copy
    code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
    <span class=\"nc\">Server</span><span class=\"p\">:</span>\n            <span
    class=\"k\">def</span> <span class=\"fm\">__init__</span><span class=\"p\">(</span>\n
    \               <span class=\"bp\">self</span><span class=\"p\">:</span> <span
    class=\"s2\">&quot;Server&quot;</span><span class=\"p\">,</span>\n                <span
    class=\"o\">*</span><span class=\"p\">,</span>\n                <span class=\"n\">auto_restart</span><span
    class=\"p\">:</span> <span class=\"nb\">bool</span> <span class=\"o\">=</span>
    <span class=\"kc\">True</span><span class=\"p\">,</span>\n                <span
    class=\"n\">directory</span><span class=\"p\">:</span> <span class=\"n\">Union</span><span
    class=\"p\">[</span><span class=\"nb\">str</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;Path&quot;</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"kc\">None</span><span class=\"p\">,</span>\n                <span
    class=\"n\">port</span><span class=\"p\">:</span> <span class=\"nb\">int</span>
    <span class=\"o\">=</span> <span class=\"mi\">8000</span><span class=\"p\">,</span>\n
    \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"k\">if</span> <span class=\"n\">directory</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                   <span class=\"kn\">from</span> <span class=\"nn\">markata</span>
    <span class=\"kn\">import</span> <span class=\"n\">Markata</span>\n\n                    <span
    class=\"n\">m</span> <span class=\"o\">=</span> <span class=\"n\">Markata</span><span
    class=\"p\">()</span>\n                    <span class=\"n\">directory</span>
    <span class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">m</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;output_dir&quot;</span><span class=\"p\">]))</span>\n\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">auto_restart</span>
    <span class=\"o\">=</span> <span class=\"n\">auto_restart</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">directory</span>
    <span class=\"o\">=</span> <span class=\"n\">directory</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">port</span>
    <span class=\"o\">=</span> <span class=\"n\">find_port</span><span class=\"p\">(</span><span
    class=\"n\">port</span><span class=\"o\">=</span><span class=\"n\">port</span><span
    class=\"p\">)</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">start_server</span><span class=\"p\">()</span>\n                <span
    class=\"n\">atexit</span><span class=\"o\">.</span><span class=\"n\">register</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">kill</span><span class=\"p\">)</span>\n\n            <span class=\"k\">def</span>
    <span class=\"nf\">start_server</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"kn\">import</span> <span class=\"nn\">subprocess</span>\n\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">cmd</span> <span class=\"o\">=</span> <span class=\"p\">[</span>\n
    \                   <span class=\"s2\">&quot;python&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;-m&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;http.server&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">port</span><span
    class=\"p\">),</span>\n                    <span class=\"s2\">&quot;--directory&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">directory</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">]</span>\n\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">proc</span> <span class=\"o\">=</span> <span
    class=\"n\">subprocess</span><span class=\"o\">.</span><span class=\"n\">Popen</span><span
    class=\"p\">(</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">cmd</span><span class=\"p\">,</span>\n                    <span
    class=\"n\">stderr</span><span class=\"o\">=</span><span class=\"n\">subprocess</span><span
    class=\"o\">.</span><span class=\"n\">PIPE</span><span class=\"p\">,</span>\n
    \                   <span class=\"n\">stdout</span><span class=\"o\">=</span><span
    class=\"n\">subprocess</span><span class=\"o\">.</span><span class=\"n\">PIPE</span><span
    class=\"p\">,</span>\n                <span class=\"p\">)</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">start_time</span>
    <span class=\"o\">=</span> <span class=\"n\">time</span><span class=\"o\">.</span><span
    class=\"n\">time</span><span class=\"p\">()</span>\n\n            <span class=\"nd\">@property</span>\n
    \           <span class=\"k\">def</span> <span class=\"nf\">uptime</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"nb\">int</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">return</span> <span class=\"nb\">round</span><span
    class=\"p\">(</span><span class=\"n\">time</span><span class=\"o\">.</span><span
    class=\"n\">time</span><span class=\"p\">()</span> <span class=\"o\">-</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">start_time</span><span
    class=\"p\">)</span>\n\n            <span class=\"nd\">@property</span>\n            <span
    class=\"k\">def</span> <span class=\"nf\">title</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"nb\">str</span><span class=\"p\">:</span>\n                <span
    class=\"k\">return</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;server
    (</span><span class=\"si\">{</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">uptime</span><span class=\"si\">}</span><span class=\"s2\">)&quot;</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"nf\">kill</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">auto_restart</span>
    <span class=\"o\">=</span> <span class=\"kc\">False</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">proc</span><span
    class=\"o\">.</span><span class=\"n\">stdout</span><span class=\"o\">.</span><span
    class=\"n\">close</span><span class=\"p\">()</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">proc</span><span class=\"o\">.</span><span
    class=\"n\">stderr</span><span class=\"o\">.</span><span class=\"n\">close</span><span
    class=\"p\">()</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">proc</span><span class=\"o\">.</span><span class=\"n\">kill</span><span
    class=\"p\">()</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">proc</span><span class=\"o\">.</span><span class=\"n\">wait</span><span
    class=\"p\">()</span>\n\n            <span class=\"k\">def</span> <span class=\"nf\">__rich__</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"n\">Panel</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">if</span> <span class=\"ow\">not</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">proc</span><span
    class=\"o\">.</span><span class=\"n\">poll</span><span class=\"p\">():</span>\n
    \                   <span class=\"k\">return</span> <span class=\"n\">Panel</span><span
    class=\"p\">(</span>\n                        <span class=\"p\">(</span>\n                            <span
    class=\"sa\">f</span><span class=\"s2\">&quot;[green]serving on port: [gold1]</span><span
    class=\"si\">{</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">port</span><span class=\"si\">}</span><span class=\"s2\"> &quot;</span>\n
    \                           <span class=\"sa\">f</span><span class=\"s2\">&quot;[green]using
    pid: [gold1]</span><span class=\"si\">{</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">proc</span><span class=\"o\">.</span><span
    class=\"n\">pid</span><span class=\"si\">}</span><span class=\"s2\"> &quot;</span>\n
    \                           <span class=\"sa\">f</span><span class=\"s2\">&quot;[green]uptime:
    [gold1]</span><span class=\"si\">{</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">uptime</span><span class=\"si\">}</span><span class=\"s2\"> &quot;</span>\n
    \                           <span class=\"sa\">f</span><span class=\"s2\">&quot;[green]link:
    [gold1] http://localhost:</span><span class=\"si\">{</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">port</span><span class=\"si\">}</span><span
    class=\"s2\">[/]&quot;</span>\n                        <span class=\"p\">),</span>\n
    \                       <span class=\"n\">border_style</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;blue&quot;</span><span class=\"p\">,</span>\n                        <span
    class=\"n\">title</span><span class=\"o\">=</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">title</span><span class=\"p\">,</span>\n
    \                       <span class=\"n\">expand</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span>\n                    <span
    class=\"p\">)</span>\n\n                <span class=\"k\">else</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"n\">Panel</span><span
    class=\"p\">(</span>\n                        <span class=\"s2\">&quot;[red]server
    died&quot;</span><span class=\"p\">,</span>\n                        <span class=\"n\">title</span><span
    class=\"o\">=</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">title</span><span class=\"p\">,</span>\n                        <span
    class=\"n\">border_style</span><span class=\"o\">=</span><span class=\"s2\">&quot;red&quot;</span><span
    class=\"p\">,</span>\n                        <span class=\"n\">expand</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \                   <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='configure' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>configure <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">configure
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">configure</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">def</span> <span class=\"nf\">get_server</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">try</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_server</span>\n                <span class=\"k\">except</span> <span
    class=\"ne\">AttributeError</span><span class=\"p\">:</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_server</span><span
    class=\"p\">:</span> <span class=\"n\">Server</span> <span class=\"o\">=</span>
    <span class=\"n\">Server</span><span class=\"p\">(</span><span class=\"n\">directory</span><span
    class=\"o\">=</span><span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;output_dir&quot;</span><span class=\"p\">]))</span>\n
    \                   <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_server</span>\n\n            <span class=\"kn\">from</span>
    <span class=\"nn\">markata</span> <span class=\"kn\">import</span> <span class=\"n\">Markata</span>\n\n
    \           <span class=\"n\">Markata</span><span class=\"o\">.</span><span class=\"n\">server</span>
    <span class=\"o\">=</span> <span class=\"nb\">property</span><span class=\"p\">(</span><span
    class=\"n\">get_server</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='run_server' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>run_server <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">run_server
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">run_server</span><span class=\"p\">()</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n            <span class=\"kn\">from</span>
    <span class=\"nn\">rich.live</span> <span class=\"kn\">import</span> <span class=\"n\">Live</span>\n\n
    \           <span class=\"kn\">from</span> <span class=\"nn\">.cli</span> <span
    class=\"kn\">import</span> <span class=\"n\">run_until_keyboard_interrupt</span>\n\n
    \           <span class=\"k\">with</span> <span class=\"n\">Live</span><span class=\"p\">(</span><span
    class=\"n\">Server</span><span class=\"p\">(),</span> <span class=\"n\">refresh_per_second</span><span
    class=\"o\">=</span><span class=\"mi\">1</span><span class=\"p\">,</span> <span
    class=\"n\">screen</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">):</span>\n                <span class=\"n\">run_until_keyboard_interrupt</span><span
    class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='cli' class='admonition-title'
    style='margin:0;padding:.5rem 1rem;'>cli <em class='small'>function</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">cli
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">cli</span><span class=\"p\">(</span><span class=\"n\">app</span><span
    class=\"p\">:</span> <span class=\"n\">typer</span><span class=\"o\">.</span><span
    class=\"n\">Typer</span><span class=\"p\">,</span> <span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">server_app</span> <span class=\"o\">=</span> <span
    class=\"n\">typer</span><span class=\"o\">.</span><span class=\"n\">Typer</span><span
    class=\"p\">()</span>\n            <span class=\"n\">app</span><span class=\"o\">.</span><span
    class=\"n\">add_typer</span><span class=\"p\">(</span><span class=\"n\">server_app</span><span
    class=\"p\">)</span>\n\n            <span class=\"nd\">@server_app</span><span
    class=\"o\">.</span><span class=\"n\">callback</span><span class=\"p\">(</span><span
    class=\"n\">invoke_without_command</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n            <span class=\"k\">def</span> <span class=\"nf\">serve</span><span
    class=\"p\">():</span>\n<span class=\"w\">                </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">                Serve the site locally.</span>\n<span class=\"sd\">
    \               &quot;&quot;&quot;</span>\n                <span class=\"n\">run_server</span><span
    class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='__init__'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'><strong>init</strong>
    <em class='small'>method</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\"><strong>init</strong> <em class='small'>source</em></p>\n</div>\n<pre
    class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button class='copy' title='copy
    code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"fm\">__init__</span><span class=\"p\">(</span>\n                <span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Server&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"o\">*</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">auto_restart</span><span class=\"p\">:</span>
    <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span class=\"kc\">True</span><span
    class=\"p\">,</span>\n                <span class=\"n\">directory</span><span
    class=\"p\">:</span> <span class=\"n\">Union</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">,</span> <span class=\"s2\">&quot;Path&quot;</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"kc\">None</span><span
    class=\"p\">,</span>\n                <span class=\"n\">port</span><span class=\"p\">:</span>
    <span class=\"nb\">int</span> <span class=\"o\">=</span> <span class=\"mi\">8000</span><span
    class=\"p\">,</span>\n            <span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
    class=\"k\">if</span> <span class=\"n\">directory</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                    <span
    class=\"kn\">from</span> <span class=\"nn\">markata</span> <span class=\"kn\">import</span>
    <span class=\"n\">Markata</span>\n\n                    <span class=\"n\">m</span>
    <span class=\"o\">=</span> <span class=\"n\">Markata</span><span class=\"p\">()</span>\n
    \                   <span class=\"n\">directory</span> <span class=\"o\">=</span>
    <span class=\"n\">Path</span><span class=\"p\">(</span><span class=\"nb\">str</span><span
    class=\"p\">(</span><span class=\"n\">m</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"p\">[</span><span class=\"s2\">&quot;output_dir&quot;</span><span
    class=\"p\">]))</span>\n\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">auto_restart</span> <span class=\"o\">=</span>
    <span class=\"n\">auto_restart</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">directory</span> <span class=\"o\">=</span>
    <span class=\"n\">directory</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">port</span> <span class=\"o\">=</span> <span
    class=\"n\">find_port</span><span class=\"p\">(</span><span class=\"n\">port</span><span
    class=\"o\">=</span><span class=\"n\">port</span><span class=\"p\">)</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">start_server</span><span class=\"p\">()</span>\n                <span
    class=\"n\">atexit</span><span class=\"o\">.</span><span class=\"n\">register</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">kill</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='start_server' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>start_server <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">start_server
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">start_server</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"kn\">import</span> <span class=\"nn\">subprocess</span>\n\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">cmd</span> <span class=\"o\">=</span> <span class=\"p\">[</span>\n
    \                   <span class=\"s2\">&quot;python&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;-m&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;http.server&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">port</span><span
    class=\"p\">),</span>\n                    <span class=\"s2\">&quot;--directory&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">directory</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">]</span>\n\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">proc</span> <span class=\"o\">=</span> <span
    class=\"n\">subprocess</span><span class=\"o\">.</span><span class=\"n\">Popen</span><span
    class=\"p\">(</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">cmd</span><span class=\"p\">,</span>\n                    <span
    class=\"n\">stderr</span><span class=\"o\">=</span><span class=\"n\">subprocess</span><span
    class=\"o\">.</span><span class=\"n\">PIPE</span><span class=\"p\">,</span>\n
    \                   <span class=\"n\">stdout</span><span class=\"o\">=</span><span
    class=\"n\">subprocess</span><span class=\"o\">.</span><span class=\"n\">PIPE</span><span
    class=\"p\">,</span>\n                <span class=\"p\">)</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">start_time</span>
    <span class=\"o\">=</span> <span class=\"n\">time</span><span class=\"o\">.</span><span
    class=\"n\">time</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='uptime' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>uptime <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">uptime
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">uptime</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">int</span><span
    class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"nb\">round</span><span
    class=\"p\">(</span><span class=\"n\">time</span><span class=\"o\">.</span><span
    class=\"n\">time</span><span class=\"p\">()</span> <span class=\"o\">-</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">start_time</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='title' class='admonition-title'
    style='margin:0;padding:.5rem 1rem;'>title <em class='small'>method</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">title
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">title</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span
    class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;server (</span><span class=\"si\">{</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">uptime</span><span class=\"si\">}</span><span
    class=\"s2\">)&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='kill'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>kill <em class='small'>method</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">kill
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">kill</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">auto_restart</span> <span class=\"o\">=</span> <span class=\"kc\">False</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">proc</span><span class=\"o\">.</span><span class=\"n\">stdout</span><span
    class=\"o\">.</span><span class=\"n\">close</span><span class=\"p\">()</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">proc</span><span class=\"o\">.</span><span class=\"n\">stderr</span><span
    class=\"o\">.</span><span class=\"n\">close</span><span class=\"p\">()</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">proc</span><span class=\"o\">.</span><span class=\"n\">kill</span><span
    class=\"p\">()</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">proc</span><span class=\"o\">.</span><span class=\"n\">wait</span><span
    class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='__rich__'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'><strong>rich</strong>
    <em class='small'>method</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\"><strong>rich</strong> <em class='small'>source</em></p>\n</div>\n<pre
    class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button class='copy' title='copy
    code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">__rich__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Panel</span><span
    class=\"p\">:</span>\n                <span class=\"k\">if</span> <span class=\"ow\">not</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">proc</span><span
    class=\"o\">.</span><span class=\"n\">poll</span><span class=\"p\">():</span>\n
    \                   <span class=\"k\">return</span> <span class=\"n\">Panel</span><span
    class=\"p\">(</span>\n                        <span class=\"p\">(</span>\n                            <span
    class=\"sa\">f</span><span class=\"s2\">&quot;[green]serving on port: [gold1]</span><span
    class=\"si\">{</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">port</span><span class=\"si\">}</span><span class=\"s2\"> &quot;</span>\n
    \                           <span class=\"sa\">f</span><span class=\"s2\">&quot;[green]using
    pid: [gold1]</span><span class=\"si\">{</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">proc</span><span class=\"o\">.</span><span
    class=\"n\">pid</span><span class=\"si\">}</span><span class=\"s2\"> &quot;</span>\n
    \                           <span class=\"sa\">f</span><span class=\"s2\">&quot;[green]uptime:
    [gold1]</span><span class=\"si\">{</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">uptime</span><span class=\"si\">}</span><span class=\"s2\"> &quot;</span>\n
    \                           <span class=\"sa\">f</span><span class=\"s2\">&quot;[green]link:
    [gold1] http://localhost:</span><span class=\"si\">{</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">port</span><span class=\"si\">}</span><span
    class=\"s2\">[/]&quot;</span>\n                        <span class=\"p\">),</span>\n
    \                       <span class=\"n\">border_style</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;blue&quot;</span><span class=\"p\">,</span>\n                        <span
    class=\"n\">title</span><span class=\"o\">=</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">title</span><span class=\"p\">,</span>\n
    \                       <span class=\"n\">expand</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span>\n                    <span
    class=\"p\">)</span>\n\n                <span class=\"k\">else</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"n\">Panel</span><span
    class=\"p\">(</span>\n                        <span class=\"s2\">&quot;[red]server
    died&quot;</span><span class=\"p\">,</span>\n                        <span class=\"n\">title</span><span
    class=\"o\">=</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">title</span><span class=\"p\">,</span>\n                        <span
    class=\"n\">border_style</span><span class=\"o\">=</span><span class=\"s2\">&quot;red&quot;</span><span
    class=\"p\">,</span>\n                        <span class=\"n\">expand</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \                   <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='get_server' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>get_server <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">get_server
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">get_server</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">):</span>\n                <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_server</span>\n                <span class=\"k\">except</span>
    <span class=\"ne\">AttributeError</span><span class=\"p\">:</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_server</span><span
    class=\"p\">:</span> <span class=\"n\">Server</span> <span class=\"o\">=</span>
    <span class=\"n\">Server</span><span class=\"p\">(</span><span class=\"n\">directory</span><span
    class=\"o\">=</span><span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;output_dir&quot;</span><span class=\"p\">]))</span>\n
    \                   <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_server</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='serve' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>serve <em class='small'>function</em></h2>\nServe the site locally.</p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">serve
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">serve</span><span class=\"p\">():</span>\n<span class=\"w\">
    \               </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span class=\"sd\">
    \               Serve the site locally.</span>\n<span class=\"sd\">                &quot;&quot;&quot;</span>\n
    \               <span class=\"n\">run_server</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n\n
    \   </section>\n</article>    </div>\n    <div></div>\n</div>\n     </body>\n</html>"
  og: "<!DOCTYPE html>\n<html lang=\"en\">\n<title>Server.Py</title>\n<meta charset=\"UTF-8\"
    />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n<meta
    name=\"description\" content=\"None ! ! ???+ source  ! ???+ source  ! ???+ source
    \ ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  !
    ???+ source  ! ???+ source  ! ?\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    <head>\n<title>Server.Py</title>\n<meta
    charset=\"UTF-8\" />\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"
    />\n<meta name=\"description\" content=\"None ! ! ???+ source  ! ???+ source  !
    ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
    source  ! ???+ source  ! ???+ source  ! ?\" />\n <link href=\"/favicon.ico\" rel=\"icon\"
    type=\"image/png\" />\n\n<link rel=\"stylesheet\" href=\"/post.css\" />\n<link
    rel=\"stylesheet\" href=\"/app.css\" />\n<script src=\"/theme.js\"></script>\n\n\n<meta
    name=\"og:author_email\" content=\"waylon@waylonwalker.com\" />\n\n    </head>\n
    \   <body>\n<article style=\"text-align: center;\">\n    <style>\n        section
    {\n            font-size: 200%;\n        }\n\n\n        .edit {\n            display:
    none;\n        }\n    </style>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Server.Py\n    </h1>\n</section></article>\n     </body>\n</html>"
  partial: "<article class='w-full'>\n<section class=\"title\">\n    <h1 id=\"title\">\n
    \       Server.Py\n    </h1>\n</section>    <section class=\"body\">\n        <p>None</p>\n<p>!!
    function <h2 id='find_port' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>find_port <em class='small'>function</em></h2>\nFind a port not in ues
    starting at given port</p>\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
    class=\"admonition-title\">find_port <em class='small'>source</em></p>\n</div>\n<pre
    class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button class='copy' title='copy
    code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">find_port</span><span class=\"p\">(</span><span class=\"n\">port</span><span
    class=\"p\">:</span> <span class=\"nb\">int</span> <span class=\"o\">=</span>
    <span class=\"mi\">8000</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"nb\">int</span><span class=\"p\">:</span>\n<span class=\"w\">            </span><span
    class=\"sd\">&quot;&quot;&quot;Find a port not in ues starting at given port&quot;&quot;&quot;</span>\n
    \           <span class=\"kn\">import</span> <span class=\"nn\">socket</span>\n\n
    \           <span class=\"k\">with</span> <span class=\"n\">socket</span><span
    class=\"o\">.</span><span class=\"n\">socket</span><span class=\"p\">(</span><span
    class=\"n\">socket</span><span class=\"o\">.</span><span class=\"n\">AF_INET</span><span
    class=\"p\">,</span> <span class=\"n\">socket</span><span class=\"o\">.</span><span
    class=\"n\">SOCK_STREAM</span><span class=\"p\">)</span> <span class=\"k\">as</span>
    <span class=\"n\">s</span><span class=\"p\">:</span>\n                <span class=\"k\">if</span>
    <span class=\"n\">s</span><span class=\"o\">.</span><span class=\"n\">connect_ex</span><span
    class=\"p\">((</span><span class=\"s2\">&quot;localhost&quot;</span><span class=\"p\">,</span>
    <span class=\"n\">port</span><span class=\"p\">))</span> <span class=\"o\">==</span>
    <span class=\"mi\">0</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">return</span> <span class=\"n\">find_port</span><span class=\"p\">(</span><span
    class=\"n\">port</span><span class=\"o\">=</span><span class=\"n\">port</span>
    <span class=\"o\">+</span> <span class=\"mi\">1</span><span class=\"p\">)</span>\n
    \               <span class=\"k\">return</span> <span class=\"n\">port</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    class <h2 id='Server' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Server
    <em class='small'>class</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\">Server <em class='small'>source</em></p>\n</div>\n<pre
    class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button class='copy' title='copy
    code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
    <span class=\"nc\">Server</span><span class=\"p\">:</span>\n            <span
    class=\"k\">def</span> <span class=\"fm\">__init__</span><span class=\"p\">(</span>\n
    \               <span class=\"bp\">self</span><span class=\"p\">:</span> <span
    class=\"s2\">&quot;Server&quot;</span><span class=\"p\">,</span>\n                <span
    class=\"o\">*</span><span class=\"p\">,</span>\n                <span class=\"n\">auto_restart</span><span
    class=\"p\">:</span> <span class=\"nb\">bool</span> <span class=\"o\">=</span>
    <span class=\"kc\">True</span><span class=\"p\">,</span>\n                <span
    class=\"n\">directory</span><span class=\"p\">:</span> <span class=\"n\">Union</span><span
    class=\"p\">[</span><span class=\"nb\">str</span><span class=\"p\">,</span> <span
    class=\"s2\">&quot;Path&quot;</span><span class=\"p\">]</span> <span class=\"o\">=</span>
    <span class=\"kc\">None</span><span class=\"p\">,</span>\n                <span
    class=\"n\">port</span><span class=\"p\">:</span> <span class=\"nb\">int</span>
    <span class=\"o\">=</span> <span class=\"mi\">8000</span><span class=\"p\">,</span>\n
    \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"k\">if</span> <span class=\"n\">directory</span>
    <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \                   <span class=\"kn\">from</span> <span class=\"nn\">markata</span>
    <span class=\"kn\">import</span> <span class=\"n\">Markata</span>\n\n                    <span
    class=\"n\">m</span> <span class=\"o\">=</span> <span class=\"n\">Markata</span><span
    class=\"p\">()</span>\n                    <span class=\"n\">directory</span>
    <span class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
    class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">m</span><span
    class=\"o\">.</span><span class=\"n\">config</span><span class=\"p\">[</span><span
    class=\"s2\">&quot;output_dir&quot;</span><span class=\"p\">]))</span>\n\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">auto_restart</span>
    <span class=\"o\">=</span> <span class=\"n\">auto_restart</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">directory</span>
    <span class=\"o\">=</span> <span class=\"n\">directory</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">port</span>
    <span class=\"o\">=</span> <span class=\"n\">find_port</span><span class=\"p\">(</span><span
    class=\"n\">port</span><span class=\"o\">=</span><span class=\"n\">port</span><span
    class=\"p\">)</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">start_server</span><span class=\"p\">()</span>\n                <span
    class=\"n\">atexit</span><span class=\"o\">.</span><span class=\"n\">register</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">kill</span><span class=\"p\">)</span>\n\n            <span class=\"k\">def</span>
    <span class=\"nf\">start_server</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"kn\">import</span> <span class=\"nn\">subprocess</span>\n\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">cmd</span> <span class=\"o\">=</span> <span class=\"p\">[</span>\n
    \                   <span class=\"s2\">&quot;python&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;-m&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;http.server&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">port</span><span
    class=\"p\">),</span>\n                    <span class=\"s2\">&quot;--directory&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">directory</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">]</span>\n\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">proc</span> <span class=\"o\">=</span> <span
    class=\"n\">subprocess</span><span class=\"o\">.</span><span class=\"n\">Popen</span><span
    class=\"p\">(</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">cmd</span><span class=\"p\">,</span>\n                    <span
    class=\"n\">stderr</span><span class=\"o\">=</span><span class=\"n\">subprocess</span><span
    class=\"o\">.</span><span class=\"n\">PIPE</span><span class=\"p\">,</span>\n
    \                   <span class=\"n\">stdout</span><span class=\"o\">=</span><span
    class=\"n\">subprocess</span><span class=\"o\">.</span><span class=\"n\">PIPE</span><span
    class=\"p\">,</span>\n                <span class=\"p\">)</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">start_time</span>
    <span class=\"o\">=</span> <span class=\"n\">time</span><span class=\"o\">.</span><span
    class=\"n\">time</span><span class=\"p\">()</span>\n\n            <span class=\"nd\">@property</span>\n
    \           <span class=\"k\">def</span> <span class=\"nf\">uptime</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"nb\">int</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">return</span> <span class=\"nb\">round</span><span
    class=\"p\">(</span><span class=\"n\">time</span><span class=\"o\">.</span><span
    class=\"n\">time</span><span class=\"p\">()</span> <span class=\"o\">-</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">start_time</span><span
    class=\"p\">)</span>\n\n            <span class=\"nd\">@property</span>\n            <span
    class=\"k\">def</span> <span class=\"nf\">title</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"nb\">str</span><span class=\"p\">:</span>\n                <span
    class=\"k\">return</span> <span class=\"sa\">f</span><span class=\"s2\">&quot;server
    (</span><span class=\"si\">{</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">uptime</span><span class=\"si\">}</span><span class=\"s2\">)&quot;</span>\n\n
    \           <span class=\"k\">def</span> <span class=\"nf\">kill</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">auto_restart</span>
    <span class=\"o\">=</span> <span class=\"kc\">False</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">proc</span><span
    class=\"o\">.</span><span class=\"n\">stdout</span><span class=\"o\">.</span><span
    class=\"n\">close</span><span class=\"p\">()</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">proc</span><span class=\"o\">.</span><span
    class=\"n\">stderr</span><span class=\"o\">.</span><span class=\"n\">close</span><span
    class=\"p\">()</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">proc</span><span class=\"o\">.</span><span class=\"n\">kill</span><span
    class=\"p\">()</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">proc</span><span class=\"o\">.</span><span class=\"n\">wait</span><span
    class=\"p\">()</span>\n\n            <span class=\"k\">def</span> <span class=\"nf\">__rich__</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">)</span> <span
    class=\"o\">-&gt;</span> <span class=\"n\">Panel</span><span class=\"p\">:</span>\n
    \               <span class=\"k\">if</span> <span class=\"ow\">not</span> <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">proc</span><span
    class=\"o\">.</span><span class=\"n\">poll</span><span class=\"p\">():</span>\n
    \                   <span class=\"k\">return</span> <span class=\"n\">Panel</span><span
    class=\"p\">(</span>\n                        <span class=\"p\">(</span>\n                            <span
    class=\"sa\">f</span><span class=\"s2\">&quot;[green]serving on port: [gold1]</span><span
    class=\"si\">{</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">port</span><span class=\"si\">}</span><span class=\"s2\"> &quot;</span>\n
    \                           <span class=\"sa\">f</span><span class=\"s2\">&quot;[green]using
    pid: [gold1]</span><span class=\"si\">{</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">proc</span><span class=\"o\">.</span><span
    class=\"n\">pid</span><span class=\"si\">}</span><span class=\"s2\"> &quot;</span>\n
    \                           <span class=\"sa\">f</span><span class=\"s2\">&quot;[green]uptime:
    [gold1]</span><span class=\"si\">{</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">uptime</span><span class=\"si\">}</span><span class=\"s2\"> &quot;</span>\n
    \                           <span class=\"sa\">f</span><span class=\"s2\">&quot;[green]link:
    [gold1] http://localhost:</span><span class=\"si\">{</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">port</span><span class=\"si\">}</span><span
    class=\"s2\">[/]&quot;</span>\n                        <span class=\"p\">),</span>\n
    \                       <span class=\"n\">border_style</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;blue&quot;</span><span class=\"p\">,</span>\n                        <span
    class=\"n\">title</span><span class=\"o\">=</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">title</span><span class=\"p\">,</span>\n
    \                       <span class=\"n\">expand</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span>\n                    <span
    class=\"p\">)</span>\n\n                <span class=\"k\">else</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"n\">Panel</span><span
    class=\"p\">(</span>\n                        <span class=\"s2\">&quot;[red]server
    died&quot;</span><span class=\"p\">,</span>\n                        <span class=\"n\">title</span><span
    class=\"o\">=</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">title</span><span class=\"p\">,</span>\n                        <span
    class=\"n\">border_style</span><span class=\"o\">=</span><span class=\"s2\">&quot;red&quot;</span><span
    class=\"p\">,</span>\n                        <span class=\"n\">expand</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \                   <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='configure' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>configure <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">configure
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">configure</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"k\">def</span> <span class=\"nf\">get_server</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">):</span>\n
    \               <span class=\"k\">try</span><span class=\"p\">:</span>\n                    <span
    class=\"k\">return</span> <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">_server</span>\n                <span class=\"k\">except</span> <span
    class=\"ne\">AttributeError</span><span class=\"p\">:</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_server</span><span
    class=\"p\">:</span> <span class=\"n\">Server</span> <span class=\"o\">=</span>
    <span class=\"n\">Server</span><span class=\"p\">(</span><span class=\"n\">directory</span><span
    class=\"o\">=</span><span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;output_dir&quot;</span><span class=\"p\">]))</span>\n
    \                   <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_server</span>\n\n            <span class=\"kn\">from</span>
    <span class=\"nn\">markata</span> <span class=\"kn\">import</span> <span class=\"n\">Markata</span>\n\n
    \           <span class=\"n\">Markata</span><span class=\"o\">.</span><span class=\"n\">server</span>
    <span class=\"o\">=</span> <span class=\"nb\">property</span><span class=\"p\">(</span><span
    class=\"n\">get_server</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='run_server' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>run_server <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">run_server
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">run_server</span><span class=\"p\">()</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n            <span class=\"kn\">from</span>
    <span class=\"nn\">rich.live</span> <span class=\"kn\">import</span> <span class=\"n\">Live</span>\n\n
    \           <span class=\"kn\">from</span> <span class=\"nn\">.cli</span> <span
    class=\"kn\">import</span> <span class=\"n\">run_until_keyboard_interrupt</span>\n\n
    \           <span class=\"k\">with</span> <span class=\"n\">Live</span><span class=\"p\">(</span><span
    class=\"n\">Server</span><span class=\"p\">(),</span> <span class=\"n\">refresh_per_second</span><span
    class=\"o\">=</span><span class=\"mi\">1</span><span class=\"p\">,</span> <span
    class=\"n\">screen</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">):</span>\n                <span class=\"n\">run_until_keyboard_interrupt</span><span
    class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!! function <h2 id='cli' class='admonition-title'
    style='margin:0;padding:.5rem 1rem;'>cli <em class='small'>function</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">cli
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">cli</span><span class=\"p\">(</span><span class=\"n\">app</span><span
    class=\"p\">:</span> <span class=\"n\">typer</span><span class=\"o\">.</span><span
    class=\"n\">Typer</span><span class=\"p\">,</span> <span class=\"n\">markata</span><span
    class=\"p\">:</span> <span class=\"s2\">&quot;Markata&quot;</span><span class=\"p\">)</span>
    <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
    \           <span class=\"n\">server_app</span> <span class=\"o\">=</span> <span
    class=\"n\">typer</span><span class=\"o\">.</span><span class=\"n\">Typer</span><span
    class=\"p\">()</span>\n            <span class=\"n\">app</span><span class=\"o\">.</span><span
    class=\"n\">add_typer</span><span class=\"p\">(</span><span class=\"n\">server_app</span><span
    class=\"p\">)</span>\n\n            <span class=\"nd\">@server_app</span><span
    class=\"o\">.</span><span class=\"n\">callback</span><span class=\"p\">(</span><span
    class=\"n\">invoke_without_command</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
    class=\"p\">)</span>\n            <span class=\"k\">def</span> <span class=\"nf\">serve</span><span
    class=\"p\">():</span>\n<span class=\"w\">                </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span
    class=\"sd\">                Serve the site locally.</span>\n<span class=\"sd\">
    \               &quot;&quot;&quot;</span>\n                <span class=\"n\">run_server</span><span
    class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='__init__'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'><strong>init</strong>
    <em class='small'>method</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\"><strong>init</strong> <em class='small'>source</em></p>\n</div>\n<pre
    class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button class='copy' title='copy
    code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"fm\">__init__</span><span class=\"p\">(</span>\n                <span
    class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">&quot;Server&quot;</span><span
    class=\"p\">,</span>\n                <span class=\"o\">*</span><span class=\"p\">,</span>\n
    \               <span class=\"n\">auto_restart</span><span class=\"p\">:</span>
    <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span class=\"kc\">True</span><span
    class=\"p\">,</span>\n                <span class=\"n\">directory</span><span
    class=\"p\">:</span> <span class=\"n\">Union</span><span class=\"p\">[</span><span
    class=\"nb\">str</span><span class=\"p\">,</span> <span class=\"s2\">&quot;Path&quot;</span><span
    class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"kc\">None</span><span
    class=\"p\">,</span>\n                <span class=\"n\">port</span><span class=\"p\">:</span>
    <span class=\"nb\">int</span> <span class=\"o\">=</span> <span class=\"mi\">8000</span><span
    class=\"p\">,</span>\n            <span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
    class=\"k\">if</span> <span class=\"n\">directory</span> <span class=\"ow\">is</span>
    <span class=\"kc\">None</span><span class=\"p\">:</span>\n                    <span
    class=\"kn\">from</span> <span class=\"nn\">markata</span> <span class=\"kn\">import</span>
    <span class=\"n\">Markata</span>\n\n                    <span class=\"n\">m</span>
    <span class=\"o\">=</span> <span class=\"n\">Markata</span><span class=\"p\">()</span>\n
    \                   <span class=\"n\">directory</span> <span class=\"o\">=</span>
    <span class=\"n\">Path</span><span class=\"p\">(</span><span class=\"nb\">str</span><span
    class=\"p\">(</span><span class=\"n\">m</span><span class=\"o\">.</span><span
    class=\"n\">config</span><span class=\"p\">[</span><span class=\"s2\">&quot;output_dir&quot;</span><span
    class=\"p\">]))</span>\n\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">auto_restart</span> <span class=\"o\">=</span>
    <span class=\"n\">auto_restart</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">directory</span> <span class=\"o\">=</span>
    <span class=\"n\">directory</span>\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">port</span> <span class=\"o\">=</span> <span
    class=\"n\">find_port</span><span class=\"p\">(</span><span class=\"n\">port</span><span
    class=\"o\">=</span><span class=\"n\">port</span><span class=\"p\">)</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">start_server</span><span class=\"p\">()</span>\n                <span
    class=\"n\">atexit</span><span class=\"o\">.</span><span class=\"n\">register</span><span
    class=\"p\">(</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">kill</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='start_server' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>start_server <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">start_server
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">start_server</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"kn\">import</span> <span class=\"nn\">subprocess</span>\n\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">cmd</span> <span class=\"o\">=</span> <span class=\"p\">[</span>\n
    \                   <span class=\"s2\">&quot;python&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;-m&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"s2\">&quot;http.server&quot;</span><span class=\"p\">,</span>\n
    \                   <span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">port</span><span
    class=\"p\">),</span>\n                    <span class=\"s2\">&quot;--directory&quot;</span><span
    class=\"p\">,</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">directory</span><span class=\"p\">,</span>\n
    \               <span class=\"p\">]</span>\n\n                <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">proc</span> <span class=\"o\">=</span> <span
    class=\"n\">subprocess</span><span class=\"o\">.</span><span class=\"n\">Popen</span><span
    class=\"p\">(</span>\n                    <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">cmd</span><span class=\"p\">,</span>\n                    <span
    class=\"n\">stderr</span><span class=\"o\">=</span><span class=\"n\">subprocess</span><span
    class=\"o\">.</span><span class=\"n\">PIPE</span><span class=\"p\">,</span>\n
    \                   <span class=\"n\">stdout</span><span class=\"o\">=</span><span
    class=\"n\">subprocess</span><span class=\"o\">.</span><span class=\"n\">PIPE</span><span
    class=\"p\">,</span>\n                <span class=\"p\">)</span>\n                <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">start_time</span>
    <span class=\"o\">=</span> <span class=\"n\">time</span><span class=\"o\">.</span><span
    class=\"n\">time</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    method <h2 id='uptime' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>uptime <em class='small'>method</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">uptime
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">uptime</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">int</span><span
    class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"nb\">round</span><span
    class=\"p\">(</span><span class=\"n\">time</span><span class=\"o\">.</span><span
    class=\"n\">time</span><span class=\"p\">()</span> <span class=\"o\">-</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">start_time</span><span
    class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='title' class='admonition-title'
    style='margin:0;padding:.5rem 1rem;'>title <em class='small'>method</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">title
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">title</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span
    class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"sa\">f</span><span
    class=\"s2\">&quot;server (</span><span class=\"si\">{</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">uptime</span><span class=\"si\">}</span><span
    class=\"s2\">)&quot;</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='kill'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'>kill <em class='small'>method</em></h2></p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">kill
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">kill</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
    class=\"p\">:</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">auto_restart</span> <span class=\"o\">=</span> <span class=\"kc\">False</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">proc</span><span class=\"o\">.</span><span class=\"n\">stdout</span><span
    class=\"o\">.</span><span class=\"n\">close</span><span class=\"p\">()</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">proc</span><span class=\"o\">.</span><span class=\"n\">stderr</span><span
    class=\"o\">.</span><span class=\"n\">close</span><span class=\"p\">()</span>\n
    \               <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">proc</span><span class=\"o\">.</span><span class=\"n\">kill</span><span
    class=\"p\">()</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">proc</span><span class=\"o\">.</span><span class=\"n\">wait</span><span
    class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n<p>!! method <h2 id='__rich__'
    class='admonition-title' style='margin:0;padding:.5rem 1rem;'><strong>rich</strong>
    <em class='small'>method</em></h2></p>\n<div class=\"admonition source is-collapsible
    collapsible-open\">\n<p class=\"admonition-title\"><strong>rich</strong> <em class='small'>source</em></p>\n</div>\n<pre
    class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button class='copy' title='copy
    code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">__rich__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Panel</span><span
    class=\"p\">:</span>\n                <span class=\"k\">if</span> <span class=\"ow\">not</span>
    <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">proc</span><span
    class=\"o\">.</span><span class=\"n\">poll</span><span class=\"p\">():</span>\n
    \                   <span class=\"k\">return</span> <span class=\"n\">Panel</span><span
    class=\"p\">(</span>\n                        <span class=\"p\">(</span>\n                            <span
    class=\"sa\">f</span><span class=\"s2\">&quot;[green]serving on port: [gold1]</span><span
    class=\"si\">{</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">port</span><span class=\"si\">}</span><span class=\"s2\"> &quot;</span>\n
    \                           <span class=\"sa\">f</span><span class=\"s2\">&quot;[green]using
    pid: [gold1]</span><span class=\"si\">{</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">proc</span><span class=\"o\">.</span><span
    class=\"n\">pid</span><span class=\"si\">}</span><span class=\"s2\"> &quot;</span>\n
    \                           <span class=\"sa\">f</span><span class=\"s2\">&quot;[green]uptime:
    [gold1]</span><span class=\"si\">{</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">uptime</span><span class=\"si\">}</span><span class=\"s2\"> &quot;</span>\n
    \                           <span class=\"sa\">f</span><span class=\"s2\">&quot;[green]link:
    [gold1] http://localhost:</span><span class=\"si\">{</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">port</span><span class=\"si\">}</span><span
    class=\"s2\">[/]&quot;</span>\n                        <span class=\"p\">),</span>\n
    \                       <span class=\"n\">border_style</span><span class=\"o\">=</span><span
    class=\"s2\">&quot;blue&quot;</span><span class=\"p\">,</span>\n                        <span
    class=\"n\">title</span><span class=\"o\">=</span><span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">title</span><span class=\"p\">,</span>\n
    \                       <span class=\"n\">expand</span><span class=\"o\">=</span><span
    class=\"kc\">True</span><span class=\"p\">,</span>\n                    <span
    class=\"p\">)</span>\n\n                <span class=\"k\">else</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"n\">Panel</span><span
    class=\"p\">(</span>\n                        <span class=\"s2\">&quot;[red]server
    died&quot;</span><span class=\"p\">,</span>\n                        <span class=\"n\">title</span><span
    class=\"o\">=</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
    class=\"n\">title</span><span class=\"p\">,</span>\n                        <span
    class=\"n\">border_style</span><span class=\"o\">=</span><span class=\"s2\">&quot;red&quot;</span><span
    class=\"p\">,</span>\n                        <span class=\"n\">expand</span><span
    class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n
    \                   <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='get_server' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>get_server <em class='small'>function</em></h2></p>\n<div class=\"admonition
    source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">get_server
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">get_server</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
    class=\"p\">):</span>\n                <span class=\"k\">try</span><span class=\"p\">:</span>\n
    \                   <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_server</span>\n                <span class=\"k\">except</span>
    <span class=\"ne\">AttributeError</span><span class=\"p\">:</span>\n                    <span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">_server</span><span
    class=\"p\">:</span> <span class=\"n\">Server</span> <span class=\"o\">=</span>
    <span class=\"n\">Server</span><span class=\"p\">(</span><span class=\"n\">directory</span><span
    class=\"o\">=</span><span class=\"nb\">str</span><span class=\"p\">(</span><span
    class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">config</span><span
    class=\"p\">[</span><span class=\"s2\">&quot;output_dir&quot;</span><span class=\"p\">]))</span>\n
    \                   <span class=\"k\">return</span> <span class=\"bp\">self</span><span
    class=\"o\">.</span><span class=\"n\">_server</span>\n</pre></div>\n\n</pre>\n\n<p>!!
    function <h2 id='serve' class='admonition-title' style='margin:0;padding:.5rem
    1rem;'>serve <em class='small'>function</em></h2>\nServe the site locally.</p>\n<div
    class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">serve
    <em class='small'>source</em></p>\n</div>\n<pre class='wrapper'>\n\n<div class='copy-wrapper'>\n\n<button
    class='copy' title='copy code to clipboard' onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"><svg
    version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
    x=\"0px\" y=\"0px\" viewBox=\"0 0 115.77 122.88\" style=\"enable-background:new
    0 0 115.77 122.88\" xml:space=\"preserve\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
    class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
    v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
    c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
    c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
    h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
    M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
    v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
    M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
    c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
    v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"/></g></svg></button>\n</div>\n
    \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
    <span class=\"nf\">serve</span><span class=\"p\">():</span>\n<span class=\"w\">
    \               </span><span class=\"sd\">&quot;&quot;&quot;</span>\n<span class=\"sd\">
    \               Serve the site locally.</span>\n<span class=\"sd\">                &quot;&quot;&quot;</span>\n
    \               <span class=\"n\">run_server</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n\n\n
    \   </section>\n</article>"
  raw.md: ''
published: true
slug: markata/cli/server
title: Server.Py


---

None


!! function <h2 id='find_port' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>find_port <em class='small'>function</em></h2>
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


!! class <h2 id='Server' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Server <em class='small'>class</em></h2>

???+ source "Server <em class='small'>source</em>"

```python

        class Server:
            def __init__(
                self: "Server",
                *,
                auto_restart: bool = True,
                directory: Union[str, "Path"] = None,
                port: int = 8000,
            ) -> None:
                if directory is None:
                    from markata import Markata

                    m = Markata()
                    directory = Path(str(m.config["output_dir"]))

                self.auto_restart = auto_restart
                self.directory = directory
                self.port = find_port(port=port)
                self.start_server()
                atexit.register(self.kill)

            def start_server(self) -> None:
                import subprocess

                self.cmd = [
                    "python",
                    "-m",
                    "http.server",
                    str(self.port),
                    "--directory",
                    self.directory,
                ]

                self.proc = subprocess.Popen(
                    self.cmd,
                    stderr=subprocess.PIPE,
                    stdout=subprocess.PIPE,
                )
                self.start_time = time.time()

            @property
            def uptime(self) -> int:
                return round(time.time() - self.start_time)

            @property
            def title(self) -> str:
                return f"server ({self.uptime})"

            def kill(self) -> None:
                self.auto_restart = False
                self.proc.stdout.close()
                self.proc.stderr.close()
                self.proc.kill()
                self.proc.wait()

            def __rich__(self) -> Panel:
                if not self.proc.poll():
                    return Panel(
                        (
                            f"[green]serving on port: [gold1]{self.port} "
                            f"[green]using pid: [gold1]{self.proc.pid} "
                            f"[green]uptime: [gold1]{self.uptime} "
                            f"[green]link: [gold1] http://localhost:{self.port}[/]"
                        ),
                        border_style="blue",
                        title=self.title,
                        expand=True,
                    )

                else:
                    return Panel(
                        "[red]server died",
                        title=self.title,
                        border_style="red",
                        expand=True,
                    )
```


!! function <h2 id='configure' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>configure <em class='small'>function</em></h2>

???+ source "configure <em class='small'>source</em>"

```python

        def configure(markata: "Markata") -> None:
            def get_server(self):
                try:
                    return self._server
                except AttributeError:
                    self._server: Server = Server(directory=str(self.config["output_dir"]))
                    return self._server

            from markata import Markata

            Markata.server = property(get_server)
```


!! function <h2 id='run_server' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>run_server <em class='small'>function</em></h2>

???+ source "run_server <em class='small'>source</em>"

```python

        def run_server() -> None:
            from rich.live import Live

            from .cli import run_until_keyboard_interrupt

            with Live(Server(), refresh_per_second=1, screen=True):
                run_until_keyboard_interrupt()
```


!! function <h2 id='cli' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>cli <em class='small'>function</em></h2>

???+ source "cli <em class='small'>source</em>"

```python

        def cli(app: typer.Typer, markata: "Markata") -> None:
            server_app = typer.Typer()
            app.add_typer(server_app)

            @server_app.callback(invoke_without_command=True)
            def serve():
                """
                Serve the site locally.
                """
                run_server()
```


!! method <h2 id='__init__' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>__init__ <em class='small'>method</em></h2>

???+ source "__init__ <em class='small'>source</em>"

```python

        def __init__(
                self: "Server",
                *,
                auto_restart: bool = True,
                directory: Union[str, "Path"] = None,
                port: int = 8000,
            ) -> None:
                if directory is None:
                    from markata import Markata

                    m = Markata()
                    directory = Path(str(m.config["output_dir"]))

                self.auto_restart = auto_restart
                self.directory = directory
                self.port = find_port(port=port)
                self.start_server()
                atexit.register(self.kill)
```


!! method <h2 id='start_server' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>start_server <em class='small'>method</em></h2>

???+ source "start_server <em class='small'>source</em>"

```python

        def start_server(self) -> None:
                import subprocess

                self.cmd = [
                    "python",
                    "-m",
                    "http.server",
                    str(self.port),
                    "--directory",
                    self.directory,
                ]

                self.proc = subprocess.Popen(
                    self.cmd,
                    stderr=subprocess.PIPE,
                    stdout=subprocess.PIPE,
                )
                self.start_time = time.time()
```


!! method <h2 id='uptime' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>uptime <em class='small'>method</em></h2>

???+ source "uptime <em class='small'>source</em>"

```python

        def uptime(self) -> int:
                return round(time.time() - self.start_time)
```


!! method <h2 id='title' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>title <em class='small'>method</em></h2>

???+ source "title <em class='small'>source</em>"

```python

        def title(self) -> str:
                return f"server ({self.uptime})"
```


!! method <h2 id='kill' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>kill <em class='small'>method</em></h2>

???+ source "kill <em class='small'>source</em>"

```python

        def kill(self) -> None:
                self.auto_restart = False
                self.proc.stdout.close()
                self.proc.stderr.close()
                self.proc.kill()
                self.proc.wait()
```


!! method <h2 id='__rich__' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>__rich__ <em class='small'>method</em></h2>

???+ source "__rich__ <em class='small'>source</em>"

```python

        def __rich__(self) -> Panel:
                if not self.proc.poll():
                    return Panel(
                        (
                            f"[green]serving on port: [gold1]{self.port} "
                            f"[green]using pid: [gold1]{self.proc.pid} "
                            f"[green]uptime: [gold1]{self.uptime} "
                            f"[green]link: [gold1] http://localhost:{self.port}[/]"
                        ),
                        border_style="blue",
                        title=self.title,
                        expand=True,
                    )

                else:
                    return Panel(
                        "[red]server died",
                        title=self.title,
                        border_style="red",
                        expand=True,
                    )
```


!! function <h2 id='get_server' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>get_server <em class='small'>function</em></h2>

???+ source "get_server <em class='small'>source</em>"

```python

        def get_server(self):
                try:
                    return self._server
                except AttributeError:
                    self._server: Server = Server(directory=str(self.config["output_dir"]))
                    return self._server
```


!! function <h2 id='serve' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>serve <em class='small'>function</em></h2>
    Serve the site locally.
???+ source "serve <em class='small'>source</em>"

```python

        def serve():
                """
                Serve the site locally.
                """
                run_server()
```

