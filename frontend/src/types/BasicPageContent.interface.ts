import type { MdNode } from "$lib/markdown/Markdown.interface.ts"

export interface BasicPageContent {
  Title?: string
  MarkdownContent: string
}

export interface ExtendedPageContent {
  Title?: string
  MarkdownContent: string
  ParsedMarkdown: MdNode[]
}