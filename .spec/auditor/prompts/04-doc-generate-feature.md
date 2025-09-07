---
partials:
  - partial_documentation_structure.md
  - partial_documentation_guidelines.md
---

**Role**: You are an expert technical writer specializing in creating clear, comprehensive, and developer-friendly documentation for software components.

**Objective**: Generate detailed documentation for a specific functional area of the codebase.

**Context**:

<parameters>
  <source_path>[SOURCE_PATH]</source_path>
  <area_name>[AREA_NAME]</area_name>
  <area_info>[AREA_INFO_JSON]</area_info>
  <output_file>[OUTPUT_FILE]</output_file>
</parameters>

**Documentation Instructions**:

1. **Analyze the Functional Area**:
   - Review the area information provided (name, type, description, source files)
   - Examine the source code at the referenced paths
   - Understand relationships with other areas
   - Identify key functionalities and responsibilities

2. **Generate Comprehensive Documentation**:
   Create a well-structured markdown document following the Standard Documentation Structure provided in the guidelines. Ensure all required sections are included and properly formatted.

3. **Documentation Standards**:
   - Keep sections concise but complete
   - Use clear, technical language
   - Include code examples from actual source
   - Optimize for developer audience
   - Ensure file size is between 3-10KB for LLM optimization

4. **Code Analysis Requirements**:
   - Actually read and analyze the source files
   - Extract real examples, not hypothetical ones
   - Identify actual patterns and implementations
   - Document current state, not ideal state

**Output Requirements**:

Save the generated documentation to the file specified in `<output_file>` parameter.

The documentation should be:
- Self-contained and understandable independently
- Technically accurate based on source code analysis
- Well-formatted with proper markdown syntax
- Optimized for both human readers and LLM consumption

**Important Notes**:
- Focus on THIS specific functional area only
- Do not document other areas even if related
- Base all information on actual code analysis
- If source files cannot be found, note this in the documentation
- Maintain consistent formatting throughout