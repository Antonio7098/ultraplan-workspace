# Accessibility Contract

## Purpose

This contract defines the minimum accessibility requirements for frontend and UI changes.

It governs:

- semantic HTML and ARIA
- keyboard interaction
- focus management
- forms and errors
- color/contrast and visual states
- motion
- screen reader support
- accessible testing and review

## Scope

This contract applies when a change:

- adds or changes UI components, pages, flows, forms, dialogs, navigation, charts, or interactive controls
- changes design-system primitives or patterns
- changes loading/error/empty states
- adds animation or motion

## Requirement Index

| ID | Title | Applies To | Severity If Violated |
| --- | --- | --- | --- |
| A11Y-SEM-001 | UI must use semantic structure before ARIA | all UI | High |
| A11Y-KBD-001 | Interactive UI must work by keyboard | controls/dialogs/nav | Blocker |
| A11Y-FOCUS-001 | Focus management must be deliberate | route/dialog/dynamic UI | High |
| A11Y-FORM-001 | Forms must expose labels, errors, and status accessibly | forms | High |
| A11Y-VIS-001 | Visual states must not rely only on color | UI states | Medium |
| A11Y-MOTION-001 | Motion must respect reduced-motion needs | animation | Medium |
| A11Y-TEST-001 | Critical UI must have accessibility checks | frontend tests/review | Medium |

## Core Principles

1. Prefer semantic HTML over custom ARIA.
2. If it is clickable, it must be keyboard-operable.
3. Focus is part of the user journey.
4. Errors, loading, and success states must be perceivable.
5. Color cannot be the only carrier of meaning.
6. Custom components inherit native accessibility responsibilities.
7. Accessibility belongs in design-system primitives, not just individual pages.

## Requirements

### A11Y-SEM-001: UI Must Use Semantic Structure Before ARIA

**Rule**
Use native semantic elements and accessible names before adding custom roles or ARIA.

**Required**
- use buttons for actions, links for navigation, labels for inputs, headings for structure
- provide accessible names for icons and controls
- preserve logical heading and landmark structure

**Forbidden**
- clickable div/span controls when native elements would work
- ARIA roles that contradict native semantics
- icon-only controls without accessible labels

**Evidence**
- component review shows semantic elements and accessible names

### A11Y-KBD-001: Interactive UI Must Work By Keyboard

**Rule**
All interactive controls and flows must be usable without a mouse.

**Required**
- support Tab/Shift+Tab navigation
- support Enter/Space activation where expected
- support Escape/arrow-key behaviour for common patterns where relevant
- avoid keyboard traps

**Forbidden**
- mouse-only controls
- hidden focusable elements in inactive UI
- overlays/dialogs that trap users without exit

**Evidence**
- tests or manual review cover keyboard path for important UI

### A11Y-FOCUS-001: Focus Management Must Be Deliberate

**Rule**
Focus must move predictably during route changes, dialogs, errors, and dynamic content updates.

**Required**
- move focus into modal/dialog content and restore it on close
- make validation errors discoverable
- avoid unexpected focus stealing
- keep focus visible

**Forbidden**
- opening dialogs without focus management
- removing focused elements without recovery
- hiding focus outlines without replacement

**Evidence**
- modal/route/form tests or review notes cover focus behaviour

### A11Y-FORM-001: Forms Must Expose Labels, Errors, And Status Accessibly

**Rule**
Form controls must have labels, validation, error messaging, and submission state that assistive technology can perceive.

**Required**
- associate labels with inputs
- associate errors with fields
- expose pending/success/failure state accessibly
- preserve user input on validation failure where practical

**Forbidden**
- placeholder-only labels
- errors shown only visually with no programmatic association
- disabling submit without explaining why where user action is needed

**Evidence**
- form tests cover validation and error state rendering

### A11Y-VIS-001: Visual States Must Not Rely Only On Color

**Rule**
Important state must be communicated through text, icons, shape, position, or programmatic attributes in addition to color.

**Required**
- provide visible and accessible state indicators
- keep contrast suitable for text and controls
- distinguish disabled, focused, selected, invalid, and loading states clearly

**Forbidden**
- color-only error/success distinction
- low-contrast text or controls in normal states

**Evidence**
- design-system components preserve visible state affordances

### A11Y-MOTION-001: Motion Must Respect Reduced-Motion Needs

**Rule**
Animations and transitions must not block usability and should respect reduced-motion preferences.

**Required**
- avoid essential information conveyed only through motion
- respect reduced-motion settings for non-essential animation
- avoid flashing or aggressive motion

**Forbidden**
- motion required to understand or operate the UI
- distracting loading animations with no reduced-motion treatment

**Evidence**
- motion utilities/components include reduced-motion handling where relevant

### A11Y-TEST-001: Critical UI Must Have Accessibility Checks

**Rule**
Important UI primitives, forms, dialogs, and flows must have automated or manual accessibility review.

**Required**
- test design-system primitives and complex interactions
- include accessible queries in component tests where practical
- run automated accessibility checks where supported

**Forbidden**
- relying only on snapshot tests for accessible behaviour
- shipping custom controls without keyboard/focus review

**Evidence**
- tests or review notes cover accessibility-critical behaviour

## Review Rejection Criteria

Reject a UI change if it:

- adds mouse-only controls
- uses custom controls without keyboard/focus behaviour
- hides focus outlines without replacement
- adds forms without labels or accessible errors
- uses color-only state communication
- opens dialogs without focus handling
- puts business-critical UI behind inaccessible interactions

---