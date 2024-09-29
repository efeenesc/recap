import { reportStore } from "$lib/stores/ReportStore.ts";
import { GetReports, GetReportById } from "$lib/wailsjs/go/app/AppMethods.js";
import type { db } from "$lib/wailsjs/go/models.ts";
import type { ExtendedReport } from "../types/ExtendedReport.interface.ts";

const processReports = (reports: db.Report[] | null): ExtendedReport[] => {
    if (!reports) return []; // Return an empty array if reports is null
    return reports.map<ExtendedReport>((v: any) => {
        v.Time = new Date(v.Timestamp * 1000).toLocaleTimeString();
        return v;
    });
};

export const getReports = async (limit: number) => {
    const res = await GetReports(limit);
    const processed = processReports(res);
    reportStore.update((val) => (val = processed));
    return res;
};

export const getReportById = async (
    id: number
): Promise<ExtendedReport | null> => {
    const res = await GetReportById(id);
    if (!res) return null; // Return null if no report is found
    const processed = processReports([res]);
    return processed[0];
};

export const addReportToStore = async (
    id: number
): Promise<ExtendedReport | null> => {
    const res = await getReportById(id);
    if (!res) return null; // Return null if no report is found
    const newReport = res;

    reportStore.update((prev) => {
        return [newReport, ...prev];
    });

    return newReport;
};
