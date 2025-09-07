/**
 * Review Profile System Types
 */

export interface ProfileMatches {
  area_type?: string[];
  name_contains?: string[];
  technologies?: string[];
}

export interface Profile {
  path: string;
  type: 'framework' | 'pattern' | 'language' | 'system';
  name: string;
  keywords: string[];
  extends?: string[];
  priority?: number;
  matches?: ProfileMatches;
  content: string;
}

export interface ProfileChain {
  profiles: Profile[];
  getMergedGuidelines(): string;
  getNames(): string[];
}