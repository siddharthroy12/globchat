import { getAuthHeaders } from "./auth.svelte";

export type User = {
  id: number;
  email: string;
  created_at: string;
  username: string;
  image: string;
  messages: number;
  is_admin: boolean;
};

export type UserQueryResult = {
  total: number;
  count: number;
  users: User[]; // Changed from 'messages' to 'users'
};

export type UserQueryResponse = {
  users: User[];
  search?: string;
  pagination: {
    total: number;
    count: number;
    page: number;
    page_size: number;
    total_pages: number;
    has_next: boolean;
    has_prev: boolean;
  };
};

export async function queryUsers( // Changed function name from 'queryMessages'
  search?: string,
  pageSize: number = 20,
  page: number = 0 // Changed from 'pageIndex' to 'page' to match Go handler
): Promise<UserQueryResponse> {
  const params = new URLSearchParams();

  if (search) {
    params.append("search", search);
  }

  params.append("page_size", pageSize.toString()); // Changed from 'pageSize' to 'page_size'
  params.append("page", page.toString()); // Changed from 'pageIndex' to 'page'

  const res = await fetch(`/api/v1/users/query?${params.toString()}`, {
    method: "GET",
    headers: getAuthHeaders(),
  });

  if (!res.ok) {
    throw new Error(`HTTP error! status: ${res.status}`);
  }

  const json = await res.json();

  return json; // Changed from json["result"] to just json since Go handler doesn't wrap in "result"
}
