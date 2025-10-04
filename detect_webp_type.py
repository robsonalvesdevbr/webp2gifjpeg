#!/usr/bin/env python3
"""
Detect if a WebP file is animated or static
Exit codes:
  0 - Static WebP
  1 - Animated WebP
  2 - Error (file not found, not WebP, etc.)
"""
import sys
from pathlib import Path

try:
    from PIL import Image
except ImportError:
    print("Error: Pillow library not installed", file=sys.stderr)
    sys.exit(2)


def detect_webp_type(input_path: str) -> int:
    """
    Detect if a WebP file is animated or static

    Returns:
        0 for static WebP
        1 for animated WebP
        2 for errors
    """
    try:
        with Image.open(input_path) as im:
            # Verify it's WebP
            if im.format != 'WEBP':
                print(f"Error: File is not WebP format: {im.format}", file=sys.stderr)
                return 2

            # Check if animated
            is_animated = getattr(im, "is_animated", False)

            if is_animated:
                print("animated")  # stdout for parsing
                return 1
            else:
                print("static")
                return 0

    except Exception as e:
        print(f"Error detecting WebP type: {e}", file=sys.stderr)
        return 2


if __name__ == "__main__":
    if len(sys.argv) != 2:
        print(f"Usage: {sys.argv[0]} <input.webp>", file=sys.stderr)
        sys.exit(2)

    input_path = sys.argv[1]

    if not Path(input_path).exists():
        print(f"Error: File not found: {input_path}", file=sys.stderr)
        sys.exit(2)

    sys.exit(detect_webp_type(input_path))
