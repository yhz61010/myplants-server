# Design System: The Digital Conservatory

## 1. Overview & Creative North Star
This design system is anchored by the Creative North Star: **"The Digital Conservatory."** 

Unlike a standard social network or e-commerce grid, this system treats the interface as an architectural glasshouse. It balances the structural precision of modern typography with the organic fluidity of nature. We move beyond "template" layouts by embracing **intentional asymmetry, overlapping editorial elements, and depth created through light rather than lines.** 

The goal is to evoke a sense of serenity and expertise. Every interaction should feel like a quiet walk through a botanical garden—intentional, breathable, and deeply tactile.

---

### 2. Colors & Tonal Architecture
We utilize a palette of "Living Greens" and "Earthy Neutrals" to create a grounded, premium atmosphere. 

*   **Primary (`#006e1c`):** Our "Chlorophyll." Used sparingly for high-intent actions and brand moments.
*   **Secondary (`#7a5649`):** Our "Terra Cotta." Provides a sophisticated, earthy anchor to the vibrant greens.
*   **Surface Hierarchy (The Nesting Rule):**
    *   **The "No-Line" Rule:** 1px solid borders are strictly prohibited for sectioning. Contrast is achieved through background shifts. For example, a `surface-container-low` section sits directly on a `surface` background to define its boundary.
    *   **Tonal Nesting:** UI depth is built by stacking surface tiers. A `surface-container-lowest` card should be placed upon a `surface-container-low` section. This creates a "soft lift" that feels organic rather than mechanical.
*   **Glass & Gradient Rule:** 
    *   For floating navigation or featured modals, use Glassmorphism: `surface` color at 70% opacity with a `24px` backdrop blur.
    *   **Signature Textures:** Main CTAs or Hero backgrounds should use a subtle linear gradient (135°) from `primary` to `primary_container`. This mimics the natural variegation found in a leaf.

---

### 3. Typography: The Editorial Voice
We use a high-contrast typographic pairing to establish an "Editorial Authority."

*   **The Display Voice (Plus Jakarta Sans):** Used for `display` and `headline` tiers. It is modern, geometric, and spacious. Use `display-lg` (3.5rem) with tight letter-spacing (-0.02em) for hero moments to create a high-end magazine feel.
*   **The Functional Voice (Manrope):** Used for `title`, `body`, and `label` tiers. Manrope’s open counters and modern proportions ensure maximum readability for plant care guides and community threads.
*   **Hierarchy as Identity:** Always lead with a large `headline-lg` or `display-sm` to anchor a page. Use `label-md` in all-caps with `0.05em` tracking for category tags (e.g., "FERNS," "SUCCULENTS") to provide a sophisticated, curated look.

---

### 4. Elevation & Depth: Layering Principle
We reject the "heavy shadow" aesthetic. Depth in this system is an atmospheric quality.

*   **Tonal Layering:** Instead of shadows, use the transition from `surface-container-lowest` to `surface-container-highest` to indicate importance.
*   **Ambient Shadows:** Where floating elements (like FABs or Modals) are necessary, shadows must be "Ambient." Use a blur radius of `32px`, a spread of `-4px`, and an opacity of `6%`. The color of the shadow must be a tinted version of `on-surface` (dark green-grey) rather than pure black.
*   **The "Ghost Border" Fallback:** If a layout requires a container edge for accessibility, use the `outline-variant` token at **15% opacity**. It should be felt, not seen.
*   **Corner Radii:** Apply `xl` (1.5rem) to primary containers and `md` (0.75rem) to smaller nested elements. This "nested rounding" mimics the soft edges of a pebble.

---

### 5. Components & Interface Elements

*   **Buttons:** 
    *   *Primary:* Use `primary` background with `on-primary` text. Shape: `full` (pill). 
    *   *Tertiary:* Use `on-surface` text with no background. On hover, apply a `surface-container-high` background.
*   **Cards (The Specimen Card):** 
    *   Cards must never have borders. Use `surface-container-lowest` for the card body. 
    *   Images should use a `xl` corner radius. 
    *   **Asymmetric Layout:** Place the plant name (`title-lg`) overlapping the edge of the image slightly to break the "boxed" feel.
*   **Lists:** 
    *   **The No-Divider Rule:** Forbid the use of horizontal lines between list items. Use `24px` of vertical white space (from the Spacing Scale) or a subtle shift to `surface-container-low` on hover to separate content.
*   **Inputs:** 
    *   Use a `surface-container-high` fill with a `none` border. On focus, transition to a `ghost-border` of the `primary` color.
*   **Chips (The Leaf Tag):** 
    *   Selection chips should use `secondary-fixed-dim` for a muted, earthy feel. Use `full` roundedness.

---

### 6. Do’s and Don’ts

#### **Do:**
*   **Do** use ample white space. If you think there is enough space, add 16px more.
*   **Do** use "Editorial Breaks." Occasionally center a piece of text or an image to break the left-aligned grid.
*   **Do** ensure photography is high-quality, featuring natural lighting and "macro" details of plants.

#### **Don’t:**
*   **Don’t** use pure black (`#000000`). Use `on-surface` or `inverse-surface` for text and shadows.
*   **Don’t** use 1px dividers. If you need to separate content, use a background color change or whitespace.
*   **Don’t** use "default" system fonts. Stick strictly to the Plus Jakarta Sans/Manrope pairing.
*   **Don’t** use harsh, high-contrast shadows. If the shadow is the first thing you see, it’s too dark.