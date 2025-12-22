# Skalkaho Logo Specification

Two overlapping mountain peaks with curved bases. Back peak recedes with reduced opacity, front peak comes forward.

## Versions

### Primary (Light Backgrounds)
Forest back peak at 60% opacity, Copper front peak at full opacity.

```svg
<svg viewBox="0 0 56 44" fill="none" xmlns="http://www.w3.org/2000/svg">
  <path d="M 0 44 Q 9 38, 18 4 Q 27 38, 36 44 Z" fill="#2d5a47" opacity="0.6"/>
  <path d="M 20 44 Q 29 38, 38 4 Q 47 38, 56 44 Z" fill="#a35a2a"/>
</svg>
```

### Dark Backgrounds
White monochrome. Back peak at 35% opacity, front peak at full opacity.

```svg
<svg viewBox="0 0 56 44" fill="none" xmlns="http://www.w3.org/2000/svg">
  <path d="M 0 44 Q 9 38, 18 4 Q 27 38, 36 44 Z" fill="#ffffff" opacity="0.35"/>
  <path d="M 20 44 Q 29 38, 38 4 Q 47 38, 56 44 Z" fill="#ffffff"/>
</svg>
```

## Wordmark

Always inline (horizontal), never stacked. Logo mark left, text right.

- Light backgrounds: Slate 900 (`#1a1d23`) text
- Dark backgrounds: White (`#ffffff`) text

## Colors

| Name | Hex | Usage |
|------|-----|-------|
| Forest | `#2d5a47` | Back peak (light bg) |
| Copper | `#a35a2a` | Front peak (light bg) |
| White | `#ffffff` | Both peaks (dark bg) |

## Clear Space

Minimum clear space around logo: height of one peak (the "4" in the viewBox).

## Minimum Size

- Icon only: 16px height
- With wordmark: 24px height

## Don'ts

- Don't use the two-color version on dark backgrounds
- Don't stack the wordmark vertically
- Don't change the opacity values
- Don't separate the peaks
- Don't rotate or skew