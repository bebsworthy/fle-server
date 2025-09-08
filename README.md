# FLE - French Language Learning Platform

A modern, AI-powered platform for learning French as a foreign language, designed with extensibility to support multiple languages in the future.

## Overview

FLE (French as a Foreign Language) is a comprehensive language learning platform that combines structured curriculum with dynamic, AI-powered activities. The system is architected to support multiple learning methodologies and allows learners to progress at their own pace through customizable learning paths.

## Technical Architecture

### Frontend
- **Framework**: React
- **Styling**: Tailwind CSS v4
- **State Management**: Jotai
- **Data Fetching**: React Query (TanStack Query)
- **Communication**: WebSocket JSON-RPC API

### Backend
- **Language**: Go
- **Web Framework**: Gorilla (HTTP and WebSocket)
- **ORM**: GORM
- **API Protocol**: JSON-RPC over WebSocket

### Database
- **PostgreSQL** for persistent storage

## Authentication & User Management

### Session-Based Access
- **Anonymous Accounts**: Visitors automatically receive a session code
- **Session Persistence**: 
  - Code visible in URL for easy sharing/bookmarking
  - Stored in localStorage for client-side persistence
- **Progress Tracking**: All learning progress tied to session identifier

### Optional Security Features
- **Email Registration**: Users can optionally add email for account recovery
- **Two-Factor Authentication**: TOTP-based 2FA for enhanced security

## Educational Model

### Hierarchical Structure

```
Cursus (Learning Path)
  └── Course (CEFR Level)
      ├── Objectives (Learning Goals)
      └── Activities (Exercises)
```

### Core Components

#### Cursus
- Complete education path from beginner to fluent speaker
- Multiple cursus available per language
- Different methodologies supported
- Features:
  - Switch between cursus at any time
  - Complete multiple cursus in parallel
  - Track progress independently per cursus

#### Course
- Aligned with CEFR (Common European Framework of Reference for Languages)
- Standard levels: A1, A2, B1, B2, C1, C2
- Granular sublevel system:
  - **A1**: A1.1, A1.2
  - **A2**: A2.1, A2.2
  - **B1**: B1.1, B1.2, B1.3, B1.4, B1.5, B1.6
  - **B2**: B2.1, B2.2, B2.3, B2.4, B2.5, B2.6
  - **C1**: C1.1, C1.2, C1.3, C1.4, C1.5, C1.6
- Composed of Objectives and Activities

#### Objectives
- Define mastery requirements for course completion
- Clear metrics for evaluation
- Examples by level:
  - **A1**: être/avoir, regular verbs, gender agreement, articles
  - **A2**: passé composé, imparfait, future simple, subjunctive introduction
- Includes vocabulary mastery targets
- Courses composed from a catalog of reusable objectives

#### Activities
- Interactive exercises supporting one or multiple objectives
- Difficulty scales with course level
- Customizable per course requirements
- Composed from activity catalog

## Activity Types

### Current Activity Examples
- **Fill-in-the-Blank Conversations**: Interactive dialogue completion
- **Template-Based Speaking**: Structured oral practice
- **Word-Picture Matching**: Visual vocabulary building
- **Categorization Games**: Grammar and vocabulary classification
- **Virtual Guided Tours**: Contextual language immersion
- **Daily Routine Comparisons**: Practical language application

### Activity Features
- **Adaptive Difficulty**: Same activity structure scales from A1 to C2
- **Dynamic Content**: LLM-powered for varied, contextual experiences
- **Multi-Objective Support**: Single activity can target multiple learning goals
- **Customizable Parameters**: Vocabulary complexity, grammar focus, length

## Key Features

### Dynamic Course System
- Courses defined by objective sets and supporting activities
- Flexible composition from catalog components
- Real-time adaptation based on learner progress

### AI Integration
- LLM-powered conversational practice
- Dynamic activity generation
- Personalized feedback and corrections
- Context-aware difficulty adjustment

### Multi-Language Architecture
- Core system designed for language agnosticism
- Course configurations per language
- Shared activity frameworks across languages
- Extensible objective catalog

## Development Status

### ✅ Backend Core Implementation Complete

The backend infrastructure has been fully implemented with production-ready features:

- **WebSocket Server**: Real-time bidirectional communication with JSON-RPC 2.0
- **Session Management**: Human-friendly session codes (e.g., "happy-panda-42")
- **Validation Framework**: Fast-fail validation with detailed error messages
- **Concurrent Connections**: Hub pattern supporting 1000+ simultaneous connections
- **Comprehensive Testing**: 94.8% session coverage, 77.5% WebSocket coverage
- **Production Features**: Graceful shutdown, panic recovery, structured logging

### Server Quick Start

```bash
# Clone and navigate to server
cd server

# Build the server
make build

# Run the server
./bin/fle-server

# Or run with custom settings
PORT=9090 LOG_LEVEL=debug ./bin/fle-server
```

### API Endpoints
- `GET /health` - Health check
- `GET /ws` - WebSocket endpoint for JSON-RPC communication

## Roadmap

### Completed ✅
- [x] Core backend infrastructure (Go/Gorilla)
- [x] WebSocket JSON-RPC API implementation
- [x] Session management system with memorable codes
- [x] Comprehensive test suite

### In Progress 🚧
- [ ] PostgreSQL schema design with GORM
- [ ] React frontend foundation
- [ ] Basic CEFR A1 course content

### Planned 📋
- [ ] First set of interactive activities
- [ ] LLM integration for dynamic content
- [ ] Progress tracking and analytics
- [ ] Email registration and 2FA
- [ ] Additional language support

## Contributing

This project is currently in the planning phase. Contribution guidelines will be established once the initial architecture is implemented.

## License

[License to be determined]