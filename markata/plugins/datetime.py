"""Default datetime plugin"""
import datetime
from pathlib import Path

import dateutil.parser
import pytz

from markata.hookspec import hook_impl

from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from markata import Markata


@hook_impl(trylast=True)
def load(markata: "Markata") -> None:
    for article in markata.iter_articles("datetime"):
        try:
            date = article.metadata["date"]
        except KeyError:
            date = "1970-01-01"
        if isinstance(date, str):
            date = dateutil.parser.parse(date)
        if isinstance(date, datetime.date):
            date = datetime.datetime(
                year=date.year,
                month=date.month,
                day=date.day,
                tzinfo=pytz.utc,
            )
        article["datetime"] = date
        article["date"] = date.date()
