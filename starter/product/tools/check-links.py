#!/usr/bin/env python3
"""Validate local Markdown links and section anchors in a product repository."""

from __future__ import annotations

import argparse
import json
import re
import sys
import unicodedata
from pathlib import Path
from urllib.parse import unquote, urlsplit

LINK = re.compile(r"(?m)(?P<image>!)?\[[^\]\n]+\]\((?P<target>[^)\n]+)\)")
ID = re.compile(r"\b[A-Z][A-Z0-9]*(?:-[A-Z0-9]+)+\b")
HEADING = re.compile(r"(?m)^#{1,6}\s+(.+?)\s*#*\s*$")
EXPLICIT_ANCHOR = re.compile(r'''<a\s+(?:[^>]*\s)?id=["']([^"']+)["']''', re.I)


def slug(value: str) -> str:
    value = re.sub(r"<[^>]+>", "", value)
    value = unicodedata.normalize("NFKD", value).encode("ascii", "ignore").decode()
    value = re.sub(r"[^\w\s-]", "", value.lower())
    return re.sub(r"[\s-]+", "-", value).strip("-")


def anchors(path: Path) -> set[str]:
    text = path.read_text(encoding="utf-8")
    result = {slug(match.group(1)) for match in HEADING.finditer(text)}
    result.update(match.group(1) for match in EXPLICIT_ANCHOR.finditer(text))
    return result


def check(root: Path) -> list[str]:
    errors: list[str] = []
    known_ids = registered_ids(root)
    for source in sorted(root.rglob("*.md")):
        text = re.sub(r"(?s)```.*?```", "", source.read_text(encoding="utf-8"))
        if "_template" not in source.parts:
            for match in ID.finditer(text):
                identifier = match.group(0)
                if (identifier in known_ids and not identifier.endswith(("-TBD", "-TEMPLATE", "-XXX"))
                        and not is_metadata_reference(text, match.start())
                        and not is_linked_id(text, match.start(), match.end())):
                    errors.append(f"{source.relative_to(root)}: unlinked artifact reference {match.group(0)}")
        for match in LINK.finditer(text):
            target = match.group("target").strip().strip("<>")
            parsed = urlsplit(target)
            if not target or target.startswith("mailto:") or parsed.scheme or parsed.netloc:
                continue
            if target.startswith("#"):
                if target[1:] not in anchors(source):
                    errors.append(f"{source.relative_to(root)}: missing section {target}")
                continue
            destination = (source.parent / unquote(parsed.path)).resolve()
            try:
                destination.relative_to(root.resolve())
            except ValueError:
                errors.append(f"{source.relative_to(root)}: link escapes root {target}")
                continue
            if not destination.exists():
                errors.append(f"{source.relative_to(root)}: missing target {target}")
            elif parsed.fragment and destination.is_file() and parsed.fragment not in anchors(destination):
                errors.append(f"{source.relative_to(root)}: missing section {parsed.fragment} in {parsed.path}")
    return errors


def registered_ids(root: Path) -> set[str]:
    """Return IDs whose canonical paths are known to the product registry."""
    result: set[str] = set()
    for relative in (Path(".product/artifacts.json"), Path(".product/decisions.json")):
        path = root / relative
        if not path.exists():
            continue
        try:
            document = json.loads(path.read_text(encoding="utf-8"))
        except json.JSONDecodeError:
            continue
        values = document.get("artifacts", []) if "artifacts" in document else document.get("decisions", [])
        for value in values:
            if isinstance(value, dict) and isinstance(value.get("id"), str):
                result.add(value["id"])
    return result


def is_linked_id(text: str, start: int, end: int) -> bool:
    """Whether an ID occurrence is the label of a Markdown link."""
    before = text[:start].rfind("[")
    close = text.find("]", end)
    return before >= 0 and close >= 0 and text[before + 1 : close].find(text[start:end]) >= 0 and text[close + 1 :].lstrip().startswith("(")


def is_metadata_reference(text: str, start: int) -> bool:
    """Do not treat YAML/table identity metadata as navigational prose."""
    line_start = text.rfind("\n", 0, start) + 1
    line = text[line_start : text.find("\n", start) if "\n" in text[start:] else len(text)]
    return bool(re.search(r"\b(id|parent ids|parents|children|depends on|related)\b", line, re.I))


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("root", nargs="?", default=".", help="repository or product root")
    root = Path(parser.parse_args().root).resolve()
    if not root.is_dir():
        print(f"error: root does not exist: {root}", file=sys.stderr)
        return 2
    errors = check(root)
    for error in errors:
        print(f"ERROR links: {error}")
    print(f"Link check: {'PASS' if not errors else 'FAIL'} ({len(errors)} errors)")
    return 0 if not errors else 1


if __name__ == "__main__":
    raise SystemExit(main())
