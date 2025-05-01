"""
The `markata.plugins.analytics` plugin generates analytics and contribution visualizations
for your Markata site. It creates a contributions heatmap similar to GitHub's contribution
graph and provides post statistics.

# Installation

This plugin is built-in and enabled by default through the 'default' plugin.
If you want to be explicit, you can add it to your list of plugins:

```toml
hooks = [
    "markata.plugins.analytics",
]
```

# Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

```toml
disabled_hooks = [
    "markata.plugins.analytics",
]
```

# Configuration

Configure analytics behavior in your `markata.toml`:

```toml
[markata.analytics]
# Maximum scale for contribution heatmap
contributions_max_post_scale = 5

# Color map for the heatmap (any matplotlib colormap)
contributions_cmap = "rocket"

# Filter to apply when generating analytics
filter = ""  # e.g. "draft != True" to exclude drafts
```

# Functionality

## Contribution Heatmap

Generates a heatmap visualization showing:
- Post frequency over time
- Relative activity levels using configurable color scales
- Year-over-year contribution patterns

## Post Statistics

Provides analytics including:
- Total number of posts
- Posts per month/year
- Word count statistics
- Reading time estimates

## Output Files

The plugin generates these files in your output directory:
- `contributions.svg`: The contribution heatmap
- `analytics.json`: Raw analytics data

## Dependencies

This plugin depends on:
- matplotlib (for heatmap generation)
- The `datetime` plugin for post date information
"""

from pathlib import Path
from typing import TYPE_CHECKING

import pydantic

from markata.hookspec import hook_impl
from markata.hookspec import register_attr

if TYPE_CHECKING:
    from markata import Markata


class AnalyticsConfig(pydantic.BaseModel):
    contributions_max_post_scale: int = 5
    contributions_cmap: str = "rocket"
    filter: str = ""


class Config(pydantic.BaseModel):
    analytics: AnalyticsConfig = AnalyticsConfig()


@hook_impl
@register_attr("post_models")
def config_model(markata: "Markata") -> None:
    markata.config_models.append(Config)


@hook_impl
def pre_render(markata: "Markata") -> None:
    output_dir = Path(markata.config.output_dir)
    output_dir.mkdir(exist_ok=True, parents=True)

    post_dates = markata.map("date", filter=markata.config.analytics.filter)
    post_dates.sort()

    analytics_key = markata.make_hash("analytics", markata.config.analytics, post_dates)
    with markata.cache as cache:
        contributions = cache.get(analytics_key)
    if (
        contributions is not None
        and (output_dir / "total_posts_over_time.png").exists()
    ):
        print("Contribution graph already rendered, skipping analytics render")
        # if conributions have not changed, don't render again
        return

    # Import matplotlib first to set backend
    import matplotlib

    matplotlib.use("Agg")  # Use non-interactive backend

    import pandas as pd
    import seaborn as sns
    from matplotlib import pyplot as plt

    # Contribution Graph as Heatmap
    df = pd.DataFrame(post_dates, columns=["date"])

    df["date"] = pd.to_datetime(df["date"])

    df["year"] = df["date"].dt.year
    df["week"] = df["date"].dt.isocalendar().week
    df["day_of_week"] = df["date"].dt.weekday

    all_weeks = range(1, 53)
    all_days = range(0, 7)

    day_labels = ["", "MON", "", "WED", "", "FRI", ""]
    month_labels = [
        "",
        "Jan",
        "Feb",
        "Mar",
        "Apr",
        "May",
        "Jun",
        "Jul",
        "Aug",
        "Sep",
        "Oct",
        "Nov",
        "Dec",
    ]

    # Determine global vmax for all years
    overall_max_value = df.groupby(["day_of_week", "week"]).size().max()
    vmax = max(5, overall_max_value // 2)
    vmax = min(markata.config.analytics.contributions_max_post_scale, vmax)

    plt.style.use("dark_background")

    for year in df["year"].unique():
        first_day_of_year = pd.Timestamp(f"{year}-01-01").weekday()
        start_date = pd.Timestamp(f"{year}-01-01") - pd.Timedelta(
            days=first_day_of_year + 1
        )
        padded_data = pd.DataFrame(
            {
                "date": pd.date_range(
                    start=start_date, periods=first_day_of_year + 1, freq="D"
                )
            }
        )
        padded_data["year"] = padded_data["date"].dt.year
        padded_data["week"] = padded_data["date"].dt.isocalendar().week
        padded_data["day_of_week"] = padded_data["date"].dt.weekday
        df = pd.concat([padded_data, df])

        fig, ax = plt.subplots(figsize=(15, 3))
        yearly_data = df[df["year"] == year]
        heatmap_data = (
            yearly_data.groupby(["day_of_week", "week"]).size().unstack(fill_value=0)
        )
        heatmap_data = heatmap_data.reindex(
            index=all_days, columns=all_weeks, fill_value=0
        )

        sns.heatmap(
            heatmap_data,
            cmap=markata.config.analytics.contributions_cmap,
            linewidths=3,
            linecolor="black",
            square=True,
            vmax=vmax,
            cbar=False,
            ax=ax,
        )
        ax.set_yticklabels(day_labels, rotation=0, color="white")
        ax.set_xticks([i * 4 for i in range(13)])
        ax.set_xticklabels(month_labels, rotation=0, color="white")
        ax.set_xlabel("")
        ax.set_ylabel("")
        ax.tick_params(left=False, bottom=False)

        plt.savefig(
            output_dir / f"contributions_{year}.png",
            facecolor="black",
            bbox_inches="tight",
            pad_inches=0.1,
        )
        plt.savefig(
            output_dir / f"contributions_{year}.svg",
            facecolor="black",
            bbox_inches="tight",
            pad_inches=0.1,
        )
        plt.close()

    # Total Posts Over Time
    cumulative_posts = df.groupby("date").size().cumsum()
    plt.figure(figsize=(10, 5))
    plt.plot(cumulative_posts.index, cumulative_posts.values, color="#c4204f")
    plt.text(
        cumulative_posts.index[-1],
        cumulative_posts.values[-1],
        f"{cumulative_posts.values[-1]} posts",
        color="white",
        fontsize=12,
        ha="right",
    )
    for year in df["year"].unique():
        year_start = pd.Timestamp(f"{year}-01-01")
        plt.axvline(x=year_start, color="#130c24", linestyle="--")
        plt.text(
            year_start,
            -cumulative_posts.values.max() * 0.1,
            str(year),
            color="white",
            fontsize=10,
            ha="center",
        )
    plt.title("Total Number of Posts Over Time", color="white")
    plt.axis("off")
    plt.savefig(
        output_dir / "total_posts_over_time.png",
        facecolor="black",
        bbox_inches="tight",
        pad_inches=0.1,
    )
    plt.savefig(
        output_dir / "total_posts_over_time.svg",
        facecolor="black",
        bbox_inches="tight",
        pad_inches=0.1,
    )
    plt.close()

    with markata.cache as cache:
        cache.set(analytics_key, "done")
