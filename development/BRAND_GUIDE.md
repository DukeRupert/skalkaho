# Skalkaho Color Palette

---

## Inspiration

The palette draws from the landscape around Skalkaho Pass and the Bitterroot Valley:

- **Granite peaks** of the Sapphire Mountains
- **Evergreen forests** of lodgepole pine and Douglas fir
- **Copper mining heritage** of the region (Marcus Daly, Anaconda)
- **River stones** of the Bitterroot River
- **Morning sky** over the valley
- **Weathered wood** of old barns and cabins

The goal: grounded, natural, professional—not "tech startup."

---

## Primary Palette

### Slate (Primary)

The foundation. Used for primary text, headers, and key UI elements.

| Name | Hex | RGB | Usage |
|------|-----|-----|-------|
| **Slate 900** | `#1a1d23` | 26, 29, 35 | Primary text, headers |
| **Slate 700** | `#3d4450` | 61, 68, 80 | Secondary text, icons |
| **Slate 500** | `#64707d` | 100, 112, 125 | Tertiary text, captions |
| **Slate 300** | `#b0b8c4` | 176, 184, 196 | Borders, disabled states |
| **Slate 100** | `#e8ebef` | 232, 235, 239 | Dividers, subtle backgrounds |
| **Slate 50** | `#f4f5f7` | 244, 245, 247 | Page backgrounds |

### Forest (Secondary)

Evergreen. Used for secondary actions, categories, success-adjacent contexts.

| Name | Hex | RGB | Usage |
|------|-----|-----|-------|
| **Forest 900** | `#1a332a` | 26, 51, 42 | Dark accents |
| **Forest 700** | `#2d5a47` | 45, 90, 71 | Secondary buttons, icons |
| **Forest 500** | `#3d7a5f` | 61, 122, 95 | Hover states |
| **Forest 300** | `#7fb09a` | 127, 176, 154 | Light accents |
| **Forest 100** | `#d4e8df` | 212, 232, 223 | Subtle backgrounds |
| **Forest 50** | `#eef5f2` | 238, 245, 242 | Tinted backgrounds |

### Copper (Accent)

Warmth and action. Used for primary CTAs, highlights, active states.

| Name | Hex | RGB | Usage |
|------|-----|-----|-------|
| **Copper 900** | `#6b3d1f` | 107, 61, 31 | Dark accent (rare) |
| **Copper 700** | `#a35a2a` | 163, 90, 42 | Primary buttons, links |
| **Copper 500** | `#c97a42` | 201, 122, 66 | Hover states |
| **Copper 300** | `#e0a879` | 224, 168, 121 | Light accents |
| **Copper 100** | `#f5e0ce` | 245, 224, 206 | Subtle highlights |
| **Copper 50** | `#faf3ec` | 250, 243, 236 | Warm backgrounds |

---

## Semantic Colors

### Success

| Name | Hex | RGB | Usage |
|------|-----|-----|-------|
| **Success 700** | `#2d6a4f` | 45, 106, 79 | Success text, icons |
| **Success 500** | `#40916c` | 64, 145, 108 | Success borders |
| **Success 100** | `#d8f3dc` | 216, 243, 220 | Success backgrounds |

### Warning

| Name | Hex | RGB | Usage |
|------|-----|-----|-------|
| **Warning 700** | `#b86e00` | 184, 110, 0 | Warning text, icons |
| **Warning 500** | `#e09200` | 224, 146, 0 | Warning borders |
| **Warning 100** | `#fff3cd` | 255, 243, 205 | Warning backgrounds |

### Error

| Name | Hex | RGB | Usage |
|------|-----|-----|-------|
| **Error 700** | `#b42a2a` | 180, 42, 42 | Error text, icons |
| **Error 500** | `#dc3c3c` | 220, 60, 60 | Error borders |
| **Error 100** | `#fce8e8` | 252, 232, 232 | Error backgrounds |

### Info

| Name | Hex | RGB | Usage |
|------|-----|-----|-------|
| **Info 700** | `#2a6496` | 42, 100, 150 | Info text, icons |
| **Info 500** | `#4a90c2` | 74, 144, 194 | Info borders |
| **Info 100** | `#e3f2fc` | 227, 242, 252 | Info backgrounds |

---

## Backgrounds

| Name | Hex | RGB | Usage |
|------|-----|-----|-------|
| **White** | `#ffffff` | 255, 255, 255 | Cards, inputs, modals |
| **Stone** | `#f4f5f7` | 244, 245, 247 | Page background |
| **Warm Stone** | `#f9f7f5` | 249, 247, 245 | Alternative warm background |

---

## Usage Guidelines

### Text Hierarchy

| Level | Color | Example |
|-------|-------|---------|
| Primary text | Slate 900 | Headings, body copy, prices |
| Secondary text | Slate 700 | Labels, descriptions |
| Tertiary text | Slate 500 | Captions, timestamps, hints |
| Disabled text | Slate 300 | Inactive elements |

### Interactive Elements

| Element | Default | Hover | Active |
|---------|---------|-------|--------|
| Primary button | Copper 700 bg | Copper 500 bg | Copper 900 bg |
| Secondary button | Forest 700 bg | Forest 500 bg | Forest 900 bg |
| Text link | Copper 700 text | Copper 500 text | Copper 900 text |
| Ghost button | Slate 700 border | Slate 500 bg | Slate 700 bg |

### Borders & Dividers

| Context | Color |
|---------|-------|
| Input borders (default) | Slate 300 |
| Input borders (focus) | Copper 500 |
| Card borders | Slate 100 |
| Dividers | Slate 100 |
| Table borders | Slate 100 |

### Data & Categories

When color-coding categories or types, use the Forest and Copper families to differentiate:

| Type | Color suggestion |
|------|------------------|
| Materials | Forest 700 |
| Labor | Copper 700 |
| Category headers | Slate 900 on Slate 50 |

---

## Contrast & Accessibility

All text combinations meet WCAG 2.1 AA standards:

| Combination | Contrast Ratio | Pass |
|-------------|----------------|------|
| Slate 900 on White | 14.5:1 | ✓ AAA |
| Slate 700 on White | 8.2:1 | ✓ AAA |
| Slate 500 on White | 4.6:1 | ✓ AA |
| Copper 700 on White | 4.8:1 | ✓ AA |
| White on Copper 700 | 4.8:1 | ✓ AA |
| Forest 700 on White | 5.9:1 | ✓ AA |
| White on Forest 700 | 5.9:1 | ✓ AA |

**Note:** Always test final implementations. Slate 500 is the lightest acceptable for body text on white.

---

## Tailwind CSS Configuration

```javascript
// tailwind.config.js
module.exports = {
  theme: {
    extend: {
      colors: {
        slate: {
          50: '#f4f5f7',
          100: '#e8ebef',
          300: '#b0b8c4',
          500: '#64707d',
          700: '#3d4450',
          900: '#1a1d23',
        },
        forest: {
          50: '#eef5f2',
          100: '#d4e8df',
          300: '#7fb09a',
          500: '#3d7a5f',
          700: '#2d5a47',
          900: '#1a332a',
        },
        copper: {
          50: '#faf3ec',
          100: '#f5e0ce',
          300: '#e0a879',
          500: '#c97a42',
          700: '#a35a2a',
          900: '#6b3d1f',
        },
        success: {
          100: '#d8f3dc',
          500: '#40916c',
          700: '#2d6a4f',
        },
        warning: {
          100: '#fff3cd',
          500: '#e09200',
          700: '#b86e00',
        },
        error: {
          100: '#fce8e8',
          500: '#dc3c3c',
          700: '#b42a2a',
        },
        info: {
          100: '#e3f2fc',
          500: '#4a90c2',
          700: '#2a6496',
        },
      },
    },
  },
}
```

---

## Sample Combinations

### Primary Button

```
Background: copper-700 (#a35a2a)
Text: white (#ffffff)
Border: none
Hover: copper-500 (#c97a42)
```

### Secondary Button

```
Background: forest-700 (#2d5a47)
Text: white (#ffffff)
Border: none
Hover: forest-500 (#3d7a5f)
```

### Input Field

```
Background: white (#ffffff)
Text: slate-900 (#1a1d23)
Border: slate-300 (#b0b8c4)
Focus border: copper-500 (#c97a42)
Placeholder: slate-500 (#64707d)
```

### Card

```
Background: white (#ffffff)
Border: slate-100 (#e8ebef)
Header text: slate-900 (#1a1d23)
Body text: slate-700 (#3d4450)
```

### Page Layout

```
Page background: stone (#f4f5f7)
Card background: white (#ffffff)
Header background: slate-900 (#1a1d23)
Header text: white (#ffffff)
```

---

## Color Don'ts

- **Don't use pure black** (`#000000`) — Use Slate 900 instead
- **Don't use pure gray** — The slate scale has subtle cool undertones
- **Don't mix warm and cool grays** — Stick to the Slate scale
- **Don't use Copper for errors** — It's too close; use the dedicated Error red
- **Don't use low-contrast text** — Slate 300 is for borders, not readable text
- **Don't over-accent** — Copper is for primary actions only; too much feels cheap

---

## Revision History

| Date | Version | Notes |
|------|---------|-------|
| 2024-12-22 | 0.1 | Initial color palette |