# Plain post CSS loading

On stock `post.html`, card-free, feed-free articles can avoid an initial
request for feed-card CSS. Generator-side detection protects embedded cards,
feeds, and custom templates. Do not optimize by hiding or limiting posts.

Use a clean build and compare generated stylesheet links in article vs feed
pages; verify transition navigation both ways, mobile LCP and layout shifts.
The narrow optimization is not a universal critical CSS extractor.
