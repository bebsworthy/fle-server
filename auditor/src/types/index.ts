/**
 * Type definitions for the Codebase Auditor tool
 */

export interface PromptWithMetadata {
  content: string;
  name: string;
  partials: string[];
}

export interface AgentConfig {
  tool: 'claude' | 'gemini';
  prompts: { name: string; content: string; partials?: string[] }[];
  variables?: Record<string, any>;
  outputPath: string;
  workingDir: string;
  dryRun: boolean | string;
  description?: string;
}

export interface AgentResult {
  success: boolean;
  output?: string;
}

export interface GlobalOptions {
  cwd?: string;
  dryRun?: boolean | string;
  raw?: boolean;
  agentOutput?: string;
}

export interface DocumentationSummary {
  total_areas: number;
  overall_quality_score: number;
  overall_importance_score: number;
  mandatory_count: number;
  very_important_count: number;
  important_count: number;
  nice_to_have_count: number;
  not_important_count: number;
  undocumented_critical: number;
}

export interface FunctionalArea {
  name: string;
  area_type: string;
  doc_quality_score: number;
  doc_importance_score: number;
  doc_reason: string;
}

export interface DocumentationQuality {
  timestamp: string;
  modules_file: string;
  total_areas: number;
  overall_quality_score: number;
  overall_importance_score: number;
  functional_areas: FunctionalArea[];
  summary?: DocumentationSummary;
}

export interface Technology {
  name: string;
  version?: string;
  type: 'language' | 'framework' | 'library' | 'tool' | 'database' | 'service';
}

export interface Module {
  name: string;
  path: string;
  type: 'frontend' | 'backend' | 'library' | 'service' | 'mobile' | 'cli' | 'shared';
  main_technology: string;
  package_manager?: string;
  dependencies?: string[];
}

export interface ModuleInfo {
  name: string;
  description: string;
  module?: string;  // For monorepos
  area_type: string;
  technologies: Technology[];  // Per-area technologies
  source_references: string[];
  documentation_references: string[];
  confidence_level: 'high' | 'medium' | 'low';
  metrics: {
    file_count: number;
    estimated_loc: number;
    primary_language: string;
  };
  complexity_score: number;
  complexity_factors: string[];
  related_areas?: string[];
  // Legacy fields for compatibility
  type?: string;
  source_files?: string[];
  confidence?: number;
}

export interface GuidelineDocument {
  file_path: string;
  description: string;
}

export interface ModulesData {
  timestamp: string;
  modules?: Module[];  // For monorepos
  functional_areas: ModuleInfo[];
  guideline_documents?: GuidelineDocument[];
  project_name?: string;
}