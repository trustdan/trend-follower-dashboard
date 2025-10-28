# Documentation Index
# Trading Engine v3

**Excel-based trading platform enforcing disciplined trend-following through systematic constraints.**

This directory contains all project documentation, organized by purpose.

---

## 📖 Quick Links

**New to the project?** Start here:
1. [WHY.md](project/WHY.md) - Core philosophy (5 min read)
2. [PLAN.md](project/PLAN.md) - Step-by-step implementation plan
3. [DEVELOPMENT_PHILOSOPHY.md](dev/DEVELOPMENT_PHILOSOPHY.md) - How we build

**Working with Claude Code?** Read:
- [CLAUDE.md](dev/CLAUDE.md) - Claude Code guidance for this project

**Writing features?** See:
- [BDD_GUIDE.md](dev/BDD_GUIDE.md) - Behavior-driven development approach

---

## 📁 Documentation Structure

```
docs/
├── README.md (this file)         # Documentation index
├── project/                      # Core project documents
│   ├── WHY.md                    # Philosophy and purpose
│   └── PLAN.md                   # Step-by-step implementation plan
├── dev/                          # Development guides
│   ├── DEVELOPMENT_PHILOSOPHY.md # How we build this system
│   ├── BDD_GUIDE.md              # Testing approach
│   ├── CLAUDE.md                 # Claude Code project guidance
│   └── CLAUDE_RULES.md           # Detailed development rules
├── milestones/                   # Milestone completion reports
│   ├── M17-M18_*.md              # M17-M18 completion docs
│   ├── M19_COMPLETION_SUMMARY.md # VBA implementation
│   └── M20_COMPLETION_SUMMARY.md # Windows integration package
├── json-schemas/                 # JSON API specifications
│   └── JSON_API_SPECIFICATION.md
└── HTTP_CLI_PARITY.md            # Transport parity testing
```

---

## 🎯 Project Documents

### [WHY.md](project/WHY.md) ⭐ **START HERE**
**Purpose:** Explains why this system exists and what problem it solves

**Key Concepts:**
- Discipline over flexibility
- The cost of impulsive trading
- Why constraints are features, not limitations
- Ed Seykota's philosophy applied

**Read this first** to understand the system's purpose. Everything else makes sense after this.

**Time:** 5 minutes

---

### [PLAN.md](project/PLAN.md)
**Purpose:** Complete step-by-step development plan (formerly Trading-Engine-v3_Step-by-Step-Plan.md)

**Contents:**
- Architecture overview (engine-first, CLI by default, HTTP optional)
- Data model (SQLite schema)
- Development phases (A through E)
- Milestones M1-M21 with exit criteria
- BDD test scenarios
- Timeline and risk mitigation

**Use this for:**
- Understanding the overall architecture
- Knowing what milestone comes next
- Checking exit criteria for completed work
- Planning development work

**Status:** Living document, updated as milestones complete

---

## 🛠️ Development Guides

### [DEVELOPMENT_PHILOSOPHY.md](dev/DEVELOPMENT_PHILOSOPHY.md)
**Purpose:** How we build this system

**Key Principles:**
- Gherkin first, always
- Question every feature request
- Anti-patterns to reject immediately
- Code quality standards
- Testing philosophy

**Read before coding** anything.

---

### [BDD_GUIDE.md](dev/BDD_GUIDE.md)
**Purpose:** Behavior-driven development approach for this project

**Contents:**
- What BDD is and why we use it
- Gherkin scenario format
- Feature file organization
- Step definition patterns
- Testing workflow

**Use this when:**
- Writing new features
- Creating test scenarios
- Implementing step definitions

---

### [CLAUDE.md](dev/CLAUDE.md)
**Purpose:** Guidance for Claude Code when working with this repository

**Contents:**
- Project overview and philosophy
- Current status and architecture
- Development commands
- Critical development rules
- Repository structure
- Standard patterns

**Audience:** Claude Code (but useful for human developers too)

---

### [CLAUDE_RULES.md](dev/CLAUDE_RULES.md)
**Purpose:** Detailed development rules and anti-patterns

**Contents:**
- Feature request evaluation criteria
- Anti-patterns to reject
- Code quality checklist
- Error message standards
- Git workflow

**Use this for:** Detailed guidance on making development decisions

---

## 📊 Milestone Reports

Completion summaries for each development milestone. These document what was built, how it was tested, and what's next.

### M17-M18: JSON Contract Validation
- [M17-M18_FINAL_SUMMARY.md](milestones/M17-M18_FINAL_SUMMARY.md) - Complete summary
- [M17-M18_COMPLETE.md](milestones/M17-M18_COMPLETE.md) - Completion checklist
- [M17-M18_PROGRESS.md](milestones/M17-M18_PROGRESS.md) - Progress tracking
- [M17-M18_ISSUES_TO_FIX.md](milestones/M17-M18_ISSUES_TO_FIX.md) - Issues found and fixed

**Milestone:** Defined and validated JSON contracts for CLI and HTTP

---

### [M19: VBA Implementation](milestones/M19_COMPLETION_SUMMARY.md)
**Milestone:** Created VBA modules as text exports (`.bas` files)

**Deliverables:**
- TFTypes.bas - Type definitions
- TFHelpers.bas - JSON parsing and utilities
- TFEngine.bas - Engine communication layer
- TFTests.bas - VBA unit tests (14 tests)

**Status:** ✅ Complete (2025-10-27)

---

### [M20: Windows Integration Package](milestones/M20_COMPLETION_SUMMARY.md)
**Milestone:** Complete Windows deployment package

**Deliverables:**
- Cross-compiled Windows binary (tf-engine.exe, 12 MB)
- VBA import automation (windows-import-vba.bat)
- Database initialization (windows-init-database.bat)
- Automated test runner (run-tests.bat, 11 tests)
- Comprehensive testing guide (WINDOWS_TESTING.md, 23 KB)
- Excel workbook template specification (14 KB)
- Test data files (21 JSON samples)

**Status:** ✅ Complete (2025-10-27), ready for M21

---

## 🔧 Technical Specifications

### [JSON_API_SPECIFICATION.md](json-schemas/JSON_API_SPECIFICATION.md)
**Purpose:** Complete JSON contract specifications for all engine commands

**Contents:**
- Request/response schemas for all commands
- Example JSON for success and error cases
- Field descriptions and validation rules

**Use this for:**
- Implementing VBA parsers
- Validating engine outputs
- Understanding API contracts

---

### [HTTP_CLI_PARITY.md](HTTP_CLI_PARITY.md)
**Purpose:** Testing that CLI and HTTP return identical JSON

**Contents:**
- Parity testing approach
- Test scenarios
- Validation results

---

## 🧭 Navigation Guide

### "I'm new to the project"
1. Read [WHY.md](project/WHY.md) (5 min)
2. Skim [PLAN.md](project/PLAN.md) for architecture (10 min)
3. Read [DEVELOPMENT_PHILOSOPHY.md](dev/DEVELOPMENT_PHILOSOPHY.md) (10 min)

### "I'm implementing a feature"
1. Write Gherkin scenario first ([BDD_GUIDE.md](dev/BDD_GUIDE.md))
2. Check anti-patterns ([DEVELOPMENT_PHILOSOPHY.md](dev/DEVELOPMENT_PHILOSOPHY.md))
3. Implement with discipline over flexibility
4. Write tests matching Gherkin

### "I'm working with VBA"
1. Read [M19_COMPLETION_SUMMARY.md](milestones/M19_COMPLETION_SUMMARY.md)
2. Check `../excel/VBA_MODULES_README.md` for VBA specifics
3. Reference [JSON_API_SPECIFICATION.md](json-schemas/JSON_API_SPECIFICATION.md) for contracts

### "I'm testing on Windows"
1. Read [M20_COMPLETION_SUMMARY.md](milestones/M20_COMPLETION_SUMMARY.md)
2. Follow `../windows/WINDOWS_TESTING.md` step-by-step (M21)
3. Use `../windows/README.md` for package info

### "I'm coming back after time away"
Context switching checklist from [CLAUDE.md](dev/CLAUDE.md):
- [ ] Re-read [WHY.md](project/WHY.md) (5 min)
- [ ] Review [PLAN.md](project/PLAN.md) current phase
- [ ] Check recent commits
- [ ] Run tests
- [ ] Read relevant Gherkin scenarios

---

## 📝 Documentation Standards

**Markdown Format:**
- Use clear headings (## for sections, ### for subsections)
- Include table of contents for docs > 500 lines
- Use code blocks with language tags
- Include examples for every concept

**Naming Conventions:**
- SCREAMING_SNAKE_CASE.md for major docs (WHY.md, PLAN.md)
- Milestone_Summary_Format.md for milestone reports
- Technical_Specification_Format.md for specs

**Cross-References:**
- Use relative links: `[WHY.md](project/WHY.md)`
- Link to source code with line numbers: `internal/domain/sizing.go:42`
- Always provide context for links

**Living Documents:**
- Mark documents as "Living" if they're updated frequently
- Include "Last Updated" date for living documents
- Archive old versions in milestones/ folder

---

## 🔄 Document Lifecycle

**Project Documents** (project/)
- Core philosophy and architecture
- Updated rarely, carefully
- Require project lead approval for changes

**Development Guides** (dev/)
- Process and standards
- Updated as we learn better approaches
- Propose changes via discussion first

**Milestone Reports** (milestones/)
- Historical record
- **Never modified** after sign-off
- Dated snapshots of completed work

**Technical Specs** (json-schemas/)
- API contracts and schemas
- Updated when implementation changes
- Version-controlled

---

## 🎯 Current Status

**Phase:** M20 Complete, M21 Ready
**Milestone:** M20 - Windows Integration Package ✅
**Next:** M21 - Windows Integration Validation (manual testing)

**Recent Updates:**
- 2025-10-27: M20 completed (Windows package ready)
- 2025-10-27: M19 completed (VBA modules created)
- 2025-10-27: M17-M18 completed (JSON contracts validated)

**Current Focus:** Preparing for M21 Windows manual testing

---

## 📚 External References

**Core Influences:**
- Ed Seykota's trading philosophy
- Van Tharp's position sizing methods
- Behavior-driven development (BDD)
- Domain-driven design (DDD)

**Technical Stack:**
- Go (backend engine)
- SQLite (storage)
- Excel + VBA (frontend)
- Gherkin + godog (BDD testing)

---

## 🤝 Contributing

**Before Contributing:**
1. Read [WHY.md](project/WHY.md) - Understand the purpose
2. Read [DEVELOPMENT_PHILOSOPHY.md](dev/DEVELOPMENT_PHILOSOPHY.md) - Understand the approach
3. Check anti-patterns - Know what to reject

**Development Workflow:**
1. Write Gherkin scenario first (see [BDD_GUIDE.md](dev/BDD_GUIDE.md))
2. Get agreement on behavior
3. Implement code
4. Write tests matching Gherkin
5. Document in milestone report

**Questions to Ask:**
- Does this support discipline or undermine it?
- Would Ed Seykota approve?
- Does it make impulsivity easier or harder?
- Is this solving a real problem or adding complexity?

**If unsure:** Read [WHY.md](project/WHY.md) again. The answer is there.

---

## 📞 Support

**Documentation Issues:**
- Found a broken link? Fix it and note in commit message
- Found unclear documentation? Propose improvement
- Found missing documentation? Check if it belongs in a different folder

**Technical Issues:**
- Check milestone reports for known issues
- Check `tf-engine.log` and `TradingSystem_Debug.log` for correlation IDs
- Reference [PLAN.md](project/PLAN.md) for architecture questions

---

**Remember:** This is a discipline enforcement system. Documentation serves that mission by being clear, accurate, and focused.

Code serves discipline. Discipline does not serve code.
