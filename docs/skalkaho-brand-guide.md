# Skalkaho — Brand Guide

**Direction: Dark Authority**
*For use by the coding agent implementing the UI. All decisions here are final unless explicitly overridden by the developer.*

---

## 1. Design Philosophy

Skalkaho is a professional tool for contractors who take their work seriously. The interface draws from the same Bitterroot Valley palette as Western Skies Contracting — granite, sage, sunburst gold — but inverted into a dark, focused workspace. Where the marketing site is warm limestone and open air, the app is ink-dark and precise.

**Tone:** Quiet confidence. A tool that stays out of your way. Dark without being theatrical — no glowing accents, no glassmorphism, no gaming-dashboard energy.

**One rule to carry through every component:** Dark backgrounds, warm accents, sharp edges.

**Relationship to Western Skies brand:** Skalkaho shares the color tokens and typography of the parent brand but applies them in dark-mode context. The login page is the bridge between the two — familiar palette, different mood.

---

## 2. Color Tokens

Inherited from the Western Skies brand guide. Applied in dark-mode context.

```css
:root {
  --color-ink:        #1E2120;  /* Primary background. Near-black with green undertone. */
  --color-granite:    #3D4A52;  /* Secondary surfaces, primary buttons, panels. */
  --color-sage:       #7A8C6E;  /* Field labels, secondary accents. Used sparingly. */
  --color-sunburst:   #C49A3C;  /* Eyebrow accents, focus states, hover highlights. ONE role per component. */
  --color-limestone:  #F2EDE4;  /* NOT used as background in-app. Reserved for rare light contexts. */
  --color-concrete:   #D4CEC6;  /* Body text on dark backgrounds, subtitle text, borders. */
  --color-white:      #FAFAF8;  /* Headings on dark backgrounds, primary text, input values. */
  --color-body-text:  #3A3830;  /* Not used in dark mode — too low contrast. */
  --color-muted-text: #7A7468;  /* Not used in dark mode — too low contrast. */
}
```

### Dark Mode Usage Rules

| Token | Use for | Never use for |
|---|---|---|
| `--color-ink` | Page background, input backgrounds | Text on any surface |
| `--color-granite` | Panels, cards, button backgrounds, nav background | Page-level background (too light) |
| `--color-sage` | Field labels, secondary UI labels | Body text, headings, large fills |
| `--color-sunburst` | Eyebrow accents, focus rings, button hover bg, one highlight per component | Multiple competing uses in one view |
| `--color-concrete` | Body copy, subtitles, descriptions on dark bg | Headings (not bright enough) |
| `--color-white` | Headings, input values, primary text on dark bg | Borders, labels (too loud) |

### Borders & Dividers on Dark Backgrounds

```css
/* Default border */
border-color: rgba(255, 255, 255, 0.08);

/* Emphasized border (e.g., active input) */
border-color: var(--color-sunburst);

/* Error accent */
border-color: #B43C3C;
```

---

## 3. Typography

Shared with Western Skies. Same fonts, same roles.

### Font Families

```css
:root {
  --font-display: 'Cormorant Garamond', Georgia, serif;
  --font-ui:      'Barlow Condensed', system-ui, sans-serif;
  --font-body:    'Barlow', system-ui, sans-serif;
}
```

**Google Fonts import (place in `<head>`):**

```html
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link href="https://fonts.googleapis.com/css2?family=Cormorant+Garamond:ital,wght@0,400;0,600;0,700;1,400;1,600;1,700&family=Barlow:wght@300;400;500&family=Barlow+Condensed:wght@400;600;700&display=swap" rel="stylesheet">
```

### Font Roles in Dark Mode

| Font | Role | Style |
|---|---|---|
| Cormorant Garamond | Page titles, hero headlines | 700 weight, mixed case, `--color-white`. Never all-caps. |
| Barlow Condensed | Eyebrow labels, field labels, buttons, nav links | Uppercase, letter-spacing 0.16–0.24em. Labels: 600 weight, `--color-sage`. Buttons: 700 weight. |
| Barlow | Body copy, subtitles, descriptions, input values | 300–400 weight, `--color-concrete` for secondary text, `--color-white` for input values. |

### Type Scale

```css
:root {
  --text-hero:    clamp(36px, 5vw, 52px);
  --text-h2:      clamp(28px, 3.5vw, 44px);
  --text-h3:      clamp(20px, 2.5vw, 26px);
  --text-label:   11px;   /* Always uppercase + tracked */
  --text-btn:     12px;   /* Always uppercase + tracked */
  --text-body:    15px;
  --text-body-lg: 17px;
  --text-small:   13px;
  --text-caption: 10px;
}
```

---

## 4. Component Patterns

### Eyebrow Label

Used above page titles. Sunburst on dark backgrounds.

```css
.eyebrow {
  font-family: var(--font-ui);
  font-size: var(--text-label);
  font-weight: 700;
  letter-spacing: 0.24em;
  text-transform: uppercase;
  color: var(--color-sunburst);
  display: flex;
  align-items: center;
  gap: 12px;
}
.eyebrow::before {
  content: '';
  display: block;
  width: 24px;
  height: 2px;
  background: var(--color-sunburst);
  flex-shrink: 0;
}
```

### Form Inputs

Bottom-border only. Transparent background.

```css
.field-input {
  width: 100%;
  background: transparent;
  border: none;
  border-bottom: 2px solid rgba(255, 255, 255, 0.08);
  padding: 10px 0;
  font-family: var(--font-body);
  font-size: 15px;
  font-weight: 400;
  color: var(--color-white);
  outline: none;
  transition: border-color 0.3s ease;
}
.field-input:focus {
  border-bottom-color: var(--color-sunburst);
}
.field-input::placeholder {
  color: rgba(255, 255, 255, 0.35);
}
```

### Primary Button

```css
.btn-primary {
  font-family: var(--font-ui);
  font-size: var(--text-btn);
  font-weight: 700;
  letter-spacing: 0.16em;
  text-transform: uppercase;
  background: var(--color-granite);
  color: var(--color-white);
  padding: 15px 28px;
  border: 2px solid var(--color-granite);
  cursor: pointer;
  transition: background 0.25s ease, border-color 0.25s ease, color 0.25s ease;
}
.btn-primary:hover {
  background: var(--color-sunburst);
  border-color: var(--color-sunburst);
  color: var(--color-ink);
}
.btn-primary:focus-visible {
  outline: 2px solid var(--color-sunburst);
  outline-offset: 3px;
}
```

### Error Banner

Left-border accent, muted red background.

```css
.error-banner {
  background: rgba(180, 60, 60, 0.12);
  border-left: 3px solid #B43C3C;
  padding: 12px 16px;
  font-family: var(--font-body);
  font-size: 13px;
  font-weight: 400;
  color: #E8A0A0;
}
```

### Accent Rule

Sunburst-to-sage gradient. Use at bottom of pages or as section dividers.

```css
.accent-rule {
  height: 3px;
  background: linear-gradient(
    90deg,
    var(--color-sunburst) 0%,
    var(--color-sage) 55%,
    transparent 100%
  );
}
```

### Grain Texture

Subtle noise overlay for tactile warmth on dark surfaces.

```css
.has-grain::before {
  content: '';
  position: absolute;
  inset: 0;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='200' height='200'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='.75' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='200' height='200' filter='url(%23n)' opacity='.4'/%3E%3C/svg%3E");
  background-size: 200px 200px;
  opacity: 0.03;
  pointer-events: none;
  z-index: 0;
}
```

---

## 5. Motion

### Entrance Animations

Staggered fade-up for content reveals. Expo easing — fast out, smooth deceleration.

```css
@keyframes fadeUp {
  from { opacity: 0; transform: translateY(16px); }
  to   { opacity: 1; transform: translateY(0); }
}
/* Apply with staggered animation-delay: 0.1s, 0.18s, 0.26s, etc. */
/* Easing: cubic-bezier(0.16, 1, 0.3, 1) */
/* Duration: 0.5s */
```

### Transition Defaults

- **Color/background transitions:** 0.25s ease
- **Border transitions:** 0.3s ease
- **Never animate:** width, height, padding, margin — use transform and opacity only
- **Never use:** bounce or elastic easing

---

## 6. Layout Principles

- **No rounded corners** — buttons, cards, inputs are all sharp-edged
- **No drop shadows** on cards — use border or background contrast to separate surfaces
- **Left-align body text** — never center paragraphs
- **Asymmetric layouts preferred** — break the grid intentionally for emphasis
- **Split panels** for focused views (login, onboarding) — brand/context on left, action on right

---

## 7. Anti-Patterns

Do not introduce any of these:

- No glassmorphism, blur effects, or glow borders
- No gradient text
- No cyan-on-dark or purple-to-blue gradients
- No rounded rectangles with generic drop shadows
- No bounce/elastic animations
- No emoji or icon fonts — SVG only if icons are needed
- No centered body copy
- No all-caps Cormorant Garamond
- No barn red, orange, or rust accent tones
- No cards nested inside cards

---

## 8. Established Pages

Pages that have been designed and should be referenced for consistency:

| Page | File | Key Decisions |
|---|---|---|
| Login | `internal/templates/pages/login.html` | Split panel layout, ink background, bottom-border inputs, staggered fade-up entrance, grain overlay, bottom accent rule |

---

## 9. Decisions To Make

*Add items here as design questions arise.*

- [ ] Nav bar pattern for authenticated pages (granite bg with sunburst bottom border, per Western Skies?)
- [ ] Table styling for estimate line items
- [ ] Card pattern for project list / dashboard
- [ ] Modal/dialog styling (or avoid modals per design principles?)
- [ ] Success/warning/info banner variants alongside error
- [ ] Sidebar vs top-nav for app navigation
- [ ] Empty state patterns
