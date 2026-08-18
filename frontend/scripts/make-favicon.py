"""Xuất favicon PNG/ICO từ cùng hình vẽ với public/favicon.svg.

Vẽ ở độ phân giải gấp 8 lần rồi thu nhỏ để cạnh mượt, vì PIL không khử răng cưa.
"""
import pathlib
from PIL import Image, ImageDraw

OUT = pathlib.Path("/Users/manhnv/Documents/Code/ELearning/frontend/public")
SS = 8                      # hệ số siêu lấy mẫu
VB = 48                     # viewBox của bản SVG

TOP_LEFT = (0x3d, 0x94, 0xf6)      # #3d94f6
BOTTOM_RIGHT = (0x0c, 0x3f, 0x9e)  # #0c3f9e
CAP_BODY = (0xcf, 0xe4, 0xff)      # #cfe4ff
BEAD = (0x9e, 0xcb, 0xff)          # #9ecbff


def quad_bezier(p0, ctrl, p2, steps=60):
    pts = []
    for i in range(steps + 1):
        t = i / steps
        u = 1 - t
        pts.append((
            u * u * p0[0] + 2 * u * t * ctrl[0] + t * t * p2[0],
            u * u * p0[1] + 2 * u * t * ctrl[1] + t * t * p2[1],
        ))
    return pts


def render(size):
    px = size * SS
    scale = px / VB

    def s(*coords):
        """Đổi toạ độ trong viewBox 48 sang pixel."""
        return tuple(c * scale for c in coords)

    # Nền chuyển sắc theo đường chéo trên-trái xuống dưới-phải.
    bg = Image.new("RGB", (px, px))
    bg_px = bg.load()
    for y in range(px):
        for x in range(px):
            t = (x + y) / (2 * px - 2)
            bg_px[x, y] = tuple(
                round(TOP_LEFT[i] + (BOTTOM_RIGHT[i] - TOP_LEFT[i]) * t) for i in range(3)
            )

    # Bo góc bằng mặt nạ alpha.
    mask = Image.new("L", (px, px), 0)
    ImageDraw.Draw(mask).rounded_rectangle([0, 0, px - 1, px - 1], radius=11 * scale, fill=255)
    icon = Image.new("RGBA", (px, px), (0, 0, 0, 0))
    icon.paste(bg, (0, 0), mask)

    d = ImageDraw.Draw(icon)

    # Thân mũ ôm đầu: hai cạnh thẳng, đáy cong.
    body = [s(14.4, 21.3)[0:2], s(14.4, 27.4)[0:2]]
    body += [(x * scale, y * scale) for x, y in quad_bezier((14.4, 27.4), (24, 36.2), (33.6, 27.4))]
    body += [s(33.6, 27.4)[0:2], s(33.6, 21.3)[0:2], s(24, 27.6)[0:2]]
    d.polygon(body, fill=CAP_BODY)

    # Ván mũ phẳng phía trên.
    d.polygon([s(24, 10.5), s(42, 18.2), s(24, 25.9), s(6, 18.2)], fill="white")

    # Dây tua và hạt cườm.
    d.line([s(39.4, 19.3), s(39.4, 27.7)], fill="white", width=round(1.7 * scale))
    r = 2.4 * scale
    cx, cy = 39.4 * scale, 30.2 * scale
    d.ellipse([cx - r, cy - r, cx + r, cy + r], fill=BEAD)

    return icon.resize((size, size), Image.LANCZOS)


sizes = [512, 180, 64, 48, 32, 16]
images = {n: render(n) for n in sizes}

images[512].save(OUT / "logo-512.png")
images[180].save(OUT / "apple-touch-icon.png")
images[32].save(OUT / "favicon-32.png")
images[16].save(OUT / "favicon-16.png")

# .ico gói nhiều kích thước cho trình duyệt cũ và thanh tác vụ Windows.
images[64].save(OUT / "favicon.ico", format="ICO",
                sizes=[(16, 16), (32, 32), (48, 48), (64, 64)])

for f in sorted(OUT.iterdir()):
    print(f"  {f.name:24} {f.stat().st_size:>7,} bytes")
