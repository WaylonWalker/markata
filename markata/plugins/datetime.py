"""Default datetime plugin"""
import datetime
from pathlib import Path

import dateutil
import pytz
from feedgen.feed import FeedGenerator
from more_itertools import flatten
from tqdm import tqdm

from markata.hookspec import hook_impl


@hook_impl(trylast=True)
def load(markata):
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
