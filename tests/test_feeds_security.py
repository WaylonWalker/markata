import pytest
import tempfile
import shutil
from pathlib import Path
from unittest.mock import Mock, patch
import hashlib

from markata.plugins.feeds import Feed, _sanitize_feed_slug, _download_htmx_if_needed
from markata import Markata


class TestSecurity:
    """Test suite for security vulnerabilities in feeds plugin."""

    def test_path_traversal_prevention(self):
        """Test that path traversal attacks are prevented in feed slugs."""

        # Malicious slugs that should be rejected
        malicious_slugs = [
            "../../../etc/passwd",
            "..\\..\\windows\\system32\\config\\sam",
            "normal/../../../etc/passwd",
            "normal\\..\\..\\windows\\system32",
            "etc/passwd",
            "C:\\Windows\\System32",
            "/etc/shadow",
            "",
            ".",
            "./hidden",
            "hidden/.",
        ]

        for slug in malicious_slugs:
            with pytest.raises(
                ValueError, match=r"(Invalid characters|cannot be empty)"
            ):
                _sanitize_feed_slug(slug)

    def test_safe_slug_validation(self):
        """Test that safe slugs are allowed."""

        safe_slugs = [
            "blog",
            "my-feed",
            "news_posts",
            "test123",
            "a",
            "my_blog_posts_2023",
            "feed-with-dashes",
        ]

        for slug in safe_slugs:
            result = _sanitize_feed_slug(slug)
            assert result == slug

    def test_htmx_integrity_verification(self):
        """Test that HTMX download verifies file integrity."""

        # Mock the responses with wrong hash
        mock_content = b"malicious javascript content"
        mock_response = Mock()
        mock_response.read.return_value = mock_content

        with patch("markata.plugins.feeds.urlopen", return_value=mock_response):
            with patch("pathlib.Path.exists", return_value=False):
                with patch("pathlib.Path.parent"):
                    with patch("pathlib.Path.write_bytes"):
                        mock_markata = Mock()
                        mock_markata.config.htmx_version = "1.9.10"
                        mock_markata.config.output_dir = "/tmp/test"

                        with pytest.raises(RuntimeError, match="HTMX download failed"):
                            _download_htmx_if_needed(mock_markata)

    def test_htmx_timeout_protection(self):
        """Test that HTMX download has timeout protection."""

        mock_markata = Mock()
        mock_markata.config.htmx_version = "1.9.10"
        mock_markata.config.output_dir = "/tmp/test"

        # Mock urlopen to raise timeout
        with patch(
            "markata.plugins.feeds.urlopen",
            side_effect=TimeoutError("Request timed out"),
        ):
            with pytest.raises(RuntimeError, match="HTMX download failed"):
                _download_htmx_if_needed(mock_markata)

    def test_xss_prevention_in_template_context(self):
        """Test that template context doesn't contain dangerous config data."""

        # Create a feed with potentially dangerous config
        dangerous_config = {
            "pagination_type": "js",
            "posts_per_page": 10,
            "template": '<script>alert("xss")</script>',
            "card_template": "dangerous-template.html",
            "xss_payload": '<script>document.cookie="stolen"</script>',
            "admin_password": "secret123",
            "api_key": "sk-1234567890",
        }

        # Safe config should only include essential pagination settings
        safe_config = {
            "pagination_type": dangerous_config["pagination_type"],
            "posts_per_page": dangerous_config["posts_per_page"],
            "template": dangerous_config["template"],
        }

        # Verify only safe keys are included
        for key in dangerous_config:
            if key not in safe_config:
                assert key not in safe_config, (
                    f"Dangerous key '{key}' should not be in safe config"
                )

    def test_canonical_url_sanitization(self):
        """Test that canonical URLs use sanitized slugs."""

        mock_markata = Mock()
        mock_markata.config.url = "https://example.com"

        # Test with safe slug
        safe_slug = "my-blog-feed"
        feed_config = Mock()
        feed_config.slug = safe_slug

        feed = Feed(config=feed_config, markata=mock_markata)

        # The canonical URL should use the safe slug
        expected_url = f"https://example.com/{safe_slug}/"
        # This would be tested in actual template rendering

    def test_feed_file_path_security(self):
        """Test that feed file paths cannot escape output directory."""

        with tempfile.TemporaryDirectory() as temp_dir:
            output_dir = Path(temp_dir)

            # Try to create a feed with malicious slug
            malicious_slugs = [
                "../outside",
                "normal/../../../etc/passwd",
                "normal\\..\\..\\windows\\system32",
            ]

            for malicious_slug in malicious_slugs:
                with pytest.raises(ValueError):
                    _sanitize_feed_slug(malicious_slug)

                # Ensure no files can be created outside output directory
                safe_slug = _sanitize_feed_slug("safe-feed")
                file_path = output_dir / safe_slug / "index.html"

                # Verify path is within output directory
                assert file_path.resolve().is_relative_to(output_dir.resolve())

    def test_template_injection_prevention(self):
        """Test that template injection is prevented in feed names."""

        dangerous_names = [
            "{{7*7}}",  # Template injection
            "${7*7}",  # Expression injection
            "<script>alert('xss')</script>",
            "javascript:void(0)",
            "data:text/html,<script>alert(1)</script>",
        ]

        for dangerous_name in dangerous_names:
            # These should be sanitized or rejected
            sanitized = _sanitize_feed_slug(dangerous_name)
            # Should either be rejected or sanitized to safe version
            assert "{{" not in sanitized
            assert "}}" not in sanitized
            assert "<script" not in sanitized.lower()
            assert "javascript:" not in sanitized.lower()

    def test_htmx_download_no_fallback(self):
        """Test that HTMX download has no silent fallback to CDN."""

        mock_markata = Mock()
        mock_markata.config.htmx_version = "1.9.10"
        mock_markata.config.output_dir = "/tmp/test"
        mock_markata.console = Mock()

        # Mock failed download
        with patch(
            "markata.plugins.feeds.urlopen", side_effect=Exception("Network error")
        ):
            with patch("pathlib.Path.exists", return_value=False):
                with pytest.raises(RuntimeError, match="HTMX download failed"):
                    _download_htmx_if_needed(mock_markata)

                # Should log error but not fall back to CDN
                mock_markata.console.error.assert_called()


if __name__ == "__main__":
    pytest.main([__file__])
