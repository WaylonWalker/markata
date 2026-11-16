"""Contract tests for deterministic, complete site Go test sharding."""

import unittest

from site_test_shards import discover_names, run_regex, select_shard


class SiteTestShardsTest(unittest.TestCase):
    def test_discover_runnable_names_from_go_list(self):
        listing = """TestZeta
TestAlpha
ExampleUsage
FuzzFrontmatter
BenchmarkBuild
ok  github.com/WaylonWalker/markata/tests 0.012s
"""
        self.assertEqual(
            discover_names(listing),
            ["ExampleUsage", "FuzzFrontmatter", "TestAlpha", "TestZeta"],
        )

    def test_split_is_complete_unique_and_deterministic(self):
        names = [f"TestCase{index:02d}" for index in range(43)]
        shards = [select_shard(names, index, 3) for index in range(3)]
        self.assertEqual(sorted(name for shard in shards for name in shard), sorted(names))
        self.assertEqual(max(map(len, shards)) - min(map(len, shards)), 1)
        self.assertEqual(shards[1], select_shard(list(reversed(names)), 1, 3))

    def test_regex_matches_only_selected_top_level_tests(self):
        names = ["TestAlpha", "ExampleBeta", "FuzzGamma"]
        regex = run_regex(names)
        import re

        for name in names:
            self.assertRegex(name, regex)
        for name in ["TestAlphabeta", "TestOther", "BenchmarkBuild", ""]:
            self.assertNotRegex(name, regex)

    def test_fail_closed(self):
        for listing in ("", "ok  github.com/WaylonWalker/markata/tests 0.0s\n"):
            with self.subTest(listing=listing), self.assertRaisesRegex(ValueError, "no runnable"):
                discover_names(listing)
        with self.assertRaisesRegex(ValueError, "duplicate"):
            discover_names("TestA\nTestA\n")
        with self.assertRaisesRegex(ValueError, "fewer runnable"):
            select_shard(["TestA"], 0, 2)
        for index, count in ((-1, 3), (3, 3), (0, 0)):
            with self.subTest(index=index, count=count), self.assertRaises(ValueError):
                select_shard(["TestA", "TestB", "TestC"], index, count)
        with self.assertRaisesRegex(ValueError, "empty shard"):
            run_regex([])


if __name__ == "__main__":
    unittest.main()
