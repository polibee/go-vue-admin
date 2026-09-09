from __future__ import annotations

import argparse
from http.server import SimpleHTTPRequestHandler, ThreadingHTTPServer
from pathlib import Path


class SpaHandler(SimpleHTTPRequestHandler):
    def translate_path(self, path: str) -> str:
        resolved = Path(super().translate_path(path))
        if resolved.is_file():
            return str(resolved)

        return str(Path(self.directory or ".") / "index.html")


def main() -> None:
    parser = argparse.ArgumentParser(description="Serve a single-page app with history fallback")
    parser.add_argument("--directory", default="dist")
    parser.add_argument("--bind", default="127.0.0.1")
    parser.add_argument("--port", type=int, default=4173)
    args = parser.parse_args()

    handler = lambda *handler_args, **handler_kwargs: SpaHandler(  # noqa: E731
        *handler_args,
        directory=args.directory,
        **handler_kwargs,
    )
    server = ThreadingHTTPServer((args.bind, args.port), handler)
    print(f"Serving {Path(args.directory).resolve()} at http://{args.bind}:{args.port}/", flush=True)
    server.serve_forever()


if __name__ == "__main__":
    main()
