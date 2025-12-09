---
date: 2025-12-09
description: "Slugify migration for projects moving from markata<0.5.0 into markata>=0.5.0
  to run this script install markata>=0.5.0 and run the following. Then make sure\u2026"
published: false
slug: markata/scripts/migrate-to-slugify
title: migrate_to_slugify.py


---

---

Slugify migration for projects moving from markata<0.5.0 into markata>=0.5.0

to run this script install markata>=0.5.0 and run the following.

``` bash
python -m markata.scripts.migrate_to_slugify
```

Then make sure that you do not explicity turn off slugify and your site is
going to be on to better urls.

``` toml
[markata]
slugify=false
```

---