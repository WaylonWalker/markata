"""SQLite schema integration checks; no Go SQLite driver is selected yet.

Run: python3 -m unittest pkg.editorsync.test_schema
"""
import hashlib
import sqlite3
import unittest
from pathlib import Path

SCHEMA = Path(__file__).with_name("schema.sql").read_text(encoding="utf-8")
CHECKSUM = hashlib.sha256(b"# Hello\n").hexdigest()


class SQLiteSchemaTests(unittest.TestCase):
    def setUp(self):
        self.db = sqlite3.connect(":memory:")
        self.db.execute("PRAGMA foreign_keys=ON")
        self.db.executescript(SCHEMA)

    def tearDown(self):
        self.db.close()

    def _insert_one(self):
        self.db.execute("INSERT INTO editor_sites (site_id, last_seq) VALUES (?, ?)", ("family", 1))
        self.db.execute(
            """INSERT INTO editor_documents
            (site_id, document_id, owner_subject, visibility, revision, changed_seq,
             content_hash, source, deleted)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)""",
            ("family", "post-1", "alice", "private", 1, 1, CHECKSUM, b"# Hello\n", 0),
        )
        self.db.execute(
            """INSERT INTO editor_document_versions
            (site_id, document_id, revision, changed_seq, owner_subject, visibility,
             content_hash, source, deleted)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)""",
            ("family", "post-1", 1, 1, "alice", "private", CHECKSUM, b"# Hello\n", 0),
        )
        self.db.execute(
            """INSERT INTO editor_mutation_receipts
            (site_id, subject_id, mutation_id, payload_digest, document_id,
             base_revision, accepted_revision, changed_seq)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?)""",
            ("family", "alice", "m-1", CHECKSUM, "post-1", 0, 1, 1),
        )
        self.db.execute(
            """INSERT INTO editor_changes
            (site_id, seq, document_id, revision, owner_subject, visibility, deleted)
            VALUES (?, ?, ?, ?, ?, ?, ?)""",
            ("family", 1, "post-1", 1, "alice", "private", 0),
        )
        self.db.commit()

    def test_migration_is_idempotent(self):
        self.db.executescript(SCHEMA)
        tables = {r[0] for r in self.db.execute(
            "SELECT name FROM sqlite_master WHERE type='table'"
        )}
        self.assertTrue({
            "editor_sites", "editor_documents", "editor_document_versions",
            "editor_mutation_receipts", "editor_changes", "editor_conflicts",
        }.issubset(tables))
        self.assertEqual(self.db.execute("PRAGMA foreign_key_check").fetchall(), [])

    def test_one_durable_mutation_and_replay_uniqueness(self):
        self._insert_one()
        rows = self.db.execute(
            "SELECT changed_seq, accepted_revision FROM editor_mutation_receipts "
            "WHERE site_id='family' AND subject_id='alice' AND mutation_id='m-1'"
        ).fetchall()
        self.assertEqual(rows, [(1, 1)])
        with self.assertRaises(sqlite3.IntegrityError):
            self.db.execute(
                """INSERT INTO editor_mutation_receipts
                (site_id, subject_id, mutation_id, payload_digest, document_id,
                 base_revision, accepted_revision, changed_seq)
                VALUES ('family','alice','m-1',?,?,0,1,1)""",
                (hashlib.sha256(b"changed").hexdigest(), "post-1"),
            )

    def test_invalid_visibility_and_missing_site_rejected(self):
        with self.assertRaises(sqlite3.IntegrityError):
            self.db.execute(
                """INSERT INTO editor_documents
                (site_id,document_id,owner_subject,visibility,revision,
                 changed_seq,content_hash,source)
                VALUES ('missing','post','alice','private',1,1,?,?)""",
                (CHECKSUM, b"data"),
            )
        self.db.execute("INSERT INTO editor_sites(site_id) VALUES ('family')")
        with self.assertRaises(sqlite3.IntegrityError):
            self.db.execute(
                """INSERT INTO editor_documents
                (site_id,document_id,owner_subject,visibility,revision,
                 changed_seq,content_hash,source)
                VALUES ('family','post','alice','admin',1,1,?,?)""",
                (CHECKSUM, b"data"),
            )

    def test_immutable_snapshot_revision_remains_readable(self):
        self._insert_one()
        new_hash = hashlib.sha256(b"# Later\n").hexdigest()
        self.db.execute(
            """UPDATE editor_documents
            SET revision=2, changed_seq=2, content_hash=?, source=?
            WHERE site_id='family' AND document_id='post-1'""",
            (new_hash, b"# Later\n"),
        )
        self.db.execute(
            """INSERT INTO editor_document_versions
            (site_id, document_id, revision, changed_seq, owner_subject,
             visibility, content_hash, source, deleted)
            VALUES ('family','post-1',2,2,'alice','private',?,?,0)""",
            (new_hash, b"# Later\n"),
        )
        self.db.commit()
        historic = self.db.execute(
            """SELECT revision, source FROM editor_document_versions
            WHERE site_id='family' AND changed_seq<=? ORDER BY revision DESC LIMIT 1""",
            (1,),
        ).fetchone()
        self.assertEqual(historic, (1, b"# Hello\n"))
        self.assertEqual(
            self.db.execute("SELECT revision FROM editor_documents").fetchone()[0], 2
        )

    def test_transaction_rollback_keeps_no_partial_receipt(self):
        try:
            with self.db:
                self.db.execute("INSERT INTO editor_sites(site_id) VALUES ('family')")
                self.db.execute("INSERT INTO editor_sites(site_id) VALUES ('family')")
        except sqlite3.IntegrityError:
            pass
        self.assertEqual(self.db.execute("SELECT COUNT(*) FROM editor_sites").fetchone()[0], 0)


if __name__ == "__main__":
    unittest.main()
