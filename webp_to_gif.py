#!/usr/bin/env python3
"""
WebP to GIF converter using PIL/Pillow
Supports animated WebP files
"""
import sys
from pathlib import Path

try:
    from PIL import Image
except ImportError:
    print("Error: Pillow library not installed. Run: pip install Pillow", file=sys.stderr)
    sys.exit(1)


def convert_webp_to_gif(input_path: str, output_path: str) -> bool:
    """
    Convert WebP (including animated) to GIF

    Args:
        input_path: Path to input WebP file
        output_path: Path to output GIF file

    Returns:
        True if successful, False otherwise
    """
    try:
        # Open WebP file
        with Image.open(input_path) as im:
            # Check if animated
            is_animated = getattr(im, "is_animated", False)

            if is_animated:
                # Extract all frames
                frames = []
                durations = []

                for frame_num in range(im.n_frames):
                    im.seek(frame_num)
                    # Convert to RGB (GIF doesn't support RGBA well)
                    frame = im.convert("RGB")
                    frames.append(frame)

                    # Get frame duration (in milliseconds)
                    duration = im.info.get('duration', 100)
                    durations.append(duration)

                # Save as animated GIF
                frames[0].save(
                    output_path,
                    format='GIF',
                    append_images=frames[1:],
                    save_all=True,
                    duration=durations,
                    loop=0,  # Loop forever
                    optimize=False
                )
            else:
                # Single frame WebP
                im.convert("RGB").save(output_path, format='GIF')

            return True

    except Exception as e:
        print(f"Error converting {input_path}: {e}", file=sys.stderr)
        return False


def main():
    if len(sys.argv) != 3:
        print(f"Usage: {sys.argv[0]} <input.webp> <output.gif>", file=sys.stderr)
        sys.exit(1)

    input_path = sys.argv[1]
    output_path = sys.argv[2]

    # Validate input file exists
    if not Path(input_path).exists():
        print(f"Error: Input file not found: {input_path}", file=sys.stderr)
        sys.exit(1)

    # Convert
    success = convert_webp_to_gif(input_path, output_path)

    if success:
        sys.exit(0)
    else:
        sys.exit(1)


if __name__ == "__main__":
    main()
