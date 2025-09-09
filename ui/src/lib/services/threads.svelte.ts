import { getAuthHeaders } from "./auth.svelte";

export type Thread = {
  id: number;
  lat: number;
  long: number;
  user_id: number;
  message: string;
  replies: number;
  username: string;
  user_image: string;
  created_at: string;
};

export async function fetchThreads(
  minLat: number,
  maxLat: number,
  minLong: number,
  maxLong: number,
  mine: boolean
): Promise<Thread[] | null> {
  let url = `/api/v1/threads?minLat=${minLat}&maxLat=${maxLat}&minLong=${minLong}&maxLong=${maxLong}`;
  if (mine) {
    url = `/api/v1/threads?mine`;
  }
  const res = await fetch(url, {
    headers: getAuthHeaders(),
  });

  const json = await res.json();

  if (json["error"]) {
    return null; // Null means too many items
  }

  return json["threads"];
}

export async function fetchRandomThread(): Promise<Thread> {
  const res = await fetch(`/api/v1/randomthread`, {
    headers: getAuthHeaders(),
  });

  const json = await res.json();

  return json["thread"];
}

export async function fetchThread(id: number): Promise<Thread> {
  const res = await fetch(`/api/v1/threads/${id}`, {
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

  if (json["error"]) {
    throw Error(json["error"]);
  }

  return json["thread"];
}

export async function deleteThread(threadId: number) {
  const res = await fetch(`/api/v1/threads?threadId=${threadId}`, {
    method: "DELETE",
    headers: getAuthHeaders(),
  });
}
