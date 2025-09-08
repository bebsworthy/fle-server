# Activity Architecture

## Overview

Activities in FLE are dynamic, LLM-powered learning experiences rather than static exercises. Each activity is a combination of three key components:

1. **Activity Handler**: Backend logic that orchestrates the LLM and manages the activity lifecycle
2. **UI Widget**: Frontend component that presents the activity and captures student responses
3. **Configuration**: JSON-based templates that define how the activity should behave for different levels and objectives

This architecture enables infinite variety in exercises while maintaining consistent interaction patterns and evaluation methods.

## Core Concepts

### Activities as Templates, Not Content

Traditional language learning platforms store thousands of pre-written exercises. FLE takes a fundamentally different approach: we store *instructions for generating exercises*. When a student needs to practice the verb "être", the system doesn't retrieve a fixed set of sentences - it generates new ones tailored to their level, recent mistakes, and learning objectives.

This means:
- Every practice session can be unique
- Exercises automatically adapt to student level
- Content stays fresh and engaging
- New vocabulary or cultural contexts can be incorporated instantly

### The Role of LLMs

Language models serve three critical functions in our activity system:

1. **Content Generation**: Creating exercises, scenarios, and examples appropriate to the student's level
2. **Interactive Participation**: Playing roles in conversations, providing dynamic responses
3. **Evaluation and Feedback**: Assessing student responses and providing personalized guidance

The LLM acts as both a content creator and a language tutor, guided by carefully crafted prompts and constraints.

## Activity Lifecycle

### Generation Phase

When a student begins an activity:

1. The system selects an appropriate activity configuration based on their current course and objectives
2. The activity handler constructs a prompt combining:
   - The base template from the configuration
   - The student's current CEFR level
   - Specific learning objectives to target
   - Constraints (vocabulary limits, grammar focus, themes)
3. The LLM generates the activity content in a structured format
4. The handler validates and prepares the content for the UI

### Interaction Phase

The student interacts with the activity through a purpose-built UI widget:

- **One-shot activities** present a single challenge and wait for a response
- **Multi-step activities** maintain context across multiple exchanges
- **Guided activities** provide scaffolding that adapts based on student performance

Each interaction type has corresponding UI patterns - from simple text inputs to drag-and-drop interfaces to chat-like conversation flows.

### Evaluation Phase

Evaluation happens differently depending on the activity type:

- **Simple evaluations** use pattern matching or exact comparison
- **Complex evaluations** send the student's response back to the LLM for assessment
- **Progressive evaluations** happen at checkpoints during multi-step activities

The evaluation produces not just a score but actionable feedback that helps the student understand their mistakes and improve.

## Activity Types

### One-Shot Activities

These are discrete exercises with a single correct answer or set of acceptable answers:

- **Fill in the blanks**: Generate sentences with missing words
- **Translation**: Translate phrases between languages
- **Conjugation**: Provide the correct verb form
- **Multiple choice**: Select from generated options

One-shot activities are ideal for drilling specific grammar points or vocabulary. They can be evaluated immediately and definitively.

### Multi-Step Activities

These maintain state across multiple interactions:

- **Conversations**: Free-form dialogue with an AI persona
- **Story completion**: Progressively build a narrative
- **Task sequences**: Complete related exercises that build on each other

Multi-step activities better simulate real language use and allow for more natural practice of communication skills.

### Guided Activities

These provide structured support while maintaining interactivity:

- **Guided conversations**: Conversations with suggested responses at each turn
- **Scaffolded writing**: Progressive prompts that help construct longer texts
- **Supported reading**: Interactive comprehension with vocabulary help

Guided activities bridge the gap between structured exercises and free-form practice, providing support that can be gradually reduced as students improve.

## Evaluation Strategies

### Exact Matching

For exercises with clear correct answers (verb conjugations, specific vocabulary), the system performs simple string comparison. This is fast, deterministic, and suitable for foundational skills.

### Pattern Matching

For exercises with multiple acceptable answers, the system checks against patterns or lists of alternatives. This allows flexibility while maintaining objective evaluation.

### LLM Evaluation

For open-ended responses, conversations, and complex language production, the LLM evaluates based on multiple criteria:

- Grammatical accuracy
- Vocabulary appropriateness
- Communication effectiveness
- Adherence to task requirements
- Pronunciation (when audio is involved)

The LLM provides both quantitative scores and qualitative feedback, explaining what was done well and what needs improvement.

### Progressive Evaluation

In multi-step activities, evaluation happens continuously:

- Each student response is assessed in context
- The system tracks improvement over the course of the activity
- Difficulty and support adjust based on performance
- Final evaluation considers the entire interaction

## Support Systems

### Scaffolding

Activities can provide various levels of support:

- **Word banks**: Suggested vocabulary for constructing responses
- **Sentence templates**: Partially completed structures to fill in
- **Grammar hints**: Reminders about relevant rules
- **Translation aids**: On-demand translation of difficult words

Scaffolding can be static (always present) or adaptive (adjusting based on student struggle).

### Adaptive Difficulty

Activities adjust their difficulty in real-time:

- **Vocabulary complexity** increases or decreases based on comprehension
- **Sentence length** adapts to student capability
- **Grammar complexity** progresses as students master simpler forms
- **Response time** expectations adjust to student pace

This ensures students are always challenged but not overwhelmed.

### Feedback Mechanisms

Feedback is provided at multiple levels:

- **Immediate corrections** for clear errors
- **Gentle guidance** for partial mistakes
- **Explanations** of grammar rules when needed
- **Encouragement** and progress acknowledgment
- **Suggestions** for improvement

The tone and detail of feedback adapts to the student's level and learning style.

## Configuration Philosophy

### Why JSON Configuration?

Activity configurations are stored as JSON rather than code because:

1. **Editability**: Teachers and content creators can modify activities without programming knowledge
2. **Versioning**: Different configurations can exist for the same activity type
3. **A/B Testing**: Multiple configurations can be tested to find the most effective
4. **Portability**: Configurations can be shared, exported, and imported easily
5. **Separation of Concerns**: Logic stays in code, content strategy stays in configuration

### Configuration Components

Each activity configuration includes:

- **Metadata**: Name, level, objectives, estimated duration
- **LLM Instructions**: System prompts, constraints, output formats
- **Evaluation Criteria**: How to assess student responses
- **UI Settings**: Which features to enable, how to display content
- **Support Options**: What help to provide and when

### Reusability and Composition

Configurations are designed to be:

- **Reusable**: The same configuration can generate countless unique instances
- **Composable**: Configurations can reference shared components (personas, word lists, evaluation rubrics)
- **Extensible**: New parameters can be added without breaking existing configurations
- **Hierarchical**: Course-level settings can override or supplement activity-level settings

## Integration with Learning Objectives

### Objective Targeting

Each activity configuration specifies which learning objectives it addresses. This enables:

- Automatic activity selection based on what needs practice
- Progress tracking against specific skills
- Balanced practice across multiple objectives
- Focused remediation on weak areas

### Mastery Tracking

As students complete activities:

1. Their responses are analyzed for demonstration of specific skills
2. Mastery levels for relevant objectives are updated
3. The system identifies which objectives need more practice
4. Future activity selection prioritizes gaps in mastery

### Curriculum Alignment

Activities are tagged with CEFR levels and can be filtered by:

- Language skills (reading, writing, speaking, listening)
- Grammar topics
- Vocabulary themes
- Cultural contexts
- Practical situations

This ensures activities align with standard language learning frameworks while maintaining flexibility.

## Future Extensibility

The architecture is designed to accommodate:

### New Activity Types

Adding a new activity type requires:
1. Creating a handler class with generation and evaluation logic
2. Building a corresponding UI widget
3. Defining configuration templates

The existing infrastructure handles everything else.

### Multiple Languages

The system is language-agnostic:
- Activity handlers work with any language the LLM supports
- Configurations can be translated or adapted
- UI widgets are language-independent
- Evaluation criteria can be language-specific

### Advanced Features

The architecture can support:
- **Voice interaction** through speech-to-text and text-to-speech
- **Video scenarios** with comprehension activities
- **Collaborative activities** between multiple students
- **Real-world integrations** (ordering from actual restaurant menus, reading current news)
- **Gamification layers** on top of existing activities

## Summary

The FLE activity architecture represents a paradigm shift in language learning technology. By treating activities as dynamic templates rather than static content, leveraging LLMs for generation and evaluation, and maintaining clear separation between logic, presentation, and configuration, we create a system that is both powerful and maintainable.

This approach enables:
- Infinite content variety without infinite content creation
- Personalized learning experiences at scale
- Rapid iteration and improvement of teaching methods
- Easy extension to new languages and activity types

Most importantly, it keeps the focus on actual language learning rather than completing predetermined exercises, making the platform more engaging and effective for students at all levels.