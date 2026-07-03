# Frontend Design Reasoning Template

## Purpose

This template gives AI coding agents a structured reasoning process to complete before designing or changing a frontend surface.

The goal is not to make the interface “pretty.” The goal is to make the interface truthful, usable, coherent, accessible, distinctive, and maintainable.

Use this when creating or changing a page, view, workflow, navigation pattern, dashboard, form, landing page, data display, empty state, onboarding flow, or major component. For implementation structure, component ownership, state placement, API boundaries, and code organisation, use the separate frontend architecture/state templates.

This template is specifically for design reasoning:

- What is the user trying to understand or do?
- What should the screen communicate first?
- What should feel easy, obvious, trustworthy, and high quality?
- What visual system should govern the design?
- What must be avoided so the result does not become generic AI UI?
- What design decisions are being made intentionally rather than by default?

---

## 0. Design Change Summary

### Surface name

```text
<page / route / workflow / component / screen / state>
```

### User/product goal

Describe the real outcome this interface should help the user achieve.

```text
<The user needs to...>
```

### Current design task

Describe the concrete design change.

```text
<What needs to be designed, redesigned, or improved?>
```

### Context

```text
- Product area: <...>
- User type/persona: <...>
- Stage of journey: <first-time / returning / expert / admin / buyer / internal user>
- Device context: <desktop / mobile / tablet / responsive>
- Frequency of use: <one-off / occasional / daily / high-volume>
```

### Non-goals

List what should not be solved in this pass.

```text
- <not included>
- <not included>
- <not included>
```

---

## 1. First-Principles UX Breakdown

### Core user job

What is the smallest truthful description of what the user is trying to do?

```text
<The user is trying to decide / find / compare / create / edit / understand / recover from...>
```

### Primary question the screen must answer

```text
<What should the user understand within the first few seconds?>
```

### Primary action

```text
<The main action the design should make obvious>
```

### Secondary actions

```text
- <secondary action>
- <secondary action>
```

### User mental model

How does the user likely think about this task before seeing the interface?

```text
<What concepts, language, order, or grouping will feel natural to them?>
```

### User anxiety / hesitation

What might make the user pause, mistrust, or abandon the flow?

```text
- <unclear value>
- <fear of making mistake>
- <unknown consequence>
- <too much effort>
- <missing context>
```

### Success criteria

```text
The design succeeds if:
- <user can complete the task>
- <user understands the state of the system>
- <user trusts what they are seeing>
- <user knows what to do next>
```

---

## 2. Information Architecture and Content Priority

### Information inventory

List every piece of information the screen could show.

```text
- <information item>
- <information item>
- <information item>
```

### Priority ranking

Rank the information by user value, not by implementation convenience.

| Priority | Information | Why it matters | Shown where/how |
| --- | --- | --- | --- |
| 1 | <...> | <...> | <...> |
| 2 | <...> | <...> | <...> |
| 3 | <...> | <...> | <...> |

### Progressive disclosure

What can be hidden until needed?

```text
Always visible:
- <...>

Visible on interaction / expansion / detail view:
- <...>

Moved out of this surface:
- <...>
```

### Content grouping

What belongs together because the user reasons about it together?

```text
Group 1: <name>
- <item>
- <item>

Group 2: <name>
- <item>
- <item>
```

### Content language

```text
Plain-language labels:
- <technical/internal term> → <user-facing term>

Microcopy needed:
- Empty state: <...>
- Error explanation: <...>
- Confirmation: <...>
- Warning: <...>
```

---

## 3. Visual Hierarchy

### Desired attention order

Describe the exact order in which the user's eye should move.

```text
1. <first thing>
2. <second thing>
3. <third thing>
4. <primary action>
5. <secondary/supporting details>
```

### Hierarchy mechanisms

Which visual tools create that order?

```text
- Size: <...>
- Weight: <...>
- Spacing: <...>
- Contrast: <...>
- Position: <...>
- Motion: <...>
- Color: <...>
```

### Primary CTA treatment

```text
Primary CTA:
- Label: <...>
- Placement: <...>
- Visual weight: <...>
- Disabled/pending state: <...>
- Why it is primary: <...>
```

### Secondary/destructive action treatment

```text
Secondary actions:
- <action>: <treatment>

Destructive actions:
- <action>: <confirmation / placement / visual treatment>
```

### Scan pattern

Which scanning pattern does this surface support?

```text
[ ] F-pattern — dense reading / lists / dashboards / documentation
[ ] Z-pattern — simple marketing / landing / decision page
[ ] Task flow — step-by-step form or wizard
[ ] Hub-and-spoke — dashboard / overview with drill-down
[ ] Comparison — table/cards/side-by-side decision
[ ] Other: <...>

Reason:
<why this pattern fits>
```

---

## 4. Layout and Composition

### Layout model

```text
[ ] Single column
[ ] Two-column editorial/product layout
[ ] Dashboard grid
[ ] Card grid
[ ] Split pane
[ ] Master-detail
[ ] Wizard/stepper
[ ] Table/list with filters
[ ] Bento/product showcase
[ ] Other: <...>
```

### Why this layout fits

```text
<Explain why this layout matches the user's task and content priority.>
```

### Spatial rhythm

Define the spacing logic before styling.

```text
- Page padding: <...>
- Section gap: <...>
- Group gap: <...>
- Element gap: <...>
- Dense areas allowed? <yes/no/where>
```

### Alignment and grid

```text
- Main content width: <...>
- Grid columns: <...>
- Alignment rule: <left / centered / mixed / baseline>
- Breakpoints affected: <...>
```

### Whitespace decision

```text
Where should whitespace create calm and clarity?
- <...>

Where should density be preserved because the user needs comparison or scanning?
- <...>
```

---

## 5. Interaction and Flow Reasoning

### User flow

```text
1. User arrives and sees <...>
2. User understands <...>
3. User interacts with <...>
4. System responds by <...>
5. User reaches <...>
```

### Interaction inventory

```text
- Click/tap: <...>
- Hover: <...>
- Focus: <...>
- Keyboard: <...>
- Drag/drop: <...>
- Search/filter/sort: <...>
- Pagination/infinite scroll: <...>
```

### Feedback design

For every interaction, define how the interface responds.

| Interaction | Immediate feedback | Final feedback | Failure feedback |
| --- | --- | --- | --- |
| <...> | <...> | <...> | <...> |

### Friction check

```text
Can the task be completed with fewer steps?
[ ] Yes — simplify: <...>
[ ] No — current steps are necessary because: <...>

Is any required input avoidable?
[ ] Yes — remove/default/infer: <...>
[ ] No — reason: <...>
```

### Confirmation and reversibility

```text
- Can the user undo? <yes/no>
- Is confirmation needed? <yes/no>
- Is the consequence clear before action? <yes/no>
- Is the result clear after action? <yes/no>
```

---

## 6. State Design at the UI Level

This section describes what the user sees. Use the separate state placement template to decide where state belongs in code.

### Screen states required

```text
[ ] Initial/default
[ ] Loading
[ ] Skeleton
[ ] Empty
[ ] Partial data
[ ] Error
[ ] Permission denied
[ ] Offline
[ ] Pending/submitting
[ ] Success
[ ] Disabled
[ ] Read-only
[ ] Validation error
```

### State-by-state design

| State | User question | UI answer | Action available |
| --- | --- | --- | --- |
| Loading | Is something happening? | <...> | <...> |
| Empty | Why is there nothing here? | <...> | <...> |
| Error | What went wrong and can I recover? | <...> | <...> |
| Success | Did it work and what next? | <...> | <...> |

### Edge states

```text
- Very long content: <...>
- Very short content: <...>
- Missing fields: <...>
- Slow network: <...>
- Duplicate items: <...>
- Conflicting statuses: <...>
```

---

## 7. Forms and Input Design

Complete this section only if the design includes user input.

### Form purpose

```text
<What decision, data, or action does this form support?>
```

### Field necessity check

For each field, ask whether the product genuinely needs it now.

| Field | Required? | Why needed? | Can it be inferred/defaulted? |
| --- | --- | --- | --- |
| <...> | <yes/no> | <...> | <...> |

### Field grouping and order

```text
Group 1:
- <field>
- <field>

Group 2:
- <field>
- <field>
```

### Validation behaviour

```text
- Validate on blur: <...>
- Validate on submit: <...>
- Inline error copy: <...>
- Summary error copy: <...>
- Recovery guidance: <...>
```

### Trust and risk

```text
Does the user need reassurance before submitting?
[ ] No
[ ] Yes — reassurance needed: <privacy / consequence / reversibility / pricing / permissions>
```

---

## 8. Data Display and Decision Support

Complete this section for dashboards, tables, analytics, search results, lists, and comparison screens.

### User decision

```text
<What decision should the data help the user make?>
```

### Data hierarchy

```text
Primary metric/data:
- <...>

Supporting data:
- <...>

Diagnostic/detail data:
- <...>
```

### Comparison needs

```text
Does the user need to compare items?
[ ] No
[ ] Yes — compare by: <status / time / score / owner / value / category>
```

### Sorting/filtering/search

```text
- Default sort: <...>
- Filters needed: <...>
- Search needed: <yes/no>
- Saved views needed: <yes/no>
```

### Data confidence

```text
Does the interface need to communicate uncertainty, freshness, or source?
[ ] No
[ ] Yes — show: <timestamp / source / confidence / sample size / caveat>
```

---

## 9. Brand, Tone, and Visual Identity

### Intended feeling

Choose a few, then explain how the interface earns them.

```text
[ ] Calm
[ ] Fast
[ ] Premium
[ ] Technical
[ ] Trustworthy
[ ] Playful
[ ] Editorial
[ ] Minimal
[ ] Dense/professional
[ ] Bold
[ ] Human
[ ] Other: <...>

Reason:
<...>
```

### Brand physics

Define concrete visual rules rather than vague style words.

```text
- Typography personality: <...>
- Type scale: <...>
- Spacing rhythm: <...>
- Border radius: <...>
- Shadow/elevation: <...>
- Icon style: <...>
- Illustration/image style: <...>
- Surface/background treatment: <...>
```

### Avoiding generic AI UI

Check for common generic defaults.

```text
- [ ] Default purple/indigo gradient without brand reason
- [ ] Generic rounded SaaS cards everywhere
- [ ] Same hero + three feature cards pattern
- [ ] Stock icons with no visual system
- [ ] Weak typography hierarchy
- [ ] Random glassmorphism/neubrutalism without product reason
- [ ] Too many gradients, glows, and vague decorative blobs
- [ ] Layout chosen because it is common, not because it fits the task
- [ ] Empty marketing adjectives instead of specific product value
```

### Distinctiveness decision

```text
What will make this interface feel specific to this product?
- <specific design choice>
- <specific design choice>
- <specific design choice>
```

### Restraint decision

```text
What tempting visual idea are we intentionally not using?
<...>

Why?
<...>
```

---

## 10. Typography System

### Typography role map

| Role | Example | Size/weight/line-height | Reason |
| --- | --- | --- | --- |
| Page title | <...> | <...> | <...> |
| Section title | <...> | <...> | <...> |
| Body | <...> | <...> | <...> |
| Label | <...> | <...> | <...> |
| Metadata | <...> | <...> | <...> |
| CTA | <...> | <...> | <...> |

### Readability check

```text
- Body text readable at target viewport? <yes/no>
- Line length appropriate? <yes/no>
- Line height comfortable? <yes/no>
- Important labels scannable? <yes/no>
- Dense data still legible? <yes/no>
```

### Typography trade-off

```text
<What are we prioritising: elegance, density, clarity, personality, speed, authority?>
```

---

## 11. Color, Contrast, and Semantic Meaning

### Color roles

| Role | Token/intention | Used for | Notes |
| --- | --- | --- | --- |
| Background | <...> | <...> | <...> |
| Surface | <...> | <...> | <...> |
| Primary | <...> | <...> | <...> |
| Accent | <...> | <...> | <...> |
| Success | <...> | <...> | <...> |
| Warning | <...> | <...> | <...> |
| Error | <...> | <...> | <...> |
| Muted | <...> | <...> | <...> |

### Color hierarchy

```text
What should color draw attention to?
- <...>

What should color intentionally de-emphasize?
- <...>
```

### Accessibility check

```text
- Text contrast checked? <yes/no>
- UI control contrast checked? <yes/no>
- Focus state visible? <yes/no>
- Meaning does not rely on color alone? <yes/no>
- Dark/light mode implications considered? <yes/no>
```

---

## 12. Motion and Microinteractions

### Motion purpose

Only add motion when it clarifies state, continuity, feedback, or hierarchy.

```text
Motion is used to:
[ ] Confirm action
[ ] Show continuity between states
[ ] Direct attention
[ ] Reduce perceived wait
[ ] Communicate hierarchy
[ ] Add brand character
[ ] Not used in this design
```

### Microinteraction map

| Moment | Trigger | Feedback | Duration/easing | Reduced motion alternative |
| --- | --- | --- | --- | --- |
| <...> | <...> | <...> | <...> | <...> |

### Motion restraint

```text
What should not animate?
<...>

Why?
<avoid distraction / performance / accessibility / unnecessary polish>
```

---

## 13. Responsive and Adaptive Design

### Breakpoint behaviour

| Viewport | Layout | Navigation | Content priority changes |
| --- | --- | --- | --- |
| Mobile | <...> | <...> | <...> |
| Tablet | <...> | <...> | <...> |
| Desktop | <...> | <...> | <...> |
| Wide | <...> | <...> | <...> |

### Mobile task check

```text
Can the primary task be completed comfortably on mobile?
[ ] Yes
[ ] No — adaptation needed: <...>
```

### Touch target check

```text
- Primary controls large enough? <yes/no>
- Spacing prevents mis-taps? <yes/no>
- Sticky actions needed? <yes/no>
- Horizontal scroll avoided unless intentional? <yes/no>
```

---

## 14. Accessibility and Inclusive Design

### Keyboard and focus

```text
- Logical tab order: <...>
- Initial focus behaviour: <...>
- Focus after action: <...>
- Escape/back behaviour: <...>
```

### Screen reader semantics

```text
- Landmark structure: <...>
- Headings: <...>
- Labels: <...>
- Status announcements: <...>
- Error announcements: <...>
```

### Inclusive content

```text
- Plain language used? <yes/no>
- Jargon explained? <yes/no>
- Error copy avoids blame? <yes/no>
- Time/date/currency formats clear? <yes/no>
```

### Accessibility risks

```text
- <risk>: <mitigation>
- <risk>: <mitigation>
```

---

## 15. Trust, Safety, Privacy, and Consequence Design

### Trust signals

```text
Does the user need trust signals?
[ ] No
[ ] Yes — use: <source, timestamp, permission explanation, customer proof, audit trail, preview, undo>
```

### Consequence clarity

```text
Before the user acts, do they understand:
- What will happen? <yes/no>
- Who/what is affected? <yes/no>
- Whether it can be undone? <yes/no>
- Whether data leaves the system? <yes/no>
```

### Sensitive data display

```text
- Sensitive information shown? <yes/no>
- Masking/truncation needed? <yes/no>
- Permission-based visibility needed? <yes/no>
- Copy-to-clipboard risk? <yes/no>
```

---

## 16. Performance Perception

This is about perceived frontend experience, not low-level code optimisation.

### Perceived speed

```text
- What appears immediately? <...>
- What can stream/progressively reveal? <...>
- What needs a skeleton? <...>
- What can wait until interaction? <...>
```

### Layout stability

```text
- Images/media dimensions reserved? <yes/no>
- Skeleton matches final layout? <yes/no>
- Late-loading content avoids jump? <yes/no>
```

### Density/performance trade-off

```text
Are we showing too much at once?
[ ] No
[ ] Yes — reduce/chunk/virtualize/progressively disclose: <...>
```

---

## 17. Design System Fit

### Existing patterns to reuse

```text
- <pattern/component/token>: <why>
- <pattern/component/token>: <why>
```

### New pattern pressure

```text
Is a new design pattern needed?
[ ] No — existing pattern fits
[ ] Partially — extend existing pattern
[ ] Yes — new pattern is justified
```

### If new pattern, why is it earned?

```text
- [ ] Repeated across multiple surfaces
- [ ] Existing pattern creates poor UX
- [ ] New behaviour cannot be expressed cleanly with current system
- [ ] Product needs a stronger branded pattern
- [ ] Accessibility requires a different pattern
- [ ] Other: <...>
```

### System risk

```text
Could this create visual inconsistency or one-off design debt?
[ ] No
[ ] Yes — mitigation: <document token/pattern/usage rule>
```

---

## 18. Options Considered

Consider at least two options unless the design change is trivial.

### Option A: Minimal direct design

```text
Description:
<...>

Pros:
- <...>

Cons:
- <...>

Risks:
- <...>
```

### Option B: Stronger UX restructuring

```text
Description:
<...>

Pros:
- <...>

Cons:
- <...>

Risks:
- <...>
```

### Option C: More distinctive/branded design

Only include if brand expression is important.

```text
Description:
<...>

Pros:
- <...>

Cons:
- <...>

Risks:
- <...>
```

### Chosen option

```text
Chosen option: <A/B/C>
Reason:
<Why this is the most honest design for the current user goal.>
```

---

## 19. Design QA Checklist

### Clarity

```text
- [ ] The screen has one obvious primary purpose
- [ ] The primary action is visually clear
- [ ] The user can understand the current state
- [ ] Labels use user language rather than internal language
- [ ] Empty/error/success states are designed, not forgotten
```

### Hierarchy

```text
- [ ] Important information is visually dominant
- [ ] Secondary information is available but not competing
- [ ] Spacing groups related things together
- [ ] The page is scannable
```

### Interaction

```text
- [ ] Every interactive element has feedback
- [ ] Loading/pending states prevent confusion
- [ ] Destructive actions are clear and recoverable where possible
- [ ] Keyboard interaction is considered
```

### Accessibility

```text
- [ ] Contrast is sufficient
- [ ] Focus states are visible
- [ ] Form fields and errors are labelled
- [ ] Motion is reduced when requested
- [ ] Meaning is not conveyed by color alone
```

### Distinctiveness

```text
- [ ] Visual choices come from product/brand/user needs
- [ ] Generic AI design defaults have been challenged
- [ ] Typography, spacing, and color feel intentional
- [ ] Decorative elements have a reason
```

---

## 20. Final Design Decision

```text
Decision: <Proceed / Redesign / Split design from implementation / Block>

Reason:
<short explanation>

Main user benefit:
<...>

Main trade-off:
<...>

Design complexity introduced:
<none / low / medium / high>

Design complexity removed:
<none / low / medium / high>

Follow-up needed in code-focused templates:
[ ] Frontend architecture / ownership
[ ] Component composition
[ ] State placement
[ ] API/data contract
[ ] Performance implementation
[ ] Testing
```

---

# Final Standard

The agent should not merely ask:

> What would look nice here?

It should ask:

> What does the user need to understand, feel, trust, and do — and what visual system best supports that?

If the answer is simple, design simply.

If the screen needs stronger hierarchy, interaction, or brand expression, make those decisions deliberately.

If the issue is code ownership, state placement, or component architecture, stop and use the code-focused frontend reasoning templates instead.
