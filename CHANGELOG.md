# Markata Changelog

## 0.2.0

* feat: auto_description plugin is all new
* deprecated: long_description has been deprecated by auto_description
* fix: covers plugin which would previously skip every time.
* feat: `markata clean` cleans up your cache and output from the command line
* fix: `publish_source` plugin will now ignore any non yaml serializable values 

### All New auto_description plugin

* Cache is busted on plugin change 
* plugin is configurable
* plugin now has docs

## Configuration

Open up your `markata.toml` file and add new entries for your
auto_descriptions.  You can have multiple desriptions, each one will be named
after the key you give it in your config.

``` toml
[markata]
hooks=[
   "markata.plugins.auto_description",
   ]

[markata.auto_description.description]
len=160
[markata.auto_description.long_description]
len=250
[markata.auto_description.super_description]
len=500
```

In the above we will end up with three different descritpions, 
(`description`, `long_description`, and `super_description`) each will be the
first number of characters from the document as specified in the config.

## 0.1.0

* fix: pyinstument plugin no longer overrides the main cli callback
* feat: default is to run the profiler if pyinstrument is installed
* fix: --profile is now under the build command
* feat: --pretty/--no-pretty will make pretty tracebacks, and skip over framework code (closes #4 )
* fix: links are now absolute so they work from github, thanks [MR Destructive](https://github.com/Mr-Destructive)
* deprecate: `article['content_hash']` has been removed with preference for simply hashing on `article['content']`
* create `heading_link` plugin
* doc: How to create your home page. https://markata.dev/home-page/

### New cli help

After the pyinstrument plugin was fixed --version and  --to-json are back, and --profile is now under the build command.

![image](https://user-images.githubusercontent.com/22648375/150662983-547aebbd-c18c-4c17-8985-a6dc01cd29c7.png)

### New Heading Link Plugin

The new heading link plugin makes it easier to share the exact part of an article you want with someone, by giving clickable links to the id of the heading it's next to.

![image](https://user-images.githubusercontent.com/22648375/151718782-08a7cd26-41c1-4f00-a12c-0a208c593e9c.png)

## 0.0.1

Initial Release 🎉

