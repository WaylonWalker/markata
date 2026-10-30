"""Verify FastAPI's LSP configuration and browser origin boundary."""
from pathlib import Path
import sys
import unittest
from unittest.mock import AsyncMock, patch

from fastapi.testclient import TestClient
from starlette.websockets import WebSocketDisconnect

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))
from editor_test_support import load_editor


class BrowserEndpointTest(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.demo = load_editor()

    def test_configuration_ignores_client_workspace(self):
        with TestClient(self.demo.app, client=("127.0.0.1", 50000)) as client:
            response = client.get("/lsp/config?rootUri=file:///outside")
            self.assertEqual(response.status_code, 200)
            self.assertEqual(response.json()["rootUri"], self.demo.LSP_WORKSPACE.as_uri())
            self.assertTrue(response.json()["documentUri"].endswith("/.markata/browser-draft.md"))
            self.assertEqual(response.headers["cache-control"], "no-store")
            for name in ["lsp.js", "lsp.css"]:
                self.assertEqual(client.get("/editor/" + name).status_code, 200)
            self.assertEqual(client.get("/editor/not-an-asset").status_code, 404)

    def test_foreign_or_missing_origin_never_starts_process(self):
        with patch.object(self.demo.asyncio, "create_subprocess_exec", new_callable=AsyncMock) as spawn:
            with TestClient(self.demo.app, client=("127.0.0.1", 50000)) as client:
                for headers in [{}, {"origin": "https://foreign.invalid"},
                                {"host": "evil.invalid", "origin": "http://evil.invalid"}]:
                    with self.assertRaises(WebSocketDisconnect):
                        with client.websocket_connect("/lsp", headers=headers):
                            self.fail("Accepted invalid origin")
            spawn.assert_not_awaited()

    def test_remote_peer_cannot_spoof_local_origin(self):
        with patch.object(self.demo.asyncio, "create_subprocess_exec", new_callable=AsyncMock) as spawn:
            with TestClient(self.demo.app, client=("198.51.100.9", 50000)) as client:
                with self.assertRaises(WebSocketDisconnect):
                    with client.websocket_connect(
                        "/lsp", headers={"host": "localhost", "origin": "http://localhost"}
                    ):
                        self.fail("A remote client spoofed localhost headers")
            spawn.assert_not_awaited()

    def test_missing_binary_closes_only_language_connection(self):
        with patch.object(self.demo, "lsp_command", return_value=["/nonexistent/markata"]):
            with TestClient(self.demo.app, client=("127.0.0.1", 50000)) as client:
                with client.websocket_connect("/lsp", headers={"host": "localhost", "origin": "http://localhost"}) as socket:
                    with self.assertRaises(WebSocketDisconnect):
                        socket.receive_json()
                self.assertEqual(client.get("/lsp/config").status_code, 200)


if __name__ == "__main__":
    unittest.main()
