import type { MdNode } from "$lib/markdown/Markdown.interface.ts"
import { db } from "$lib/wailsjs/go/models.ts"

export interface ExtendedReport extends db.Report {
  Date: string
  Time: string
  Selected: boolean
  ParsedMarkdown: MdNode[]
  Visible: boolean
}

export interface DatedReport {
  [key: string]: ExtendedReport[]
}