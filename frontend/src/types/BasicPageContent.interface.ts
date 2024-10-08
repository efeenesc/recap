import type { MdNode } from "$lib/markdown/Parser.ts"

export interface BasicPageContent {
  Title?: string
  MarkdownContent: string
}

export interface ExtendedPageContent {
  Title?: string
  MarkdownContent: string
  ParsedMarkdown: MdNode[]
}