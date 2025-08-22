import { getAuthHeaders } from "./auth.svelte";

export type Thread = {
  id: number;
  lat: number;
  long: number;
  user_id: number;
  user_name: string;
  user_image: string;
  created_at: string;
};

export async function fetchThreads(
  lat: number,
  long: number,
  km: number
): Promise<Thread[]> {
  const res = await fetch(`/api/v1/threads?lat=${lat}&long=${long}&km=${km}`, {
    headers: getAuthHeaders(),
  });

  const json = await res.json();

  return json["threads"];
}

export async function fetchRandomThread(): Promise<Thread> {
  const res = await fetch(`/api/v1/threads/random`, {
    headers: getAuthHeaders(),
  });

  const json = await res.json();

  return json["thread"];
}

export async function createThread(
  lat: number,
  long: number,
  message: string
): Promise<Thread> {
  const res = await fetch(`/api/v1/threads`, {
    method: "POST",
    body: JSON.stringify({
      lat,
      long,
      message,
    }),
    headers: getAuthHeaders(),
  });

  const json = await res.json();

  return json["thread"];
}

export async function deleteThread(threadId: number) {
  const res = await fetch(`/api/v1/threads?threadId=${threadId}`, {
    method: "DELETE",
    headers: getAuthHeaders(),
  });
}
