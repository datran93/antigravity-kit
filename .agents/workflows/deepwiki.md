---
description: Generate comprehensive, interactive repository wikis and deep knowledge bases, similar to Windsurf Deepwiki.
---

# 📖 Deepwiki Workflow (Automated Knowledge Base)

This workflow guides you to transform a codebase into a structured, navigable wiki. It combines architectural overview with deep implemention details, diagrams, and verifiable citations.

## 🚀 Execution Phase

### Phase 1: Wiki Architecture Design 🏗️
- Use `@mcp:skill-router` to activate the `wiki-architect` skill.
- Scan the repository to detect project type, layers, and tech stack.
- Generate a hierarchical **Wiki Catalogue** (JSON) mapping out:
    - **Onboarding**: Senior-level insights and learning paths.
    - **Deep Dives**: System → Components → Methods.

### Phase 2: Knowledge Synthesis 🧠
- Iteratively process each section of the catalogue using the `wiki-page-writer` skill.
- For each page:
    - Trace actual code paths (no guessing).
    - Create dark-mode Mermaid diagrams (min 2 per page).
    - Cite source files using `file_path:line_number`.

### Phase 3: Repository Integration 🖇️
- Create a `WIKI/` directory in the root.
- Generate a `SUMMARY.md` or `INDEX.md` as the table of contents.
- Link all generated markdown pages together for seamless navigation.

### Phase 4: Final Validation ✅
- Verify all relative links between wiki pages.
- Check that all cited files and lines are accurate.
- Ensure Mermaid diagrams render correctly in the target environment (e.g., GitHub/GitLab).

### Phase 5: Mission Success 🏁
- Present the generated Wiki to the user.
- Offer to explain specific "First Principles" or design tradeoffs discovered.

## 🔴 Critical Constraints
1. **Verifiable Depth**: Every claim MUST have a source. No hand-waving "this likely does X".
2. **First Principles**: Always explain *WHY* an architectural pattern was chosen before *WHAT* it does.
3. **Structured Navigation**: The Wiki must be as easy to navigate as Wikipedia.

---

## 📌 Usage Example
`/deepwiki "Generate a complete onboarding wiki for this repository"`
`/deepwiki "Create a deep technical reference for the internal messaging system"`
