# Data Model Documentation

## Overview

The FLE platform uses a file-based JSON storage system optimized for an LLM-driven architecture. Rather than storing static content, the model focuses on configurations, user progress tracking, and interaction history. This approach minimizes storage requirements while maximizing flexibility and personalization.

## Storage Strategy

The data model employs a hierarchical file system with three primary divisions:

1. **System Data**: Platform-wide configurations and registries
2. **Content Data**: Language-specific learning configurations and references
3. **User Data**: Individual learner profiles, progress, and history

Each user session operates in isolation with its own data directory, enabling horizontal scaling and simple backup/recovery. Static configurations are loaded once at startup and cached, while user data is accessed on-demand.

## Core Entities

### System-Level Entities

#### Language Registry
Defines available languages for learning with their specific characteristics and configuration paths. Each language entry includes:
- Language code (ISO 639-1 standard)
- Display names in native and interface languages
- Activation status
- Language-specific settings (gender systems, conjugation complexity, writing systems)
- Path to language-specific content

#### Activity Type Registry
Maps activity types to their backend handlers and frontend widgets. Each activity type defines:
- Unique identifier for the activity pattern
- Backend handler class name for execution logic
- Frontend widget component for UI rendering
- Streaming support capabilities
- Available evaluation methods
- Required LLM capabilities

#### Prompt Library
Centralized repository of reusable prompt components that can be referenced by activity configurations. Contains:
- Persona definitions with personality traits and behavior guidelines
- Evaluation rubrics with scoring criteria
- Constraint sets defining vocabulary and grammar limits
- Template fragments for common interaction patterns

#### Evaluation Rubrics
Standardized criteria for assessing student performance across different activity types and levels. Each rubric specifies:
- Evaluation dimensions (grammar, vocabulary, fluency, etc.)
- Scoring scales and thresholds
- Feedback generation guidelines
- Level-appropriate expectations

### Content Entities

#### Cursus (Learning Paths)
Complete educational journey from beginner to fluency. Each cursus contains:
- Unique identifier and language association
- Methodology approach (immersive, grammar-focused, conversational)
- Ordered sequence of courses
- Progression rules and prerequisites
- Difficulty curve configuration

#### Course
Individual learning module aligned with CEFR levels. Each course includes:
- CEFR level and sublevel designation
- Learning objectives to be mastered
- Activity distribution strategy
- Estimated completion time
- Prerequisites from previous courses
- Advancement criteria

#### Learning Objectives
Discrete, measurable learning goals that students must master. Each objective defines:
- Unique identifier and human-readable code
- Category classification (grammar, vocabulary, pronunciation, culture)
- CEFR level range where applicable
- Mastery criteria with specific thresholds
- Related concepts and common errors
- Assessment contexts required for demonstration

#### Activity Configuration
Templates for generating dynamic learning activities. Each configuration specifies:
- Activity type reference
- Target learning objectives
- CEFR level and difficulty parameters
- LLM generation instructions and constraints
- Evaluation method and criteria
- UI presentation settings
- Support and scaffolding options
- Success thresholds

#### Vocabulary Lists
Level-appropriate word collections used as constraints for activity generation. Contains:
- Words grouped by CEFR level
- Thematic categorization
- Frequency rankings
- Usage contexts and collocations
- Words to explicitly exclude at each level

#### Grammar Rules
Reference definitions for grammatical concepts used in constraint sets and evaluation. Includes:
- Rule identifiers and descriptions
- CEFR level associations
- Example patterns
- Common exceptions
- Related objectives

### User Data Entities

#### User Profile
Core identity and preference information for each learner. Contains:
- Session identifier (primary key)
- Optional email and authentication details
- Interface and target language preferences
- Learning style preferences
- Difficulty and pacing preferences
- Identified strength and weakness areas

#### Enrollment Record
Tracks user participation in specific learning paths. Maintains:
- Active cursus associations
- Current course position
- Enrollment timestamps
- Overall progress percentage
- Completion status

#### Progress Tracking
Detailed record of advancement through learning objectives. Tracks:
- Per-objective mastery levels (0.0 to 1.0 scale)
- Context-specific performance (affirmative, negative, interrogative)
- Practice recency and frequency
- Success rates and attempt counts
- Review scheduling status

#### Skill Assessment
Aggregate performance metrics across skill categories. Measures:
- Grammar proficiency
- Vocabulary breadth
- Conversational ability
- Listening comprehension
- Reading comprehension
- Writing capability
- Pronunciation accuracy

#### Activity History
Complete record of all learning activities attempted. Each entry includes:
- Activity configuration reference
- Instance identifier for deduplication
- Timestamp and duration
- Performance score
- Objective progression increments
- LLM token usage
- Generation parameters for reproducibility
- User behavior metrics (response time, hints used)

#### LLM Interaction Log
Detailed record of all LLM-mediated conversations and evaluations. Stores:
- Complete message exchanges
- System prompts and constraints used
- User responses with metadata
- Evaluation results and feedback
- Token consumption metrics
- Response timing information

#### Learning Session
Time-bounded periods of platform engagement. Records:
- Session start and end times
- Activities completed during session
- Objectives practiced
- Device and context information
- Aggregate performance metrics

#### Daily Statistics
Aggregated metrics for analytics and streak tracking. Compiles:
- Time spent learning
- Activities completed
- Objectives mastered
- Streak maintenance
- Performance trends

## Relationships

### Hierarchical Relationships

- **Language → Cursus**: Each cursus belongs to exactly one language
- **Cursus → Course**: Cursus contain ordered sequences of courses
- **Course → Objectives**: Courses specify required and optional objectives
- **Course → Activities**: Courses define activity distribution strategies

### Reference Relationships

- **Activity Configuration → Objectives**: Activities target specific learning objectives
- **Activity Configuration → Activity Type**: Each configuration references its type
- **User Progress → Objectives**: Progress tracked per objective
- **Activity History → Activity Configuration**: History references the template used

### Temporal Relationships

- **User → Enrollments**: Users can have multiple sequential or parallel enrollments
- **Enrollment → Course Progress**: Each enrollment tracks progress through courses
- **Activity History → LLM Interactions**: Activities may generate interaction logs

## Data Lifecycle

### Static Data
System and content configurations are:
- Loaded once at application startup
- Cached in memory for performance
- Updated only through deployment
- Versioned through git

### User Data
Individual learner data is:
- Created on first session access
- Updated continuously during learning
- Backed up periodically
- Retained based on activity recency

### Generated Content
LLM-generated activities are:
- Created on-demand per request
- Not stored in full (only parameters)
- Logged for analysis and debugging
- Regeneratable from seeds if needed

## Key Design Decisions

### File-Based Storage
Chosen over database for:
- Simplified deployment and maintenance
- Easy debugging and manual inspection
- Natural git versioning for content
- Straightforward backup and recovery
- No schema migration complexity

### Session-Based Identity
Using session codes rather than traditional accounts because:
- Removes friction from starting learning
- Enables easy sharing of progress
- Supports privacy-conscious users
- Allows gradual security upgrade

### Configuration Over Content
Storing generation instructions rather than exercises because:
- Infinite variety from finite configuration
- Easy updates affect all future content
- Reduced storage requirements
- Natural personalization capability
- Simplified content management

### Separate User Directories
Each user in their own directory for:
- Horizontal scaling simplicity
- Independent backup/recovery
- Clear data isolation
- Easy GDPR compliance

### Progress as Mastery
Tracking mastery levels rather than completion because:
- Better represents actual learning
- Enables spaced repetition
- Supports multiple practice paths
- More meaningful analytics

## Performance Considerations

### Caching Strategy
- System configurations cached at startup
- Active user sessions cached in memory
- LRU eviction for inactive sessions
- Content configurations cached indefinitely

### File Access Patterns
- Read-heavy for configurations (cached)
- Write-heavy for user progress (direct to file)
- Append-only for history logs
- Periodic aggregation for statistics

### Scalability Approach
- User data shardable by session prefix
- Static content served from CDN
- Stateless application servers
- File system replication for redundancy

## Migration Path

When ready to move to a database:
- File structure maps directly to tables
- JSON fields become JSONB columns
- Directories become foreign key relationships
- Session codes become primary keys
- History logs become time-series data

This design provides a simple starting point that can evolve into a more sophisticated system as the platform grows, while maintaining clear data organization and relationships throughout the journey.