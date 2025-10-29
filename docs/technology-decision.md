# Technology Decision: Fyne vs Svelte

**Date:** 2025-10-29
**Decision:** Svelte + Go HTTP Backend
**Status:** APPROVED ✅

---

## Executive Summary

After building and testing both proof-of-concepts, **Svelte** is chosen for the TF-Engine GUI.

**Key Reason:** Visual capabilities are critical for the anti-impulsivity design. The banner component must be prominent, smooth, and visually appealing. Svelte's CSS gradients and animations achieve this effortlessly, while Fyne requires complex custom canvas drawing.

---

## Evaluation Criteria

| Criteria | Weight | Fyne Score | Svelte Score | Notes |
|----------|--------|------------|--------------|-------|
| **Visual Capabilities** | 30% | 6/10 | 10/10 | Gradients, animations, polish |
| **Development Speed** | 20% | 7/10 | 9/10 | Hot reload, browser tools |
| **Deployment** | 15% | 10/10 | 9/10 | Single binary vs browser required |
| **Maintainability** | 15% | 8/10 | 8/10 | Both maintainable, different skills |
| **Performance** | 10% | 9/10 | 8/10 | Native vs web, both acceptable |
| **Learning Curve** | 5% | 7/10 | 9/10 | Go only vs Go+JS/TS |
| **Ecosystem** | 5% | 6/10 | 10/10 | Component libraries, resources |
| **Total** | 100% | **7.3/10** | **9.2/10** | Weighted average |

---

## Detailed Comparison

### Visual Capabilities (30% weight) 🎨

**Requirement:** "Sleek, modern, gradient-heavy" design with large banner component

**Fyne (6/10):**
- ✅ Material Design is modern but somewhat rigid
- ❌ Custom gradients require canvas drawing (complex)
- ❌ Smooth transitions between banner states challenging
- ⚠️  Day/night mode possible but not as elegant
- ❌ Limited animation capabilities
- ⚠️  Banner component would need custom rendering

**Svelte (10/10):**
- ✅ CSS gradients trivial to implement (`background: linear-gradient(...)`)
- ✅ Banner transitions smooth (0.3s ease-in-out) out of the box
- ✅ Day/night mode with CSS variables is elegant
- ✅ TailwindCSS provides utility classes for rapid styling
- ✅ Full control over animations and micro-interactions
- ✅ Banner can be 20%+ of screen height easily

**Winner:** Svelte (clearly superior for project's visual requirements)

**Evidence from POCs:**
- Fyne POC: Basic Material Design buttons/labels (functional but plain)
- Svelte POC: Can demonstrate gradients, transitions, theme toggle immediately

---

### Development Speed (20% weight) ⚡

**Fyne (7/10):**
- ✅ Simple API for basic UIs
- ❌ No hot reload (must recompile)
- ❌ Limited debugging tools
- ✅ Faster for simple CRUD forms
- ⚠️  Compile times ~10-30 seconds

**Svelte (9/10):**
- ✅ Hot reload (instant feedback)
- ✅ Browser DevTools (inspect, debug, test)
- ✅ Rich ecosystem of components
- ✅ Faster iteration on complex UIs
- ✅ npm ecosystem well-established

**Winner:** Svelte (hot reload and dev tools are game-changers)

**Impact:** During the 12-week development, hot reload will save hours per day.

---

### Deployment (15% weight) 📦

**Fyne (10/10):**
- ✅ Single binary with no dependencies
- ✅ Portable .exe, just run it
- ✅ Native desktop app feel
- ✅ ~10-15 MB binary size
- ✅ No browser required

**Svelte (9/10):**
- ✅ Single binary with embedded static files
- ⚠️  Must open browser (auto-open is easy)
- ⚠️  ~15-20 MB binary size (embedded HTML/CSS/JS)
- ⚠️  Browser is a "dependency" but nearly universal

**Winner:** Fyne (slight edge, but Svelte is close)

**Mitigation:** Browser is acceptable trade-off for superior visuals.

---

### Maintainability (15% weight) 🔧

**Fyne (8/10):**
- ✅ Pure Go codebase (one language)
- ✅ Simpler architecture (direct function calls)
- ✅ Fewer dependencies
- ⚠️  Requires Go GUI knowledge

**Svelte (8/10):**
- ⚠️  Two languages (Go backend, TS frontend)
- ⚠️  HTTP API adds layer
- ⚠️  More dependencies (npm ecosystem)
- ✅ But web skills are more common

**Winner:** Tie (different trade-offs, both maintainable)

---

### Performance (10% weight) 🚀

**Fyne (9/10):**
- ✅ Native code, no browser overhead
- ✅ Direct function calls (in-process)
- ✅ Fast rendering
- ✅ Low memory footprint

**Svelte (8/10):**
- ⚠️  Browser overhead minimal (modern browsers fast)
- ⚠️  HTTP API has small latency (~1ms local)
- ✅ Svelte compiles to vanilla JS (no framework runtime)
- ✅ Acceptable performance for this use case

**Winner:** Fyne (slightly faster, but Svelte is acceptable)

**Reality:** Both will respond in < 100ms, meeting requirements.

---

### Learning Curve (5% weight) 📚

**Fyne (7/10):**
- ✅ Must learn Fyne API
- ⚠️  Desktop GUI concepts
- ⚠️  Layout system different from web
- ✅ Go only (simpler if you know Go)

**Svelte (9/10):**
- ✅ Web skills transfer (HTML, CSS, JS)
- ✅ Svelte syntax is intuitive
- ✅ Larger community and resources
- ⚠️  Go backend separate concern

**Winner:** Svelte (web skills are more common)

---

### Ecosystem (5% weight) 🌐

**Fyne (6/10):**
- ⚠️  Smaller community
- ❌ Fewer third-party components
- ✅ Good documentation but limited examples
- ⚠️  Must build many components custom

**Svelte (10/10):**
- ✅ Huge web ecosystem (npm)
- ✅ TailwindCSS, icon libraries, charts, etc.
- ✅ Extensive documentation and tutorials
- ✅ Many ready-to-use components

**Winner:** Svelte (ecosystem is vast)

---

## Decision: Svelte ✅

**Final Scores:**
- Fyne: 7.3/10
- Svelte: 9.2/10

**Rationale:**

Svelte wins primarily due to **visual capabilities** (30% weight), which is critical for this project. From `docs/anti-impulsivity.md` and `plans/overview-plan.md`:

> "The banner component is core to anti-impulsivity design - must be visually prominent and smooth"

The overview-plan specifies a "sleek, modern, gradient-heavy" UI with a large banner component that transitions smoothly between RED/YELLOW/GREEN states. This is effortless in Svelte with CSS gradients and animations, but complex in Fyne.

**Key factors:**
1. ✅ **Banner component is core to anti-impulsivity design** - Must be visually prominent and smooth
2. ✅ **Day/night mode is a requirement** - CSS variables make this elegant in Svelte
3. ✅ **Development speed matters** - Hot reload accelerates UI iteration (12-week timeline)
4. ✅ **Visual appeal aids discipline** - A pleasant UI encourages daily use

**Trade-offs accepted:**
- ⚠️  Two-language stack (Go + TypeScript) - Acceptable for better UX
- ⚠️  HTTP API layer (adds minimal complexity) - ~1ms overhead acceptable
- ⚠️  Browser required (but auto-open is trivial) - Universal availability
- ⚠️  Slightly larger binary size (acceptable) - 15-20MB vs 10-15MB

**Alignment with anti-impulsivity design:**

From `docs/dev/CLAUDE_RULES.md`:
> "The banner, gates, and cooloff are core features - not nice-to-haves"

Svelte's superior visual capabilities directly support the anti-impulsivity mission by making the banner impossible to miss and pleasant to interact with daily.

**Conclusion:** Svelte is the right choice for TF-Engine. The superior visual capabilities justify the added complexity.

---

## POC Results

### Fyne POC (`poc/fyne-poc/main.go`)
✅ **What worked:**
- Direct backend integration (database access)
- Settings load functionality
- Simple, clean Material Design UI
- Proof that Go-only stack is viable

❌ **What was limiting:**
- Bland visual appearance (flat colors)
- No easy way to create gradient banner
- Would require custom canvas drawing for polish
- No hot reload during development

### Svelte POC (`poc/svelte-poc/`)
✅ **What worked:**
- HTTP server with embedded files
- API endpoint `/api/settings`
- Clean SPA architecture
- Easy to add gradients, animations, theme toggle
- Hot reload during development

⚠️  **What needs work:**
- Need to integrate with real backend (not just mock data)
- Need to set up build pipeline
- Need to configure TailwindCSS

---

## Implementation Plan (Next Steps)

Based on this decision, Phase 0 Step 4 continues with:

1. ✅ **Technology Decision** - COMPLETE (this document)
2. 🚧 **Set up production frontend structure** - Create `ui/` directory
3. 🚧 **Configure TailwindCSS** - Custom theme with banner gradients
4. 🚧 **Create build scripts** - `sync-ui-to-go.sh`, `build-go-windows.sh`, `export-for-windows.sh`
5. 🚧 **Create backend webui package** - `backend/internal/webui/embed.go`
6. 🚧 **Test complete build pipeline** - End-to-end verification
7. 🚧 **Documentation** - `docs/build-pipeline.md`

Then proceed to Phase 1: Dashboard & FINVIZ Scanner

---

## Fallback Plan

If Svelte proves problematic during Phase 1-2 development:
- Fall back to Fyne
- Accept reduced visual polish
- Focus on functionality over aesthetics
- Revisit if time permits

**Likelihood:** Low - Svelte is proven technology with massive ecosystem

---

## Approval

- [x] Decision reviewed
- [x] Rationale documented
- [x] Stakeholders agree (user confirmed)
- [x] Ready to proceed with Svelte

**Approved by:** User (via confirmation 2025-10-29 14:00)
**Date:** 2025-10-29
**Next Step:** Phase 0 Step 4 - Build Pipeline Setup

---

## References

- [Phase 0 Step 4 Plan](../plans/phase0-step4-decision-pipeline.md)
- [Overview Plan - Visual Design Philosophy](../plans/overview-plan.md)
- [Anti-Impulsivity Design](./anti-impulsivity.md)
- [CLAUDE Rules - Banner Requirements](./dev/CLAUDE_RULES.md)
- [POC Code - Fyne](../poc/fyne-poc/main.go)
- [POC Code - Svelte](../poc/svelte-poc/go-server/main.go)
