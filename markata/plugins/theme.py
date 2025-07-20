from typing import Dict

from colour import Color as BaseColor
from pydantic import GetCoreSchemaHandler
from pydantic_core import core_schema
from rich.console import Console
from rich.console import ConsoleOptions
from rich.console import RenderResult
from rich.repr import RichReprResult
from rich.text import Text


def print_color_swatch(name: str, hex_color: str):
    text_color = "#ffffff"
    hex_color = str(hex_color)
    if Color(hex_color).get_luminance() > 0.5:
        text_color = "#000000"

    text = Text(f" {name.ljust(46)} ", style=f"on {hex_color} {text_color}")
    console = Console()
    console.print(text, justify="left", end=" ")


class Color(BaseColor):
    def __init__(self, value=None, **kwargs):
        if value is not None:
            # Handle Tailwind-style like "blue-500"
            if isinstance(value, str):
                if "-" in value:
                    name, index = value.split("-")
                    value = tailwind_v4_colors[name][int(index)]
                elif value in tailwind_v4_colors:
                    val = tailwind_v4_colors[value]
                    if isinstance(val, dict):
                        value = val.get(500, list(val.values())[0])  # fallback
                    else:
                        value = val
            super().__init__(value)
        else:
            # Support keyword-based construction like hex=, hsl=, rgb=
            super().__init__(**kwargs)

    @classmethod
    def __get_pydantic_core_schema__(
        cls, source_type, handler: GetCoreSchemaHandler
    ) -> core_schema.CoreSchema:
        # Parse from a string and return a Color
        return core_schema.no_info_plain_validator_function(
            cls._validate_color,
            serialization=core_schema.plain_serializer_function_ser_schema(str),
        )

    @classmethod
    def _validate_color(cls, value):
        if isinstance(value, cls):
            return value
        if isinstance(value, BaseColor):
            return cls(value.hex_l)
        if isinstance(value, str):
            return cls(value)
        raise TypeError(f"Cannot convert {value!r} to {cls.__name__}")

    def __str__(self):
        return self.hex_l

    def __repr__(self):
        return self.hex_l

    def __rich_repr__(self) -> RichReprResult:
        text_color = "#000000" if self.get_luminance() > 0.5 else "#ffffff"
        swatch = Text(f" {self.hex.ljust(46)} ", style=f"on {self.hex_l} {text_color}")
        yield "swatch", swatch

    def __rich_console__(
        self, console: Console, options: ConsoleOptions
    ) -> RenderResult:
        text_color = "#000000" if self.get_luminance() > 0.5 else "#ffffff"
        swatch = Text(f" {self.hex.ljust(46)} ", style=f"on {self.hex_l} {text_color}")
        yield swatch

    def __rich__(self):
        text_color = "#000000" if self.get_luminance() > 0.5 else "#ffffff"
        return Text(f" {self.hex} ", style=f"on {self.hex_l} {text_color}")

    def __add__(self, other):
        return self.mix_hsl(other)

    def __radd__(self, other):
        return self.mix_hsl(other)

    def __sub__(self, other):
        return self.mix_hsl(other, ratio=0.5)

    def __rsub__(self, other):
        return self.mix_hsl(other, ratio=0.5)

    def __mul__(self, other):
        return self.mix_hsl(other)

    def mix(self, other, ratio=0.5):
        base = self.get_hsl()
        overlay = other.get_hsl()

        h = (1 - ratio) * base[0] + ratio * overlay[0]
        s = (1 - ratio) * base[1] + ratio * overlay[1]
        l = (1 - ratio) * base[2] + ratio * overlay[2]
        return Color(hsl=(h, s, l))  # .hex_l

    def blend(self, other, ratio=0.5):
        return self.mix(other, ratio)

    def _rgb_math(self, other, op):
        a = self.get_rgb()
        b = other.get_rgb()
        result = [op(x, y) for x, y in zip(a, b)]
        # Clamp between 0 and 1
        clamped = [min(max(c, 0), 1) for c in result]
        return Color(rgb=tuple(clamped))

    def _rgb_scalar_math(self, scalar, op):
        a = self.get_rgb()
        result = [op(x, scalar) for x in a]
        clamped = [min(max(c, 0), 1) for c in result]
        return Color(rgb=tuple(clamped))

    def __add__(self, other):
        return self._rgb_math(other, lambda a, b: a + b)

    def __sub__(self, other):
        return self._rgb_math(other, lambda a, b: a - b)

    def __mul__(self, other):
        if isinstance(other, Color):
            return self._rgb_math(other, lambda a, b: a * b)
        else:
            return self._rgb_scalar_math(other, lambda a, b: a * b)

    def __truediv__(self, other):
        if isinstance(other, Color):
            return self._rgb_math(other, lambda a, b: a / b if b != 0 else 0)
        else:
            return self._rgb_scalar_math(other, lambda a, b: a / b if b != 0 else 0)

    def __pow__(self, exponent):
        return self._rgb_scalar_math(exponent, lambda a, b: a**b)

    def print_color_swatch(self):
        text_color = "#ffffff"
        if self.get_luminance() > 0.5:
            text_color = "#000000"

        text = Text(f" {self.hex.ljust(46)} ", style=f"on {self.hex_l} {text_color}")
        console = Console()
        console.print(text, justify="left", end=" ")

    def compliment(self):
        h, s, l = self.get_hsl()
        return Color(hsl=((h + 0.5) % 1, s, l))

    def lighten(self, amount=0.1):
        h, s, l = self.get_hsl()
        return Color(hsl=(h, s, clamp(l + amount)))

    def darken(self, amount=0.1):
        h, s, l = self.get_hsl()

        return Color(hsl=(h, s, clamp(l - amount)))


def clamp(value, min_value=0, max_value=1):
    return max(min(value, max_value), min_value)


tailwind_v4_colors = {
    "black": "#000000",
    "white": "#ffffff",
    "transparent": "transparent",
    "slate": {
        50: "#f8fafc",
        100: "#f1f5f9",
        200: "#e2e8f0",
        300: "#cbd5e1",
        400: "#94a3b8",
        500: "#64748b",
        600: "#475569",
        700: "#334155",
        800: "#1e293b",
        900: "#0f172a",
        950: "#020617",
    },
    "gray": {
        50: "#f9fafb",
        100: "#f3f4f6",
        200: "#e5e7eb",
        300: "#d1d5db",
        400: "#9ca3af",
        500: "#6b7280",
        600: "#4b5563",
        700: "#374151",
        800: "#1f2937",
        900: "#111827",
        950: "#030712",
    },
    "zinc": {
        50: "#fafafa",
        100: "#f4f4f5",
        200: "#e4e4e7",
        300: "#d4d4d8",
        400: "#a1a1aa",
        500: "#71717a",
        600: "#52525b",
        700: "#3f3f46",
        800: "#27272a",
        900: "#18181b",
        950: "#09090b",
    },
    "neutral": {
        50: "#fafafa",
        100: "#f5f5f5",
        200: "#e5e5e5",
        300: "#d4d4d4",
        400: "#a3a3a3",
        500: "#737373",
        600: "#525252",
        700: "#404040",
        800: "#262626",
        900: "#171717",
        950: "#0a0a0a",
    },
    "stone": {
        50: "#fafaf9",
        100: "#f5f5f4",
        200: "#e7e5e4",
        300: "#d6d3d1",
        400: "#a8a29e",
        500: "#78716c",
        600: "#57534e",
        700: "#44403c",
        800: "#292524",
        900: "#1c1917",
        950: "#0c0a09",
    },
    "red": {
        50: "#fef2f2",
        100: "#fee2e2",
        200: "#fecaca",
        300: "#fca5a5",
        400: "#f87171",
        500: "#ef4444",
        600: "#dc2626",
        700: "#b91c1c",
        800: "#991b1b",
        900: "#7f1d1d",
        950: "#450a0a",
    },
    "orange": {
        50: "#fff7ed",
        100: "#ffedd5",
        200: "#fed7aa",
        300: "#fdba74",
        400: "#fb923c",
        500: "#f97316",
        600: "#ea580c",
        700: "#c2410c",
        800: "#9a3412",
        900: "#7c2d12",
        950: "#431407",
    },
    "amber": {
        50: "#fffbeb",
        100: "#fef3c7",
        200: "#fde68a",
        300: "#fcd34d",
        400: "#fbbf24",
        500: "#f59e0b",
        600: "#d97706",
        700: "#b45309",
        800: "#92400e",
        900: "#78350f",
        950: "#451a03",
    },
    "yellow": {
        50: "#fefce8",
        100: "#fef9c3",
        200: "#fef08a",
        300: "#fde047",
        400: "#facc15",
        500: "#eab308",
        600: "#ca8a04",
        700: "#a16207",
        800: "#854d0e",
        900: "#713f12",
        950: "#422006",
    },
    "lime": {
        50: "#f7fee7",
        100: "#ecfccb",
        200: "#d9f99d",
        300: "#bef264",
        400: "#a3e635",
        500: "#84cc16",
        600: "#65a30d",
        700: "#4d7c0f",
        800: "#3f6212",
        900: "#365314",
        950: "#1a2e05",
    },
    "green": {
        50: "#f0fdf4",
        100: "#dcfce7",
        200: "#bbf7d0",
        300: "#86efac",
        400: "#4ade80",
        500: "#22c55e",
        600: "#16a34a",
        700: "#15803d",
        800: "#166534",
        900: "#14532d",
        950: "#052e16",
    },
    "emerald": {
        50: "#ecfdf5",
        100: "#d1fae5",
        200: "#a7f3d0",
        300: "#6ee7b7",
        400: "#34d399",
        500: "#10b981",
        600: "#059669",
        700: "#047857",
        800: "#065f46",
        900: "#064e3b",
        950: "#022c22",
    },
    "teal": {
        50: "#f0fdfa",
        100: "#ccfbf1",
        200: "#99f6e4",
        300: "#5eead4",
        400: "#2dd4bf",
        500: "#14b8a6",
        600: "#0d9488",
        700: "#0f766e",
        800: "#115e59",
        900: "#134e4a",
        950: "#042f2e",
    },
    "cyan": {
        50: "#ecfeff",
        100: "#cffafe",
        200: "#a5f3fc",
        300: "#67e8f9",
        400: "#22d3ee",
        500: "#06b6d4",
        600: "#0891b2",
        700: "#0e7490",
        800: "#155e75",
        900: "#164e63",
        950: "#083344",
    },
    "sky": {
        50: "#f0f9ff",
        100: "#e0f2fe",
        200: "#bae6fd",
        300: "#7dd3fc",
        400: "#38bdf8",
        500: "#0ea5e9",
        600: "#0284c7",
        700: "#0369a1",
        800: "#075985",
        900: "#0c4a6e",
        950: "#082f49",
    },
    "blue": {
        50: "#eff6ff",
        100: "#dbeafe",
        200: "#bfdbfe",
        300: "#93c5fd",
        400: "#60a5fa",
        500: "#3b82f6",
        600: "#2563eb",
        700: "#1d4ed8",
        800: "#1e40af",
        900: "#1e3a8a",
        950: "#172554",
    },
    "indigo": {
        50: "#eef2ff",
        100: "#e0e7ff",
        200: "#c7d2fe",
        300: "#a5b4fc",
        400: "#818cf8",
        500: "#6366f1",
        600: "#4f46e5",
        700: "#4338ca",
        800: "#3730a3",
        900: "#312e81",
        950: "#1e1b4b",
    },
    "violet": {
        50: "#f5f3ff",
        100: "#ede9fe",
        200: "#ddd6fe",
        300: "#c4b5fd",
        400: "#a78bfa",
        500: "#8b5cf6",
        600: "#7c3aed",
        700: "#6d28d9",
        800: "#5b21b6",
        900: "#4c1d95",
        950: "#2e1065",
    },
    "purple": {
        50: "#faf5ff",
        100: "#f3e8ff",
        200: "#e9d5ff",
        300: "#d8b4fe",
        400: "#c084fc",
        500: "#a855f7",
        600: "#9333ea",
        700: "#7e22ce",
        800: "#6b21a8",
        900: "#581c87",
        950: "#3b0764",
    },
    "fuchsia": {
        50: "#fdf4ff",
        100: "#fae8ff",
        200: "#f5d0fe",
        300: "#f0abfc",
        400: "#e879f9",
        500: "#d946ef",
        600: "#c026d3",
        700: "#a21caf",
        800: "#86198f",
        900: "#701a75",
        950: "#4a044e",
    },
    "pink": {
        50: "#fdf2f8",
        100: "#fce7f3",
        200: "#fbcfe8",
        300: "#f9a8d4",
        400: "#f472b6",
        500: "#ec4899",
        600: "#db2777",
        700: "#be185d",
        800: "#9d174d",
        900: "#831843",
        950: "#500724",
    },
    "rose": {
        50: "#fff1f2",
        100: "#ffe4e6",
        200: "#fecdd3",
        300: "#fda4af",
        400: "#fb7185",
        500: "#f43f5e",
        600: "#e11d48",
        700: "#be123c",
        800: "#9f1239",
        900: "#881337",
        950: "#4c0519",
    },
}

THEME_DEFAULTS: Dict[str, Dict[str, str]] = {
    "tokyo-night": {
        "light": {
            "text": "gray-900",
            "muted": "gray-500",
            "heading": "black",
            "accent": "indigo-600",
            "accent_alt": "purple-600",
            "background": "white",
            "surface": "gray-50",
            "code_bg": "gray-100",
            "blockquote_bg": "gray-100",
            "blockquote_border": "indigo-300",
            "link_hover": "black",
            "selection_bg": "indigo-100",
            "selection_text": "gray-900",
            "border": "gray-200",
        },
        "dark": {
            "text": "gray-100",
            "muted": "gray-400",
            "heading": "white",
            "accent": "indigo-400",
            "accent_alt": "purple-400",
            "background": "#1a1b26",
            "surface": "#222436",
            "code_bg": "#2f3549",
            "blockquote_bg": "#1f2335",
            "blockquote_border": "indigo-500",
            "link_hover": "white",
            "selection_bg": "#2f3549",
            "selection_text": "white",
            "border": "#3b4261",
        },
    },
    "catppuccin": {
        "light": {
            "text": "rose-900",
            "muted": "rose-500",
            "heading": "rose-800",
            "accent": "pink-500",
            "accent_alt": "purple-400",
            "background": "rose-50",
            "surface": "rose-100",
            "code_bg": "rose-100",
            "blockquote_bg": "rose-200",
            "blockquote_border": "pink-400",
            "link_hover": "pink-800",
            "selection_bg": "rose-300",
            "selection_text": "rose-900",
            "border": "rose-300",
            "code_theme": "stata-light",
        },
        "dark": {
            "text": "rose-200",
            "muted": "rose-400",
            "heading": "rose-100",
            "accent": "pink-400",
            "accent_alt": "violet-300",
            "background": "#1e1e28",
            "surface": "#2a2a38",
            "code_bg": "#2c2c3a",
            "blockquote_bg": "#2b2b3a",
            "blockquote_border": "pink-500",
            "link_hover": "white",
            "selection_bg": "#403d52",
            "selection_text": "rose-50",
            "border": "#4e4e5a",
            "code_theme": "dracula",
        },
    },
    "everforest": {
        "light": {
            "text": "green-900",
            "muted": "green-500",
            "heading": "green-800",
            "accent": "green-600",
            "accent_alt": "lime-500",
            "background": "green-50",
            "surface": "green-100",
            "code_bg": "green-100",
            "blockquote_bg": "green-200",
            "blockquote_border": "green-400",
            "link_hover": "green-800",
            "selection_bg": "green-200",
            "selection_text": "green-900",
            "border": "green-300",
            "code_theme": "stata-light",
        },
        "dark": {
            "text": "green-100",
            "muted": "green-400",
            "heading": "green-300",
            "accent": "green-400",
            "accent_alt": "lime-400",
            "background": "#2b3339",
            "surface": "#374045",
            "code_bg": "#3b444a",
            "blockquote_bg": "#3d484f",
            "blockquote_border": "green-500",
            "link_hover": "white",
            "selection_bg": "#475258",
            "selection_text": "white",
            "border": "#517d90",
            "code_theme": "stata-dark",
        },
    },
    "gruvbox": {
        "light": {
            "text": "orange-900",
            "muted": "orange-400",
            "heading": "yellow-900",
            "accent": "orange-600",
            "accent_alt": "yellow-500",
            "background": "white",
            "surface": "orange-50",
            "code_bg": "orange-100",
            "blockquote_bg": "orange-200",
            "blockquote_border": "orange-300",
            "link_hover": "orange-800",
            "selection_bg": "orange-200",
            "selection_text": "orange-900",
            "border": "orange-300",
        },
        "dark": {
            "text": "orange-100",
            "muted": "orange-400",
            "heading": "yellow-100",
            "accent": "orange-400",
            "accent_alt": "yellow-400",
            "background": "#282828",
            "surface": "#3c3836",
            "code_bg": "#504945",
            "blockquote_bg": "#3a3634",
            "blockquote_border": "orange-500",
            "link_hover": "white",
            "selection_bg": "#665c54",
            "selection_text": "orange-50",
            "border": "#7c6f64",
        },
    },
    "kanagwa": {
        "light": {
            "text": "slate-900",
            "muted": "slate-400",
            "heading": "slate-800",
            "accent": "blue-600",
            "accent_alt": "indigo-500",
            "background": "slate-50",
            "surface": "slate-100",
            "code_bg": "slate-100",
            "blockquote_bg": "slate-200",
            "blockquote_border": "blue-300",
            "link_hover": "blue-800",
            "selection_bg": "blue-100",
            "selection_text": "slate-900",
            "border": "slate-300",
        },
        "dark": {
            "text": "slate-100",
            "muted": "slate-400",
            "heading": "slate-50",
            "accent": "blue-400",
            "accent_alt": "indigo-400",
            "background": "#1f2335",
            "surface": "#2a2e3e",
            "code_bg": "#3a3f52",
            "blockquote_bg": "#2e3440",
            "blockquote_border": "blue-500",
            "link_hover": "white",
            "selection_bg": "#394260",
            "selection_text": "white",
            "border": "#4b5162",
        },
    },
    "nord": {
        "light": {
            "text": "cyan-900",
            "muted": "cyan-400",
            "heading": "cyan-800",
            "accent": "cyan-600",
            "accent_alt": "blue-500",
            "background": "cyan-200",
            "surface": "cyan-100",
            "code_bg": "cyan-50",
            "blockquote_bg": "cyan-200",
            "blockquote_border": "cyan-300",
            "link_hover": "cyan-800",
            "selection_bg": "cyan-200",
            "selection_text": "cyan-900",
            "border": "cyan-300",
            "code_theme": "solarized-light",
        },
        "dark": {
            "text": "cyan-100",
            "muted": "cyan-400",
            "heading": "cyan-50",
            "accent": "cyan-400",
            "accent_alt": "blue-300",
            "background": "#2e3440",
            "surface": "#3b4252",
            "code_bg": "#434c5e",
            "blockquote_bg": "#4c566a",
            "blockquote_border": "cyan-500",
            "link_hover": "white",
            "selection_bg": "#5e81ac",
            "selection_text": "cyan-50",
            "border": "#6b7d97",
            "code_theme": "nord-darker",
        },
    },
    "synthwave-84": {
        "light": {
            "text": "purple-900",
            "muted": "pink-500",
            "heading": "fuchsia-800",
            "accent": "pink-500",
            "accent_alt": "fuchsia-500",
            "background": "pink-50",
            "surface": "pink-100",
            "code_bg": "pink-100",
            "blockquote_bg": "pink-200",
            "blockquote_border": "pink-400",
            "link_hover": "purple-800",
            "selection_bg": "fuchsia-200",
            "selection_text": "purple-900",
            "border": "pink-300",
            "code_theme": "monokai",
        },
        "dark": {
            "text": "#ff00ff",
            "muted": "#c060c0",
            "heading": "#ff66ff",
            "accent": "pink-400",
            "accent_alt": "fuchsia-400",
            "background": "#2d0036",
            "surface": "#440055",
            "code_bg": "#3d0047",
            "blockquote_bg": "#520066",
            "blockquote_border": "pink-500",
            "link_hover": "white",
            "selection_bg": "#8800aa",
            "selection_text": "#ffffff",
            "border": "#ff00ff",
            "code_theme": "monokai",
        },
    },
}
