"""Default datetime plugin"""
import datetime
from typing import TYPE_CHECKING

import dateutil
import pytz

from markata.hookspec import hook_impl

if TYPE_CHECKING:
    from markata import Markata


@hook_impl(trylast=True)
def load(markata: "Markata") -> None:
    for article in markata.iter_articles("datetime"):
        try:
            date = article.metadata["date"]
        except KeyError:
            date = None
        if isinstance(date, str):
            date = dateutil.parser.parse(date)
        if isinstance(date, datetime.date):
            date = datetime.datetime(
                year=date.year,
                month=date.month,
                day=date.day,
                tzinfo=pytz.utc,
            )

        article["today"] = datetime.date.today()
        article["now"] = datetime.datetime.now()
        article["datetime"] = date
        # if date is not None:
        #     article["date"] = date.date()
