import { db } from "$lib/wailsjs/go/models.ts"

export interface ExtendedReport extends db.Report {
  Date: string
  Time: string
  Selected: boolean
}

export interface DatedReport {
  [key: string]: ExtendedReport[]
}