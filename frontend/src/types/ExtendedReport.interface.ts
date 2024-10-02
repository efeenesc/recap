import type { MdNode } from "$lib/markdown/MarkdownParser.ts"
import { db } from "$lib/wailsjs/go/models.ts"

export interface ExtendedReport extends db.Report {
  Date: string
  Time: string
  Selected: boolean
  ParsedMarkdown: MdNode[]
}

export interface DatedReport {
  [key: string]: ExtendedReport[]
}