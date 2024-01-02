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
  \               run_server()\n```\n"
date: 0001-01-01
description: None ! ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  !
  ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
  source  ! ?
html: "<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n<title>Server.Py</title>\n<meta
  charset=\"utf-8\"/>\n<meta content=\"width=device-width, initial-scale=1\" name=\"viewport\"/>\n<meta
  content=\"None ! ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
  source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source
  \ ! ?\" name=\"description\"/>\n<link href=\"/static/favicon.ico\" rel=\"icon\"
  type=\"image/png\"/>\n<script>\n        function setTheme(theme) {\n            document.documentElement.setAttribute(\"data-theme\",
  theme);\n        }\n\n        function detectColorSchemeOnLoad() {\n            //local
  storage is used to override OS theme settings\n            if (localStorage.getItem(\"theme\"))
  {\n                if (localStorage.getItem(\"theme\") == \"dark\") {\n                    setTheme(\"dark\");\n
  \               } else if (localStorage.getItem(\"theme\") == \"light\") {\n                    setTheme(\"light\");\n
  \               }\n            } else if (!window.matchMedia) {\n                //matchMedia
  method not supported\n                setTheme(\"light\");\n                return
  false;\n            } else if (window.matchMedia(\"(prefers-color-scheme: dark)\").matches)
  {\n                //OS theme setting detected as dark\n                setTheme(\"dark\");\n
  \           } else {\n                setTheme(\"light\");\n            }\n        }\n
  \       detectColorSchemeOnLoad();\n        document.addEventListener(\n            \"DOMContentLoaded\",\n
  \           function () {\n                //identify the toggle switch HTML element\n
  \               const toggleSwitch = document.querySelector(\n                    '#theme-switch
  input[type=\"checkbox\"]',\n                );\n\n                //function that
  changes the theme, and sets a localStorage variable to track the theme between page
  loads\n                function switchTheme(e) {\n                    if (e.target.checked)
  {\n                        localStorage.setItem(\"theme\", \"dark\");\n                        document.documentElement.setAttribute(\"data-theme\",
  \"dark\");\n                        toggleSwitch.checked = true;\n                    }
  else {\n                        localStorage.setItem(\"theme\", \"light\");\n                        document.documentElement.setAttribute(\"data-theme\",
  \"light\");\n                        toggleSwitch.checked = false;\n                    }\n
  \               }\n\n                //listener for changing themes\n                toggleSwitch.addEventListener(\"change\",
  switchTheme, false);\n\n                //pre-check the dark-theme checkbox if dark-theme
  is set\n                if (document.documentElement.getAttribute(\"data-theme\")
  == \"dark\") {\n                    toggleSwitch.checked = true;\n                }\n
  \           },\n            false,\n        );\n    </script>\n<style>\n      :root
  {\n        --color-bg: #1f2022;\n        --color-bg-2: ;\n        --color-bg-code:
  #1f2022;\n        --color-text: #eefbfe;\n        --color-link: #fb30c4; \n        --color-accent:
  #e1bd00c9;\n        --overlay-brightness: .85;\n        --body-width: 800px;\n      }\n
  \     [data-theme=\"dark\"] {\n        --color-bg: #1f2022;\n        --color-bg-2:
  ;\n        --color-bg-code: #1f2022;\n        --color-text: #eefbfe;\n        --color-link:
  #fb30c4; \n        --color-accent: #e1bd00c9;\n        --overlay-brightness: .85;\n
  \       --body-width: 800px;\n      }\n      [data-theme=\"light\"] {\n        --color-bg:
  #eefbfe;\n        --color-bg-2: ;\n        --color-bg-code: #eefbfe;\n        --color-text:
  #1f2022;\n        --color-link: #fb30c4; \n        --color-accent: #ffeb00;\n        --overlay-brightness:
  .95;\n      }\n\n        html {\n            font-family: \"Space Mono\", monospace;\n
  \           background: var(--color-bg);\n            color: var(--color-text);\n
  \       }\n\n        a {\n            color: var(--color-link);\n        }\n\n        main
  a {\n            max-width: 100%;\n        }\n\n        .heading-permalink {\n            font-size:
  .7em;\n        }\n\n        body {\n            max-width: var(--body-width);\n
  \           margin: 5rem auto;\n            padding: 0 .5rem;\n            font-size:
  1rem;\n            line-height: 1.56;\n        }\n\n        blockquote {\n            background:
  var(--color-bg);\n            filter: brightness(var(--overlay-brightness));\n            border-left:
  4px solid var(--color-accent);\n            border-radius: 4px;\n            box-shadow:\n
  \               -0.8rem 0rem 1rem -1rem #f1fa8c,\n                0.2rem 0rem 1rem
  rgb(0, 0, 0, .4);\n            padding-left: 1rem;\n            margin: 1rem;\n
  \       }\n\n        li.post {\n            list-style-type: None;\n            padding:
  .2rem 0;\n        }\n\n        pre.wrapper {\n            padding: 0;\n            box-shadow:
  0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n            display: flex;\n            flex-direction:
  column;\n            position: relative;\n            margin: 2rem;\n        }\n\n
  \       pre {\n            margin: 0;\n            padding: 1rem;\n            min-width:
  -webkit-fill-available;\n            max-width: fit-content;\n            overflow-x:
  auto;\n        }\n\n        pre .filepath {\n            margin: 0;\n            padding-left:
  1rem;\n            border-radius: 4px 4px 0 0;\n            background: black;\n
  \           display: flex;\n            justify-content: space-between;\n            align-items:
  center;\n        }\n\n        pre .filepath p {\n            margin: 0\n        }\n\n
  \       pre .filepath .right {\n            display: flex;\n            gap: .2rem;\n
  \           align-items: center;\n        }\n\n        pre::-webkit-scrollbar {\n
  \           height: 4px;\n            background-color: transparent;\n        }\n\n
  \       pre::-webkit-scrollbar-thumb {\n            background-color: #d3d3d32e;\n
  \           border-radius: 2px;\n        }\n\n        pre::-webkit-scrollbar-track
  {\n            background-color: transparent;\n        }\n\n        .copy-wrapper
  {\n            background: none;\n            position: absolute;\n            width:
  100%;\n            z-index: 100;\n            display: flex;\n            justify-content:
  flex-end;\n        }\n\n        button.copy {\n            z-index: 100;\n            background:
  none;\n            fill: #ffffff45;\n            border: none;\n            width:
  32px;\n            align-self: flex-end;\n            top: 0;\n            right:
  0;\n            margin: 0.5rem 0.2rem;\n\n        }\n\n        button.copy:hover
  {\n            fill: white\n        }\n\n        a.help {\n            fill: #ffffff45;\n
  \       }\n\n        a.help:hover {\n            fill: white;\n        }\n\n        a.help
  svg {\n            height: 24px;\n            width: 24px;\n        }\n\n        .highlight
  {\n            background: var(--color-bg-code);\n            color: var(--color-text);\n
  \           filter: brightness(var(--overlay-brightness));\n            border-radius:
  0 0 4px 4px;\n        }\n\n        .highlight .c {\n            color: #8b8b8b\n
  \       }\n\n        /* Comment */\n        .highlight .err {\n            color:
  #960050;\n            background-color: #1e0010\n        }\n\n        /* Error */\n
  \       .highlight .k {\n            color: #c678dd\n        }\n\n        /* Keyword
  */\n        .highlight .l {\n            color: #ae81ff\n        }\n\n        /*
  Literal */\n        .highlight .n {\n            color: #abb2bf\n        }\n\n        /*
  Name */\n        .highlight .o {\n            color: #c678dd\n        }\n\n        /*
  Operator */\n        .highlight .p {\n            color: #abb2bf\n        }\n\n
  \       /* Punctuation */\n        .highlight .ch {\n            color: #8b8b8b\n
  \       }\n\n        /* Comment.Hashbang */\n        .highlight .cm {\n            color:
  #8b8b8b\n        }\n\n        /* Comment.Multiline */\n        .highlight .cp {\n
  \           color: #8b8b8b\n        }\n\n        /* Comment.Preproc */\n        .highlight
  .cpf {\n            color: #8b8b8b\n        }\n\n        /* Comment.PreprocFile
  */\n        .highlight .c1 {\n            color: #8b8b8b\n        }\n\n        /*
  Comment.Single */\n        .highlight .cs {\n            color: #8b8b8b\n        }\n\n
  \       /* Comment.Special */\n        .highlight .gd {\n            color: #c678dd\n
  \       }\n\n        /* Generic.Deleted */\n        .highlight .ge {\n            font-style:
  italic\n        }\n\n        /* Generic.Emph */\n        .highlight .gi {\n            color:
  #a6e22e\n        }\n\n        /* Generic.Inserted */\n        .highlight .gs {\n
  \           font-weight: bold\n        }\n\n        /* Generic.Strong */\n        .highlight
  .gu {\n            color: #8b8b8b\n        }\n\n        /* Generic.Subheading */\n
  \       .highlight .kc {\n            color: #c678dd\n        }\n\n        /* Keyword.Constant
  */\n        .highlight .kd {\n            color: #c678dd\n        }\n\n        /*
  Keyword.Declaration */\n        .highlight .kn {\n            color: #c678dd\n        }\n\n
  \       /* Keyword.Namespace */\n        .highlight .kp {\n            color: #c678dd\n
  \       }\n\n        /* Keyword.Pseudo */\n        .highlight .kr {\n            color:
  #c678dd\n        }\n\n        /* Keyword.Reserved */\n        .highlight .kt {\n
  \           color: #c678dd\n        }\n\n        /* Keyword.Type */\n        .highlight
  .ld {\n            color: #e6db74\n        }\n\n        /* Literal.Date */\n        .highlight
  .m {\n            color: #ae81ff\n        }\n\n        /* Literal.Number */\n        .highlight
  .s {\n            color: #e6db74\n        }\n\n        /* Literal.String */\n        .highlight
  .na {\n            color: #a6e22e\n        }\n\n        /* Name.Attribute */\n        .highlight
  .nb {\n            color: #98c379\n        }\n\n        /* Name.Builtin */\n        .highlight
  .nc {\n            color: #abb2bf\n        }\n\n        /* Name.Class */\n        .highlight
  .no {\n            color: #c678dd\n        }\n\n        /* Name.Constant */\n        .highlight
  .nd {\n            color: #abb2bf\n        }\n\n        /* Name.Decorator */\n        .highlight
  .ni {\n            color: #abb2bf\n        }\n\n        /* Name.Entity */\n        .highlight
  .ne {\n            color: #a6e22e\n        }\n\n        /* Name.Exception */\n        .highlight
  .nf {\n            color: #61afef\n        }\n\n        /* Name.Function */\n        .highlight
  .nl {\n            color: #abb2bf\n        }\n\n        /* Name.Label */\n        .highlight
  .nn {\n            color: #abb2bf\n        }\n\n        /* Name.Namespace */\n        .highlight
  .nx {\n            color: #a6e22e\n        }\n\n        /* Name.Other */\n        .highlight
  .py {\n            color: #abb2bf\n        }\n\n        /* Name.Property */\n        .highlight
  .nt {\n            color: #c678dd\n        }\n\n        /* Name.Tag */\n        .highlight
  .nv {\n            color: #abb2bf\n        }\n\n        /* Name.Variable */\n        .highlight
  .ow {\n            color: #c678dd\n        }\n\n        /* Operator.Word */\n        .highlight
  .w {\n            color: #abb2bf\n        }\n\n        /* Text.Whitespace */\n        .highlight
  .mb {\n            color: #ae81ff\n        }\n\n        /* Literal.Number.Bin */\n
  \       .highlight .mf {\n            color: #ae81ff\n        }\n\n        /* Literal.Number.Float
  */\n        .highlight .mh {\n            color: #ae81ff\n        }\n\n        /*
  Literal.Number.Hex */\n        .highlight .mi {\n            color: #ae81ff\n        }\n\n
  \       /* Literal.Number.Integer */\n        .highlight .mo {\n            color:
  #ae81ff\n        }\n\n        /* Literal.Number.Oct */\n        .highlight .sa {\n
  \           color: #e6db74\n        }\n\n        /* Literal.String.Affix */\n        .highlight
  .sb {\n            color: #e6db74\n        }\n\n        /* Literal.String.Backtick
  */\n        .highlight .sc {\n            color: #e6db74\n        }\n\n        /*
  Literal.String.Char */\n        .highlight .dl {\n            color: #e6db74\n        }\n\n
  \       /* Literal.String.Delimiter */\n        .highlight .sd {\n            color:
  #98c379\n        }\n\n        /* Literal.String.Doc */\n        .highlight .s2 {\n
  \           color: #98c379\n        }\n\n        /* Literal.String.Double */\n        .highlight
  .se {\n            color: #ae81ff\n        }\n\n        /* Literal.String.Escape
  */\n        .highlight .sh {\n            color: #e6db74\n        }\n\n        /*
  Literal.String.Heredoc */\n        .highlight .si {\n            color: #e6db74\n
  \       }\n\n        /* Literal.String.Interpol */\n        .highlight .sx {\n            color:
  #e6db74\n        }\n\n        /* Literal.String.Other */\n        .highlight .sr
  {\n            color: #e6db74\n        }\n\n        /* Literal.String.Regex */\n
  \       .highlight .s1 {\n            color: #e6db74\n        }\n\n        /* Literal.String.Single
  */\n        .highlight .ss {\n            color: #e6db74\n        }\n\n        /*
  Literal.String.Symbol */\n        .highlight .bp {\n            color: #abb2bf\n
  \       }\n\n        /* Name.Builtin.Pseudo */\n        .highlight .fm {\n            color:
  #61afef\n        }\n\n        /* Name.Function.Magic */\n        .highlight .vc
  {\n            color: #abb2bf\n        }\n\n        /* Name.Variable.Class */\n
  \       .highlight .vg {\n            color: #abb2bf\n        }\n\n        /* Name.Variable.Global
  */\n        .highlight .vi {\n            color: #abb2bf\n        }\n\n        /*
  Name.Variable.Instance */\n        .highlight .vm {\n            color: #abb2bf\n
  \       }\n\n        /* Name.Variable.Magic */\n        .highlight .il {\n            color:
  #ae81ff\n        }\n\n        /* Literal.Number.Integer.Long */\n\n        /* Tab
  style starts here */\n        .tabbed-set {\n            position: relative;\n            display:
  flex;\n            flex-wrap: wrap;\n            margin: 1em 0;\n            border-radius:
  0.1rem;\n        }\n\n        .tabbed-set>input {\n            display: none;\n
  \       }\n\n        .tabbed-set label {\n            width: auto;\n            padding:
  0.9375em 1.25em 0.78125em;\n            font-weight: 700;\n            font-size:
  0.84em;\n            white-space: nowrap;\n            border-bottom: 0.15rem solid
  transparent;\n            border-top-left-radius: 0.1rem;\n            border-top-right-radius:
  0.1rem;\n            cursor: pointer;\n            transition: background-color
  250ms, color 250ms;\n        }\n\n        .tabbed-set .tabbed-content {\n            width:
  100%;\n            display: none;\n            box-shadow: 0 -.05rem #ddd;\n        }\n\n
  \       .tabbed-set input {\n            position: absolute;\n            opacity:
  0;\n        }\n\n        /* fonts */\n        h1 {\n            font-weight: 700;\n
  \       }\n\n        h1#title a {\n            font-size: 16px;\n        }\n\n        h1,\n
  \       h2,\n        h3,\n        h4,\n        h5,\n        h6 {\n            margin-top:
  3rem;\n        }\n\n        h1 {\n            font-size: 2.5em;\n            margin-top:
  5rem;\n        }\n\n        h2 {\n            font-size: 1.63rem;\n            margin-top:
  5rem;\n        }\n\n\n\n        p {\n            font-size: 21px;\n            font-style:
  normal;\n            font-variant: normal;\n            font-weight: 400;\n            line-height:
  1.5;\n        }\n\n        @media only screen and (max-width: 700px) {\n            p
  {\n                font-size: 18px;\n            }\n        }\n\n        @media
  only screen and (max-width: 600px) {\n            p {\n                font-size:
  16px;\n            }\n        }\n\n        @media only screen and (max-width: 500px)
  {\n            p {\n                font-size: 14px;\n            }\n        }\n\n
  \       @media only screen and (max-width: 400px) {\n            p {\n                font-size:
  12px;\n            }\n        }\n\n\n        pre {\n            font-style: normal;\n
  \           font-variant: normal;\n            font-weight: 400;\n            line-height:
  18.5714px;\n            */\n        }\n\n        a {\n            font-weight: 600;\n
  \           text-decoration-color: var(--color-accent);\n            color: var(--color-link);\n
  \           padding: .3rem .5rem;\n            display: inline-block;\n        }\n\n
  \       .admonition,\n        details {\n            box-shadow: 0.2rem 0rem 1rem
  rgb(0, 0, 0, .4);\n            margin: 5rem 0;\n            border: 1px solid transparent;\n
  \           border-radius: 4px;\n            text-align: left;\n            padding:
  0;\n            border: 0;\n\n        }\n\n        .admonition {\n            padding-bottom:
  1rem;\n        }\n\n        details[open] {\n            padding-bottom: .5rem;\n
  \       }\n\n        .admonition p {\n            padding: .2rem .6rem;\n        }\n\n
  \       .admonition-title,\n        .details-title,\n        summary {\n            background:
  var(--color-bg-2);\n            padding: 0;\n            margin: 0;\n            position:
  sticky;\n            top: 0;\n            z-index: 10;\n        }\n\n        summary:hover
  {\n            cursor: pointer;\n        }\n\n        summary.admonition-title,\n
  \       summary.details-title {\n            padding: .5rem;\n            padding-left:
  1rem;\n        }\n\n        .note {\n            border-left: 4px solid #f1fa8c;\n
  \           box-shadow:\n                -0.8rem 0rem 1rem -1rem #f1fa8c,\n                0.2rem
  0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n        .note>.admonition-title {\n            border-bottom:
  1px solid #3c3d2d;\n        }\n\n        .abstract {\n            border-left: 4px
  solid #8be9fd;\n            box-shadow:\n                -0.8rem 0rem 1rem -1rem
  #8be9fd,\n                0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n        .abstract>.admonition-title
  {\n            border-bottom: 1px solid #2c3a3f;\n        }\n\n        .info {\n
  \           border-left: 4px solid;\n            box-shadow:\n                -0.8rem
  0rem 1rem -1rem #8bb0fd,\n                0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n
  \       .info>.admonition-title {\n            border-bottom: 1px solid #2c313f;\n
  \       }\n\n        .tip {\n            border-left: 4px solid #008080;\n            box-shadow:\n
  \               -0.8rem 0rem 1rem -1rem #008080,\n                0.2rem 0rem 1rem
  rgb(0, 0, 0, .4);\n        }\n\n        .tip>.admonition-title {\n            border-bottom:
  1px solid #1b2a2b;\n        }\n\n        .success {\n            border-left: 4px
  solid #50fa7b;\n            box-shadow:\n                -0.8rem 0rem 1rem -1rem
  #50fa7b,\n                0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n        .success>.admonition-title
  {\n            border-bottom: 1px solid #263e2b;\n        }\n\n        .question
  {\n            border-left: 4px solid #a7fcbd;\n            box-shadow:\n                -0.8rem
  0rem 1rem -1rem #a7fcbd,\n                0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n
  \       .question>.admonition-title {\n            border-bottom: 1px solid #303e35;\n
  \       }\n\n        .warning {\n            border-left: 4px solid #ffb86c;\n            box-shadow:\n
  \               -0.8rem 0rem 1rem -1rem #ffb86c,\n                0.2rem 0rem 1rem
  rgb(0, 0, 0, .4);\n        }\n\n        .warning>.admonition-title {\n            border-bottom:
  1px solid #3f3328;\n        }\n\n        .failure {\n            border-left: 4px
  solid #b23b3b;\n            box-shadow:\n                -0.8rem 0rem 1rem -1rem
  #b23b3b,\n                0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n        .failure>.admonition-title
  {\n            border-bottom: 1px solid #34201f;\n        }\n\n        .danger {\n
  \           border-left: 4px solid #ff5555;\n            box-shadow:\n                -0.8rem
  0rem 1rem -1rem #ff5555,\n                0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n
  \       .danger>.admonition-title {\n            border-bottom: 1px solid #402523;\n
  \       }\n\n        .bug {\n            border-left: 4px solid #b2548a;\n            box-shadow:\n
  \               -0.8rem 0rem 1rem -1rem #b2548a,\n                0.2rem 0rem 1rem
  rgb(0, 0, 0, .4);\n        }\n\n        .bug>.admonition-title {\n            border-bottom:
  1px solid #32232c;\n        }\n\n        .example {\n            border-left: 4px
  solid #bd93f9;\n            box-shadow:\n                -0.8rem 0rem 1rem -1rem
  #bd93f9,\n                0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n        .example>.admonition-title
  {\n            border-bottom: 1px solid #332d3e;\n        }\n\n        .source {\n
  \           border-left: 4px solid #bd93f9;\n            box-shadow:\n                -0.8rem
  0rem 1rem -1rem #bd93f9,\n                0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n
  \       .source>.admonition-title {\n            border-bottom: 1px solid #332d3e;\n
  \       }\n\n        .quote {\n            border-left: 4px solid #999;\n            box-shadow:\n
  \               -0.8rem 0rem 1rem -1rem #999,\n                0.2rem 0rem 1rem
  rgb(0, 0, 0, .4);\n        }\n\n        .quote>.admonition-title {\n            border-bottom:
  1px solid #2d2e2f;\n        }\n\n        table {\n            margin: 1rem 0;\n
  \           border-collapse: collapse;\n            border-spacing: 0;\n            display:
  block;\n            max-width: -moz-fit-content;\n            max-width: fit-content;\n
  \           overflow-x: auto;\n            white-space: nowrap;\n        }\n\n        table
  thead th {\n            border: solid 1px var(--color-text);\n            padding:
  10px;\n            text-align: left;\n        }\n\n        table tbody td {\n            border:
  solid 1px var(--color-text);\n            padding: 10px;\n        }\n\n        .theme-switch
  {\n            z-index: 10;\n            display: inline-block;\n            height:
  34px;\n            position: relative;\n            width: 60px;\n\n            display:
  flex;\n            justify-content: flex-end;\n            margin-right: 1rem;\n
  \           margin-left: auto;\n            position: fixed;\n            right:
  1rem;\n            top: 1rem;\n        }\n\n        .theme-switch input {\n            display:
  none;\n\n        }\n\n        .slider {\n            background-color: #ccc;\n            bottom:
  0;\n            cursor: pointer;\n            left: 0;\n            position: absolute;\n
  \           right: 0;\n            top: 0;\n            transition: .4s;\n        }\n\n
  \       .slider:before {\n            background-color: #fff;\n            bottom:
  4px;\n            content: \"\";\n            height: 26px;\n            left: 4px;\n
  \           position: absolute;\n            transition: .4s;\n            width:
  26px;\n        }\n\n        input:checked+.slider {\n            background-color:
  #343434;\n        }\n\n        input:checked+.slider:before {\n            background-color:
  #848484;\n        }\n\n        input:checked+.slider:before {\n            transform:
  translateX(26px);\n        }\n\n        .slider.round {\n            border-radius:
  34px;\n        }\n\n        .slider.round:before {\n            border-radius: 50%;\n
  \       }\n\n        main p img {\n            width: 100%;\n            width:
  -moz-available;\n            width: -webkit-fill-available;\n            width:
  fill-available;\n        }\n\n        details>* {\n            margin: 1rem;\n        }\n\n
  \       .admonition>* {\n            margin: 1rem;\n        }\n\n        p.admonition-title,\n
  \       summary {\n            margin: 0;\n            padding-left: 1.2rem;\n        }\n\n
  \       .small {\n            font-size: .9rem;\n            color: #888;\n        }\n\n
  \       admonition+admonition {\n            margin-top: 20rem;\n        }\n\n        ::-webkit-scrollbar
  {\n            height: 12px;\n            background-color: transparent;\n        }\n\n
  \       ::-webkit-scrollbar-thumb {\n            background-color: #d3d3d32e;\n
  \           border-radius: 6px;\n        }\n\n        ::-webkit-scrollbar-track
  {\n            background-color: transparent;\n        }\n    </style>\n<script>\n
  \       if (\"serviceWorker\" in navigator) {\n            navigator.serviceWorker.register(\"/service-worker.js\");\n
  \           navigator.serviceWorker.addEventListener(\"controllerchange\", () =>
  {\n                console.log(\"new worker\");\n                window.location.reload();\n
  \           });\n        }\n    </script>\n<meta content=\"waylon@waylonwalker.com\"
  name=\"og:author_email\"/>\n<meta content=\"waylon@waylonwalker.com\" name=\"og:author_email\"/>\n<meta
  content=\"Waylon Walker\" name=\"og:author\" property=\"og:author\"/><meta content=\"waylon@waylonwalaker.com\"
  name=\"og:author_email\" property=\"og:author_email\"/><meta content=\"website\"
  name=\"og:type\" property=\"og:type\"/><meta content=\"None ! ! ???+ source  ! ???+
  source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source
  \ ! ???+ source  ! ???+ source  ! ???+ source  ! ?\" name=\"description\" property=\"description\"/><meta
  content=\"None ! ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
  source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source
  \ ! ?\" name=\"og:description\" property=\"og:description\"/><meta content=\"None
  ! ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+
  source  ! ???+ source  ! ???+ source  ! ???+ source  ! ???+ source  ! ?\" name=\"twitter:description\"
  property=\"twitter:description\"/><meta content=\"Server.Py | Markata\" name=\"og:title\"
  property=\"og:title\"/><meta content=\"Server.Py | Markata\" name=\"twitter:title\"
  property=\"twitter:title\"/><meta content=\"https://markata.dev//markata/cli/server-og.png\"
  name=\"og:image\" property=\"og:image\"/><meta content=\"https://markata.dev//markata/cli/server-og.png\"
  name=\"twitter:image\" property=\"twitter:image\"/><meta content=\"1600\" name=\"og:image:width\"
  property=\"og:image:width\"/><meta content=\"900\" name=\"og:image:width\" property=\"og:image:width\"/><meta
  content=\"summary_large_image\" name=\"twitter:card\" property=\"twitter:card\"/><meta
  content=\"Markata\" name=\"og:site_name\" property=\"og:site_name\"/><meta content=\"@_waylonwalker\"
  name=\"twitter:creator\" property=\"twitter:creator\"/><meta content=\"Server.Py\"
  name=\"title\" property=\"title\"/><meta content=\"markata 0.8.0.dev6\" name=\"generator\"
  property=\"generator\"/><link href=\"https://markata.dev//markata/cli/server/\"
  rel=\"canonical\"/><meta content=\"https://markata.dev//markata/cli/server/\" name=\"og:url\"
  property=\"og:url\"/></head>\n<body><nav>\n<a href=\"/\">markata</a>\n<a href=\"https://github.com/WaylonWalker/markata\">GitHub</a>\n</nav>\n<div>\n<label
  class=\"theme-switch\" for=\"checkbox-theme\" id=\"theme-switch\" title=\"light/dark
  mode toggle\">\n<input id=\"checkbox-theme\" type=\"checkbox\"/>\n<div class=\"slider
  round\"></div>\n</label>\n</div>\n<section class=\"title\">\n<h1 id=\"title\">\n
  \           Server.Py \n            \n        </h1>\n</section>\n<main><p>None</p>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"find_port\" style=\"margin:0;padding:.5rem
  1rem;\">find_port <em class=\"small\">function</em></h2>\nFind a port not in ues
  starting at given port\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
  class=\"admonition-title\">\"find_port <em class=\"small\">source</em>\"</p>\n</div>\n<pre
  class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">find_port</span><span class=\"p\">(</span><span class=\"n\">port</span><span
  class=\"p\">:</span> <span class=\"nb\">int</span> <span class=\"o\">=</span> <span
  class=\"mi\">8000</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
  <span class=\"nb\">int</span><span class=\"p\">:</span>\n<span class=\"w\">            </span><span
  class=\"sd\">\"\"\"Find a port not in ues starting at given port\"\"\"</span>\n
  \           <span class=\"kn\">import</span> <span class=\"nn\">socket</span>\n\n
  \           <span class=\"k\">with</span> <span class=\"n\">socket</span><span class=\"o\">.</span><span
  class=\"n\">socket</span><span class=\"p\">(</span><span class=\"n\">socket</span><span
  class=\"o\">.</span><span class=\"n\">AF_INET</span><span class=\"p\">,</span> <span
  class=\"n\">socket</span><span class=\"o\">.</span><span class=\"n\">SOCK_STREAM</span><span
  class=\"p\">)</span> <span class=\"k\">as</span> <span class=\"n\">s</span><span
  class=\"p\">:</span>\n                <span class=\"k\">if</span> <span class=\"n\">s</span><span
  class=\"o\">.</span><span class=\"n\">connect_ex</span><span class=\"p\">((</span><span
  class=\"s2\">\"localhost\"</span><span class=\"p\">,</span> <span class=\"n\">port</span><span
  class=\"p\">))</span> <span class=\"o\">==</span> <span class=\"mi\">0</span><span
  class=\"p\">:</span>\n                    <span class=\"k\">return</span> <span
  class=\"n\">find_port</span><span class=\"p\">(</span><span class=\"n\">port</span><span
  class=\"o\">=</span><span class=\"n\">port</span> <span class=\"o\">+</span> <span
  class=\"mi\">1</span><span class=\"p\">)</span>\n                <span class=\"k\">return</span>
  <span class=\"n\">port</span>\n</pre></div>\n\n</pre>\n<p>!! class </p><h2 class=\"admonition-title\"
  id=\"Server\" style=\"margin:0;padding:.5rem 1rem;\">Server <em class=\"small\">class</em></h2>\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"Server
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
  <span class=\"nc\">Server</span><span class=\"p\">:</span>\n            <span class=\"k\">def</span>
  <span class=\"fm\">__init__</span><span class=\"p\">(</span>\n                <span
  class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">\"Server\"</span><span
  class=\"p\">,</span>\n                <span class=\"o\">*</span><span class=\"p\">,</span>\n
  \               <span class=\"n\">auto_restart</span><span class=\"p\">:</span>
  <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span class=\"kc\">True</span><span
  class=\"p\">,</span>\n                <span class=\"n\">directory</span><span class=\"p\">:</span>
  <span class=\"n\">Union</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
  class=\"p\">,</span> <span class=\"s2\">\"Path\"</span><span class=\"p\">]</span>
  <span class=\"o\">=</span> <span class=\"kc\">None</span><span class=\"p\">,</span>\n
  \               <span class=\"n\">port</span><span class=\"p\">:</span> <span class=\"nb\">int</span>
  <span class=\"o\">=</span> <span class=\"mi\">8000</span><span class=\"p\">,</span>\n
  \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
  class=\"p\">:</span>\n                <span class=\"k\">if</span> <span class=\"n\">directory</span>
  <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
  \                   <span class=\"kn\">from</span> <span class=\"nn\">markata</span>
  <span class=\"kn\">import</span> <span class=\"n\">Markata</span>\n\n                    <span
  class=\"n\">m</span> <span class=\"o\">=</span> <span class=\"n\">Markata</span><span
  class=\"p\">()</span>\n                    <span class=\"n\">directory</span> <span
  class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
  class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">m</span><span
  class=\"o\">.</span><span class=\"n\">config</span><span class=\"p\">[</span><span
  class=\"s2\">\"output_dir\"</span><span class=\"p\">]))</span>\n\n                <span
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
  \               <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">cmd</span>
  <span class=\"o\">=</span> <span class=\"p\">[</span>\n                    <span
  class=\"s2\">\"python\"</span><span class=\"p\">,</span>\n                    <span
  class=\"s2\">\"-m\"</span><span class=\"p\">,</span>\n                    <span
  class=\"s2\">\"http.server\"</span><span class=\"p\">,</span>\n                    <span
  class=\"nb\">str</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">port</span><span class=\"p\">),</span>\n                    <span
  class=\"s2\">\"--directory\"</span><span class=\"p\">,</span>\n                    <span
  class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">directory</span><span
  class=\"p\">,</span>\n                <span class=\"p\">]</span>\n\n                <span
  class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">proc</span>
  <span class=\"o\">=</span> <span class=\"n\">subprocess</span><span class=\"o\">.</span><span
  class=\"n\">Popen</span><span class=\"p\">(</span>\n                    <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">cmd</span><span class=\"p\">,</span>\n                    <span
  class=\"n\">stderr</span><span class=\"o\">=</span><span class=\"n\">subprocess</span><span
  class=\"o\">.</span><span class=\"n\">PIPE</span><span class=\"p\">,</span>\n                    <span
  class=\"n\">stdout</span><span class=\"o\">=</span><span class=\"n\">subprocess</span><span
  class=\"o\">.</span><span class=\"n\">PIPE</span><span class=\"p\">,</span>\n                <span
  class=\"p\">)</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">start_time</span> <span class=\"o\">=</span> <span class=\"n\">time</span><span
  class=\"o\">.</span><span class=\"n\">time</span><span class=\"p\">()</span>\n\n
  \           <span class=\"nd\">@property</span>\n            <span class=\"k\">def</span>
  <span class=\"nf\">uptime</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">int</span><span
  class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"nb\">round</span><span
  class=\"p\">(</span><span class=\"n\">time</span><span class=\"o\">.</span><span
  class=\"n\">time</span><span class=\"p\">()</span> <span class=\"o\">-</span> <span
  class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">start_time</span><span
  class=\"p\">)</span>\n\n            <span class=\"nd\">@property</span>\n            <span
  class=\"k\">def</span> <span class=\"nf\">title</span><span class=\"p\">(</span><span
  class=\"bp\">self</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
  <span class=\"nb\">str</span><span class=\"p\">:</span>\n                <span class=\"k\">return</span>
  <span class=\"sa\">f</span><span class=\"s2\">\"server (</span><span class=\"si\">{</span><span
  class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">uptime</span><span
  class=\"si\">}</span><span class=\"s2\">)\"</span>\n\n            <span class=\"k\">def</span>
  <span class=\"nf\">kill</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
  class=\"p\">:</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">auto_restart</span> <span class=\"o\">=</span> <span class=\"kc\">False</span>\n
  \               <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">proc</span><span
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
  class=\"sa\">f</span><span class=\"s2\">\"[green]serving on port: [gold1]</span><span
  class=\"si\">{</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">port</span><span class=\"si\">}</span><span class=\"s2\"> \"</span>\n
  \                           <span class=\"sa\">f</span><span class=\"s2\">\"[green]using
  pid: [gold1]</span><span class=\"si\">{</span><span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">proc</span><span class=\"o\">.</span><span
  class=\"n\">pid</span><span class=\"si\">}</span><span class=\"s2\"> \"</span>\n
  \                           <span class=\"sa\">f</span><span class=\"s2\">\"[green]uptime:
  [gold1]</span><span class=\"si\">{</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">uptime</span><span class=\"si\">}</span><span class=\"s2\"> \"</span>\n
  \                           <span class=\"sa\">f</span><span class=\"s2\">\"[green]link:
  [gold1] http://localhost:</span><span class=\"si\">{</span><span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">port</span><span class=\"si\">}</span><span
  class=\"s2\">[/]\"</span>\n                        <span class=\"p\">),</span>\n
  \                       <span class=\"n\">border_style</span><span class=\"o\">=</span><span
  class=\"s2\">\"blue\"</span><span class=\"p\">,</span>\n                        <span
  class=\"n\">title</span><span class=\"o\">=</span><span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">title</span><span class=\"p\">,</span>\n                        <span
  class=\"n\">expand</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
  class=\"p\">,</span>\n                    <span class=\"p\">)</span>\n\n                <span
  class=\"k\">else</span><span class=\"p\">:</span>\n                    <span class=\"k\">return</span>
  <span class=\"n\">Panel</span><span class=\"p\">(</span>\n                        <span
  class=\"s2\">\"[red]server died\"</span><span class=\"p\">,</span>\n                        <span
  class=\"n\">title</span><span class=\"o\">=</span><span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">title</span><span class=\"p\">,</span>\n                        <span
  class=\"n\">border_style</span><span class=\"o\">=</span><span class=\"s2\">\"red\"</span><span
  class=\"p\">,</span>\n                        <span class=\"n\">expand</span><span
  class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n                    <span
  class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!! function </p><h2 class=\"admonition-title\"
  id=\"configure\" style=\"margin:0;padding:.5rem 1rem;\">configure <em class=\"small\">function</em></h2>\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"configure
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">configure</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">)</span>
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
  class=\"p\">[</span><span class=\"s2\">\"output_dir\"</span><span class=\"p\">]))</span>\n
  \                   <span class=\"k\">return</span> <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">_server</span>\n\n            <span class=\"kn\">from</span>
  <span class=\"nn\">markata</span> <span class=\"kn\">import</span> <span class=\"n\">Markata</span>\n\n
  \           <span class=\"n\">Markata</span><span class=\"o\">.</span><span class=\"n\">server</span>
  <span class=\"o\">=</span> <span class=\"nb\">property</span><span class=\"p\">(</span><span
  class=\"n\">get_server</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"run_server\" style=\"margin:0;padding:.5rem
  1rem;\">run_server <em class=\"small\">function</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"run_server
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
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
  class=\"p\">()</span>\n</pre></div>\n\n</pre>\n<p>!! function </p><h2 class=\"admonition-title\"
  id=\"cli\" style=\"margin:0;padding:.5rem 1rem;\">cli <em class=\"small\">function</em></h2>\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"cli
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">cli</span><span class=\"p\">(</span><span class=\"n\">app</span><span
  class=\"p\">:</span> <span class=\"n\">typer</span><span class=\"o\">.</span><span
  class=\"n\">Typer</span><span class=\"p\">,</span> <span class=\"n\">markata</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
  \           <span class=\"n\">server_app</span> <span class=\"o\">=</span> <span
  class=\"n\">typer</span><span class=\"o\">.</span><span class=\"n\">Typer</span><span
  class=\"p\">()</span>\n            <span class=\"n\">app</span><span class=\"o\">.</span><span
  class=\"n\">add_typer</span><span class=\"p\">(</span><span class=\"n\">server_app</span><span
  class=\"p\">)</span>\n\n            <span class=\"nd\">@server_app</span><span class=\"o\">.</span><span
  class=\"n\">callback</span><span class=\"p\">(</span><span class=\"n\">invoke_without_command</span><span
  class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n            <span
  class=\"k\">def</span> <span class=\"nf\">serve</span><span class=\"p\">():</span>\n<span
  class=\"w\">                </span><span class=\"sd\">\"\"\"</span>\n<span class=\"sd\">
  \               Serve the site locally.</span>\n<span class=\"sd\">                \"\"\"</span>\n
  \               <span class=\"n\">run_server</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"__init__\" style=\"margin:0;padding:.5rem
  1rem;\"><strong>init</strong> <em class=\"small\">method</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"<strong>init</strong>
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"fm\">__init__</span><span class=\"p\">(</span>\n                <span
  class=\"bp\">self</span><span class=\"p\">:</span> <span class=\"s2\">\"Server\"</span><span
  class=\"p\">,</span>\n                <span class=\"o\">*</span><span class=\"p\">,</span>\n
  \               <span class=\"n\">auto_restart</span><span class=\"p\">:</span>
  <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span class=\"kc\">True</span><span
  class=\"p\">,</span>\n                <span class=\"n\">directory</span><span class=\"p\">:</span>
  <span class=\"n\">Union</span><span class=\"p\">[</span><span class=\"nb\">str</span><span
  class=\"p\">,</span> <span class=\"s2\">\"Path\"</span><span class=\"p\">]</span>
  <span class=\"o\">=</span> <span class=\"kc\">None</span><span class=\"p\">,</span>\n
  \               <span class=\"n\">port</span><span class=\"p\">:</span> <span class=\"nb\">int</span>
  <span class=\"o\">=</span> <span class=\"mi\">8000</span><span class=\"p\">,</span>\n
  \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
  class=\"p\">:</span>\n                <span class=\"k\">if</span> <span class=\"n\">directory</span>
  <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
  \                   <span class=\"kn\">from</span> <span class=\"nn\">markata</span>
  <span class=\"kn\">import</span> <span class=\"n\">Markata</span>\n\n                    <span
  class=\"n\">m</span> <span class=\"o\">=</span> <span class=\"n\">Markata</span><span
  class=\"p\">()</span>\n                    <span class=\"n\">directory</span> <span
  class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
  class=\"nb\">str</span><span class=\"p\">(</span><span class=\"n\">m</span><span
  class=\"o\">.</span><span class=\"n\">config</span><span class=\"p\">[</span><span
  class=\"s2\">\"output_dir\"</span><span class=\"p\">]))</span>\n\n                <span
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
  class=\"n\">kill</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"start_server\" style=\"margin:0;padding:.5rem
  1rem;\">start_server <em class=\"small\">method</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"start_server
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">start_server</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
  class=\"p\">:</span>\n                <span class=\"kn\">import</span> <span class=\"nn\">subprocess</span>\n\n
  \               <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">cmd</span>
  <span class=\"o\">=</span> <span class=\"p\">[</span>\n                    <span
  class=\"s2\">\"python\"</span><span class=\"p\">,</span>\n                    <span
  class=\"s2\">\"-m\"</span><span class=\"p\">,</span>\n                    <span
  class=\"s2\">\"http.server\"</span><span class=\"p\">,</span>\n                    <span
  class=\"nb\">str</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">port</span><span class=\"p\">),</span>\n                    <span
  class=\"s2\">\"--directory\"</span><span class=\"p\">,</span>\n                    <span
  class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">directory</span><span
  class=\"p\">,</span>\n                <span class=\"p\">]</span>\n\n                <span
  class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">proc</span>
  <span class=\"o\">=</span> <span class=\"n\">subprocess</span><span class=\"o\">.</span><span
  class=\"n\">Popen</span><span class=\"p\">(</span>\n                    <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">cmd</span><span class=\"p\">,</span>\n                    <span
  class=\"n\">stderr</span><span class=\"o\">=</span><span class=\"n\">subprocess</span><span
  class=\"o\">.</span><span class=\"n\">PIPE</span><span class=\"p\">,</span>\n                    <span
  class=\"n\">stdout</span><span class=\"o\">=</span><span class=\"n\">subprocess</span><span
  class=\"o\">.</span><span class=\"n\">PIPE</span><span class=\"p\">,</span>\n                <span
  class=\"p\">)</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">start_time</span> <span class=\"o\">=</span> <span class=\"n\">time</span><span
  class=\"o\">.</span><span class=\"n\">time</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"uptime\" style=\"margin:0;padding:.5rem
  1rem;\">uptime <em class=\"small\">method</em></h2>\n<div class=\"admonition source
  is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"uptime <em class=\"small\">source</em>\"</p>\n</div>\n<pre
  class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">uptime</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">int</span><span
  class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"nb\">round</span><span
  class=\"p\">(</span><span class=\"n\">time</span><span class=\"o\">.</span><span
  class=\"n\">time</span><span class=\"p\">()</span> <span class=\"o\">-</span> <span
  class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">start_time</span><span
  class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!! method </p><h2 class=\"admonition-title\"
  id=\"title\" style=\"margin:0;padding:.5rem 1rem;\">title <em class=\"small\">method</em></h2>\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"title
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">title</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"nb\">str</span><span
  class=\"p\">:</span>\n                <span class=\"k\">return</span> <span class=\"sa\">f</span><span
  class=\"s2\">\"server (</span><span class=\"si\">{</span><span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">uptime</span><span class=\"si\">}</span><span
  class=\"s2\">)\"</span>\n</pre></div>\n\n</pre>\n<p>!! method </p><h2 class=\"admonition-title\"
  id=\"kill\" style=\"margin:0;padding:.5rem 1rem;\">kill <em class=\"small\">method</em></h2>\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"kill
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">kill</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
  class=\"p\">:</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">auto_restart</span> <span class=\"o\">=</span> <span class=\"kc\">False</span>\n
  \               <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">proc</span><span
  class=\"o\">.</span><span class=\"n\">stdout</span><span class=\"o\">.</span><span
  class=\"n\">close</span><span class=\"p\">()</span>\n                <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">proc</span><span class=\"o\">.</span><span
  class=\"n\">stderr</span><span class=\"o\">.</span><span class=\"n\">close</span><span
  class=\"p\">()</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">proc</span><span class=\"o\">.</span><span class=\"n\">kill</span><span
  class=\"p\">()</span>\n                <span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">proc</span><span class=\"o\">.</span><span class=\"n\">wait</span><span
  class=\"p\">()</span>\n</pre></div>\n\n</pre>\n<p>!! method </p><h2 class=\"admonition-title\"
  id=\"__rich__\" style=\"margin:0;padding:.5rem 1rem;\"><strong>rich</strong> <em
  class=\"small\">method</em></h2>\n<div class=\"admonition source is-collapsible
  collapsible-open\">\n<p class=\"admonition-title\">\"<strong>rich</strong> <em class=\"small\">source</em>\"</p>\n</div>\n<pre
  class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">__rich__</span><span class=\"p\">(</span><span class=\"bp\">self</span><span
  class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"n\">Panel</span><span
  class=\"p\">:</span>\n                <span class=\"k\">if</span> <span class=\"ow\">not</span>
  <span class=\"bp\">self</span><span class=\"o\">.</span><span class=\"n\">proc</span><span
  class=\"o\">.</span><span class=\"n\">poll</span><span class=\"p\">():</span>\n
  \                   <span class=\"k\">return</span> <span class=\"n\">Panel</span><span
  class=\"p\">(</span>\n                        <span class=\"p\">(</span>\n                            <span
  class=\"sa\">f</span><span class=\"s2\">\"[green]serving on port: [gold1]</span><span
  class=\"si\">{</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">port</span><span class=\"si\">}</span><span class=\"s2\"> \"</span>\n
  \                           <span class=\"sa\">f</span><span class=\"s2\">\"[green]using
  pid: [gold1]</span><span class=\"si\">{</span><span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">proc</span><span class=\"o\">.</span><span
  class=\"n\">pid</span><span class=\"si\">}</span><span class=\"s2\"> \"</span>\n
  \                           <span class=\"sa\">f</span><span class=\"s2\">\"[green]uptime:
  [gold1]</span><span class=\"si\">{</span><span class=\"bp\">self</span><span class=\"o\">.</span><span
  class=\"n\">uptime</span><span class=\"si\">}</span><span class=\"s2\"> \"</span>\n
  \                           <span class=\"sa\">f</span><span class=\"s2\">\"[green]link:
  [gold1] http://localhost:</span><span class=\"si\">{</span><span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">port</span><span class=\"si\">}</span><span
  class=\"s2\">[/]\"</span>\n                        <span class=\"p\">),</span>\n
  \                       <span class=\"n\">border_style</span><span class=\"o\">=</span><span
  class=\"s2\">\"blue\"</span><span class=\"p\">,</span>\n                        <span
  class=\"n\">title</span><span class=\"o\">=</span><span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">title</span><span class=\"p\">,</span>\n                        <span
  class=\"n\">expand</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
  class=\"p\">,</span>\n                    <span class=\"p\">)</span>\n\n                <span
  class=\"k\">else</span><span class=\"p\">:</span>\n                    <span class=\"k\">return</span>
  <span class=\"n\">Panel</span><span class=\"p\">(</span>\n                        <span
  class=\"s2\">\"[red]server died\"</span><span class=\"p\">,</span>\n                        <span
  class=\"n\">title</span><span class=\"o\">=</span><span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">title</span><span class=\"p\">,</span>\n                        <span
  class=\"n\">border_style</span><span class=\"o\">=</span><span class=\"s2\">\"red\"</span><span
  class=\"p\">,</span>\n                        <span class=\"n\">expand</span><span
  class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n                    <span
  class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!! function </p><h2 class=\"admonition-title\"
  id=\"get_server\" style=\"margin:0;padding:.5rem 1rem;\">get_server <em class=\"small\">function</em></h2>\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"get_server
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
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
  class=\"p\">[</span><span class=\"s2\">\"output_dir\"</span><span class=\"p\">]))</span>\n
  \                   <span class=\"k\">return</span> <span class=\"bp\">self</span><span
  class=\"o\">.</span><span class=\"n\">_server</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"serve\" style=\"margin:0;padding:.5rem
  1rem;\">serve <em class=\"small\">function</em></h2>\nServe the site locally.\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"serve
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
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
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">serve</span><span class=\"p\">():</span>\n<span class=\"w\">
  \               </span><span class=\"sd\">\"\"\"</span>\n<span class=\"sd\">                Serve
  the site locally.</span>\n<span class=\"sd\">                \"\"\"</span>\n                <span
  class=\"n\">run_server</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n</main>\n<footer>\xA9
  2024</footer>\n</body></html>"
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
