#!/usr/bin/env python3
"""
WebP to JPEG converter using PIL/Pillow
Supports only static WebP files
"""
import sys
from pathlib import Path

try:
    from PIL import Image
except ImportError:
    print("Error: Pillow library not installed. Run: pip install Pillow", file=sys.stderr)
    sys.exit(1)


def convert_webp_to_jpeg(input_path: str, output_path: str, quality: int = 100) -> bool:
    """
    Convert static WebP to JPEG

    Args:
        input_path: Path to input WebP file
        output_path: Path to output JPEG file
        quality: JPEG quality (1-100, default 100)

    Returns:
        True if successful, False otherwise
    """
    try:
        with Image.open(input_path) as im:
            # Verify it's WebP
            if im.format != 'WEBP':
                print(f"Error: File is not WebP format: {im.format}", file=sys.stderr)
                return False

            # Check if animated (should not be)
            is_animated = getattr(im, "is_animated", False)
            if is_animated:
                print("Error: Animated WebP detected. Use GIF conversion instead.", file=sys.stderr)
                return False

            # Convert to RGB (JPEG doesn't support transparency)
            # If WebP has alpha channel, composite on white background
            if im.mode in ('RGBA', 'LA', 'PA'):
                # Create white background
                background = Image.new('RGB', im.size, (255, 255, 255))
                # Paste with alpha channel as mask
                background.paste(im, mask=im.split()[-1] if im.mode == 'RGBA' else None)
                rgb_im = background
            elif im.mode != 'RGB':
                rgb_im = im.convert('RGB')
            else:
                rgb_im = im

            # Save as JPEG with specified quality
            # optimize=True enables Huffman table optimization
            # Keep EXIF data if present
            exif_data = im.info.get('exif', None)

            save_kwargs = {
                'format': 'JPEG',
                'quality': quality,
                'optimize': True,
            }

            if exif_data:
                save_kwargs['exif'] = exif_data

            rgb_im.save(output_path, **save_kwargs)

            return True

    except Exception as e:
        print(f"Error converting {input_path}: {e}", file=sys.stderr)
        return False


def main():
    if len(sys.argv) < 3 or len(sys.argv) > 4:
        print(f"Usage: {sys.argv[0]} <input.webp> <output.jpg> [quality]", file=sys.stderr)
        sys.exit(1)

    input_path = sys.argv[1]
    output_path = sys.argv[2]
    quality = int(sys.argv[3]) if len(sys.argv) == 4 else 100

    # Validate quality
    if not 1 <= quality <= 100:
        print(f"Error: Quality must be between 1 and 100, got {quality}", file=sys.stderr)
        sys.exit(1)

    # Validate input file exists
    if not Path(input_path).exists():
        print(f"Error: Input file not found: {input_path}", file=sys.stderr)
        sys.exit(1)

    # Convert
    success = convert_webp_to_jpeg(input_path, output_path, quality)

    if success:
        sys.exit(0)
    else:
        sys.exit(1)


if __name__ == "__main__":
    main()
