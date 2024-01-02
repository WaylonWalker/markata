---
content: "Setup Logging hook sets up the RichHandler for pretty console logs, and
  file\nlogs to the configured markata's configured `log_dir`, or `output_dir/_logs`
  if\n`log_dir` is not configured.  The log file will be named after the\n`<levelname>.log`\n\n#
  The log files\n\nThere will be 6 log files created based on log level and file type.\n\n```\nmarkout/_logs\n\u251C\u2500\u2500
  debug\n\u2502   \u2514\u2500\u2500 index.html\n\u251C\u2500\u2500 debug.log\n\u251C\u2500\u2500
  info\n\u2502   \u2514\u2500\u2500 index.html\n\u251C\u2500\u2500 info.log\n\u251C\u2500\u2500
  warning\n\u2502   \u2514\u2500\u2500 index.html\n\u2514\u2500\u2500 warning.log\n```\n\n#
  Configuration\n\nEnsure that setup_logging is in your hooks.  You can check if `setup_logging`\nis
  in your hooks by running `markata list --hooks` from your terminal and\nchecking
  the output, or creating an instance of `Markata()` and checking the\n`Markata().hooks`
  attribute.  If its missing or you wan to be more explicit,\nyou can add `setup_logging`
  to your `markata.toml` `[markata.hooks]`.\n\n``` toml\n[markata]\n\n# make sure
  its in your list of hooks\nhooks=[\n   \"markata.plugins.setup_logging\",\n   ]\n```\n\n#
  Log Template\n``` toml\n[markata]\n\n# make sure its in your list of hooks\nhooks=[\n
  \  \"markata.plugins.setup_logging\",\n   ]\n\n# point log template to the path
  of your logging template\nlog_template='templates/log_template.html'\n```\n\nYou
  can see the latest default `log_template` on\n[GitHub](https://github.com/WaylonWalker/markata/blob/main/markata/plugins/default_log_template.html)\n\n#
  Disable Logging\n\nIf you do not want logging, you can explicityly disable it by
  adding it to your\n`[markata.disabled_hooks]` array in your `[markata.toml]`\n\n```
  toml\n[markata]\n\n# make sure its in your list of hooks\ndisabled_hooks=[\n   \"markata.plugins.setup_logging\",\n
  \  ]\n```\n\n\n!! class <h2 id='SilentUndefined' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>SilentUndefined <em class='small'>class</em></h2>\n\n???+ source \"SilentUndefined
  <em class='small'>source</em>\"\n\n```python\n\n        class SilentUndefined(Undefined):\n
  \           def _fail_with_undefined_error(self, *args, **kwargs):\n                return
  \"\"\n```\n\n\n!! function <h2 id='has_rich_handler' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>has_rich_handler <em class='small'>function</em></h2>\n    Returns a boolean
  whether or not there is a RichHandler attached to the\n    root logger.\n???+ source
  \"has_rich_handler <em class='small'>source</em>\"\n\n```python\n\n        def has_rich_handler()
  -> bool:\n            \"\"\"\n            Returns a boolean whether or not there
  is a RichHandler attached to the\n            root logger.\n            \"\"\"\n\n
  \           logger = logging.getLogger()\n            return bool([h for h in logger.handlers
  if isinstance(h, RichHandler)])\n```\n\n\n!! function <h2 id='has_file_handler'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>has_file_handler <em
  class='small'>function</em></h2>\n\n???+ source \"has_file_handler <em class='small'>source</em>\"\n\n```python\n\n
  \       def has_file_handler(log_file: Path) -> bool:\n            logger = logging.getLogger()\n
  \           existing_logger_files = [\n                handler.baseFilename\n                for
  handler in logger.handlers\n                if isinstance(handler, logging.FileHandler)\n
  \           ]\n            return str(log_file.absolute()) in existing_logger_files\n```\n\n\n!!
  function <h2 id='setup_log' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>setup_log <em class='small'>function</em></h2>\n\n???+ source \"setup_log
  <em class='small'>source</em>\"\n\n```python\n\n        def setup_log(markata: \"Markata\",
  level: int = logging.INFO) -> Path:\n            path = setup_html_log(markata,
  level)\n            setup_text_log(markata, level)\n            return path\n```\n\n\n!!
  class <h2 id='LoggingConfig' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>LoggingConfig <em class='small'>class</em></h2>\n\n???+ source \"LoggingConfig
  <em class='small'>source</em>\"\n\n```python\n\n        class LoggingConfig(pydantic.BaseModel):\n
  \           output_dir: pydantic.DirectoryPath = Path(\"markout\")\n            log_dir:
  Optional[Path] = None\n            template: Optional[Path] = Path(__file__).parent
  / \"default_log_template.html\"\n\n            @pydantic.validator(\"log_dir\",
  pre=True, always=True)\n            def validate_log_dir(cls, v, *, values):\n                if
  v is None:\n                    return values[\"output_dir\"] / \"_logs\"\n                return
  Path(v)\n```\n\n\n!! class <h2 id='Config' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>Config <em class='small'>class</em></h2>\n\n???+ source \"Config <em class='small'>source</em>\"\n\n```python\n\n
  \       class Config(pydantic.BaseModel):\n            logging: LoggingConfig =
  LoggingConfig()\n```\n\n\n!! function <h2 id='config_model' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>config_model <em class='small'>function</em></h2>\n\n???+
  source \"config_model <em class='small'>source</em>\"\n\n```python\n\n        def
  config_model(markata: \"Markata\") -> None:\n            markata.config_models.append(Config)\n```\n\n\n!!
  function <h2 id='setup_text_log' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>setup_text_log <em class='small'>function</em></h2>\n    sets up a plain
  text log in markata's configured `log_dir`, or\n    `output_dir/_logs` if `log_dir`
  is not configured.  The log file will be\n    named after the `<levelname>.log`\n???+
  source \"setup_text_log <em class='small'>source</em>\"\n\n```python\n\n        def
  setup_text_log(markata: \"Markata\", level: int = logging.INFO) -> Path:\n            \"\"\"\n
  \           sets up a plain text log in markata's configured `log_dir`, or\n            `output_dir/_logs`
  if `log_dir` is not configured.  The log file will be\n            named after the
  `<levelname>.log`\n            \"\"\"\n            log_file = markata.config.logging.log_dir
  / (\n                logging.getLevelName(level).lower() + \".log\"\n            )\n\n
  \           if has_file_handler(log_file):\n                return log_file\n\n
  \           if not log_file.parent.exists():\n                log_file.parent.mkdir(parents=True)\n
  \           fh = logging.FileHandler(log_file)\n            fh.setLevel(level)\n
  \           fh_formatter = logging.Formatter(\n                \"%(asctime)s %(name)-12s
  %(levelname)-8s %(message)s\",\n            )\n            fh.setFormatter(fh_formatter)\n
  \           logging.getLogger(\"\").addHandler(fh)\n\n            return log_file\n```\n\n\n!!
  function <h2 id='setup_html_log' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>setup_html_log <em class='small'>function</em></h2>\n    sets up an html
  log in markata's configured `log_dir`, or\n    `output_dir/_logs` if `log_dir` is
  not configured.  The log file will be\n    named after the `<levelname>/index.html`.
  \ The goal of this is to give\n???+ source \"setup_html_log <em class='small'>source</em>\"\n\n```python\n\n
  \       def setup_html_log(markata: \"Markata\", level: int = logging.INFO) -> Path:\n
  \           \"\"\"\n            sets up an html log in markata's configured `log_dir`,
  or\n            `output_dir/_logs` if `log_dir` is not configured.  The log file
  will be\n            named after the `<levelname>/index.html`.  The goal of this
  is to give\n            \"\"\"\n\n            log_file = (\n                markata.config.logging.log_dir\n
  \               / logging.getLevelName(level).lower()\n                / \"index.html\"\n
  \           )\n\n            if has_file_handler(log_file):\n                return
  log_file\n\n            log_file.parent.mkdir(parents=True, exist_ok=True)\n\n            if
  not log_file.exists():\n                template = Template(\n                    markata.config.logging.template.read_text(),
  undefined=SilentUndefined\n                )\n                log_header = template.render(\n
  \                   title=markata.config.title + \" logs\",\n                    config=markata.config,\n
  \               )\n                log_file.write_text(log_header)\n            with
  open(log_file, \"a\") as f:\n                command = Path(sys.argv[0]).name +
  \" \" + \" \".join(sys.argv[1:])\n                f.write(\n                    f\"\"\"\n
  \                   <div style=\"\n                    width: 100%;\n                    height:
  20px;\n                    margin-top: 5rem;\n                    border-bottom:
  1px solid goldenrod;\n                    text-align: center\">\n                        <span
  style=\"padding: 0 10px;\">\n                            {datetime.datetime.now()}
  running \"{command}\"\n                        </span>\n                    </div>\n
  \           \"\"\",\n                )\n            fh = logging.FileHandler(log_file)\n
  \           fh.setLevel(level)\n            fh_formatter = logging.Formatter(\n
  \               \"\"\"\n                <li>\n                    <p>\n                        <span
  class=\"time\">%(asctime)s</span>\n                        <span class=\"name %(name)s\">%(name)-12s</span>\n
  \                       <span class=\"levelname %(levelname)s\">%(levelname)-8s</span>\n
  \                   </p>\n                    <p class=\"message\">%(message)s</p>\n
  \               </li>\n                \"\"\",\n            )\n            fh.setFormatter(fh_formatter)\n
  \           logging.getLogger(\"\").addHandler(fh)\n\n            return log_file\n```\n\n\n!!
  function <h2 id='configure' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>configure <em class='small'>function</em></h2>\n\n???+ source \"configure
  <em class='small'>source</em>\"\n\n```python\n\n        def configure(markata: \"Markata\")
  -> None:\n            setup_log(markata, logging.DEBUG)\n            setup_log(markata,
  logging.INFO)\n            setup_log(markata, logging.WARNING)\n\n            if
  not has_rich_handler():\n                console = RichHandler(\n                    rich_tracebacks=True,\n
  \               )\n                console.setLevel(logging.INFO)\n                formatter
  = logging.Formatter(\"%(message)s\")\n                console.setFormatter(formatter)\n
  \               logging.getLogger(\"\").addHandler(console)\n```\n\n\n!! method
  <h2 id='_fail_with_undefined_error' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>_fail_with_undefined_error <em class='small'>method</em></h2>\n\n???+ source
  \"_fail_with_undefined_error <em class='small'>source</em>\"\n\n```python\n\n        def
  _fail_with_undefined_error(self, *args, **kwargs):\n                return \"\"\n```\n\n\n!!
  method <h2 id='validate_log_dir' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>validate_log_dir <em class='small'>method</em></h2>\n\n???+ source \"validate_log_dir
  <em class='small'>source</em>\"\n\n```python\n\n        def validate_log_dir(cls,
  v, *, values):\n                if v is None:\n                    return values[\"output_dir\"]
  / \"_logs\"\n                return Path(v)\n```\n"
date: 0001-01-01
description: Setup Logging hook sets up the RichHandler for pretty console logs, and
  file There will be 6 log files created based on log level and file type. Ensure
  that set
html: "<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n<title>Setup_Logging.Py</title>\n<meta
  charset=\"utf-8\"/>\n<meta content=\"width=device-width, initial-scale=1\" name=\"viewport\"/>\n<meta
  content=\"Setup Logging hook sets up the RichHandler for pretty console logs, and
  file There will be 6 log files created based on log level and file type. Ensure
  that set\" name=\"description\"/>\n<link href=\"/static/favicon.ico\" rel=\"icon\"
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
  name=\"og:type\" property=\"og:type\"/><meta content=\"Setup Logging hook sets up
  the RichHandler for pretty console logs, and file There will be 6 log files created
  based on log level and file type. Ensure that set\" name=\"description\" property=\"description\"/><meta
  content=\"Setup Logging hook sets up the RichHandler for pretty console logs, and
  file There will be 6 log files created based on log level and file type. Ensure
  that set\" name=\"og:description\" property=\"og:description\"/><meta content=\"Setup
  Logging hook sets up the RichHandler for pretty console logs, and file There will
  be 6 log files created based on log level and file type. Ensure that set\" name=\"twitter:description\"
  property=\"twitter:description\"/><meta content=\"Setup_Logging.Py | Markata\" name=\"og:title\"
  property=\"og:title\"/><meta content=\"Setup_Logging.Py | Markata\" name=\"twitter:title\"
  property=\"twitter:title\"/><meta content=\"https://markata.dev//markata/plugins/setup-logging-og.png\"
  name=\"og:image\" property=\"og:image\"/><meta content=\"https://markata.dev//markata/plugins/setup-logging-og.png\"
  name=\"twitter:image\" property=\"twitter:image\"/><meta content=\"1600\" name=\"og:image:width\"
  property=\"og:image:width\"/><meta content=\"900\" name=\"og:image:width\" property=\"og:image:width\"/><meta
  content=\"summary_large_image\" name=\"twitter:card\" property=\"twitter:card\"/><meta
  content=\"Markata\" name=\"og:site_name\" property=\"og:site_name\"/><meta content=\"@_waylonwalker\"
  name=\"twitter:creator\" property=\"twitter:creator\"/><meta content=\"Setup_Logging.Py\"
  name=\"title\" property=\"title\"/><meta content=\"markata 0.8.0.dev6\" name=\"generator\"
  property=\"generator\"/><link href=\"https://markata.dev//markata/plugins/setup-logging/\"
  rel=\"canonical\"/><meta content=\"https://markata.dev//markata/plugins/setup-logging/\"
  name=\"og:url\" property=\"og:url\"/></head>\n<body><nav>\n<a href=\"/\">markata</a>\n<a
  href=\"https://github.com/WaylonWalker/markata\">GitHub</a>\n</nav>\n<div>\n<label
  class=\"theme-switch\" for=\"checkbox-theme\" id=\"theme-switch\" title=\"light/dark
  mode toggle\">\n<input id=\"checkbox-theme\" type=\"checkbox\"/>\n<div class=\"slider
  round\"></div>\n</label>\n</div>\n<section class=\"title\">\n<h1 id=\"title\">\n
  \           Setup_Logging.Py \n            \n        </h1>\n</section>\n<main><p>Setup
  Logging hook sets up the RichHandler for pretty console logs, and file\nlogs to
  the configured markata's configured <code>log_dir</code>, or <code>output_dir/_logs</code>
  if\n<code>log_dir</code> is not configured.  The log file will be named after the\n<code>&lt;levelname&gt;.log</code></p>\n<h1
  id=\"the-log-files\">The log files <a class=\"header-anchor\" href=\"#the-log-files\"><svg
  aria-hidden=\"true\" class=\"heading-permalink\" fill=\"currentColor\" focusable=\"false\"
  height=\"1em\" viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
  d=\"M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>There will be 6 log files
  created based on log level and file type.</p>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
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
  \       \n<div class=\"highlight\"><pre><span></span>markout/_logs\n\u251C\u2500\u2500
  debug\n\u2502   \u2514\u2500\u2500 index.html\n\u251C\u2500\u2500 debug.log\n\u251C\u2500\u2500
  info\n\u2502   \u2514\u2500\u2500 index.html\n\u251C\u2500\u2500 info.log\n\u251C\u2500\u2500
  warning\n\u2502   \u2514\u2500\u2500 index.html\n\u2514\u2500\u2500 warning.log\n</pre></div>\n\n</pre>\n<h1
  id=\"configuration\">Configuration <a class=\"header-anchor\" href=\"#configuration\"><svg
  aria-hidden=\"true\" class=\"heading-permalink\" fill=\"currentColor\" focusable=\"false\"
  height=\"1em\" viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
  d=\"M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>Ensure that setup_logging
  is in your hooks.  You can check if <code>setup_logging</code>\nis in your hooks
  by running <code>markata list --hooks</code> from your terminal and\nchecking the
  output, or creating an instance of <code>Markata()</code> and checking the\n<code>Markata().hooks</code>
  attribute.  If its missing or you wan to be more explicit,\nyou can add <code>setup_logging</code>
  to your <code>markata.toml</code> <code>[markata.hooks]</code>.</p>\n<pre class=\"wrapper\">\n\n<div
  class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
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
  \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata]</span>\n\n<span
  class=\"c1\"># make sure its in your list of hooks</span>\n<span class=\"n\">hooks</span><span
  class=\"o\">=</span><span class=\"p\">[</span>\n<span class=\"w\">   </span><span
  class=\"s2\">\"markata.plugins.setup_logging\"</span><span class=\"p\">,</span>\n<span
  class=\"w\">   </span><span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n<h1 id=\"log-template\">Log
  Template <a class=\"header-anchor\" href=\"#log-template\"><svg aria-hidden=\"true\"
  class=\"heading-permalink\" fill=\"currentColor\" focusable=\"false\" height=\"1em\"
  viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path d=\"M9.199
  13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<pre class=\"wrapper\">\n\n<div
  class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
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
  \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata]</span>\n\n<span
  class=\"c1\"># make sure its in your list of hooks</span>\n<span class=\"n\">hooks</span><span
  class=\"o\">=</span><span class=\"p\">[</span>\n<span class=\"w\">   </span><span
  class=\"s2\">\"markata.plugins.setup_logging\"</span><span class=\"p\">,</span>\n<span
  class=\"w\">   </span><span class=\"p\">]</span>\n\n<span class=\"c1\"># point log
  template to the path of your logging template</span>\n<span class=\"n\">log_template</span><span
  class=\"o\">=</span><span class=\"s1\">'templates/log_template.html'</span>\n</pre></div>\n\n</pre>\n<p>You
  can see the latest default <code>log_template</code> on\n<a href=\"https://github.com/WaylonWalker/markata/blob/main/markata/plugins/default_log_template.html\">GitHub</a></p>\n<h1
  id=\"disable-logging\">Disable Logging <a class=\"header-anchor\" href=\"#disable-logging\"><svg
  aria-hidden=\"true\" class=\"heading-permalink\" fill=\"currentColor\" focusable=\"false\"
  height=\"1em\" viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
  d=\"M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>If you do not want logging,
  you can explicityly disable it by adding it to your\n<code>[markata.disabled_hooks]</code>
  array in your <code>[markata.toml]</code></p>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
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
  \       \n<div class=\"highlight\"><pre><span></span><span class=\"k\">[markata]</span>\n\n<span
  class=\"c1\"># make sure its in your list of hooks</span>\n<span class=\"n\">disabled_hooks</span><span
  class=\"o\">=</span><span class=\"p\">[</span>\n<span class=\"w\">   </span><span
  class=\"s2\">\"markata.plugins.setup_logging\"</span><span class=\"p\">,</span>\n<span
  class=\"w\">   </span><span class=\"p\">]</span>\n</pre></div>\n\n</pre>\n<p>!!
  class </p><h2 class=\"admonition-title\" id=\"SilentUndefined\" style=\"margin:0;padding:.5rem
  1rem;\">SilentUndefined <em class=\"small\">class</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"SilentUndefined
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
  <span class=\"nc\">SilentUndefined</span><span class=\"p\">(</span><span class=\"n\">Undefined</span><span
  class=\"p\">):</span>\n            <span class=\"k\">def</span> <span class=\"nf\">_fail_with_undefined_error</span><span
  class=\"p\">(</span><span class=\"bp\">self</span><span class=\"p\">,</span> <span
  class=\"o\">*</span><span class=\"n\">args</span><span class=\"p\">,</span> <span
  class=\"o\">**</span><span class=\"n\">kwargs</span><span class=\"p\">):</span>\n
  \               <span class=\"k\">return</span> <span class=\"s2\">\"\"</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"has_rich_handler\" style=\"margin:0;padding:.5rem
  1rem;\">has_rich_handler <em class=\"small\">function</em></h2>\nReturns a boolean
  whether or not there is a RichHandler attached to the\nroot logger.\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"has_rich_handler
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
  <span class=\"nf\">has_rich_handler</span><span class=\"p\">()</span> <span class=\"o\">-&gt;</span>
  <span class=\"nb\">bool</span><span class=\"p\">:</span>\n<span class=\"w\">            </span><span
  class=\"sd\">\"\"\"</span>\n<span class=\"sd\">            Returns a boolean whether
  or not there is a RichHandler attached to the</span>\n<span class=\"sd\">            root
  logger.</span>\n<span class=\"sd\">            \"\"\"</span>\n\n            <span
  class=\"n\">logger</span> <span class=\"o\">=</span> <span class=\"n\">logging</span><span
  class=\"o\">.</span><span class=\"n\">getLogger</span><span class=\"p\">()</span>\n
  \           <span class=\"k\">return</span> <span class=\"nb\">bool</span><span
  class=\"p\">([</span><span class=\"n\">h</span> <span class=\"k\">for</span> <span
  class=\"n\">h</span> <span class=\"ow\">in</span> <span class=\"n\">logger</span><span
  class=\"o\">.</span><span class=\"n\">handlers</span> <span class=\"k\">if</span>
  <span class=\"nb\">isinstance</span><span class=\"p\">(</span><span class=\"n\">h</span><span
  class=\"p\">,</span> <span class=\"n\">RichHandler</span><span class=\"p\">)])</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"has_file_handler\" style=\"margin:0;padding:.5rem
  1rem;\">has_file_handler <em class=\"small\">function</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"has_file_handler
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
  <span class=\"nf\">has_file_handler</span><span class=\"p\">(</span><span class=\"n\">log_file</span><span
  class=\"p\">:</span> <span class=\"n\">Path</span><span class=\"p\">)</span> <span
  class=\"o\">-&gt;</span> <span class=\"nb\">bool</span><span class=\"p\">:</span>\n
  \           <span class=\"n\">logger</span> <span class=\"o\">=</span> <span class=\"n\">logging</span><span
  class=\"o\">.</span><span class=\"n\">getLogger</span><span class=\"p\">()</span>\n
  \           <span class=\"n\">existing_logger_files</span> <span class=\"o\">=</span>
  <span class=\"p\">[</span>\n                <span class=\"n\">handler</span><span
  class=\"o\">.</span><span class=\"n\">baseFilename</span>\n                <span
  class=\"k\">for</span> <span class=\"n\">handler</span> <span class=\"ow\">in</span>
  <span class=\"n\">logger</span><span class=\"o\">.</span><span class=\"n\">handlers</span>\n
  \               <span class=\"k\">if</span> <span class=\"nb\">isinstance</span><span
  class=\"p\">(</span><span class=\"n\">handler</span><span class=\"p\">,</span> <span
  class=\"n\">logging</span><span class=\"o\">.</span><span class=\"n\">FileHandler</span><span
  class=\"p\">)</span>\n            <span class=\"p\">]</span>\n            <span
  class=\"k\">return</span> <span class=\"nb\">str</span><span class=\"p\">(</span><span
  class=\"n\">log_file</span><span class=\"o\">.</span><span class=\"n\">absolute</span><span
  class=\"p\">())</span> <span class=\"ow\">in</span> <span class=\"n\">existing_logger_files</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"setup_log\" style=\"margin:0;padding:.5rem
  1rem;\">setup_log <em class=\"small\">function</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"setup_log
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
  <span class=\"nf\">setup_log</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">,</span>
  <span class=\"n\">level</span><span class=\"p\">:</span> <span class=\"nb\">int</span>
  <span class=\"o\">=</span> <span class=\"n\">logging</span><span class=\"o\">.</span><span
  class=\"n\">INFO</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
  <span class=\"n\">Path</span><span class=\"p\">:</span>\n            <span class=\"n\">path</span>
  <span class=\"o\">=</span> <span class=\"n\">setup_html_log</span><span class=\"p\">(</span><span
  class=\"n\">markata</span><span class=\"p\">,</span> <span class=\"n\">level</span><span
  class=\"p\">)</span>\n            <span class=\"n\">setup_text_log</span><span class=\"p\">(</span><span
  class=\"n\">markata</span><span class=\"p\">,</span> <span class=\"n\">level</span><span
  class=\"p\">)</span>\n            <span class=\"k\">return</span> <span class=\"n\">path</span>\n</pre></div>\n\n</pre>\n<p>!!
  class </p><h2 class=\"admonition-title\" id=\"LoggingConfig\" style=\"margin:0;padding:.5rem
  1rem;\">LoggingConfig <em class=\"small\">class</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"LoggingConfig
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
  <span class=\"nc\">LoggingConfig</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
  class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
  \           <span class=\"n\">output_dir</span><span class=\"p\">:</span> <span
  class=\"n\">pydantic</span><span class=\"o\">.</span><span class=\"n\">DirectoryPath</span>
  <span class=\"o\">=</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
  class=\"s2\">\"markout\"</span><span class=\"p\">)</span>\n            <span class=\"n\">log_dir</span><span
  class=\"p\">:</span> <span class=\"n\">Optional</span><span class=\"p\">[</span><span
  class=\"n\">Path</span><span class=\"p\">]</span> <span class=\"o\">=</span> <span
  class=\"kc\">None</span>\n            <span class=\"n\">template</span><span class=\"p\">:</span>
  <span class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"n\">Path</span><span
  class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"n\">Path</span><span
  class=\"p\">(</span><span class=\"vm\">__file__</span><span class=\"p\">)</span><span
  class=\"o\">.</span><span class=\"n\">parent</span> <span class=\"o\">/</span> <span
  class=\"s2\">\"default_log_template.html\"</span>\n\n            <span class=\"nd\">@pydantic</span><span
  class=\"o\">.</span><span class=\"n\">validator</span><span class=\"p\">(</span><span
  class=\"s2\">\"log_dir\"</span><span class=\"p\">,</span> <span class=\"n\">pre</span><span
  class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
  class=\"n\">always</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
  class=\"p\">)</span>\n            <span class=\"k\">def</span> <span class=\"nf\">validate_log_dir</span><span
  class=\"p\">(</span><span class=\"bp\">cls</span><span class=\"p\">,</span> <span
  class=\"n\">v</span><span class=\"p\">,</span> <span class=\"o\">*</span><span class=\"p\">,</span>
  <span class=\"n\">values</span><span class=\"p\">):</span>\n                <span
  class=\"k\">if</span> <span class=\"n\">v</span> <span class=\"ow\">is</span> <span
  class=\"kc\">None</span><span class=\"p\">:</span>\n                    <span class=\"k\">return</span>
  <span class=\"n\">values</span><span class=\"p\">[</span><span class=\"s2\">\"output_dir\"</span><span
  class=\"p\">]</span> <span class=\"o\">/</span> <span class=\"s2\">\"_logs\"</span>\n
  \               <span class=\"k\">return</span> <span class=\"n\">Path</span><span
  class=\"p\">(</span><span class=\"n\">v</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
  class </p><h2 class=\"admonition-title\" id=\"Config\" style=\"margin:0;padding:.5rem
  1rem;\">Config <em class=\"small\">class</em></h2>\n<div class=\"admonition source
  is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"Config <em class=\"small\">source</em>\"</p>\n</div>\n<pre
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
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">class</span>
  <span class=\"nc\">Config</span><span class=\"p\">(</span><span class=\"n\">pydantic</span><span
  class=\"o\">.</span><span class=\"n\">BaseModel</span><span class=\"p\">):</span>\n
  \           <span class=\"n\">logging</span><span class=\"p\">:</span> <span class=\"n\">LoggingConfig</span>
  <span class=\"o\">=</span> <span class=\"n\">LoggingConfig</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"config_model\" style=\"margin:0;padding:.5rem
  1rem;\">config_model <em class=\"small\">function</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"config_model
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
  <span class=\"nf\">config_model</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
  \           <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config_models</span><span
  class=\"o\">.</span><span class=\"n\">append</span><span class=\"p\">(</span><span
  class=\"n\">Config</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"setup_text_log\" style=\"margin:0;padding:.5rem
  1rem;\">setup_text_log <em class=\"small\">function</em></h2>\nsets up a plain text
  log in markata's configured <code>log_dir</code>, or\n<code>output_dir/_logs</code>
  if <code>log_dir</code> is not configured.  The log file will be\nnamed after the
  <code>&lt;levelname&gt;.log</code>\n<div class=\"admonition source is-collapsible
  collapsible-open\">\n<p class=\"admonition-title\">\"setup_text_log <em class=\"small\">source</em>\"</p>\n</div>\n<pre
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
  <span class=\"nf\">setup_text_log</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">,</span>
  <span class=\"n\">level</span><span class=\"p\">:</span> <span class=\"nb\">int</span>
  <span class=\"o\">=</span> <span class=\"n\">logging</span><span class=\"o\">.</span><span
  class=\"n\">INFO</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
  <span class=\"n\">Path</span><span class=\"p\">:</span>\n<span class=\"w\">            </span><span
  class=\"sd\">\"\"\"</span>\n<span class=\"sd\">            sets up a plain text
  log in markata's configured `log_dir`, or</span>\n<span class=\"sd\">            `output_dir/_logs`
  if `log_dir` is not configured.  The log file will be</span>\n<span class=\"sd\">
  \           named after the `&lt;levelname&gt;.log`</span>\n<span class=\"sd\">
  \           \"\"\"</span>\n            <span class=\"n\">log_file</span> <span class=\"o\">=</span>
  <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"o\">.</span><span class=\"n\">logging</span><span class=\"o\">.</span><span
  class=\"n\">log_dir</span> <span class=\"o\">/</span> <span class=\"p\">(</span>\n
  \               <span class=\"n\">logging</span><span class=\"o\">.</span><span
  class=\"n\">getLevelName</span><span class=\"p\">(</span><span class=\"n\">level</span><span
  class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">lower</span><span
  class=\"p\">()</span> <span class=\"o\">+</span> <span class=\"s2\">\".log\"</span>\n
  \           <span class=\"p\">)</span>\n\n            <span class=\"k\">if</span>
  <span class=\"n\">has_file_handler</span><span class=\"p\">(</span><span class=\"n\">log_file</span><span
  class=\"p\">):</span>\n                <span class=\"k\">return</span> <span class=\"n\">log_file</span>\n\n
  \           <span class=\"k\">if</span> <span class=\"ow\">not</span> <span class=\"n\">log_file</span><span
  class=\"o\">.</span><span class=\"n\">parent</span><span class=\"o\">.</span><span
  class=\"n\">exists</span><span class=\"p\">():</span>\n                <span class=\"n\">log_file</span><span
  class=\"o\">.</span><span class=\"n\">parent</span><span class=\"o\">.</span><span
  class=\"n\">mkdir</span><span class=\"p\">(</span><span class=\"n\">parents</span><span
  class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">)</span>\n            <span
  class=\"n\">fh</span> <span class=\"o\">=</span> <span class=\"n\">logging</span><span
  class=\"o\">.</span><span class=\"n\">FileHandler</span><span class=\"p\">(</span><span
  class=\"n\">log_file</span><span class=\"p\">)</span>\n            <span class=\"n\">fh</span><span
  class=\"o\">.</span><span class=\"n\">setLevel</span><span class=\"p\">(</span><span
  class=\"n\">level</span><span class=\"p\">)</span>\n            <span class=\"n\">fh_formatter</span>
  <span class=\"o\">=</span> <span class=\"n\">logging</span><span class=\"o\">.</span><span
  class=\"n\">Formatter</span><span class=\"p\">(</span>\n                <span class=\"s2\">\"</span><span
  class=\"si\">%(asctime)s</span><span class=\"s2\"> </span><span class=\"si\">%(name)-12s</span><span
  class=\"s2\"> </span><span class=\"si\">%(levelname)-8s</span><span class=\"s2\">
  </span><span class=\"si\">%(message)s</span><span class=\"s2\">\"</span><span class=\"p\">,</span>\n
  \           <span class=\"p\">)</span>\n            <span class=\"n\">fh</span><span
  class=\"o\">.</span><span class=\"n\">setFormatter</span><span class=\"p\">(</span><span
  class=\"n\">fh_formatter</span><span class=\"p\">)</span>\n            <span class=\"n\">logging</span><span
  class=\"o\">.</span><span class=\"n\">getLogger</span><span class=\"p\">(</span><span
  class=\"s2\">\"\"</span><span class=\"p\">)</span><span class=\"o\">.</span><span
  class=\"n\">addHandler</span><span class=\"p\">(</span><span class=\"n\">fh</span><span
  class=\"p\">)</span>\n\n            <span class=\"k\">return</span> <span class=\"n\">log_file</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"setup_html_log\" style=\"margin:0;padding:.5rem
  1rem;\">setup_html_log <em class=\"small\">function</em></h2>\nsets up an html log
  in markata's configured <code>log_dir</code>, or\n<code>output_dir/_logs</code>
  if <code>log_dir</code> is not configured.  The log file will be\nnamed after the
  <code>&lt;levelname&gt;/index.html</code>.  The goal of this is to give\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"setup_html_log
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
  <span class=\"nf\">setup_html_log</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">,</span>
  <span class=\"n\">level</span><span class=\"p\">:</span> <span class=\"nb\">int</span>
  <span class=\"o\">=</span> <span class=\"n\">logging</span><span class=\"o\">.</span><span
  class=\"n\">INFO</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
  <span class=\"n\">Path</span><span class=\"p\">:</span>\n<span class=\"w\">            </span><span
  class=\"sd\">\"\"\"</span>\n<span class=\"sd\">            sets up an html log in
  markata's configured `log_dir`, or</span>\n<span class=\"sd\">            `output_dir/_logs`
  if `log_dir` is not configured.  The log file will be</span>\n<span class=\"sd\">
  \           named after the `&lt;levelname&gt;/index.html`.  The goal of this is
  to give</span>\n<span class=\"sd\">            \"\"\"</span>\n\n            <span
  class=\"n\">log_file</span> <span class=\"o\">=</span> <span class=\"p\">(</span>\n
  \               <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">logging</span><span
  class=\"o\">.</span><span class=\"n\">log_dir</span>\n                <span class=\"o\">/</span>
  <span class=\"n\">logging</span><span class=\"o\">.</span><span class=\"n\">getLevelName</span><span
  class=\"p\">(</span><span class=\"n\">level</span><span class=\"p\">)</span><span
  class=\"o\">.</span><span class=\"n\">lower</span><span class=\"p\">()</span>\n
  \               <span class=\"o\">/</span> <span class=\"s2\">\"index.html\"</span>\n
  \           <span class=\"p\">)</span>\n\n            <span class=\"k\">if</span>
  <span class=\"n\">has_file_handler</span><span class=\"p\">(</span><span class=\"n\">log_file</span><span
  class=\"p\">):</span>\n                <span class=\"k\">return</span> <span class=\"n\">log_file</span>\n\n
  \           <span class=\"n\">log_file</span><span class=\"o\">.</span><span class=\"n\">parent</span><span
  class=\"o\">.</span><span class=\"n\">mkdir</span><span class=\"p\">(</span><span
  class=\"n\">parents</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
  class=\"p\">,</span> <span class=\"n\">exist_ok</span><span class=\"o\">=</span><span
  class=\"kc\">True</span><span class=\"p\">)</span>\n\n            <span class=\"k\">if</span>
  <span class=\"ow\">not</span> <span class=\"n\">log_file</span><span class=\"o\">.</span><span
  class=\"n\">exists</span><span class=\"p\">():</span>\n                <span class=\"n\">template</span>
  <span class=\"o\">=</span> <span class=\"n\">Template</span><span class=\"p\">(</span>\n
  \                   <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">logging</span><span
  class=\"o\">.</span><span class=\"n\">template</span><span class=\"o\">.</span><span
  class=\"n\">read_text</span><span class=\"p\">(),</span> <span class=\"n\">undefined</span><span
  class=\"o\">=</span><span class=\"n\">SilentUndefined</span>\n                <span
  class=\"p\">)</span>\n                <span class=\"n\">log_header</span> <span
  class=\"o\">=</span> <span class=\"n\">template</span><span class=\"o\">.</span><span
  class=\"n\">render</span><span class=\"p\">(</span>\n                    <span class=\"n\">title</span><span
  class=\"o\">=</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">title</span>
  <span class=\"o\">+</span> <span class=\"s2\">\" logs\"</span><span class=\"p\">,</span>\n
  \                   <span class=\"n\">config</span><span class=\"o\">=</span><span
  class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"p\">,</span>\n                <span class=\"p\">)</span>\n                <span
  class=\"n\">log_file</span><span class=\"o\">.</span><span class=\"n\">write_text</span><span
  class=\"p\">(</span><span class=\"n\">log_header</span><span class=\"p\">)</span>\n
  \           <span class=\"k\">with</span> <span class=\"nb\">open</span><span class=\"p\">(</span><span
  class=\"n\">log_file</span><span class=\"p\">,</span> <span class=\"s2\">\"a\"</span><span
  class=\"p\">)</span> <span class=\"k\">as</span> <span class=\"n\">f</span><span
  class=\"p\">:</span>\n                <span class=\"n\">command</span> <span class=\"o\">=</span>
  <span class=\"n\">Path</span><span class=\"p\">(</span><span class=\"n\">sys</span><span
  class=\"o\">.</span><span class=\"n\">argv</span><span class=\"p\">[</span><span
  class=\"mi\">0</span><span class=\"p\">])</span><span class=\"o\">.</span><span
  class=\"n\">name</span> <span class=\"o\">+</span> <span class=\"s2\">\" \"</span>
  <span class=\"o\">+</span> <span class=\"s2\">\" \"</span><span class=\"o\">.</span><span
  class=\"n\">join</span><span class=\"p\">(</span><span class=\"n\">sys</span><span
  class=\"o\">.</span><span class=\"n\">argv</span><span class=\"p\">[</span><span
  class=\"mi\">1</span><span class=\"p\">:])</span>\n                <span class=\"n\">f</span><span
  class=\"o\">.</span><span class=\"n\">write</span><span class=\"p\">(</span>\n                    <span
  class=\"sa\">f</span><span class=\"s2\">\"\"\"</span>\n<span class=\"s2\">                    &lt;div
  style=\"</span>\n<span class=\"s2\">                    width: 100%;</span>\n<span
  class=\"s2\">                    height: 20px;</span>\n<span class=\"s2\">                    margin-top:
  5rem;</span>\n<span class=\"s2\">                    border-bottom: 1px solid goldenrod;</span>\n<span
  class=\"s2\">                    text-align: center\"&gt;</span>\n<span class=\"s2\">
  \                       &lt;span style=\"padding: 0 10px;\"&gt;</span>\n<span class=\"s2\">
  \                           </span><span class=\"si\">{</span><span class=\"n\">datetime</span><span
  class=\"o\">.</span><span class=\"n\">datetime</span><span class=\"o\">.</span><span
  class=\"n\">now</span><span class=\"p\">()</span><span class=\"si\">}</span><span
  class=\"s2\"> running \"</span><span class=\"si\">{</span><span class=\"n\">command</span><span
  class=\"si\">}</span><span class=\"s2\">\"</span>\n<span class=\"s2\">                        &lt;/span&gt;</span>\n<span
  class=\"s2\">                    &lt;/div&gt;</span>\n<span class=\"s2\">            \"\"\"</span><span
  class=\"p\">,</span>\n                <span class=\"p\">)</span>\n            <span
  class=\"n\">fh</span> <span class=\"o\">=</span> <span class=\"n\">logging</span><span
  class=\"o\">.</span><span class=\"n\">FileHandler</span><span class=\"p\">(</span><span
  class=\"n\">log_file</span><span class=\"p\">)</span>\n            <span class=\"n\">fh</span><span
  class=\"o\">.</span><span class=\"n\">setLevel</span><span class=\"p\">(</span><span
  class=\"n\">level</span><span class=\"p\">)</span>\n            <span class=\"n\">fh_formatter</span>
  <span class=\"o\">=</span> <span class=\"n\">logging</span><span class=\"o\">.</span><span
  class=\"n\">Formatter</span><span class=\"p\">(</span>\n<span class=\"w\">                </span><span
  class=\"sd\">\"\"\"</span>\n<span class=\"sd\">                &lt;li&gt;</span>\n<span
  class=\"sd\">                    &lt;p&gt;</span>\n<span class=\"sd\">                        &lt;span
  class=\"time\"&gt;%(asctime)s&lt;/span&gt;</span>\n<span class=\"sd\">                        &lt;span
  class=\"name %(name)s\"&gt;%(name)-12s&lt;/span&gt;</span>\n<span class=\"sd\">
  \                       &lt;span class=\"levelname %(levelname)s\"&gt;%(levelname)-8s&lt;/span&gt;</span>\n<span
  class=\"sd\">                    &lt;/p&gt;</span>\n<span class=\"sd\">                    &lt;p
  class=\"message\"&gt;%(message)s&lt;/p&gt;</span>\n<span class=\"sd\">                &lt;/li&gt;</span>\n<span
  class=\"sd\">                \"\"\"</span><span class=\"p\">,</span>\n            <span
  class=\"p\">)</span>\n            <span class=\"n\">fh</span><span class=\"o\">.</span><span
  class=\"n\">setFormatter</span><span class=\"p\">(</span><span class=\"n\">fh_formatter</span><span
  class=\"p\">)</span>\n            <span class=\"n\">logging</span><span class=\"o\">.</span><span
  class=\"n\">getLogger</span><span class=\"p\">(</span><span class=\"s2\">\"\"</span><span
  class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">addHandler</span><span
  class=\"p\">(</span><span class=\"n\">fh</span><span class=\"p\">)</span>\n\n            <span
  class=\"k\">return</span> <span class=\"n\">log_file</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"configure\" style=\"margin:0;padding:.5rem
  1rem;\">configure <em class=\"small\">function</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"configure
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
  \           <span class=\"n\">setup_log</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
  class=\"p\">,</span> <span class=\"n\">logging</span><span class=\"o\">.</span><span
  class=\"n\">DEBUG</span><span class=\"p\">)</span>\n            <span class=\"n\">setup_log</span><span
  class=\"p\">(</span><span class=\"n\">markata</span><span class=\"p\">,</span> <span
  class=\"n\">logging</span><span class=\"o\">.</span><span class=\"n\">INFO</span><span
  class=\"p\">)</span>\n            <span class=\"n\">setup_log</span><span class=\"p\">(</span><span
  class=\"n\">markata</span><span class=\"p\">,</span> <span class=\"n\">logging</span><span
  class=\"o\">.</span><span class=\"n\">WARNING</span><span class=\"p\">)</span>\n\n
  \           <span class=\"k\">if</span> <span class=\"ow\">not</span> <span class=\"n\">has_rich_handler</span><span
  class=\"p\">():</span>\n                <span class=\"n\">console</span> <span class=\"o\">=</span>
  <span class=\"n\">RichHandler</span><span class=\"p\">(</span>\n                    <span
  class=\"n\">rich_tracebacks</span><span class=\"o\">=</span><span class=\"kc\">True</span><span
  class=\"p\">,</span>\n                <span class=\"p\">)</span>\n                <span
  class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">setLevel</span><span
  class=\"p\">(</span><span class=\"n\">logging</span><span class=\"o\">.</span><span
  class=\"n\">INFO</span><span class=\"p\">)</span>\n                <span class=\"n\">formatter</span>
  <span class=\"o\">=</span> <span class=\"n\">logging</span><span class=\"o\">.</span><span
  class=\"n\">Formatter</span><span class=\"p\">(</span><span class=\"s2\">\"</span><span
  class=\"si\">%(message)s</span><span class=\"s2\">\"</span><span class=\"p\">)</span>\n
  \               <span class=\"n\">console</span><span class=\"o\">.</span><span
  class=\"n\">setFormatter</span><span class=\"p\">(</span><span class=\"n\">formatter</span><span
  class=\"p\">)</span>\n                <span class=\"n\">logging</span><span class=\"o\">.</span><span
  class=\"n\">getLogger</span><span class=\"p\">(</span><span class=\"s2\">\"\"</span><span
  class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">addHandler</span><span
  class=\"p\">(</span><span class=\"n\">console</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
  method </p><h2 class=\"admonition-title\" id=\"_fail_with_undefined_error\" style=\"margin:0;padding:.5rem
  1rem;\">_fail_with_undefined_error <em class=\"small\">method</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"_fail_with_undefined_error
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
  <span class=\"nf\">_fail_with_undefined_error</span><span class=\"p\">(</span><span
  class=\"bp\">self</span><span class=\"p\">,</span> <span class=\"o\">*</span><span
  class=\"n\">args</span><span class=\"p\">,</span> <span class=\"o\">**</span><span
  class=\"n\">kwargs</span><span class=\"p\">):</span>\n                <span class=\"k\">return</span>
  <span class=\"s2\">\"\"</span>\n</pre></div>\n\n</pre>\n<p>!! method </p><h2 class=\"admonition-title\"
  id=\"validate_log_dir\" style=\"margin:0;padding:.5rem 1rem;\">validate_log_dir
  <em class=\"small\">method</em></h2>\n<div class=\"admonition source is-collapsible
  collapsible-open\">\n<p class=\"admonition-title\">\"validate_log_dir <em class=\"small\">source</em>\"</p>\n</div>\n<pre
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
  <span class=\"nf\">validate_log_dir</span><span class=\"p\">(</span><span class=\"bp\">cls</span><span
  class=\"p\">,</span> <span class=\"n\">v</span><span class=\"p\">,</span> <span
  class=\"o\">*</span><span class=\"p\">,</span> <span class=\"n\">values</span><span
  class=\"p\">):</span>\n                <span class=\"k\">if</span> <span class=\"n\">v</span>
  <span class=\"ow\">is</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
  \                   <span class=\"k\">return</span> <span class=\"n\">values</span><span
  class=\"p\">[</span><span class=\"s2\">\"output_dir\"</span><span class=\"p\">]</span>
  <span class=\"o\">/</span> <span class=\"s2\">\"_logs\"</span>\n                <span
  class=\"k\">return</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
  class=\"n\">v</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n</main>\n<footer>\xA9
  2024</footer>\n</body></html>"
published: true
slug: markata/plugins/setup-logging
title: Setup_Logging.Py


---

Setup Logging hook sets up the RichHandler for pretty console logs, and file
logs to the configured markata's configured `log_dir`, or `output_dir/_logs` if
`log_dir` is not configured.  The log file will be named after the
`<levelname>.log`

# The log files

There will be 6 log files created based on log level and file type.

```
markout/_logs
├── debug
│   └── index.html
├── debug.log
├── info
│   └── index.html
├── info.log
├── warning
│   └── index.html
└── warning.log
```

# Configuration

Ensure that setup_logging is in your hooks.  You can check if `setup_logging`
is in your hooks by running `markata list --hooks` from your terminal and
checking the output, or creating an instance of `Markata()` and checking the
`Markata().hooks` attribute.  If its missing or you wan to be more explicit,
you can add `setup_logging` to your `markata.toml` `[markata.hooks]`.

``` toml
[markata]

# make sure its in your list of hooks
hooks=[
   "markata.plugins.setup_logging",
   ]
```

# Log Template
``` toml
[markata]

# make sure its in your list of hooks
hooks=[
   "markata.plugins.setup_logging",
   ]

# point log template to the path of your logging template
log_template='templates/log_template.html'
```

You can see the latest default `log_template` on
[GitHub](https://github.com/WaylonWalker/markata/blob/main/markata/plugins/default_log_template.html)

# Disable Logging

If you do not want logging, you can explicityly disable it by adding it to your
`[markata.disabled_hooks]` array in your `[markata.toml]`

``` toml
[markata]

# make sure its in your list of hooks
disabled_hooks=[
   "markata.plugins.setup_logging",
   ]
```


!! class <h2 id='SilentUndefined' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>SilentUndefined <em class='small'>class</em></h2>

???+ source "SilentUndefined <em class='small'>source</em>"

```python

        class SilentUndefined(Undefined):
            def _fail_with_undefined_error(self, *args, **kwargs):
                return ""
```


!! function <h2 id='has_rich_handler' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>has_rich_handler <em class='small'>function</em></h2>
    Returns a boolean whether or not there is a RichHandler attached to the
    root logger.
???+ source "has_rich_handler <em class='small'>source</em>"

```python

        def has_rich_handler() -> bool:
            """
            Returns a boolean whether or not there is a RichHandler attached to the
            root logger.
            """

            logger = logging.getLogger()
            return bool([h for h in logger.handlers if isinstance(h, RichHandler)])
```


!! function <h2 id='has_file_handler' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>has_file_handler <em class='small'>function</em></h2>

???+ source "has_file_handler <em class='small'>source</em>"

```python

        def has_file_handler(log_file: Path) -> bool:
            logger = logging.getLogger()
            existing_logger_files = [
                handler.baseFilename
                for handler in logger.handlers
                if isinstance(handler, logging.FileHandler)
            ]
            return str(log_file.absolute()) in existing_logger_files
```


!! function <h2 id='setup_log' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>setup_log <em class='small'>function</em></h2>

???+ source "setup_log <em class='small'>source</em>"

```python

        def setup_log(markata: "Markata", level: int = logging.INFO) -> Path:
            path = setup_html_log(markata, level)
            setup_text_log(markata, level)
            return path
```


!! class <h2 id='LoggingConfig' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>LoggingConfig <em class='small'>class</em></h2>

???+ source "LoggingConfig <em class='small'>source</em>"

```python

        class LoggingConfig(pydantic.BaseModel):
            output_dir: pydantic.DirectoryPath = Path("markout")
            log_dir: Optional[Path] = None
            template: Optional[Path] = Path(__file__).parent / "default_log_template.html"

            @pydantic.validator("log_dir", pre=True, always=True)
            def validate_log_dir(cls, v, *, values):
                if v is None:
                    return values["output_dir"] / "_logs"
                return Path(v)
```


!! class <h2 id='Config' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>Config <em class='small'>class</em></h2>

???+ source "Config <em class='small'>source</em>"

```python

        class Config(pydantic.BaseModel):
            logging: LoggingConfig = LoggingConfig()
```


!! function <h2 id='config_model' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>config_model <em class='small'>function</em></h2>

???+ source "config_model <em class='small'>source</em>"

```python

        def config_model(markata: "Markata") -> None:
            markata.config_models.append(Config)
```


!! function <h2 id='setup_text_log' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>setup_text_log <em class='small'>function</em></h2>
    sets up a plain text log in markata's configured `log_dir`, or
    `output_dir/_logs` if `log_dir` is not configured.  The log file will be
    named after the `<levelname>.log`
???+ source "setup_text_log <em class='small'>source</em>"

```python

        def setup_text_log(markata: "Markata", level: int = logging.INFO) -> Path:
            """
            sets up a plain text log in markata's configured `log_dir`, or
            `output_dir/_logs` if `log_dir` is not configured.  The log file will be
            named after the `<levelname>.log`
            """
            log_file = markata.config.logging.log_dir / (
                logging.getLevelName(level).lower() + ".log"
            )

            if has_file_handler(log_file):
                return log_file

            if not log_file.parent.exists():
                log_file.parent.mkdir(parents=True)
            fh = logging.FileHandler(log_file)
            fh.setLevel(level)
            fh_formatter = logging.Formatter(
                "%(asctime)s %(name)-12s %(levelname)-8s %(message)s",
            )
            fh.setFormatter(fh_formatter)
            logging.getLogger("").addHandler(fh)

            return log_file
```


!! function <h2 id='setup_html_log' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>setup_html_log <em class='small'>function</em></h2>
    sets up an html log in markata's configured `log_dir`, or
    `output_dir/_logs` if `log_dir` is not configured.  The log file will be
    named after the `<levelname>/index.html`.  The goal of this is to give
???+ source "setup_html_log <em class='small'>source</em>"

```python

        def setup_html_log(markata: "Markata", level: int = logging.INFO) -> Path:
            """
            sets up an html log in markata's configured `log_dir`, or
            `output_dir/_logs` if `log_dir` is not configured.  The log file will be
            named after the `<levelname>/index.html`.  The goal of this is to give
            """

            log_file = (
                markata.config.logging.log_dir
                / logging.getLevelName(level).lower()
                / "index.html"
            )

            if has_file_handler(log_file):
                return log_file

            log_file.parent.mkdir(parents=True, exist_ok=True)

            if not log_file.exists():
                template = Template(
                    markata.config.logging.template.read_text(), undefined=SilentUndefined
                )
                log_header = template.render(
                    title=markata.config.title + " logs",
                    config=markata.config,
                )
                log_file.write_text(log_header)
            with open(log_file, "a") as f:
                command = Path(sys.argv[0]).name + " " + " ".join(sys.argv[1:])
                f.write(
                    f"""
                    <div style="
                    width: 100%;
                    height: 20px;
                    margin-top: 5rem;
                    border-bottom: 1px solid goldenrod;
                    text-align: center">
                        <span style="padding: 0 10px;">
                            {datetime.datetime.now()} running "{command}"
                        </span>
                    </div>
            """,
                )
            fh = logging.FileHandler(log_file)
            fh.setLevel(level)
            fh_formatter = logging.Formatter(
                """
                <li>
                    <p>
                        <span class="time">%(asctime)s</span>
                        <span class="name %(name)s">%(name)-12s</span>
                        <span class="levelname %(levelname)s">%(levelname)-8s</span>
                    </p>
                    <p class="message">%(message)s</p>
                </li>
                """,
            )
            fh.setFormatter(fh_formatter)
            logging.getLogger("").addHandler(fh)

            return log_file
```


!! function <h2 id='configure' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>configure <em class='small'>function</em></h2>

???+ source "configure <em class='small'>source</em>"

```python

        def configure(markata: "Markata") -> None:
            setup_log(markata, logging.DEBUG)
            setup_log(markata, logging.INFO)
            setup_log(markata, logging.WARNING)

            if not has_rich_handler():
                console = RichHandler(
                    rich_tracebacks=True,
                )
                console.setLevel(logging.INFO)
                formatter = logging.Formatter("%(message)s")
                console.setFormatter(formatter)
                logging.getLogger("").addHandler(console)
```


!! method <h2 id='_fail_with_undefined_error' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_fail_with_undefined_error <em class='small'>method</em></h2>

???+ source "_fail_with_undefined_error <em class='small'>source</em>"

```python

        def _fail_with_undefined_error(self, *args, **kwargs):
                return ""
```


!! method <h2 id='validate_log_dir' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>validate_log_dir <em class='small'>method</em></h2>

???+ source "validate_log_dir <em class='small'>source</em>"

```python

        def validate_log_dir(cls, v, *, values):
                if v is None:
                    return values["output_dir"] / "_logs"
                return Path(v)
```
