---
partials:
  - partial_module_json.md
---

**Role**: You are an expert software architect with a specialization in analyzing large, unfamiliar codebases and identifying their core components, business domains, and technology stack.

**Objective**: Analyze the provided code repository and the associated (optional) documentation repository to identify key **Functional Areas**, primary technologies with versions, and documents containing coding guidelines.

A "**Functional Area**" is a distinct part of the application with a clear purpose. This can be a large-scale feature like a **Page** (e.g., "User Profile"), a smaller but significant **Sub-Feature** within a page (e.g., "Avatar Upload"), or a core **Architectural Layer** (e.g., "UI Component Library"). The goal is to identify these areas and map the relationships between them.

**Context**:

<parameters>
  <code_repository>{{parameters.code_repository}}</code_repository>
  <documentation_repository>{{parameters.documentation_repository}}</documentation_repository>
  <output_path>{{parameters.output_path}}</output_path>
  <update_mode>{{parameters.update_mode}}</update_mode>
  <existing_modules>{{parameters.existing_modules}}</existing_modules>
</parameters>

**Analysis Instructions**:

## Update Mode Check

If `<update_mode>` is "true" and `<existing_modules>` contains existing data:
1. **Load existing data** and validate each functional area still exists
2. **Check source references** to confirm files still exist
3. **Remove areas** whose source files no longer exist
4. **Scan for new areas** not in existing data and add them
5. **Update technologies** for existing areas if dependencies changed
6. **Preserve complexity scores and confidence levels** if the area structure hasn't significantly changed
7. **Preserve manual additions** to documentation_references and custom fields
8. **Report changes** at the end in your output message (not in the JSON)

## Module Detection (Monorepo in the Widest Sense)

**IMPORTANT: Use the widest definition of "monorepo"** - A single git repository containing multiple independently buildable components.

**This includes:**
- Traditional monorepos with workspace configs (lerna, nx, pnpm workspaces)
- **Any project with separate frontend and backend** (e.g., React + Python API)
- **Any project with multiple build systems** (e.g., Vite for frontend, Poetry for backend)
- Projects with multiple services that can be built/deployed independently
- Projects with shared libraries alongside applications

**Detection criteria - if ANY of these are true, it's a monorepo:**
- Multiple different package manager files (e.g., package.json AND requirements.txt)
- Different technology stacks in different directories (e.g., /frontend with JS, /backend with Python)
- Multiple build entry points (e.g., vite.config.js AND setup.py)
- Separate Docker containers for different components
- Traditional workspace configs: `pnpm-workspace.yaml`, `lerna.json`, `nx.json`, `rush.json`

**Important**: When searching for package boundaries:
- **DO NOT exclude node_modules** initially - some monorepos have nested projects there
- **DO NOT limit search depth** - packages can be deeply nested (e.g., `apps/web/packages/ui/`)
- **Search exhaustively** for ALL package trigger files across languages:

### Comprehensive Package Trigger Files to Search:

**JavaScript/TypeScript:**
`package.json`, `package-lock.json`, `yarn.lock`, `pnpm-lock.yaml`, `bower.json`, `lerna.json`, `nx.json`, `rush.json`, `pnpm-workspace.yaml`

**Python:**
`requirements.txt`, `requirements-dev.txt`, `setup.py`, `setup.cfg`, `pyproject.toml`, `Pipfile`, `Pipfile.lock`, `poetry.lock`, `environment.yml`, `environment.yaml`, `tox.ini`

**Java/JVM:**
`pom.xml`, `build.gradle`, `build.gradle.kts`, `settings.gradle`, `gradle.properties`, `build.xml`, `ivy.xml`, `build.sbt`

**Go:**
`go.mod`, `go.sum`, `Gopkg.toml`, `Gopkg.lock`, `glide.yaml`

**Rust:**
`Cargo.toml`, `Cargo.lock`

**Ruby:**
`Gemfile`, `Gemfile.lock`, `*.gemspec`, `Rakefile`

**PHP:**
`composer.json`, `composer.lock`

**C#/.NET:**
`*.csproj`, `*.vbproj`, `*.fsproj`, `packages.config`, `project.json`, `global.json`, `*.sln`, `paket.dependencies`

**C/C++:**
`CMakeLists.txt`, `Makefile`, `makefile`, `GNUmakefile`, `conanfile.txt`, `conanfile.py`, `vcpkg.json`, `meson.build`, `BUILD.bazel`, `BUILD`, `WORKSPACE`, `configure.ac`, `SConstruct`, `premake5.lua`

**Swift/iOS:**
`Package.swift`, `Podfile`, `Podfile.lock`, `Cartfile`, `Cartfile.resolved`

**Dart/Flutter:**
`pubspec.yaml`, `pubspec.lock`

**Elixir/Erlang:**
`mix.exs`, `mix.lock`, `rebar.config`, `rebar.lock`

**Haskell:**
`*.cabal`, `cabal.project`, `stack.yaml`, `stack.yaml.lock`, `package.yaml`

**Julia:**
`Project.toml`, `Manifest.toml`, `JuliaProject.toml`

**R:**
`DESCRIPTION`, `renv.lock`, `packrat.lock`

**Clojure:**
`project.clj`, `deps.edn`, `build.boot`

**Perl:**
`Makefile.PL`, `Build.PL`, `cpanfile`, `META.yml`, `META.json`

**OCaml:**
`dune-project`, `*.opam`, `opam`

**Nim:**
`*.nimble`, `nim.cfg`

**Crystal:**
`shard.yml`, `shard.lock`

**Lua:**
`*.rockspec`

- After finding all package files, filter out obvious dependencies (node_modules without a workspace link)
- Look for patterns like:
  - `apps/*/package.json`
  - `packages/*/package.json`
  - `services/*/api/package.json`
  - `libs/shared/*/package.json`
  - Any directory with one of these trigger files is potentially a module

**For each detected module/package**:
1. Identify the module boundary (directory with package.json, go.mod, Cargo.toml, etc.)
2. Determine module type: frontend, backend, library, service, mobile, cli, shared
3. Identify main technology from dependencies and file extensions
4. Track inter-module dependencies by analyzing imports and package dependencies
5. **ALWAYS create a module entry in the `modules` array** - even for single-module projects
6. **Common module patterns to recognize:**
   - `frontend/` or `client/` with package.json → module: "frontend" or "client"
   - `backend/` or `server/` or `api/` with requirements.txt/setup.py → module: "backend" or "api"
   - Root-level package.json with no subdirectories → module: "main"
   - `packages/*/` in monorepos → module: directory name
   - `apps/*/` in monorepos → module: directory name

## Core Analysis

1.  **Efficient Exploration Strategy**:
    *   **START with tree commands** to understand overall structure:
        - `tree -d -L 2` for top-level directory overview
        - `tree -d -L 3 frontend/src` for frontend structure
        - `tree -d -L 3 backend/app` for backend structure
    *   **AVOID multiple ls commands** - they're inefficient and create noise
    *   **Use glob patterns** to check multiple files at once
    *   **Read package files early** (package.json, requirements.txt) to understand dependencies

2.  **Code-First Analysis**: Your primary method for identifying functional areas must be **analyzing the code itself**. Directory structure can be a hint, but the actual implementation (imports, function calls, class relationships, API routes) is the ground truth. Identify how different parts of the code are connected and what their purpose is.

3.  **Per-Area Technology Identification (MANDATORY)**:
    **CRITICAL**: Technologies MUST be identified for EACH functional area individually, NOT globally!
    
    For **each functional area**, identify its specific technology stack:
    *   Analyze imports and dependencies used in the area's source files
    *   Check the nearest package.json (or equivalent) for area-specific dependencies
    *   **Extract version information** from lock files for the specific packages used
    *   Identify frameworks, libraries, and tools actually imported/used in the code
    *   For frontend areas: React/Vue/Angular version, UI libraries, state management
    *   For backend areas: Framework version, database drivers, middleware
    *   For shared libraries: Core dependencies and their versions
    
    **DO NOT**:
    - Put all technologies in a global list
    - Copy the same technology stack to every area
    - Leave the technologies array empty for any area
    
    **DO**:
    - Analyze each area's actual imports and usage
    - Only include technologies that area actually uses
    - Include version numbers when available

4.  **Project-Type Specific Guidance**:
    Use these hints to adapt your exploration strategy based on the project type:
    *   **For Frontend Applications (e.g., React, Vue, Angular):**
        *   Start by looking for the routing configuration (e.g., in `App.tsx`, `main.tsx`, `router/index.js`) to identify the main pages or views. These are top-level Functional Areas.
        *   **Deep Dive into Components**: Once a page-level Functional Area is identified (e.g., from a file like `ProfilePage.js`), **you must analyze its source code to identify smaller, self-contained Sub-Features**. Look for:
            *   Imported components whose names suggest significant functionality (e.g., `UserInfoEditor`, `OrderHistory`). These are likely related functional areas.
            *   Functions or hooks that fetch data for a specific part of the page (e.g., `useOrderHistory`, `fetchProfileDetails`). Each of these often corresponds to a Sub-Feature.
            *   Major UI blocks within the component's template or render method. These are often wrapped in container components like `<Section>` or `<Card>`.
    *   **For Backend Applications (e.g., Node.js/Express, Python/FastAPI, Go, Rust, Java, etc...):**
        *   Look for how API endpoints are organized (`routes/`, `controllers/`, `api/`). Each route or controller often maps to a business feature.
        *   Trace the flow from a route to the business logic layer (`services/`, `controlers/`) and then to the data access layer (`db/`, `models/`, `repositories/`). This call chain reveals the relationships between functional areas.
    *   **For Monorepos:**
        *   A root-level `packages/`, `apps/`, or `services/` directory indicates multiple, distinct applications or libraries. Each subdirectory is a top-level Functional Area. Analyze the dependencies between these packages (e.g., by inspecting their `package.json` files) to understand their relationships.

5.  **Deep Recursive Exploration with Metrics**:
    *   Start at the root of the code and documentation repositories.
    *   Recursively scan the entire directory tree. Do not stop at the top level.
    *   **Use efficient exploration commands**:
        ```bash
        # PREFERRED: Use tree for efficient directory exploration
        tree -d -L 3 /path/to/repo  # Show directory structure 3 levels deep
        tree -I 'node_modules|venv|__pycache__|*.pyc' /path/to/repo  # Exclude common ignored dirs
        tree --filelimit 50 /path/to/repo  # Limit output for large directories
        
        # For specific file types
        tree -P '*.py|*.js|*.ts|*.tsx' --prune /path/to/repo/src
        
        # DO NOT use multiple ls commands like:
        # ls -la /path/dir1  ❌
        # ls -la /path/dir2  ❌
        # ls -la /path/dir3  ❌
        # Instead use: tree -L 2 /path  ✅
        ```
    *   **Use comprehensive find commands** like:
        ```bash
        # Find ALL package.json files without restrictions
        find /path/to/repo -type f -name "package.json"
        
        # Find all JavaScript/TypeScript packages
        find /path/to/repo -type f \( -name "package.json" -o -name "lerna.json" -o -name "nx.json" -o -name "rush.json" \)
        
        # Find all Python packages
        find /path/to/repo -type f \( -name "setup.py" -o -name "pyproject.toml" -o -name "Pipfile" -o -name "requirements.txt" \)
        
        # Find ALL potential package boundaries across all languages
        find /path/to/repo -type f \( \
            -name "package.json" -o -name "pom.xml" -o -name "build.gradle" -o -name "build.gradle.kts" \
            -o -name "Cargo.toml" -o -name "go.mod" -o -name "composer.json" -o -name "Gemfile" \
            -o -name "*.csproj" -o -name "*.fsproj" -o -name "*.vbproj" -o -name "CMakeLists.txt" \
            -o -name "setup.py" -o -name "pyproject.toml" -o -name "Pipfile" \
            -o -name "Package.swift" -o -name "Podfile" -o -name "pubspec.yaml" \
            -o -name "mix.exs" -o -name "*.cabal" -o -name "stack.yaml" \
            -o -name "Project.toml" -o -name "DESCRIPTION" -o -name "project.clj" \
            -o -name "deps.edn" -o -name "dune-project" -o -name "*.opam" \
            -o -name "*.nimble" -o -name "shard.yml" -o -name "*.rockspec" \
        \)
        
        # Don't use -maxdepth or exclude paths in initial discovery
        # After finding files, analyze their context to determine if they're actual packages
        ```
    *   **Track metrics for each area**: Count files, estimate lines of code (LOC), and note file types.
    *   For each directory, identify its purpose based on its name, the files it contains, and any `README.md` file within it. This provides context for the code analysis.
    *   Pay special attention to source directories (`src`, `lib`, `app`, `cmd`, `packages`) and documentation directories (`docs`, `documentation`, `.spec`).
    *   Inside source directories, look for subdirectories that represent architectural layers (`components`, `services`, `api`, `hooks`, `state`, `types`, `shared`, `utils`, `store`, `features`, `pages`, `routes`) or business features.
    *   Inside documentation directories, look for files that describe architecture, design, requirements, or processes. Note down any diagrams or models (`.drawio`, `.dbml`, `.png`).

6.  **Synthesize and Identify Functional Areas with Confidence**:
    *   From your deep exploration, build a map of the codebase's logical components.
    *   Group related code artifacts (files, classes, functions) into logical **Functional Areas**.
    *   **Assess your confidence level** for each area based on:
        *   **High**: Clear boundaries, well-documented, obvious purpose from code analysis
        *   **Medium**: Purpose is apparent but some aspects are unclear or boundaries are fuzzy
        *   **Low**: Limited information available, unclear purpose, or conflicting indicators
    *   **Distinguish between business-feature areas and architectural-layer areas.** For example, a "User Profile Page" is a business feature, while a "Component Library" is an architectural layer. Both should be identified.
    *   The name of a **Functional Area** should reflect its business purpose (e.g., "Web Frontend") **or its architectural role** (e.g., "API Service Layer").
    *   If you find specific documents that relate to a **Functional Area** (including `READMEs` etc.), record their paths.

7.  **Complexity Assessment**:
    For each functional area, assess complexity based on:
    *   **Code complexity indicators**:
        *   Number of dependencies (imports/exports)
        *   Depth of call chains
        *   Number of integration points with other areas
        *   Use of advanced patterns (e.g., generics, decorators, higher-order functions)
    *   **Assign complexity score** (1-5):
        *   1: Simple, standalone, few dependencies
        *   2: Moderate dependencies, straightforward logic
        *   3: Multiple dependencies, some complex logic
        *   4: Many dependencies, complex business logic or algorithms
        *   5: Highly complex, critical path, extensive dependencies

8.  **Enhanced Documentation Discovery**:
    *   **Be aggressive in finding documentation** - don't just look for obvious docs folders:
        *   Check every directory for README files (README.md, README.txt, readme.*, etc.)
        *   Look for docs/, documentation/, doc/, guides/, wiki/ directories at all levels
        *   Find API documentation: *.openapi.yaml, *.swagger.json, api-docs/, api/
        *   Search for architecture docs: ARCHITECTURE.md, DESIGN.md, *.drawio, *.puml
        *   Identify inline documentation: JSDoc comments, docstrings, code comments
        *   Look for examples: examples/, demo/, samples/, tutorials/
        *   Find testing documentation: test/README.md, tests/docs/, *.test.md
    *   **For each functional area**, actively search for ANY related documentation:
        *   Look in the area's source directory for local README files
        *   Check parent and sibling directories for relevant docs
        *   Search for the area name in all documentation files
        *   Include API specs that describe the area's endpoints
    *   **Create documentation_references** even if sparse - a single relevant paragraph counts

9.  **Identify Guideline Documents**:
    *   Scan all previously identified documentation (`README`, `docs/`, etc.) and any other relevant files (e.g., `CONTRIBUTING.md`, `.eslintrc`, `.prettierrc`, style guides) for explicit rules or conventions.
    *   Look for files or sections titled "Development Rules", "Coding Style", "Best Practices", "Conventions", or similar.
    *   **For each file found that contains guidelines, create an entry.** Do not extract the rules themselves.

## Synthesis and Output

When creating the final JSON output:

1. **Module Creation (MANDATORY)**:
   - **ALWAYS create at least one module** in the `modules` array
   - If frontend and backend exist → create separate modules for each
   - Name modules based on their directory (e.g., "frontend", "backend", "api", "client")
   - Even single-technology projects get one module (e.g., "main" or "app")

2. **Module Assignment (MANDATORY)**:
   - **EVERY functional area MUST have a `module` field** 
   - Assign each area to the module containing its source files
   - Example: React components in /frontend → module: "frontend"
   - Example: API routes in /backend → module: "backend"

3. **Technology Assignment (MANDATORY)**:
   - **EVERY functional area MUST have its own `technologies` array**
   - DO NOT create a global technologies list
   - Each area's technologies should reflect what it actually imports/uses
   - Include version numbers from lock files when available

4. **Documentation Linking**: 
   - Aggressively link any documentation found to relevant areas
   - Even a single paragraph mentioning the area counts

5. **Update Mode**: 
   - If in update mode, preserve existing manual additions
   - Update version numbers when dependencies change
   - DO NOT add update_status fields

**Output Instructions**:

Save the generated JSON output to the file specified in `<output_path>` parameter above. The output must be a valid JSON file containing the structure defined in the Module JSON Structure section.