# State Machine Reference for Hiri

> Created 2026-03-03. Based on analysis of QOR's `transition` package and Skalkaho's job lifecycle patterns.
> Purpose: Reference guide for implementing a state machine in Hiri.

## Why a State Machine?

1. **Prevent invalid transitions** — Can't go from `accepted` back to `draft`, can't skip `sent` to reach `accepted`
2. **Enforce business rules per state** — Block editing on accepted records, block deletion of sent records
3. **Enable all lifecycle states** — Make `rejected` and `expired` states reachable with proper handlers
4. **Synchronize parent ↔ child status** — Parent status auto-updates when child entity transitions happen
5. **Audit trail** — Log who changed status and when (useful for dispute resolution)
6. **Testable domain logic** — State transitions become unit-testable business rules

## Recommended Architecture (Go)

### Pattern: Function-based, not framework-based

QOR uses a heavy reflection-based state machine with GORM callbacks. **Don't copy that.** A lightweight, typed approach is better:

```go
// domain/status.go

type Status string

const (
    StatusDraft    Status = "draft"
    StatusSent     Status = "sent"
    StatusAccepted Status = "accepted"
    StatusRejected Status = "rejected"
    StatusExpired  Status = "expired"
)

// validTransitions defines the state graph
var validTransitions = map[Status][]Status{
    StatusDraft:    {StatusSent, StatusExpired},
    StatusSent:     {StatusAccepted, StatusRejected, StatusExpired},
    StatusRejected: {StatusDraft, StatusExpired},
    // StatusAccepted: terminal (no outgoing transitions)
    // StatusExpired:  terminal
}

// CanTransition checks if a transition is valid.
func CanTransition(from, to Status) bool {
    allowed, ok := validTransitions[from]
    if !ok {
        return false
    }
    for _, s := range allowed {
        if s == to {
            return true
        }
    }
    return false
}

// Transition validates and returns an error if invalid.
func Transition(from, to Status) error {
    if !CanTransition(from, to) {
        return fmt.Errorf("invalid transition from %q to %q", from, to)
    }
    return nil
}

// Helpers for business rule enforcement
func AllowsEditing(s Status) bool  { return s == StatusDraft }
func AllowsDeletion(s Status) bool { return s == StatusDraft || s == StatusRejected }
```

### Why this approach over QOR's

| QOR's approach | This approach |
|---------------|---------------|
| Reflection-based, any type | Typed `Status` constants |
| GORM callbacks for persistence | Explicit handler calls |
| Enter/Exit hooks on states | Business logic in handlers |
| Before/After hooks on events | Validation before, side-effects after |
| Automatic audit logging | Explicit audit calls (when needed) |

The QOR approach is powerful for a framework that manages many different models. For a single app with known models, explicit typed functions are simpler, more debuggable, and have zero dependencies.

## State Graph for a Quote Lifecycle

```
                    ┌──────────┐
                    │  draft   │ (initial state)
                    └────┬─────┘
                         │ send
                    ┌────▼─────┐
               ┌────│   sent   │────┐
               │    └──────────┘    │
        accept │                    │ reject
        ┌──────▼───┐          ┌─────▼────┐
        │ accepted  │          │ rejected │
        └──────────┘          └─────┬────┘
                                    │ revise
                              ┌─────▼────┐
                              │  draft   │ (back to draft)
                              └──────────┘

        Any non-accepted state ──expire──▶ expired
```

## Implementation Checklist for Hiri

### Domain Layer
- [ ] Define `Status` type with constants
- [ ] Define `validTransitions` map
- [ ] Implement `CanTransition()`, `Transition()`, `AllowsEditing()`, `AllowsDeletion()`
- [ ] Write unit tests for every valid transition (should pass)
- [ ] Write unit tests for every invalid transition (should error)
- [ ] Write unit tests for terminal states (accepted, expired have no outgoing transitions)

### Handler Layer
- [ ] Wire `Transition()` into handlers that change status
- [ ] Add guards: `AllowsEditing()` before edit operations
- [ ] Add guards: `AllowsDeletion()` before delete operations
- [ ] Ensure parent entity status updates when child status changes (e.g., estimate accepted → job accepted)
- [ ] Add reject handler if needed

### Database Layer
- [ ] Ensure CHECK constraint matches all valid statuses
- [ ] Have an `UpdateStatus` query that only changes the status field
- [ ] Consider: audit log table for status changes (can defer)

### Things to Watch For (Edge Cases from QOR Analysis)

1. **Create vs Update should both use transactions** — QOR had a bug where Create wasn't wrapped in a transaction but Update was. Always use transactions for multi-step operations.
2. **Status sync between related entities** — If estimate changes status, job should too. If this update fails, log but don't fail the primary operation (non-fatal). Alternatively, use a DB trigger.
3. **Terminal states must be truly terminal** — Once accepted/signed, there should be no way back. The state machine enforces this, but also add DB-level guards if possible.
4. **Don't use reflection** — QOR's biggest weakness. Use typed constants and explicit function calls.
5. **`expires_at` needs a background job or cron** — The state machine can validate the transition, but something needs to trigger it. Consider a periodic check or lazy evaluation (check on read).

## Lessons from Skalkaho's Current (Non-State-Machine) Implementation

What works well today:
- DB CHECK constraint catches invalid values at the persistence layer
- Status filtering in list views works
- Status badges in UI are already wired for all 5 states
- `UpdateJobStatus` query exists and is efficient (single column update)

What's broken without a state machine:
- `rejected` and `expired` states are unreachable (no handlers)
- Job and estimate status can desync (estimate sent but job stays draft)
- No validation prevents invalid transitions (accepted → draft)
- Only one business rule enforces status (`if job.Status != "draft"` for client editing)
- Line items, categories, and surcharges can be edited on accepted jobs
