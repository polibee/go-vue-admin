import unittest
from pathlib import Path

from serve_spa import SpaHandler


class SpaHandlerTest(unittest.TestCase):
    def test_unknown_history_route_falls_back_to_index(self):
        handler = object.__new__(SpaHandler)
        handler.directory = str(Path(__file__).parent / "test-spa")

        resolved = handler.translate_path("/admin/dashboard")

        self.assertEqual(resolved, str(Path(handler.directory) / "index.html"))


if __name__ == "__main__":
    unittest.main()
