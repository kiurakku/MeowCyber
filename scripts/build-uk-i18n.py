#!/usr/bin/env python3
"""Build uk-UA.json from en-US.json with Ukrainian translations."""
import json
import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
SRC = ROOT / "web" / "static" / "i18n" / "en-US.json"
DST = ROOT / "web" / "static" / "i18n" / "uk-UA.json"

BRAND_RE = [
    (re.compile(r"CyberStrikeAI", re.I), "MeowCyber"),
    (re.compile(r"CyberStrike", re.I), "MeowCyber"),
]


def brand(s: str) -> str:
    for pat, rep in BRAND_RE:
        s = pat.sub(rep, s)
    return s


def collect_strings(obj, out: set):
    if isinstance(obj, dict):
        for v in obj.values():
            collect_strings(v, out)
    elif isinstance(obj, str) and obj.strip():
        out.add(obj)


def apply_strings(obj, cache: dict):
    if isinstance(obj, dict):
        return {k: apply_strings(v, cache) for k, v in obj.items()}
    if isinstance(obj, str):
        return cache.get(obj, brand(obj))
    return obj


def main():
    data = json.loads(SRC.read_text(encoding="utf-8"))
    data["lang"] = {
        "zhCN": "中文",
        "enUS": "English",
        "ukUA": "Українська",
    }
    strings = set()
    collect_strings(data, strings)
    strings = sorted(strings, key=len, reverse=True)

    try:
        from deep_translator import GoogleTranslator
    except ImportError:
        print("Installing deep-translator...", file=sys.stderr)
        import subprocess
        subprocess.check_call([sys.executable, "-m", "pip", "install", "deep-translator", "-q"])
        from deep_translator import GoogleTranslator

    tr = GoogleTranslator(source="en", target="uk")
    cache = {}
    total = len(strings)
    for i, s in enumerate(strings):
        if "{{" in s or re.match(r"^[\W\d_]+$", s):
            cache[s] = brand(s)
            continue
        try:
            cache[s] = brand(tr.translate(s))
        except Exception as e:
            print(f"warn translate ({i+1}/{total}): {e!r} -> keep EN", file=sys.stderr)
            cache[s] = brand(s)
        if (i + 1) % 50 == 0:
            print(f"translated {i+1}/{total}", file=sys.stderr)

    out = apply_strings(data, cache)
    DST.write_text(json.dumps(out, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    print(f"Wrote {DST}")


if __name__ == "__main__":
    main()
