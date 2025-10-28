# Claude Code Rules & Context

## Instructions for Claude

When working on this project, follow these rules:

### 1. Always Start Here
Before writing any code, read:
- `WHY.md` - Understand the purpose
- `DEVELOPMENT_PHILOSOPHY.md` - Understand the approach
- `BDD_GUIDE.md` - Understand how we test

### 2. Question Feature Requests
**Every feature request should be evaluated:**

Ask yourself:
- Does this support discipline or undermine it?
- Would Ed Seykota approve?
- Does it make impulsivity easier or harder?
- Is this solving a real problem or creating complexity?

**If it makes impulsivity easier, push back.** That's your job.

### 3. Write Gherkin First
For any new functionality:

1. Write the Gherkin scenario
2. Get agreement on behavior
3. THEN write code
4. Code matches Gherkin, not the other way around

**No Gherkin = No implementation.**

### 4. Prefer Simple Over Clever
When choosing between:
- Elegant abstraction vs straightforward duplication → Choose straightforward
- Clever one-liner vs verbose clarity → Choose clarity
- Flexible architecture vs rigid simplicity → Choose simplicity
- Generic solution vs specific implementation → Choose specific

**Boring code is good code.**

### 5. Make Constraints Explicit
Every business rule should be:
- Documented in Gherkin
- Enforced in code
- Tested explicitly
- Impossible to bypass

Examples:
- Heat caps → Hard limits, no overrides
- Checklist → Must be GREEN to save
- Impulse timer → Cannot be skipped
- Position sizing → Must use ATR math

### 6. Error Messages Must Teach
Every error should:
- State what's wrong
- Show the actual values
- Show the limit/expectation
- Suggest how to fix it

Bad: "Invalid input"
Good: "Portfolio heat ($425) exceeds cap ($400). Reduce position size or close existing positions."

### 7. Fail Loudly
**Never fail silently.** When something goes wrong:
- Log it verbosely
- Show user-friendly error
- Make it impossible to ignore
- Don't let user proceed

### 8. Context Switching Checklist
When returning to this project after time away:

```
[ ] Re-read WHY.md (5 min)
[ ] Review recent commits (understand what's been done)
[ ] Check open issues (understand what's next)
[ ] Run tests (verify everything works)
[ ] Read relevant Gherkin scenarios (understand expected behavior)

NOW you can code.
```

### 9. Code Review Standards
When reviewing (or generating) code, check:

**Business Logic:**
- [ ] Does it match the Gherkin scenario?
- [ ] Does it support WHY.md principles?
- [ ] Are all edge cases handled?
- [ ] Are constraints enforced?

**Code Quality:**
- [ ] Is it simple and obvious?
- [ ] Are errors handled explicitly?
- [ ] Are names clear and descriptive?
- [ ] Is it well-commented (why, not what)?
- [ ] Are magic numbers named?

**Testing:**
- [ ] Do unit tests exist?
- [ ] Do integration tests exist?
- [ ] Can we demonstrate the Gherkin scenario works?

**Documentation:**
- [ ] Is the behavior clear from reading the code?
- [ ] Are assumptions documented?
- [ ] Are edge cases explained?

### 10. Refactoring Rules
Before refactoring:
- [ ] Tests are passing
- [ ] You understand what it currently does
- [ ] You can articulate why the refactor is better
- [ ] The Gherkin scenarios still apply

After refactoring:
- [ ] Tests still pass
- [ ] Behavior is unchanged
- [ ] Code is simpler or clearer
- [ ] Performance is same or better

**If tests break during refactoring, you changed behavior. Stop and reconsider.**

---

## Project-Specific Patterns

### Position Sizing Always Follows This Pattern:
```
1. Validate inputs (entry, ATR, K must be positive)
2. Calculate stop distance (K × ATR)
3. Calculate initial stop (entry - stop distance)
4. Calculate shares (risk ÷ stop distance, rounded down)
5. Verify actual risk ≤ specified risk
6. Return result with all components
```

Never deviate. This is the Van Tharp method.

### Heat Management Always Follows This Pattern:
```
1. Sum risk across all open positions
2. Add proposed new position risk
3. Compare to portfolio cap (equity × heat_pct)
4. Compare to bucket cap (equity × bucket_pct)
5. Reject if either exceeded
6. Return detailed breakdown
```

Never allow trades that exceed caps. No exceptions.

### Checklist Validation Always Follows This Pattern:
```
1. Count missing checks
2. If 0 missing → GREEN (go)
3. If 1 missing → YELLOW (caution)
4. If 2+ missing → RED (no-go)
5. Return banner, missing count, missing items
6. Start impulse timer only on GREEN
```

Never skip checklist. Never allow save without GREEN.

### Data Flow Always Follows This Pattern:
```
Excel (input)
  → VBA (call backend)
    → Backend (calculate)
      → Database (persist)
    ← Backend (return JSON)
  ← VBA (parse JSON)
← Excel (display result)
```

Keep VBA thin. Keep Excel dumb. Keep backend smart.

---

## Anti-Patterns to Reject

### ❌ "Let's make it configurable"
**Response:** No. Hard-code the rules. Configuration is complexity.

### ❌ "Let's add a bypass for edge cases"
**Response:** No. If it's an edge case, document it. Don't build a backdoor.

### ❌ "Let's make it more flexible"
**Response:** No. Flexibility = opportunity for impulsivity.

### ❌ "Let's add this convenience feature"
**Response:** Only if it reduces technical complexity without reducing discipline.

### ❌ "This is how other systems do it"
**Response:** We're not building other systems. We're building THIS system.

### ❌ "But what if the user wants to..."
**Response:** Read WHY.md again. This system constrains users. That's the point.

---

## Language & Tone

### When Writing Code
- Clear over clever
- Explicit over implicit
- Verbose over terse
- Simple over elegant

### When Writing Docs
- Direct and concrete
- Use examples liberally
- Explain the "why"
- Assume reader is future you, tired and confused

### When Writing Error Messages
- State the problem
- Show the numbers
- Explain the rule
- Suggest the fix

### When Discussing Features
- Question everything
- Default to "no"
- Require strong justification
- Refer back to WHY.md

---

## Decision-Making Framework

### When Uncertain, Ask:

**1. Does it support discipline?**
   - If yes → Consider it
   - If no → Reject it
   - If unclear → Read WHY.md

**2. Can it be tested with Gherkin?**
   - If yes → Good sign
   - If no → Probably too vague

**3. Is it as simple as possible?**
   - If yes → Good
   - If no → Simplify first

**4. Will you understand it in 6 months?**
   - If yes → Ship it
   - If no → Rewrite it

**5. Would you trust it with your money?**
   - If yes → Ship it
   - If no → Fix it

---

## Git Commit Messages

### Format:
```
[Type] Brief description

- Detail 1
- Detail 2
- Reasoning for change
- Reference to Gherkin scenario if applicable
```

### Types:
- `[feat]` - New feature (with Gherkin)
- `[fix]` - Bug fix
- `[refactor]` - Code improvement, no behavior change
- `[test]` - Adding/fixing tests
- `[docs]` - Documentation only
- `[chore]` - Build, dependencies, etc.

### Examples:
```
[feat] Add portfolio heat validation

- Implement heat cap checking
- Reject trades exceeding 4% cap
- Return detailed error with overage amount
- Implements: features/heat-management.feature:15

[refactor] Simplify position sizing calculation

- Extract stop distance calculation
- Remove redundant validation
- Add explanatory comments
- No behavior change - tests still pass

[fix] Correct shares calculation rounding

- Changed ceil() to floor() to avoid oversizing
- Adds test for fractional share rounding
- Fixes issue where actual risk exceeded specified risk
```

---

## Testing Philosophy

### Unit Tests Should:
- Test one thing
- Be fast (<10ms each)
- Be independent
- Have clear assertions
- Use descriptive names

### Integration Tests Should:
- Test component interaction
- Use real database (test.db)
- Clean up after themselves
- Verify actual behavior
- Match Gherkin scenarios

### Manual Tests Should:
- Be documented
- Be repeatable
- Cover happy path
- Cover error cases
- Be done before every commit

---

## Performance Standards

### Backend Should:
- Respond in <100ms for calculations
- Respond in <500ms for database operations
- Respond in <5s for web scraping
- Handle errors gracefully
- Log performance issues

### Excel Integration Should:
- Feel instant for calculations
- Show progress for slow operations
- Never freeze the UI
- Provide feedback on long operations

### Database Should:
- Use indexes for queries
- Keep schema simple
- Use transactions for writes
- Back up automatically

---

## Security Considerations

### Input Validation:
- Validate ALL inputs before use
- Reject invalid data immediately
- Never trust user input
- Sanitize before database insert

### Database:
- Use parameterized queries (prevent SQL injection)
- No raw SQL with user input
- Limit query results
- Handle connection errors

### File System:
- Validate file paths
- Don't allow path traversal
- Check file permissions
- Handle missing files gracefully

---

## The Meta-Rule

**When in doubt, read WHY.md.**

If the answer isn't there, you don't understand the question well enough yet.

Think more. Code less.

**Every line of code is a liability. Every feature is a potential failure point. Every configuration is a decision users must make.**

Minimize all three.

---

## Remember

This is not a software project that happens to be about trading.

This is a discipline system that happens to be implemented in software.

**Code serves discipline. Discipline does not serve code.**

Act accordingly.
