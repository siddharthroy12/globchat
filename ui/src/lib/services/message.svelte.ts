import { getAuthHeaders } from "./auth.svelte";

export type Message = {
  id: number;
  text: string;
  image: string;
  thread_id: number;
  user_id: number;
  username: string;
  user_image: string;
  created_at: string;
  is_first: boolean;
};

export type MessageQueryResult = {
  total: number;
  count: number;
  messages: Message[];
};

export async function createMessage(threadId: number, text: string) {
  const res = await fetch(`/api/v1/messages`, {
    method: "POST",
    body: JSON.stringify({
      thread_id: threadId,
      text: text,
      image: "",
    }),
    headers: getAuthHeaders(),
  });

  const json = await res.json();

  return json["message"];
}

export async function deleteMessage(messageId: number) {
  const res = await fetch(`/api/v1/messages?messageId=${messageId}`, {
    method: "DELETE",
    headers: getAuthHeaders(),
  });
}

export async function getMessages(
  threadId: number,
  limit: number,
  messageId?: number,
  direction: "before" | "after" = "before"
): Promise<Message[]> {
  const res = await fetch(
    `/api/v1/messages?threadId=${threadId}&limit=${limit}&messageId=${messageId}&direction=${direction}`,
    {
      headers: getAuthHeaders(),
    }
  );

  const json = await res.json();

  return json["messages"];
}

export async function getMessageById(messageId: number): Promise<Message> {
  const res = await fetch(`/api/v1/messages/${messageId}`, {
    headers: getAuthHeaders(),
  });

  const json = await res.json();

  return json["message"];
}

export async function queryMessages(
  search?: string,
  pageSize: number = 20,
  page: number = 0
): Promise<MessageQueryResult> {
  const params = new URLSearchParams();

  if (search) {
    params.append("search", search);
  }

  params.append("page_size", pageSize.toString());
  params.append("page", page.toString());

  const res = await fetch(`/api/v1/query/messages?${params.toString()}`, {
    method: "GET",
    headers: getAuthHeaders(),
  });

  if (!res.ok) {
    throw new Error(`HTTP error! status: ${res.status}`);
  }

  const json = await res.json();

  return json["result"];
}
