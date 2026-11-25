# Credential-bearing Dropper uploader transport security

Issue #2100.

## Trust boundary

The Markata editor must never receive a Dropper upstream credential. A
separate process exposes a constrained local upload capability through a Unix
socket. On POSIX targets the socket must have mode `0660`, with directory
access also scoped to trusted users/groups. The server refuses non-socket
stale paths and validates the fixed upstream endpoint and token.

## Windows behavior

Windows does not guarantee the same POSIX mode semantics for AF_UNIX socket
filesystem entries. A successful `os.Chmod(0660)` followed by `Stat`
reporting `0666` is not sufficient evidence of a secure ACL. Rather than
silently offering credential-bearing uploads to unintended local users,
`ServeDropperUploader` must return a clear unsupported-platform error on
Windows **before binding any socket**.

This restriction applies only to the credential-bearing uploader backend.
The local editor, Markdown authoring, builder, and any ordinary noncredential
features remain cross-platform. A future Windows implementation must use and
test Windows-native ACLs or an authenticated named-pipe transport before this
gate can be lifted.

## Tests

- On POSIX systems, assert actual Unix socket type and `0660` mode, then
  ensure clean server shutdown.
- On Windows, assert a clear fail-closed unsupported-platform diagnostic;
  do not skip the security contract or claim POSIX permission enforcement.
- Run cross-platform CI on Windows/macOS/Linux. Keep this fix independent of
  installer PR #2086.
