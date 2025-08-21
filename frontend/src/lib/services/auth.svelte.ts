export type UserData = {
  id: number;
  email: string;
  image: string;
  messages: number;
  username: string;
  new_account: string;
};

export enum AuthenticationStatus {
  Unknown,
  LoggedIn,
  LoggedOut,
}

let authenticationStatus = $state(AuthenticationStatus.Unknown);
let userData: UserData | null = $state(null);

export function getAuthenticationStatus() {
  return authenticationStatus;
}

export function getUserData() {
  return userData;
}

export async function updateUserImageAndUsername(
  image: string,
  username: string
) {
  const res = await fetch("/api/v1/google/login", {
    method: "POST",
    body: JSON.stringify({
      username,
      image,
    }),
  });

  if (res.status === 200) {
    userData!.username = username;
    userData!.image = image;
  }
}

export async function checkAuthenticationStatus() {
  authenticationStatus = AuthenticationStatus.Unknown;
  const res = await fetch("/api/v1/user", {
    headers: getAuthHeaders(),
  });

  if (res.status == 200) {
    const json = await res.json();
    userData = json["account"];
    authenticationStatus = AuthenticationStatus.LoggedIn;
    return;
  }
  authenticationStatus = AuthenticationStatus.LoggedOut;
}

export function logout() {
  removeToken();
  authenticationStatus = AuthenticationStatus.LoggedOut;
}

export function getAuthHeaders(): HeadersInit {
  return {
    token: getToken(),
  };
}

function saveToken(token: string) {
  localStorage.setItem("token", token);
}

function getToken() {
  return localStorage.getItem("token") ?? "";
}
function removeToken() {
  return localStorage.removeItem("token");
}

export async function loginUsingGoogle(jwt: string) {
  const res = await fetch("/api/v1/google/login", {
    method: "POST",
    body: JSON.stringify({
      token: jwt,
    }),
  });

  if (res.status === 200) {
    const json = await res.json();

    userData = json["account"];
    saveToken(json["token"]);
    authenticationStatus = AuthenticationStatus.LoggedIn;
  }
}
