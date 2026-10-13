# Feed view contract

## Built-in views

The canonical built-in feed presentation identifiers are `default`, `simple`, and `calendar`.
When no view list is configured, all three are available. `default` is mandatory; `simple` and
`calendar` may be removed globally through `feed_defaults.views` or overridden per feed with
`feeds[].views`.

## Rendering and generation

View controls MUST only link to enabled presentations. Calendar assets and its hidden source
payload MUST NOT be emitted when `calendar` is disabled. Simple HTML MUST NOT be generated when
`simple` is disabled, even if the `simple_html` output-format flag is otherwise enabled. Feed
view availability participates in incremental output hashing so configuration changes rebuild
affected feeds.

## Validation

Explicit view lists MUST contain `default`. Unknown view identifiers MUST fail validation with
the supported identifiers in the diagnostic. Omitted lists inherit; explicitly empty lists are
invalid rather than meaning "inherit".

Config layering MUST preserve an explicit empty view list so validation can
reject it; merging MUST NOT turn that list into omitted/inherited views.
Disabling Simple MUST remove stale generated `simple/` output on rebuild.

## Shorts viewer controls

The Shorts viewer MUST keep previous and next controls at a stable viewport
position as the active feed item changes between images and videos. Playback
controls MUST be positioned independently from navigation controls so showing
video-only actions does not move the previous and next controls.

The viewer MUST support `j`/`k` and the up/down arrows for next/previous,
Escape to return to the source feed grid, Space to play or pause an active
video, left/right arrows to seek five seconds, and `m` to mute or unmute.
Keyboard shortcuts MUST NOT intercept typing in form fields or editable
content.
