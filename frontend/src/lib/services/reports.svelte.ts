import { getAuthHeaders } from "./auth.svelte";

export type Report = {
  id: number;
  reason: string;
  reporter_id: number;
  message_id: number;
  created_at: string;
};

export type ReportQueryResult = {
  total: number;
  count: number;
  reports: Report[];
};

export async function createReport(
  messageId: number,
  reason: string
): Promise<string> {
  const res = await fetch(`/api/v1/reports`, {
    method: "POST",
    body: JSON.stringify({
      message_id: messageId,
      reason: reason,
    }),
    headers: getAuthHeaders(),
  });

  const json = await res.json();

  return json["message"] || json["error"];
}

export async function deleteReport(reportId: number) {
  await fetch(`/api/v1/reports?reportId=${reportId}`, {
    method: "DELETE",
    headers: getAuthHeaders(),
  });
}

export async function resolveReport(reportId: number) {
  await fetch(`/api/v1/reports/resolve?reportId=${reportId}`, {
    method: "PATCH",
    headers: getAuthHeaders(),
  });
}

export async function queryReports(
  search?: string,
  pageSize: number = 20,
  page: number = 0
): Promise<ReportQueryResult> {
  const params = new URLSearchParams();

  if (search) {
    params.append("search", search);
  }

  params.append("page_size", pageSize.toString());
  params.append("page", page.toString());

  const res = await fetch(`/api/v1/query/reports?${params.toString()}`, {
    method: "GET",
    headers: getAuthHeaders(),
  });

  if (!res.ok) {
    throw new Error(`HTTP error! status: ${res.status}`);
  }

  const json = await res.json();

  return json["result"];
}
