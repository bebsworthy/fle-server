/**
 * Type definitions for the Code Audit feature
 */

export interface CodeAuditOptions {
  file: string;
  source: string;
  output: string;
  minComplexity?: string;
  name?: string[];
  module?: string;
  focus?: 'quality' | 'architecture' | 'security' | 'performance';
}

export interface AuditFinding {
  severity: 'critical' | 'high' | 'medium' | 'low';
  category: string;
  issue: string;
  location?: string;
  recommendation: string;
  effort?: 'quick-win' | 'small' | 'medium' | 'large';
}

export interface AreaAuditResult {
  area_name: string;
  area_type: string;
  timestamp: string;
  health_score: number; // 1-10
  complexity_score: number;
  findings: {
    critical: AuditFinding[];
    high: AuditFinding[];
    medium: AuditFinding[];
    low: AuditFinding[];
  };
  metrics: {
    files_analyzed: number;
    lines_of_code: number;
    duplication_percentage?: number;
    test_coverage?: number;
    cyclomatic_complexity?: number;
  };
  strengths: string[];
  quick_wins: string[];
  action_plan: {
    priority: number;
    action: string;
    effort: string;
    impact: string;
  }[];
}

export interface AuditSummary {
  timestamp: string;
  total_areas_audited: number;
  average_health_score: number;
  critical_issues_count: number;
  high_issues_count: number;
  areas_by_health: {
    excellent: string[]; // 9-10
    good: string[];      // 7-8
    fair: string[];      // 5-6
    poor: string[];      // 3-4
    critical: string[];  // 1-2
  };
}