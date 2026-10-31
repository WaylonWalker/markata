"""Catalog coverage and real-renderer conformance for custom Markdown commands."""
import json
from pathlib import Path
import subprocess
import unittest

from markata import render

ROOT = Path(__file__).resolve().parents[2]


def load_catalog():
    result = subprocess.run(['node', '-e', """
      const fs = require('node:fs'), vm = require('node:vm');
      const context = {window:{}};
      vm.runInNewContext(fs.readFileSync(process.argv[1], 'utf8'), context);
      process.stdout.write(JSON.stringify(context.window.MarkataMarkdownCommands));
    """, str(ROOT / 'pkg/workbench/web/command-catalog.js')], capture_output=True, text=True, check=True)
    return json.loads(result.stdout)


class CommandCatalogTest(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.commands = load_catalog()
        cls.by_id = {command['id']: command for command in cls.commands}

    def snippet(self, command_id):
        command = self.by_id[command_id]
        return command['template'].replace('{{text}}', command['placeholder'])

    def test_complete_unique_catalog(self):
        self.assertEqual(len(self.commands), len(self.by_id))
        required = {
            'note', 'tip', 'warning', 'aside', 'collapse', 'collapse-open', 'vsplit',
            'highlight', 'keys', 'definition', 'mermaid', 'chart', 'lucide', 'icon-pack',
            'md-card', 'md-box', 'md-layout', 'md-columns', 'md-grid', 'md-grid-rail',
            'md-grid-pending', 'md-stack', 'md-divider', 'md-connector', 'md-section',
            'md-icon', 'md-feature', 'md-status', 'md-stat', 'md-delta', 'md-meter',
            'md-progress', 'md-sparkline', 'md-checklist', 'md-source', 'footnote', 'attributes',
            'emoji', 'superscript', 'subscript', 'figure', 'hover-card', 'hover-template',
            'wikilink', 'post-embed', 'external-embed', 'youtube', 'toc', 'glossary',
            'contribution-graph', 'wa-html',
        }
        required.update('road-' + style for style in ['winding', 'timeline', 'steps', 'snake', 'cycle', 'constellation', 'git', 'editorial', 'story', 'reveal', 'story-reveal'])
        required.update('wa-' + kind for kind in ['comparison', 'details', 'tabs', 'copy', 'qr', 'badge', 'tag', 'tooltip', 'carousel', 'animated-image'])
        required.update('admonition-' + kind for kind in ['info', 'hint', 'success', 'caution', 'danger', 'error', 'bug', 'example', 'quote', 'abstract'])
        self.assertFalse(required - self.by_id.keys())
        for command in self.commands:
            self.assertEqual(command['template'].count('{{text}}'), 1, command['id'])
            self.assertTrue(command['placeholder'], command['id'])
        for command_id in ['wa-details', 'icon-pack', 'wikilink', 'post-embed', 'external-embed', 'toc']:
            self.assertIn('Requires', self.by_id[command_id]['description'])
        self.assertEqual(self.snippet('toc'), '[[toc]]')

    def test_every_bundled_icon_is_insertable_and_renders(self):
        names = {path.stem for path in (ROOT / 'pkg/plugins/standalone_icons').glob('*.svg')}
        actual = {command['id'].removeprefix('lucide-') for command in self.commands if command['id'].startswith('lucide-')}
        self.assertEqual(actual, names)
        for name in names:
            self.assertFalse(self.by_id['lucide-' + name]['useSelection'])
        source = '\n\n'.join(self.snippet('lucide-' + name) for name in sorted(names))
        result = render(source)
        self.assertEqual(result.html.count('data-icon="lucide/'), len(names))
        self.assertNotIn(':lucide-', result.html)

    def test_standalone_custom_containers_render_semantically(self):
        command_ids = [key for key in self.by_id if key.startswith(('road-', 'md-', 'admonition-'))]
        source = '\n\n'.join(self.snippet(key) for key in command_ids)
        result = render(source)
        for key in command_ids:
            if key == 'md-grid-rail': expected = 'class="md-grid rail"'
            elif key == 'md-grid-pending': expected = 'class="md-grid pending"'
            elif key.startswith('md-'): expected = 'class="' + key + '"'
            elif key.startswith('road-'): expected = 'class="road '
            else: expected = 'admonition ' + key.removeprefix('admonition-')
            self.assertIn(expected, result.html, key)
        self.assertIn('<meter ', result.html)
        self.assertIn('<progress ', result.html)
        self.assertIn('<details class="md-source"', result.html)
        self.assertNotIn('::::', result.html)

    def test_example_copies_match_the_shared_asset(self):
        source = (ROOT / 'examples/command-palette-showcase.md').read_bytes()
        for path in ['pkg/workbench/web/command-palette-example.md',
                     'examples/wasm-workbench/command-palette-example.md']:
            with self.subTest(path=path):
                self.assertEqual((ROOT / path).read_bytes(), source)

    def test_chart_defaults_are_valid_configs(self):
        for command in self.commands:
            if command['id'] == 'chart' or command['id'].startswith('chart-'):
                with self.subTest(command=command['id']):
                    config = json.loads(command['placeholder'])
                    self.assertTrue(config['data']['datasets'])
                    self.assertIn('chartjs', self.snippet(command['id']))

    def test_showcase_renders_without_literal_shortcodes(self):
        result = render((ROOT / 'examples/command-palette-showcase.md').read_text(), unsafe=True)
        for marker in ['md-card', 'road steps', 'md-feature', '<progress ', 'md-source', 'data-icon="lucide/']:
            self.assertIn(marker, result.html)
        self.assertNotIn(':lucide-', result.html)
        self.assertTrue(result.features['mermaid'])
        self.assertTrue(result.features['chartjs'])


if __name__ == '__main__':
    unittest.main()
