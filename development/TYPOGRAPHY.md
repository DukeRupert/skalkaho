# Skalkaho Typography

---

## Inspiration

The typography draws from the same "grounded, professional, not-tech-startup" ethos as the color palette:

- **Barlow** — Primary typeface. Inspired by California highway signage and public infrastructure. Sturdy, legible, slightly condensed. Feels like equipment manuals, spec sheets, and construction documents.
- **JetBrains Mono** — Monospace for keyboard shortcuts and code. Clean, modern, excellent legibility at small sizes.

The goal: workwear typography. Functional, no-nonsense, built to last.

---

## Font Families

### Primary: Barlow

Used for all headings, body text, labels, and UI elements.

| Property | Value |
|----------|-------|
| **Family** | `'Barlow', sans-serif` |
| **Source** | [Google Fonts](https://fonts.google.com/specimen/Barlow) |
| **Weights used** | 400 (Regular), 500 (Medium), 600 (SemiBold), 700 (Bold) |
| **Features** | Tabular figures, proportional figures |

### Monospace: JetBrains Mono

Used exclusively for keyboard shortcuts, inline code, and data entry fields where character alignment matters.

| Property | Value |
|----------|-------|
| **Family** | `'JetBrains Mono', monospace` |
| **Source** | [Google Fonts](https://fonts.google.com/specimen/JetBrains+Mono) |
| **Weights used** | 400 (Regular), 500 (Medium) |

---

## Font Loading

### HTML (Google Fonts)

```html
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link href="https://fonts.googleapis.com/css2?family=Barlow:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500&display=swap" rel="stylesheet">
```

### CSS (@import)

```css
@import url('https://fonts.googleapis.com/css2?family=Barlow:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500&display=swap');
```

### Self-hosted (if preferred)

Download from Google Fonts and serve via `/static/fonts/`. Use `@font-face` declarations with `font-display: swap`.

---

## Type Scale

Based on a 1.25 ratio (Major Third) with a 16px base.

| Name | Size (rem) | Size (px) | Line Height | Use |
|------|------------|-----------|-------------|-----|
| **xs** | 0.75rem | 12px | 1.5 | Badges, fine print |
| **sm** | 0.875rem | 14px | 1.5 | Table content, secondary text |
| **base** | 1rem | 16px | 1.5 | Body text, form inputs |
| **lg** | 1.125rem | 18px | 1.4 | Lead paragraphs |
| **xl** | 1.25rem | 20px | 1.4 | Section totals, emphasis |
| **2xl** | 1.5rem | 24px | 1.3 | Page headings (h1) |
| **3xl** | 1.875rem | 30px | 1.3 | Hero headings (rare) |

---

## Text Styles

### Page Heading (h1)

The primary heading for a page or view. Used for job names, category names.

```css
.page-heading {
  font-family: 'Barlow', sans-serif;
  font-size: 1.5rem;      /* 24px */
  font-weight: 700;
  line-height: 1.3;
  letter-spacing: -0.01em;
  color: var(--slate-900);
}
```

**Tailwind:** `font-sans text-2xl font-bold leading-tight tracking-tight text-slate-900`

### Section Heading (h2)

Used for section labels like "Subcategories", "Items", "Line Items".

```css
.section-heading {
  font-family: 'Barlow', sans-serif;
  font-size: 0.875rem;    /* 14px */
  font-weight: 600;
  line-height: 1.5;
  letter-spacing: 0.025em;
  text-transform: uppercase;
  color: var(--slate-700);
}
```

**Tailwind:** `font-sans text-sm font-semibold leading-normal tracking-wide uppercase text-slate-700`

### Body Text

Default text for descriptions, instructions, and general content.

```css
.body-text {
  font-family: 'Barlow', sans-serif;
  font-size: 1rem;        /* 16px */
  font-weight: 400;
  line-height: 1.5;
  color: var(--slate-900);
}
```

**Tailwind:** `font-sans text-base font-normal leading-normal text-slate-900`

### Secondary Text

Used for supporting information, timestamps, hints.

```css
.secondary-text {
  font-family: 'Barlow', sans-serif;
  font-size: 0.875rem;    /* 14px */
  font-weight: 400;
  line-height: 1.5;
  color: var(--slate-500);
}
```

**Tailwind:** `font-sans text-sm font-normal leading-normal text-slate-500`

### Table Header

Column headers in data tables.

```css
.table-header {
  font-family: 'Barlow', sans-serif;
  font-size: 0.75rem;     /* 12px */
  font-weight: 500;
  line-height: 1.5;
  letter-spacing: 0.05em;
  text-transform: uppercase;
  color: var(--slate-500);
}
```

**Tailwind:** `font-sans text-xs font-medium leading-normal tracking-wider uppercase text-slate-500`

### Table Cell

Data within tables.

```css
.table-cell {
  font-family: 'Barlow', sans-serif;
  font-size: 0.875rem;    /* 14px */
  font-weight: 400;
  line-height: 1.5;
  color: var(--slate-900);
}
```

**Tailwind:** `font-sans text-sm font-normal leading-normal text-slate-900`

### Numeric Data

Prices, quantities, totals. Uses tabular figures for alignment.

```css
.numeric {
  font-family: 'Barlow', sans-serif;
  font-size: 0.875rem;    /* 14px */
  font-weight: 400;
  line-height: 1.5;
  font-variant-numeric: tabular-nums;
  color: var(--slate-900);
}
```

**Tailwind:** `font-sans text-sm font-normal leading-normal tabular-nums text-slate-900`

### Total / Emphasis Numeric

Category totals, job totals, final prices.

```css
.total {
  font-family: 'Barlow', sans-serif;
  font-size: 1.25rem;     /* 20px */
  font-weight: 600;
  line-height: 1.4;
  font-variant-numeric: tabular-nums;
  color: var(--slate-900);
}
```

**Tailwind:** `font-sans text-xl font-semibold leading-snug tabular-nums text-slate-900`

### Keyboard Shortcut

Inline keyboard hints displayed in the UI.

```css
.kbd {
  font-family: 'JetBrains Mono', monospace;
  font-size: 0.75rem;     /* 12px */
  font-weight: 400;
  line-height: 1;
  padding: 0.125rem 0.375rem;
  background-color: var(--slate-100);
  border: 1px solid var(--slate-300);
  border-radius: 4px;
  color: var(--slate-700);
}
```

**Tailwind:** `font-mono text-xs font-normal px-1.5 py-0.5 bg-slate-100 border border-slate-300 rounded text-slate-700`

### Form Label

Labels above or beside form inputs.

```css
.form-label {
  font-family: 'Barlow', sans-serif;
  font-size: 0.875rem;    /* 14px */
  font-weight: 500;
  line-height: 1.5;
  color: var(--slate-700);
}
```

**Tailwind:** `font-sans text-sm font-medium leading-normal text-slate-700`

### Form Input

Text inside form fields.

```css
.form-input {
  font-family: 'Barlow', sans-serif;
  font-size: 1rem;        /* 16px */
  font-weight: 400;
  line-height: 1.5;
  color: var(--slate-900);
}

.form-input::placeholder {
  color: var(--slate-500);
}
```

**Tailwind:** `font-sans text-base font-normal leading-normal text-slate-900 placeholder:text-slate-500`

### Breadcrumb

Navigation breadcrumbs.

```css
.breadcrumb {
  font-family: 'Barlow', sans-serif;
  font-size: 0.875rem;    /* 14px */
  font-weight: 400;
  line-height: 1.5;
  color: var(--slate-500);
}

.breadcrumb a {
  color: var(--copper-700);
  text-decoration: none;
}

.breadcrumb a:hover {
  color: var(--copper-500);
}
```

**Tailwind:** `font-sans text-sm font-normal leading-normal text-slate-500` (links: `text-copper-700 hover:text-copper-500 no-underline`)

### Badge / Tag

Small labels for categorization (Material, Labor, Equipment).

```css
.badge {
  font-family: 'Barlow', sans-serif;
  font-size: 0.625rem;    /* 10px */
  font-weight: 500;
  line-height: 1;
  letter-spacing: 0.05em;
  text-transform: uppercase;
  padding: 0.25rem 0.5rem;
  border-radius: 3px;
}
```

**Tailwind:** `font-sans text-[10px] font-medium leading-none tracking-wider uppercase px-2 py-1 rounded-sm`

---

## Tailwind CSS Configuration

Add to your existing `tailwind.config.js`:

```javascript
module.exports = {
  theme: {
    extend: {
      fontFamily: {
        sans: ['Barlow', 'ui-sans-serif', 'system-ui', 'sans-serif'],
        mono: ['JetBrains Mono', 'ui-monospace', 'monospace'],
      },
      fontSize: {
        'xs': ['0.75rem', { lineHeight: '1.5' }],
        'sm': ['0.875rem', { lineHeight: '1.5' }],
        'base': ['1rem', { lineHeight: '1.5' }],
        'lg': ['1.125rem', { lineHeight: '1.4' }],
        'xl': ['1.25rem', { lineHeight: '1.4' }],
        '2xl': ['1.5rem', { lineHeight: '1.3' }],
        '3xl': ['1.875rem', { lineHeight: '1.3' }],
      },
      letterSpacing: {
        tighter: '-0.02em',
        tight: '-0.01em',
        normal: '0',
        wide: '0.025em',
        wider: '0.05em',
      },
    },
  },
}
```

---

## Usage Guidelines

### Weight Hierarchy

| Weight | Value | Use |
|--------|-------|-----|
| **Bold (700)** | `font-bold` | Page headings only |
| **SemiBold (600)** | `font-semibold` | Section headings, totals, emphasis |
| **Medium (500)** | `font-medium` | Labels, table headers, badges |
| **Regular (400)** | `font-normal` | Body text, table cells, inputs |

### When to Use Monospace

- Keyboard shortcut hints (`c`, `m`, `l`, `%`, `⏎`)
- Inline code references (if any)
- **Not for:** prices, quantities, or general numeric data (use tabular-nums instead)

### Numeric Alignment

Always use `font-variant-numeric: tabular-nums` (Tailwind: `tabular-nums`) for:
- Price columns
- Quantity columns
- Any vertically-aligned numbers

This ensures digits align in columns without needing monospace.

### Text Colors by Context

| Context | Color | Tailwind |
|---------|-------|----------|
| Primary text | Slate 900 | `text-slate-900` |
| Secondary text | Slate 700 | `text-slate-700` |
| Tertiary / hints | Slate 500 | `text-slate-500` |
| Disabled | Slate 300 | `text-slate-300` |
| Links | Copper 700 | `text-copper-700` |
| Link hover | Copper 500 | `text-copper-500` |
| Error text | Error 700 | `text-error-700` |
| Success text | Success 700 | `text-success-700` |

---

## Component Examples

### Page Header

```html
<div class="mb-6">
  <nav class="text-sm text-slate-500 mb-2">
    <a href="#" class="text-copper-700 hover:text-copper-500">Jobs</a>
    <span class="mx-1">/</span>
    <a href="#" class="text-copper-700 hover:text-copper-500">Pole Barn</a>
    <span class="mx-1">/</span>
    <span>Main Barn</span>
  </nav>
  <h1 class="text-2xl font-bold tracking-tight text-slate-900">Main Barn</h1>
  <div class="text-sm text-slate-500 mt-1">
    <span>Level 1</span>
    <span class="mx-2">•</span>
    <span>Markup: inherit</span>
  </div>
</div>
```

### Section with Keyboard Hints

```html
<div class="flex items-center gap-3 mb-3">
  <h2 class="text-sm font-semibold tracking-wide uppercase text-slate-700">
    Items
  </h2>
  <span class="font-mono text-xs px-1.5 py-0.5 bg-slate-100 border border-slate-300 rounded text-slate-700">m</span>
  <span class="text-xs text-slate-500">material</span>
  <span class="font-mono text-xs px-1.5 py-0.5 bg-slate-100 border border-slate-300 rounded text-slate-700">l</span>
  <span class="text-xs text-slate-500">labor</span>
</div>
```

### Data Table

```html
<table class="w-full">
  <thead>
    <tr class="border-b border-slate-100">
      <th class="text-left text-xs font-medium tracking-wider uppercase text-slate-500 py-3 px-4">Name</th>
      <th class="text-right text-xs font-medium tracking-wider uppercase text-slate-500 py-3 px-4">Qty</th>
      <th class="text-right text-xs font-medium tracking-wider uppercase text-slate-500 py-3 px-4">Unit</th>
      <th class="text-right text-xs font-medium tracking-wider uppercase text-slate-500 py-3 px-4">Price</th>
      <th class="text-right text-xs font-medium tracking-wider uppercase text-slate-500 py-3 px-4">Total</th>
    </tr>
  </thead>
  <tbody>
    <tr class="border-b border-slate-100">
      <td class="text-sm text-slate-900 py-3 px-4">1/2"x5" Titan anchor bolt</td>
      <td class="text-sm tabular-nums text-slate-900 text-right py-3 px-4">25.00</td>
      <td class="text-sm text-slate-900 text-right py-3 px-4">ea</td>
      <td class="text-sm tabular-nums text-slate-900 text-right py-3 px-4">$5.00</td>
      <td class="text-sm tabular-nums text-slate-900 text-right py-3 px-4">$125.00</td>
    </tr>
  </tbody>
</table>
```

### Total Row

```html
<div class="flex justify-between items-center py-4 border-t border-slate-100">
  <span class="text-sm font-medium text-slate-700">Category Total</span>
  <span class="text-xl font-semibold tabular-nums text-slate-900">$1,775.00</span>
</div>
```

---

## Typography Don'ts

- **Don't use more than 700 weight** — Barlow has 800/900 but they're too heavy for this UI
- **Don't use italic for emphasis** — Use weight or color instead
- **Don't mix font families** — Barlow for everything except keyboard hints
- **Don't use all-caps for body text** — Only for section headers, table headers, badges
- **Don't use letter-spacing on body text** — Only on uppercase elements
- **Don't use monospace for prices** — Use `tabular-nums` with Barlow instead
- **Don't go below 12px (0.75rem)** — Accessibility floor

---

## Revision History

| Date | Version | Notes |
|------|---------|-------|
| 2024-12-27 | 1.0 | Initial typography spec |
