# Prompt: Comprehensive Component Review

Act as a senior software engineer performing a comprehensive review of the `[FEATURE_NAME]` components located at `[COMPONENT_PATH]`.

## Context & Instructions

Your goal is to produce a detailed, actionable review that helps improve code quality, consistency, and reusability. You will be provided with the necessary context files and paths to the relevant source code.

You MUST perform the following steps:

1.  **Analyze Provided Context**: Thoroughly review the feature's requirements, context, and component map (if available) to understand its purpose and scope.
    -   **Component Source**: `[PATH_TO_FEATURE_COMPONENTS_ROOT_DIR]`
    -   **Requirements**: `[PATH_TO_REQUIREMENTS_DOC]`
    -   **Context / README**: `[PATH_TO_CONTEXT_OR_README_DOC]`

2.  **Analyze for Duplication**: Compare the feature's components against the core component library to identify any components that are duplicates or could be refactored to use existing core components. The core library is located at:
    -   `mockup/v3/src/components/atoms`
    -   `mockup/v3/src/components/molecules`
    -   `mockup/v3/src/services/hooks`

3.  **Assess Quality and Architecture**:
    -   Evaluate the component structure, separation of concerns (e.g., UI components vs. hooks), and adherence to architectural principles (e.g., Atomic Design, Composition over Inheritance).
    -   Review custom hooks for complexity, reusability, and clarity. Identify logic that is overly complex or logic that is generic and should be promoted to a shared service.

4.  **Generate the Review Document**: Structure your output as a markdown file with the exact format specified in the template below. Be constructive, highlighting both strengths and concrete areas for improvement.

---

## Review Template (To be filled by the AI)

# Component Review: [FEATURE_NAME]

**Feature**: `[FEATURE_NAME]`
**Module**: `mockup_v3`
**Review Date**: `[CURRENT_DATE]`
**Status**: Completed

## 1. Overall Assessment

Provide a high-level summary of the feature's implementation, architecture, and overall quality.

### Strengths

-   List 2-4 key strengths, providing brief explanations. Examples: Excellent Component Composition, Robust Error Handling, Clear Documentation, Good Separation of Concerns.

### Areas for Improvement

-   List 2-4 high-level areas that need attention. Examples: Component Duplication, State Management Complexity, Missed Generalization Opportunities.

---

## 2. Component Duplication and Consolidation

Present your findings on component duplication and generalization opportunities using the tables below. This is a critical section that requires detailed analysis.

### Direct Duplications

*Components that have a direct, existing equivalent in the core library.*

| Feature Component | Existing Core Component | Recommendation |
| :--- | :--- | :--- |
| `[e.g., atoms/FeatureEmptyState.tsx]` | `[e.g., molecules/EmptyState.tsx]` | **Replace.** Use the existing core molecule directly. |

### Refactoring and Generalization Candidates

*Components that should be refactored to use core components, or new components that are generic enough to be promoted to the core library.*

| Feature Component | Core Equivalent / Recommendation | Analysis |
| :--- | :--- | :--- |
| `[e.g., atoms/FeatureSpecificBadge.tsx]` | `[e.g., atoms/Badge/]` | **Refactor.** Explain how this component should be refactored to use or compose the core component to avoid duplicating styles and logic. |
| `[e.g., hooks/useFeatureSpecificLogic.ts]` | Create core `hooks/useGenericLogic.ts` | **Generalize.** Explain why this hook's logic is not feature-specific and would be valuable in the core library. |
| `[e.g., molecules/NewGenericComponent.tsx]` | Create core `molecules/NewGenericComponent.tsx` | **Generalize.** Identify this as a new, high-value component that should be part of the core design system from the start. |

---

## 3. State Management & Hooks Review

Provide a detailed analysis of the custom hooks used in the feature. For each major hook, provide an assessment and recommendation.

### `[e.g., useFeatureManager]`

-   **Assessment**: Describe the hook's role, responsibilities, and complexity. Is it doing too much? Is its purpose clear?
-   **Recommendation**: Suggest specific, actionable improvements. Examples: Decompose the hook into smaller, more focused hooks; extract generic logic into a separate utility; clarify state variable names.

---

## 4. Action Plan Summary

Provide a numbered list of the most important, actionable steps that should be taken based on the review's findings. The list should be prioritized.

1.  **[Example]** **Promote Shared Logic**: Move the `use...` hook to a shared `services/` directory.
2.  **[Example]** **Address Duplication**: Replace `...` with the core `...` component.
3.  **[Example]** **Fix Component Location**: Move the `...` component from `[feature]/atoms` to `core/atoms` and update imports.
4.  **[Example]** **Generalize New Components**: Move `...` and `...` to the core library to enrich the design system.
5.  **[Example]** **Document Best Practices**: The composition pattern used in `...` should be documented as a reference for future development.
