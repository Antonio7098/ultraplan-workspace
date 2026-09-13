# Area reasoning template: Frontend verification and observer isolation

## Purpose

Define the proof that the browser tells the same story as the runtime while remaining passive, bounded, accessible, and reliable under client and evidence failures.

## Required inputs

- accepted Sprint 4 observation and experience decisions;
- completed lifecycle race, replay, and benchmark evidence;
- testing, performance, observability, frontend, accessibility, security, and API contracts.

## Questions to resolve

1. How will tests compare runtime, CLI, DTO, initial HTML, and live browser state?
2. Which controlled schedules cover connection before, during, and after terminal commitment?
3. How are slow readers, disconnects, reconnects, malformed data, unknown versions, and service shutdown tested?
4. Which checks prove keyboard use, focus order, semantic structure, non-colour status, reduced motion, and responsive layouts?
5. What payload, update, render, memory, and runtime-task bounds apply?
6. Which browser checks run routinely, and which belong in the final smoke suite?
7. What failure would reveal a second lifecycle model in the frontend?

## Required output

- cross-interface agreement matrix;
- controlled connection and terminal schedules;
- observer isolation and leak tests;
- schema and transport failure tests;
- accessibility and responsive test plan;
- frontend performance measurements tied to Sprint 3;
- real-browser smoke scenarios;
- Phase 1 rejection criteria.

