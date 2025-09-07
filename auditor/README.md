# Codebase Auditor

This project aims to build a comprehensive, AI-driven codebase analysis tool.

## Installation

```bash
cd auditor
npm install
```

## Project Structure

```
auditor/
├── src/                      # Source code (refactored)
│   ├── index.ts             # Main entry point
│   ├── commands/            # Command implementations
│   │   ├── explore.ts       # Codebase exploration
│   │   ├── list.ts          # List functional areas
│   │   └── doc/             # Documentation commands
│   ├── utils/               # Utility functions
│   └── types/               # TypeScript type definitions
├── prompts/                  # AI prompt templates
│   └── review-profiles/      # Technology-specific review guidelines
│       ├── javascript/       # JavaScript ecosystem profiles
│       ├── patterns/         # Architectural pattern profiles
│       └── ...              # Other technology profiles
├── audit.ts                  # Legacy entry point (redirects to src/index.ts)
└── package.json
```

## Usage

### Commands

#### Explore Command

Analyze a codebase to identify functional areas, technologies, and documentation:

```bash
tsx audit.ts [global-options] explore [options] <source_code_path>
```

#### Global Options:
- `--cwd <path>`: Working directory to launch CLI from (default: current directory)
- `--dry-run [level]`: Show what would be executed without actually running the agent (0=minimal, 1=normal, 2=verbose with prompt)
- `--raw`: Show raw output from Claude without filtering (Claude only, default: filtered)
- `--agent-output <file>`: Write raw agent output (stdout+stderr) to file for debugging

#### Explore Command Options:
- `<source_code_path>`: Path to the source code directory (required)
- `--doc <path>`: Path to documentation directory (default: none)
- `--output <path>`: Path to output directory (default: `<cwd>/audit_result`)

#### Examples:

```bash
# Minimal usage (uses all defaults)
npx tsx audit.ts explore ./src

# With documentation
npx tsx audit.ts explore --doc ./docs ./src

# With custom output directory
npx tsx audit.ts explore --output ./analysis ./src

# Full example with all options
npx tsx audit.ts explore --doc ../documentation/modules/mockup_v3/ --output audit_result ../mockup/v3/src/

# Working from a different directory (--cwd is a global flag, goes before 'explore')
npx tsx audit.ts --cwd ../../storyflow explore src/

# Combined example with global and command options
npx tsx audit.ts --cwd ../../storyflow explore --doc documentation/ --output analysis src/

# Dry-run with minimal output (default)
npx tsx audit.ts --dry-run explore ./src

# Dry-run with normal verbosity (shows prompts and variables)
npx tsx audit.ts --dry-run 1 explore ./src

# Dry-run with full verbosity (includes prompt preview)
npx tsx audit.ts --dry-run 2 explore ./src

# When using --cwd, paths are resolved relative to the new working directory
# Option 1: Use paths relative to the new cwd
npx tsx audit.ts --cwd ../../storyflow explore --output audit .

# Option 2: Use absolute paths
npx tsx audit.ts --cwd ../../storyflow explore --output /absolute/path/to/audit /absolute/path/to/source

# Option 3: Simple case - analyze current directory from different cwd
npx tsx audit.ts --cwd ../../storyflow explore .

# Capture raw agent output to a file for debugging
npx tsx audit.ts --agent-output agent-debug.log explore ./src

# Combined with other global options
npx tsx audit.ts --raw --agent-output claude-output.jsonl explore ./src
```

#### List Command

List functional areas from a previously generated modules.json file:

```bash
tsx audit.ts [global-options] list [options]
```

##### List Command Options:
- `--file <path>`: Path to modules.json file (default: `audit_result/modules.json`)
- `--detailed`: Show detailed information including area type, complexity, and metrics

##### Examples:

```bash
# List functional areas from default location
npx tsx audit.ts list

# List from specific file
npx tsx audit.ts list --file ./my-audit/modules.json

# Show detailed information
npx tsx audit.ts list --detailed

# List from different working directory (as in your example)
npx tsx audit.ts --cwd ../../storyflow list --file audit/modules.json

# Detailed view from custom location
npx tsx audit.ts --cwd ../project list --file analysis/modules.json --detailed
```

The list command displays:
- Simple view: numbered list with area names, types, and complexity ratings (★ stars)
- Detailed view: includes description, confidence level, file count, and lines of code
- Summary: counts of features, sub-features, and architectural layers

#### Doc Commands

Documentation analysis and generation with granular control:

```bash
tsx audit.ts [global-options] doc [subcommand] [options]
```

##### Doc Command (No Subcommand)

Display documentation quality summary:

```bash
npx tsx audit.ts doc
```

Shows the documentation quality summary if `documentation_quality.json` exists, or prompts to run `doc review` first.

##### Doc Review Subcommand

Review and score documentation needs for each functional area:

```bash
npx tsx audit.ts doc review [options]
```

**Options:**
- `--file <path>`: Path to modules.json file (default: `audit_result/modules.json`)
- `--output <path>`: Path to output directory (default: `audit_result`)

**Examples:**
```bash
# Review with default paths
npx tsx audit.ts doc review

# Custom paths
npx tsx audit.ts doc review --file ./analysis/modules.json --output ./doc-analysis

# From different working directory
npx tsx audit.ts --cwd ../project doc review

# Dry-run to preview what would be executed
npx tsx audit.ts --dry-run doc review
```

**Output:**
- Creates `documentation_quality.json` with quality and importance scores for each area
- Displays summary ranked by importance
- Removes existing quality file before running

##### Doc Generate Subcommand

Generate documentation incrementally for selected functional areas:

```bash
npx tsx audit.ts doc generate [options]
```

**Options:**
- `--file <path>`: Path to modules.json (default: `audit_result/modules.json`)
- `--quality <path>`: Path to documentation_quality.json (default: `audit_result/documentation_quality.json`)
- `--source <path>`: Path to source code (default: `.`)
- `--output <path>`: Path to output directory (default: `audit_result`)
- `--min-score <number>`: Minimum importance score to document (1-5, default: 3)
- `--name <names...>`: Generate documentation for specific areas by name (overrides --min-score)

##### Doc Update Subcommand

Update existing documentation or create new documentation with incremental update support:

```bash
npx tsx audit.ts doc update [options]
```

**Options:**
- `--file <path>`: Path to modules.json (default: `audit_result/modules.json`)
- `--quality <path>`: Path to documentation_quality.json (default: `audit_result/documentation_quality.json`)
- `--source <path>`: Path to source code (default: `.`)
- `--output <path>`: Path to output directory (default: `audit_result`)
- `--min-score <number>`: Minimum importance score to document (1-5, default: 3)
- `--name <names...>`: Update documentation for specific areas by name (overrides --min-score)

**Examples:**
```bash
# Update existing docs or create new ones with default score (3)
npx tsx audit.ts doc update

# Update documentation for areas with importance >= 4
npx tsx audit.ts doc update --min-score 4

# Update specific areas by name
npx tsx audit.ts doc update --name "Dashboard Page" "Search API"

# Check what would be updated without executing
npx tsx audit.ts --dry-run doc update --name "Database Layer"

# From different directory
npx tsx audit.ts --cwd ../project doc update --min-score 2
```

**Key Differences from `doc generate`:**
- **Preserves existing content**: Manual additions, migration guides, and custom examples are retained
- **Incremental updates**: Only updates sections that have changed in the code
- **Revision history**: Adds change logs documenting what was modified
- **Smart detection**: Identifies what's new, modified, or removed since last update
- **Backward compatibility**: Maintains deprecated feature documentation for migration

**When to use each command:**
- `doc generate`: First-time documentation creation or when you want fresh documentation from scratch
- `doc update`: Updating existing documentation after code changes while preserving manual additions

**Note**: Both commands intelligently handle the documentation index (README.md):
- If no index exists, it creates a new one
- If an index exists, it preserves custom sections while updating auto-generated content
- Custom content should be placed between `<!-- custom:start -->` and `<!-- custom:end -->` markers

##### Doc Generate Examples

```bash
# Generate with default score (3)
npx tsx audit.ts doc generate

# Generate for areas with importance >= 4
npx tsx audit.ts doc generate --min-score 4

# Generate for specific areas by name
npx tsx audit.ts doc generate --name "Dashboard Page" "Search API"

# Generate for a single area
npx tsx audit.ts doc generate --name "Database Layer"

# Custom paths
npx tsx audit.ts doc generate --source ./src --output ./documentation --min-score 4

# From different directory
npx tsx audit.ts --cwd ../project doc generate --min-score 2

# Dry-run to see what areas would be documented and preview prompts
npx tsx audit.ts --dry-run doc generate --min-score 3

# Dry-run with specific names
npx tsx audit.ts --dry-run doc generate --name "Story Retrieval API" "Search Feature"
```

**Process:**
1. Validates modules.json and documentation_quality.json consistency
2. Filters areas by:
   - Specific names if `--name` is provided (ignores --min-score)
   - OR importance score if `--min-score` is provided (default: 3)
3. Generates documentation for each area sequentially
4. Creates documentation index (README.md) with table of contents
5. Shows progress for each area being documented

**Note**: When using `--name`, the importance score is ignored. This allows you to generate documentation for any area regardless of its importance rating.

##### Documentation Scoring

**Quality Scores (Current State):**
- 5: Comprehensive documentation exists
- 4: Good documentation with minor gaps
- 3: Basic documentation present
- 2: Minimal documentation
- 1: No documentation found

**Importance Scores (Priority):**
- 5: MANDATORY - Critical, must be documented
- 4: VERY IMPORTANT - Highly recommended
- 3: IMPORTANT - Should be documented
- 2: NICE TO HAVE - Document when time permits
- 1: NOT IMPORTANT - Self-documenting code sufficient

**Importance Criteria:**
- Complexity (from modules.json complexity_score)
- Business criticality
- Change frequency
- Onboarding impact
- Integration points
- Security/performance impact
- Innovation/uniqueness

##### Generated Structure

```
audit_result/
├── modules.json                    # From explore command
├── documentation_quality.json      # From doc review
├── doc/                            # From doc generate
│   ├── README.md                   # Index with TOC
│   └── features/                   # Individual documentation
│       ├── feature-1.md
│       ├── feature-2.md
│       └── ...
└── audit/                          # From audit command
    ├── README.md                   # Index of audit reports
    ├── dashboard-page.md           # Individual audit reports
    ├── search-feature.md
    └── ...
```

#### Audit Command

Perform code quality and architecture audits on functional areas:

```bash
npx tsx src/index.ts audit [options]
```

##### Audit Command Options:
- `--file <path>`: Path to modules.json file (default: `audit_result/modules.json`)
- `--source <path>`: Path to source code directory (default: `.`)
- `--output <path>`: Path to output directory (default: `audit_result`)
- `--min-complexity <number>`: Minimum complexity score to audit (1-5)
- `--module <name>`: Filter areas by module (for monorepos)
- `--name <names...>`: Audit specific areas by name
- `--focus <aspect>`: Focus on specific aspect (quality|architecture|security|performance)

##### Examples:

```bash
# Audit all functional areas
npx tsx src/index.ts audit

# Audit only high-complexity areas (complexity >= 4)
npx tsx src/index.ts audit --min-complexity 4

# Audit specific areas by name
npx tsx src/index.ts audit --name "Dashboard Page" "Search API"

# Audit areas in a specific module (for monorepos)
npx tsx src/index.ts audit --module frontend

# Combine module and complexity filters
npx tsx src/index.ts audit --module frontend --min-complexity 4

# Focus on security aspects
npx tsx src/index.ts audit --focus security

# Custom paths
npx tsx src/index.ts audit --source ./src --output ./analysis

# From different directory with module filter
npx tsx src/index.ts --cwd ../../storyflow audit --module frontend --min-complexity 4

# Dry-run to preview
npx tsx src/index.ts --dry-run audit --min-complexity 3
```

##### Audit Report Contents

Each audit report includes:
- **Executive Summary**: Health score and key metrics
- **Critical Findings**: Issues requiring immediate attention
- **Code Quality Assessment**: Complexity, duplication, SOLID principles
- **Architecture Review**: Component structure, dependencies, patterns
- **Security & Performance**: Vulnerabilities and bottlenecks
- **Action Plan**: Prioritized improvements with effort estimates

##### Severity Levels

- 🔴 **Critical**: Security vulnerabilities, data loss risks, crashes
- 🟠 **High**: Performance issues, missing error handling, high complexity
- 🟡 **Medium**: Moderate complexity, inconsistent patterns
- 🟢 **Low**: Minor improvements, optimization opportunities

##### Workflow Example

```bash
# Step 1: Explore codebase
npx tsx src/index.ts explore ./src

# Step 2: Review documentation needs
npx tsx src/index.ts doc review

# Step 3: Check summary
npx tsx src/index.ts doc

# Step 4: Generate documentation for important areas
npx tsx src/index.ts doc generate --min-score 3

# Step 5: Audit code quality and architecture
npx tsx src/index.ts audit --min-complexity 3
```

## Review Profiles System

The auditor includes a sophisticated review profiles system that automatically applies technology-specific best practices and guidelines during code audits and documentation reviews. This ensures that audits are contextually aware of the technologies being used and can provide more relevant, actionable feedback.

### How It Works

1. **Automatic Detection**: When running audits, the system automatically detects technologies used in each functional area based on the exploration phase
2. **Profile Matching**: It matches detected technologies against available review profiles
3. **Guideline Application**: Matched profiles are injected into the AI prompts to provide technology-specific review criteria
4. **Visibility**: The audit command clearly shows which profiles are being applied for each area

### Profile Structure

Review profiles are markdown files located in `prompts/review-profiles/` with YAML frontmatter:

```yaml
---
type: framework          # Type: language|framework|pattern|system
name: react              # Profile identifier
keywords: [react, jsx]   # Keywords for searching/matching
priority: 8              # Priority (1-10, higher = more specific)
matches:
  technologies: [react, "@react"]  # Technology names to match
---

# Content with best practices and guidelines...
```

### Available Profile Types

- **language**: Programming languages (TypeScript, JavaScript, Python, etc.)
- **framework**: Application frameworks (React, Vue, Express, etc.)
- **pattern**: Architectural patterns (SPA, microservices, etc.)
- **system**: System-level concerns (security, performance, etc.)

### Current Profiles

The system includes comprehensive profiles for modern web development:

#### JavaScript Ecosystem
- **TypeScript** (`javascript/typescript-best-practice.md`): Strict mode, type safety, immutability patterns
- **React** (`javascript/react-best-practice.md`): React 19 best practices, hooks, performance optimization
- **TanStack Query** (`javascript/tanstack-best-practice.md`): Server state management, caching strategies
- **Tailwind CSS v4** (`javascript/tailwind-4-best-practice.md`): Utility-first CSS, theming, variant management
- **JavaScript Base** (`javascript/_base.md`): Core JavaScript best practices

#### Architectural Patterns
- **SPA** (`patterns/spa.md`): Single-page application patterns
- **Additional patterns**: Extensible for microservices, serverless, etc.

#### System Profiles
- **Fundamental** (`fundamental-best-practice.md`): SOLID principles, clean code practices

### Profile Resolution

Profiles are resolved in a specific order based on:
1. **Priority**: Higher priority profiles override lower ones
2. **Specificity**: More specific profiles (framework) override general ones (language)
3. **Inheritance**: Profiles can extend other profiles using the `extends` field

Example resolution chain for a React + TypeScript project:
```
fundamental → javascript → typescript → react → tanstack-query
```

### Visibility During Audits

When running the audit command, you'll see which profiles are applied:

```bash
📚 Loading review profiles...
✓ Loaded 14 review profiles

[1/3] Auditing: Dashboard Page
📋 Applying Review Profiles:
  🏗️ react (framework) - javascript/react-best-practice.md
  🏗️ tanstack-query (framework) - javascript/tanstack-best-practice.md  
  🔤 typescript (language) - javascript/typescript-best-practice.md
  🔤 javascript (language) - javascript/_base.md
  → 4 profiles with ~93 review guidelines
```

### Adding Custom Profiles

To add your own technology-specific profiles:

1. Create a markdown file in `prompts/review-profiles/[category]/`
2. Add appropriate frontmatter with type, name, keywords, and matches
3. Write comprehensive best practices and guidelines
4. The profile will be automatically detected and applied

Example custom profile:

```yaml
---
type: framework
name: nextjs
keywords: [nextjs, next, react, ssr, ssg]
priority: 9
matches:
  technologies: [nextjs, next, "next.js"]
---

# Next.js Best Practices

## Routing
- Use App Router for new projects...

## Data Fetching  
- Implement proper caching strategies...
```

### Profile Benefits

1. **Consistency**: Ensures all audits follow the same technology-specific standards
2. **Comprehensiveness**: ~50-100 guidelines per technology stack
3. **Maintainability**: Update best practices in one place
4. **Extensibility**: Easy to add new technologies
5. **Transparency**: Clear visibility of what's being checked

##### Output Filtering (Claude)

When using Claude CLI, the tool automatically filters the verbose JSON stream output to show only essential information:
- **Init messages**: Shows model and working directory
- **Agent messages**: Shows the agent's text responses
- **Tool use**: Shows tool name and primary parameter (e.g., file path, command)
- **Result**: Shows success/error status with truncated message

To see the raw, unfiltered JSON stream from Claude, use the `--raw` flag:

```bash
# Filtered output (default) - clean and readable
npx tsx audit.ts explore ./src

# Raw output - full JSON stream
npx tsx audit.ts --raw explore ./src
```

**Note**: Output filtering only applies to Claude. Gemini output is always shown as-is.

##### Dry-Run Mode

The `--dry-run` flag allows you to preview what would be executed without actually running the AI agent. This is useful for:
- Reviewing the prompts that would be sent to the AI
- Understanding what files would be created or modified
- Debugging command configurations
- Previewing the workflow before execution

**Verbosity Levels:**
- `--dry-run` or `--dry-run 0`: Minimal output (default)
  - Shows task description, configuration, and command to be executed
- `--dry-run 1`: Normal verbosity
  - Adds list of prompts to be concatenated
  - Shows variables to be injected
- `--dry-run 2`: Full verbosity
  - Includes first 1000 characters of the final prompt

**Example:**
```bash
# Minimal preview (default)
npx tsx audit.ts --dry-run explore ./src

# Normal verbosity - see prompts and variables
npx tsx audit.ts --dry-run 1 doc review

# Full verbosity - includes prompt preview
npx tsx audit.ts --dry-run 2 doc generate --min-score 4
```

##### Customization

- Edit `prompts/partial_documentation_importance.md` to customize importance scoring criteria
- Edit `prompts/partial_documentation_guidelines.md` to customize quality standards

The script will:
1. Prompt you to select between Claude CLI or Gemini CLI
2. Generate a comprehensive analysis of the codebase
3. Save the results to `<output_folder>/modules.json`

## Process Overview

The process is divided into two main phases:

1.  **Exploration**: The first step is to run an initial exploration prompt (`prompts/01-explore.md`) to get a high-level overview of the entire codebase. This identifies the key modules, technology stack, and overall architecture.

2.  **In-Depth Analysis**: Once the modules are identified, subsequent prompts will be used to analyze each module independently. The goal of this phase is to generate a detailed list of potential improvements, such as refactoring opportunities, performance enhancements, or documentation gaps.

## Prerequisites

You need to have either Claude CLI or Gemini CLI installed:

- **Claude CLI**: Install from [Anthropic's Claude CLI](https://github.com/anthropics/claude-cli)
- **Gemini CLI**: Install from [Google's Gemini CLI](https://github.com/google/gemini-cli)

Make sure the CLI tool is available in your system PATH.
