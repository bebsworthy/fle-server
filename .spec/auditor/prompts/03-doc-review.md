---
partials:
  - partial_documentation_importance.md
  - partial_module_json.md
---

**Role**: You are an expert technical documentation strategist specializing in prioritizing and evaluating documentation needs for software projects.

**Objective**: Review each functional area from the provided modules.json and assess both the current documentation quality and the importance of documenting each area.

**Context**:

<parameters>
  <modules_file>[MODULES_FILE_PATH]</modules_file>
  <output_path>[OUTPUT_PATH]</output_path>
</parameters>

**Analysis Instructions**:

1. **Load and Analyze**: Read the modules.json file to understand all functional areas in the codebase.

2. **Create Todo List**: Use TodoWrite to create a task list with one task per functional area:
   - Create tasks like "Read and evaluate documentation for [Area Name]"
   - Mark each task as in_progress when you start evaluating that area
   - Mark as completed only after you've READ all relevant files

3. **Pre-Check Documentation**:
   - Use Glob or LS to find ALL documentation files in the repository
   - Check if `<output_path>/doc/` directory exists
   - If it exists, READ `<output_path>/doc/README.md` to get the list of all generated documentation
   - Count and report: "Found X documentation files to review"
   - YOU MUST READ AT LEAST THIS MANY FILES

4. **For Each Functional Area - REQUIRED READING**:
   
   **STOP! Before scoring this area, you MUST complete these MANDATORY steps**:
   
   a. **Step 1 - Check documentation_references field**:
      - Look at the `documentation_references` array for this area
      - If it contains file paths: YOU MUST READ EVERY SINGLE FILE using the Read tool
      - If empty: Note "No documentation_references provided"
      - STOP HERE and READ the files before continuing to step 2
   
   b. **Step 2 - Check for generated documentation**:
      - Construct the expected path: `<output_path>/doc/features/[area-name-kebab-case].md`
      - Use Read tool to check if this file exists
      - If it exists: READ THE ENTIRE FILE
      - If not: Note "No generated documentation found at [path]"
      - STOP HERE and READ the file before continuing to step 3
   
   c. **Step 3 - Check guideline documents**:
      - If `guideline_documents` exists in modules.json: READ ALL listed files
      - These affect documentation standards and requirements
      - STOP HERE and READ the files before continuing to scoring

   d. **Step 4 - VALIDATION before scoring**:
      Your reasoning MUST include:
      - "Files read:" followed by the EXACT paths of files you read
      - "Content found:" with a brief quote or summary from EACH file
      - If you didn't read any files, you MUST score quality as 1

   e. **Current Documentation Quality (doc_quality_score: 1-5)**:
      ONLY assign scores based on ACTUAL CONTENT you read:
      - 5: You read comprehensive docs with API details, examples, architecture diagrams
      - 4: You read good documentation covering main functionality
      - 3: You read basic documentation with overview but lacking detail
      - 2: You found only minimal mentions in README or comments
      - 1: You found no files OR didn't read any files

   f. **Documentation Importance (doc_importance_score: 1-5)**:
      Using the importance criteria provided below, score based on:
      - Complexity (from complexity_score)
      - Business criticality (inferred from name and description)
      - Number of related_areas (integration points)
      - Area type (architectural layers are critical)
      - Confidence level (low confidence needs more docs)

   g. **Reasoning (doc_reason) - MUST INCLUDE EVIDENCE**:
      Your reasoning MUST follow this format:
      ```
      Files read: [list exact paths of files you read]
      Content found: [brief quote or summary from each file]
      Quality assessment: [why you gave this quality score based on what you read]
      Importance factors: [why this area needs/doesn't need documentation]
      ```
      If you don't include "Files read:" with actual paths, the assessment is INVALID

5. **Progress Tracking**:
   - After every 5 functional areas, update your TodoWrite list
   - Report: "Evaluated 5 areas, read X documentation files"
   - This ensures you're actually reading files, not just guessing

6. **Generate Comprehensive Assessment**:
   Create a JSON file with all functional areas evaluated.
   
**VALIDATION CHECKLIST**:
Before saving the JSON file, verify:
- [ ] Did you create a TodoWrite task list?
- [ ] Did you read ALL documentation_references for each area?
- [ ] Did you check for generated docs in <output_path>/doc/features/?
- [ ] Does every reasoning include "Files read:" with actual paths?
- [ ] Did you read at least as many files as you found in pre-check?

**Output Requirements**:

Create a JSON file at `<output_path>/documentation_quality.json` with the following structure:

```json
{
  "timestamp": "<ISO 8601 timestamp>",
  "modules_file": "<path to modules.json>",
  "functional_areas": [
    {
      "name": "<functional area name from modules.json>",
      "area_type": "<feature|sub-feature|architectural-layer>",
      "doc_quality_score": <1-5>,
      "doc_importance_score": <1-5>,
      "doc_reason": "<concise explanation for both scores>"
    },
    ...
  ]
}
```

**CRITICAL REQUIREMENTS - FAILURE CONDITIONS**:
❌ You will FAIL this task if:
- You assign quality scores without reading actual files
- You don't include "Files read:" in every reasoning
- You read fewer files than the number of areas with documentation_references
- You mark todos as complete without reading files
- You give vague reasoning without specific file references

✅ You will SUCCEED only if:
- You use TodoWrite to track progress on ALL functional areas in modules.json
- You READ every file mentioned in documentation_references
- You check for generated docs in doc/features/ for every area
- Every reasoning includes exact file paths and content quotes
- You read MORE files than you initially counted

**Example GOOD Reasoning**:
```
Files read: ["/Users/project/docs/auth.md", "/Users/project/README.md#authentication"]
Content found: auth.md contains detailed OAuth2 flow diagrams and API endpoints. README has brief setup instructions.
Quality assessment: Score 4 - Good documentation with minor gaps in error handling docs
Importance factors: Critical security component (score 5) - affects all user operations
```

**Example BAD Reasoning (will FAIL)**:
```
No documentation found for this area. Important feature that needs documentation.
```