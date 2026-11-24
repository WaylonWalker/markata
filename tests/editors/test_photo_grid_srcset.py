"""Regression contracts for responsive photo-grid rendering."""

from pathlib import Path
import re
import unittest


ROOT = Path(__file__).resolve().parents[2]
TEMPLATE = ROOT / "pkg/themes/default/templates/feed-photo-grid.html"


class PhotoGridResponsiveImageTests(unittest.TestCase):
    def test_both_pagination_paths_use_trusted_responsive_filter(self):
        template = TEMPLATE.read_text(encoding="utf-8")
        images = re.findall(
            r'<img\s+[^>]*class="shot-card-media u-photo"[^>]*>',
            template,
        )
        self.assertEqual(len(images), 2, "Expected both photo-grid render paths")
        for image in images:
            with self.subTest(image=image[:80]):
                self.assertIn('responsive_srcset:"240,400,600,800,1200"', image)
                self.assertIn("{% if shot_srcset %}srcset=", image)
                self.assertIn("{% endif %}", image)
                self.assertNotIn("with_size:'240' }} 240w", image)

    def test_large_cards_use_their_actual_display_width(self):
        template = TEMPLATE.read_text(encoding="utf-8")
        images = re.findall(
            r'<img\s+[^>]*class="shot-card-media u-photo"[^>]*>',
            template,
        )
        self.assertEqual(len(images), 2)
        for image in images:
            with self.subTest(image=image[:80]):
                self.assertIn("media_size == '600'", image)
                self.assertIn("(max-width: 640px) 92vw, 600px", image)
                self.assertIn("(max-width: 640px) 46vw, 280px", image)


if __name__ == "__main__":
    unittest.main()
