# Module JSON Structure

*   **JSON Structure**: The root object must contain: `timestamp`, `modules`, `functional_areas`, and `guideline_documents`.
*   **Timestamp**: (string) ISO 8601 timestamp of when the analysis was performed
*   **Modules Array**: (optional) For monorepos, an array of detected modules/packages, where each contains:
    *   `name`: (string) Module/package name
    *   `path`: (string) Absolute path to the module root
    *   `type`: (string) One of: "frontend", "backend", "library", "service", "mobile", "cli", "shared"
    *   `main_technology`: (string) Primary technology/framework used
    *   `package_manager`: (string, optional) Package manager used (npm, yarn, pnpm, pip, go, cargo, etc.)
    *   `dependencies`: (array of strings, optional) Other modules this module depends on
*   **Functional Areas Array**: A **flat array** of objects representing identified **Functional Areas**.
*   **Guideline Documents Array**: An array of files containing coding guidelines or conventions.
*   **Functional Area Object Properties**: Each functional area must contain:
    *   `name`: (string) Name of the functional area
    *   `description`: (string) Description of its purpose and functionality
    *   `module`: (string, optional) For monorepos, the module this area belongs to
    *   `area_type`: (string) One of: "feature", "sub-feature", "architectural-layer"
    *   `technologies`: (array) Technology stack specific to this area:
        *   `name`: (string) Technology name
        *   `version`: (string, optional) Version if detected
        *   `type`: (string) One of: "language", "framework", "library", "tool", "database", "service"
    *   `source_references`: (array of strings) File paths to source code
    *   `documentation_references`: (array of strings) File paths to documentation (be aggressive in finding these!)
    *   `confidence_level`: (string) One of: "high", "medium", "low"
    *   `metrics`: (object) containing:
        *   `file_count`: (number) Number of files in this area
        *   `estimated_loc`: (number) Estimated lines of code
        *   `primary_language`: (string) Main programming language used
    *   `complexity_score`: (number) 1-5 complexity rating
    *   `complexity_factors`: (array of strings) Key factors contributing to complexity
    *   `related_areas`: (array of strings, optional) Names of other functional areas it interacts with
*   **Guideline Document Properties**: Each guideline document must contain:
    *   `file_path`: (string) The absolute path to the guideline document
    *   `description`: (string) A one-sentence explanation of the guidelines contained

## Example JSON Output

```json
{
  "timestamp": "2024-01-20T10:30:00Z",
  "modules": [
    {
      "name": "web-app",
      "path": "/path/to/project/packages/web-app",
      "type": "frontend",
      "main_technology": "React",
      "package_manager": "pnpm",
      "dependencies": ["shared-ui", "api-client"]
    },
    {
      "name": "api-server",
      "path": "/path/to/project/packages/api-server",
      "type": "backend",
      "main_technology": "Node.js",
      "package_manager": "pnpm",
      "dependencies": ["shared-types"]
    },
    {
      "name": "shared-ui",
      "path": "/path/to/project/packages/shared-ui",
      "type": "library",
      "main_technology": "React",
      "package_manager": "pnpm"
    }
  ],
  "functional_areas": [
    {
      "name": "User Profile Page",
      "description": "Allows users to view and edit their profile information. It contains the Profile Editor and Order History components.",
      "module": "web-app",
      "area_type": "feature",
      "technologies": [
        {"name": "React", "version": "18.2.0", "type": "framework"},
        {"name": "TypeScript", "version": "5.2.2", "type": "language"},
        {"name": "TanStack Query", "version": "5.0.0", "type": "library"},
        {"name": "Tailwind CSS", "version": "3.3.0", "type": "library"}
      ],
      "source_references": [
        "/path/to/project/packages/web-app/src/pages/profile.tsx",
        "/path/to/project/packages/web-app/src/pages/profile.module.css"
      ],
      "documentation_references": [
        "/path/to/project/packages/web-app/README.md#user-profile",
        "/path/to/project/docs/features/user-profile.md",
        "/path/to/project/packages/web-app/src/pages/profile/README.md"
      ],
      "confidence_level": "high",
      "metrics": {
        "file_count": 8,
        "estimated_loc": 1200,
        "primary_language": "TypeScript"
      },
      "complexity_score": 3,
      "complexity_factors": [
        "Multiple sub-components",
        "API integration",
        "State management with context"
      ],
      "related_areas": ["Profile Editor", "Order History", "User API Service"]
    },
    {
      "name": "User API Service",
      "description": "Backend API endpoints for user profile management, authentication, and authorization.",
      "module": "api-server",
      "area_type": "feature",
      "technologies": [
        {"name": "Node.js", "version": "20.10.0", "type": "framework"},
        {"name": "Express", "version": "4.18.2", "type": "framework"},
        {"name": "TypeScript", "version": "5.2.2", "type": "language"},
        {"name": "PostgreSQL", "version": "15", "type": "database"},
        {"name": "Prisma", "version": "5.7.0", "type": "library"}
      ],
      "source_references": [
        "/path/to/project/packages/api-server/src/routes/user.ts",
        "/path/to/project/packages/api-server/src/services/userService.ts",
        "/path/to/project/packages/api-server/src/controllers/userController.ts"
      ],
      "documentation_references": [
        "/path/to/project/packages/api-server/README.md#user-api",
        "/path/to/project/docs/api/user-endpoints.md",
        "/path/to/project/packages/api-server/docs/user-api.openapi.yaml"
      ],
      "confidence_level": "high",
      "metrics": {
        "file_count": 12,
        "estimated_loc": 1800,
        "primary_language": "TypeScript"
      },
      "complexity_score": 3,
      "complexity_factors": [
        "Authentication middleware",
        "Database transactions",
        "Input validation",
        "Role-based access control"
      ],
      "related_areas": ["User Profile Page", "Auth Service", "Database Layer"]
    }
  ],
  "guideline_documents": [
    {
      "file_path": "/path/to/project/CONTRIBUTING.md",
      "description": "Contains instructions for contributing to the project, including code style and commit message format."
    },
    {
      "file_path": "/path/to/project/.eslintrc.json",
      "description": "ESLint configuration defining code quality rules and formatting standards."
    }
  ]
}
```