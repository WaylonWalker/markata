"""Check the documented FastAPI examples against an installed Markata wheel."""
import re
from pathlib import Path
from fastapi.testclient import TestClient

source = Path(__file__).resolve().parents[1] / 'docs/guides/python-renderer.md'
examples = re.findall(r'```python\n(.*?)\n```', source.read_text(), flags=re.S)

inline = {}
exec(examples[0], inline)
with TestClient(inline['app']) as client:
    response = client.get('/')
    assert response.status_code == 200
    assert '<h1' in response.text and '<style>' in response.text
    assert response.headers['content-type'].startswith('text/html')

linked = {}
exec(examples[2], linked)
with TestClient(linked['app']) as client:
    response = client.get('/')
    assert response.status_code == 200 and '/md-assets/css/main.css' in response.text
    for asset in linked['assets'].values():
        response = client.get('/md-assets/' + asset['id'])
        assert response.status_code == 200
        assert response.text == asset['content']
        expected = 'text/css' if asset['kind'] == 'style' else 'text/javascript'
        assert response.headers['content-type'].startswith(expected)
    assert client.get('/md-assets/missing').status_code == 404
print('Documented inline and linked-asset FastAPI examples verified.')
