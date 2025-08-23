import { getAuthHeaders } from "./auth.svelte";

export type Message = {
  id: number;
  text: string;
  image: string;
  thread_id: number;
  reported: number;
  user_id: number;
  username: string;
  user_image: string;
  created_at: string;
  is_first: boolean;
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

export async function reportMessage(messageId: number) {
  const res = await fetch(`/api/v1/messages/report?messageId=${messageId}`, {
    method: "POST",
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
