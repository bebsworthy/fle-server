---
type: framework
name: jotai
keywords: [jotai, state-management, atoms, react, state, atomic, reactive]
priority: 8
matches:
  technologies: [jotai, "jotai/utils"]
---

## **Jotai Best Practices: A Guideline for Scalable State Management**

### **1. The Guiding Philosophy: Embrace Atomicity**

Our primary goal is to write code that is **performant, scalable, and easy to maintain**. Jotai's atomic "bottom-up" approach is the key to achieving this. Unlike traditional state management libraries, Jotai avoids unnecessary re-renders by default.

These guidelines exist to help us leverage this power consistently. When a rule seems unclear, default to what best serves **granularity and clarity** for a future developer.

---

### **2. Core Principles & Best Practices (The "Do's")**

#### **Use Small, Focused Atoms**

**Rule:** **Define** state in the smallest possible units (atoms). Resist the urge to group unrelated state into a single, large object.

**Why:** This is the core of Jotai. Small atoms ensure that a component only re-renders when the specific piece of state it cares about actually changes. This leads to highly optimized performance automatically, without manual memoization.

**Example:**

```javascript
// ❌ Bad: A monolithic atom causes over-rendering.
// A component using user.name will re-render if settings.theme changes.
const globalStateAtom = atom({
  user: { name: 'John Doe' },
  settings: { theme: 'dark' },
  products: [],
});

// ✅ Good: State is granular and independent.
// Components subscribe to only what they need.
export const userNameAtom = atom('John Doe');
export const themeAtom = atom('dark');
export const productsAtom = atom([]);
```

#### **Isolate Business Logic in Writable Atoms**

**Rule:** **Encapsulate** business logic, validation, and state mutations within write-only or read-write atoms.

**Why:** This separates your business logic from your UI components, making it reusable, testable, and easier to manage. Components become cleaner, focusing only on dispatching actions and rendering state.

**Example:**

```javascript
// ❌ Bad: Business logic is mixed inside the component.
function AddTodo() {
  const [todos, setTodos] = useAtom(todosAtom);
  const [newTodoText, setNewTodoText] = useState('');

  const handleAddTodo = () => {
    if (newTodoText.trim() === '') return; // Validation logic
    const newTodo = { id: Date.now(), text: newTodoText, completed: false };
    setTodos((current) => [...current, newTodo]); // Update logic
    setNewTodoText('');
  };
  // ...
}

// ✅ Good: Logic is centralized in a reusable "action" atom.
export const addTodoAtom = atom(
  null, // write-only
  (get, set, text: string) => {
    if (text.trim() === '') return; // Validation logic
    const newTodo = { id: Date.now(), text, completed: false };
    const currentTodos = get(todosAtom);
    set(todosAtom, [...currentTodos, newTodo]); // Centralized update
  }
);

function AddTodo() {
  const addTodo = useSetAtom(addTodoAtom); // Clean, declarative API
  const [newTodoText, setNewTodoText] = useState('');

  const handleAddTodo = () => {
    addTodo(newTodoText);
    setNewTodoText('');
  };
  // ...
}
```

#### **Create Custom Hooks for Clean APIs**

**Rule:** **Combine** related atoms and actions into a custom hook (e.g., `useTodos`) to provide a clean, declarative API for your components.

**Why:** This abstracts away the implementation details of your state management. Components don't need to know about `useAtom`, `useSetAtom`, or `useAtomValue`. They simply use the well-named functions and values exposed by the hook, making the code easier to read and refactor.

**Example:**

```javascript
// src/hooks/useTodos.ts
import { useAtomValue, useSetAtom } from 'jotai';
import { todosAtom, addTodoAtom } from '../state/todos';

// ✅ Good: A custom hook provides a clean interface to the state.
export function useTodos() {
  return {
    todos: useAtomValue(todosAtom),
    addTodo: useSetAtom(addTodoAtom),
  };
}

// src/components/AddTodo.tsx
import { useTodos } from '../hooks/useTodos';

function AddTodo() {
  const { addTodo } = useTodos(); // Consuming component is clean and declarative
  // ...
}
```

#### **Use `jotai/utils` for Complex Scenarios**

**Rule:** **Leverage** the powerful utilities from `jotai/utils` to handle common complex patterns like dynamic lists, performance optimizations, and robust async state.

**Why:** These utilities are purpose-built to solve specific, challenging problems in a highly optimized way. Reinventing this logic yourself can lead to bugs and performance issues.

*   **`atomFamily`**: Use when you need to create a collection of atoms dynamically from a list of data (e.g., a list of users, where each user has their own state).
*   **`selectAtom`**: Use to subscribe a component to a small slice of a large object atom. The component will only re-render when that specific slice changes.
*   **`focusAtom`**: Use for read/write access to a nested property within a larger state object (e.g., binding form inputs to nested state).
*   **`splitAtom`**: Use for rendering dynamic lists. It is a critical performance utility that takes an atom of an array and returns an array of atoms, preventing re-renders of the entire list when a single item changes.
*   **`loadable`**: Use for robust async state management, giving you explicit control over loading, data, and error states without relying on React Suspense.

---

### **3. Common Anti-Patterns (The "Don'ts")**

#### **The Monolithic Atom Fallacy**

**Anti-Pattern:** Defining the entire application state in a single, large atom, mimicking a Redux store.

**Harm:** This completely negates Jotai's primary benefit of automatic, granular re-renders. Every component that uses any piece of the state will re-render whenever *any other piece* of the state changes.

**Preferred Alternative:** Break the state down into the smallest reasonable, independent atoms.

#### **Defining Atoms Inside Components**

**Anti-Pattern:** Creating a new atom inside the render function of a React component.

**Harm:** This causes an infinite render loop. Jotai tracks atoms by referential equality. Creating a new atom on every render tells Jotai the component depends on a "new" piece of state, which triggers a re-render, which creates another new atom, and so on.

**Preferred Alternative:** Always define atoms at the module level (outside any component). If an atom must be created dynamically, it **must** be memoized with `useMemo`.

```javascript
// ❌ Bad: Creates an infinite render loop.
function MyComponent({ initialValue }) {
  const valueAtom = atom(initialValue);
  const [value] = useAtom(valueAtom);
  // ...
}

// ✅ Good: Atom is defined once at the module level.
const valueAtom = atom(0);

// ✅ Also Good (for dynamic atoms): Memoize the atom creation.
function MyComponent({ initialValue }) {
  const valueAtom = useMemo(() => atom(initialValue), [initialValue]);
  // ...
}
```

#### **Heavy Computations in Derived Atoms**

**Anti-Pattern:** Placing an expensive, synchronous calculation directly inside a derived atom's `read` function.

**Harm:** The `read` function of a synchronous atom executes during the render phase. A heavy computation here will block the main thread, causing the UI to freeze and become unresponsive.

**Preferred Alternative:** Move the heavy computation into an async, write-only "action" atom. Trigger this action, perform the work outside the render cycle, and store the result in a separate "result" atom.

#### **Over-subscribing in Large Components**

**Anti-Pattern:** Creating a large, monolithic component that subscribes to many different, unrelated atoms.

**Harm:** This negates the benefits of Jotai's granular subscription model. The component becomes a performance bottleneck, re-rendering whenever *any* of its many state dependencies change.

**Preferred Alternative:** Follow the principle of small, focused components. Break down large components into smaller children, each subscribing only to the one or two atoms it absolutely needs to perform its function.