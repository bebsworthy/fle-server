---
type: framework
name: ant-design
keywords: [antd, ant-design, ui, components, design-system, react, enterprise, forms, tables]
priority: 8
matches:
  technologies: [antd, ant-design, "@ant-design/icons", "antd/es", "antd/lib"]
---

# Ant Design Best Practices: A Developer's Guide

## **The Guiding Philosophy: Why These Rules Exist**

Our primary goal is to build applications that are **maintainable, consistent, and provide a high-quality user experience**. Ant Design is a powerful and comprehensive UI library, but its power comes with complexity. These guidelines are designed to help us navigate that complexity, avoid common pitfalls, and leverage the library's strengths. When a rule seems unclear, we should default to the choice that best aligns with Ant Design's core design values: **Natural, Certain, Meaningful, and Growing.**

---

## **Architectural & Layout Best Practices**

### **1. Use the Grid System for All Layouts**

**Use** Ant Design's built-in 24-point `Row` and `Col` components for structuring page content. Do not mix it with other grid systems (e.g., Bootstrap) to avoid CSS conflicts.

*   **Why?** This ensures a consistent, predictable, and responsive layout that adheres to the "Certainty" design value. It leverages the framework's primary tool for layout, ensuring future compatibility and stability.

*   **Example:**

    ```jsx
    // ❌ Bad: Mixing layout systems or avoiding the grid.
    <div style={{ display: 'flex', width: '100%' }}>
      <div style={{ width: '50%' }}>...</div>
      <div style={{ width: '50%' }}>...</div>
    </div>

    // ✅ Good: Clean, responsive, and idiomatic.
    import { Row, Col } from 'antd';

    <Row gutter={16}>
      <Col span={12}>...</Col>
      <Col span={12}>...</Col>
    </Row>
    ```

### **2. Structure Pages with the `Layout` Component**

**Use** the `Layout` component suite (`Layout`, `Header`, `Sider`, `Content`, `Footer`) for the overall structure of your application's pages.

*   **Why?** These components are built on Flexbox and are designed to work seamlessly together, providing a robust foundation for common enterprise application interfaces.

---

## **Form Best Practices**

### **1. Let the `Form` Component Manage State**

**Do not** use React's `useState` to manage the state of inputs inside a named `Form.Item`. Provide default values to the form via the `initialValues` prop on the `<Form>` component itself.

*   **Why?** The Ant Design `Form` is a high-performance state management system. When a `Form.Item` has a `name`, it automatically controls its children. Fighting this behavior leads to bugs and unpredictable UI.

### **2. Use Declarative Rules for Validation**

**Use** the `rules` prop on `Form.Item` for all standard validation. For performance-critical inputs with asynchronous validation, **set** `validateTrigger` to `onBlur` and consider using `validateDebounce`.

*   **Why?** This keeps validation logic declarative and co-located with the field. Optimizing the trigger prevents excessive API calls and improves user experience.

### **3. Organize Complex Forms with `Tabs`**

**Group** long forms (more than 15 items) into `Tabs`. Avoid multi-column layouts for a single form.

*   **Why?** This reduces cognitive load and creates a clear, predictable reading order for the user, aligning with the "Natural" design principle.

### **4. Use `Form.List` for Dynamic Fields**

**Use** the `Form.List` component to manage a variable number of form fields.

*   **Why?** It provides a clean, structured API for managing dynamic form data, abstracting away the complexity of managing field names and keys.

---

## **Data Display: The `Table` Component**

### **1. Always Provide a `rowKey`**

**Always** specify the `rowKey` prop on the `Table` component.

*   **Why?** React needs a unique key for each item in a list for efficient DOM updates. For the `Table`, `rowKey` is critical for features like row selection and expansion to work without bugs.

### **2. Handle Data Operations on the Server**

**For any non-trivial dataset, perform pagination, sorting, and filtering on the server.** The `Table` should be a "controlled" component that reflects server state.

*   **Why?** Client-side sorting on a single page of data is functionally incorrect and misleading. The `onChange` handler must trigger an API call to fetch the correctly sorted and filtered data from the backend.

### **3. Customize Rendering with `columns.render`**

**Use** the `render` function in a column definition to customize a cell's content, such as formatting dates or adding action buttons.

*   **Why?** This is the designated API for custom cell rendering, keeping display logic cleanly organized within the column configuration.

### **4. Use Virtualization for Large, Unpaginated Lists**

**Set** the `virtual` prop to `true` when displaying very long lists of data without pagination.

*   **Why?** This prevents rendering thousands of DOM nodes at once, which can cause the UI to freeze. Virtual scrolling dramatically improves performance.

---

## **Performance & Theming**

### **1. Use Selective Imports for Tree Shaking**

**Always** import components and their styles from their specific ES module path.

*   **Why?** Using default barrel file imports (`import { Button } from 'antd';`) prevents tree shaking, resulting in the *entire* library being included in your final bundle. This is the single biggest cause of performance issues.

*   **Example:**

    ```javascript
    // ❌ Bad: Pulls in the entire library, massive bundle size.
    import { Button } from 'antd';

    // ✅ Good: Only the code for Button and its styles are bundled.
    import Button from 'antd/es/button';
    import 'antd/es/button/style/css';
    ```

### **2. Prune the Icon Library**

**Avoid** including the entire `@ant-design/icons` package if you only use a few icons.

*   **Why?** The full icon package is very large. For applications needing only a small subset, it's more performant to import SVGs directly or use a lighter library.

### **3. Customize Styles with `ConfigProvider`**

**Use** the `ConfigProvider` and its `theme` prop to apply global style changes. Avoid brittle CSS overrides with `!important`.

*   **Why?** This is the intended, maintainable way to customize the look and feel in v5+. It allows for dynamic theme switching and ensures consistency.

### **4. Code-Split Non-Critical Components**

**Use** `React.lazy()` to dynamically import components that are not required for the initial page render (e.g., modals, complex components below the fold).

*   **Why?** This reduces the initial JavaScript payload, leading to faster page load times.

---

## **Accessibility (a11y)**

### **1. Use `App.useApp` for Simple Dialogs**

**Prefer** using `App.useApp()` for simple, non-disruptive feedback like success messages.

*   **Why?** This avoids the accessibility complexity of managing focus for a full `Modal` component when a simpler `message` or `notification` would suffice, aligning with the "Meaningful" design principle.

---

## **Common Anti-Patterns to Avoid**

### **1. The "Import All" Fallacy**

*   **Antipattern:** Using `import { Component } from 'antd';`.
*   **Harm:** Disables tree shaking, leading to a massive final bundle size.
*   **Solution:** **Always use selective imports:** `import Button from 'antd/es/button';`

### **2. The "Portal-per-Item" Mistake**

*   **Antipattern:** Rendering a component that uses a Portal (like `Popover` or `Tooltip`) inside every item of a large list.
*   **Harm:** Creates thousands of hidden DOM nodes, leading to catastrophic rendering performance.
*   **Solution:** Build a custom component with a single, shared portal whose content is updated on interaction.

### **3. Client-Side Sorting of Server-Side Data**

*   **Antipattern:** Applying a client-side `sorter` function to a `Table` that is using server-side pagination.
*   **Harm:** Misleads the user by only sorting the data on the current page, not the entire dataset.
*   **Solution:** Handle all sorting on the server via API calls triggered by the `onChange` event.

### **4. Styling Portals with CSS Cascade**

*   **Antipattern:** Trying to style a `Modal` or `Select` dropdown with CSS selectors from its parent.
*   **Harm:** These components render into a Portal outside the parent's DOM tree, so CSS cascade does not apply.
*   **Solution:** Use the `popupClassName` prop to target the portal element directly, or use `ConfigProvider` to apply theme tokens.
