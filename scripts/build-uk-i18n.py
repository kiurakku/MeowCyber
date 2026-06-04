#!/usr/bin/env python3
"""Build uk-UA.json from en-US.json (batched translation + cache)."""
import json
import re
import sys
import time
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
SRC = ROOT / "web" / "static" / "i18n" / "en-US.json"
DST = ROOT / "web" / "static" / "i18n" / "uk-UA.json"
CACHE = Path(__file__).resolve().parent / "uk-i18n-cache.json"
BATCH = 40

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


def should_skip(s: str) -> bool:
    if "{{" in s:
        return True
    if re.match(r"^[\W\d_\s]+$", s):
        return True
    return False


def main():
    data = json.loads(SRC.read_text(encoding="utf-8"))
    data["lang"] = {"zhCN": "中文", "enUS": "English", "ukUA": "Українська"}

    strings = sorted({s for s in set() if True} | set(), key=len, reverse=True)
    tmp = set()
    collect_strings(data, tmp)
    strings = sorted(tmp, key=len, reverse=True)

    cache = {}
    if CACHE.exists():
        cache = json.loads(CACHE.read_text(encoding="utf-8"))
        print(f"Loaded cache: {len(cache)} entries", file=sys.stderr)

    try:
        from deep_translator import GoogleTranslator
    except ImportError:
        import subprocess
        subprocess.check_call([sys.executable, "-m", "pip", "install", "deep-translator", "-q"])
        from deep_translator import GoogleTranslator

    tr = GoogleTranslator(source="en", target="uk")
    pending = [s for s in strings if s not in cache]
    total = len(pending)
    print(f"Translating {total} strings ({len(cache)} cached)", file=sys.stderr)

    for i in range(0, total, BATCH):
        batch = pending[i : i + BATCH]
        to_tr = []
        for s in batch:
            if should_skip(s):
                cache[s] = brand(s)
            else:
                to_tr.append(s)
        if to_tr:
            try:
                out = tr.translate_batch(to_tr)
                for src, dst in zip(to_tr, out):
                    cache[src] = brand(dst) if dst else brand(src)
            except Exception:
                for s in to_tr:
                    try:
                        cache[s] = brand(tr.translate(s))
                    except Exception:
                        cache[s] = brand(s)
                    time.sleep(0.05)
        CACHE.write_text(json.dumps(cache, ensure_ascii=False), encoding="utf-8")
        done = min(i + BATCH, total)
        if done % 200 < BATCH or done == total:
            print(f"progress {done}/{total}", file=sys.stderr)
        time.sleep(0.2)

    out = apply_strings(data, cache)
    DST.write_text(json.dumps(out, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    print(f"Wrote {DST}", file=sys.stderr)


if __name__ == "__main__":
    main()
