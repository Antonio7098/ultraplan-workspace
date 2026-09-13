# Area reasoning template: Performance evidence and regression policy

## Purpose

Define how Aren records, compares, interprets, and retains performance results. This area owns evidence identity, schema, variance, failure truth, and the path toward future regression gates.

## Required inputs

- accepted Sprint 3 workload and measurement decisions;
- performance engineering mandate;
- current repository and continuous-integration capabilities;
- selected evaluation, observability, testing, and documentation evidence.

## Questions to resolve

1. Which source, environment, scenario, parameter, and sample identities make one result attributable?
2. Which detailed facts are machine-readable and which summaries belong in human reports?
3. How are partial, cancelled, timed-out, or malformed measurements represented without false success?
4. How much repetition characterizes ordinary variance on the comparison host?
5. Which comparisons use `benchstat`, and how are scenario-runner distributions compared?
6. What evidence belongs in local output, continuous-integration artifacts, and the promoted baseline?
7. What evidence would justify a later warning or blocking threshold?

## Required output

- versioned result and baseline schemas;
- terminal measurement outcomes and stable failure meanings;
- statistical comparison procedure;
- variance report method;
- retention and artifact policy;
- pull-request reporting policy;
- future regression-gate entry conditions;
- baseline interpretation limits.

