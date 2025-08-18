export function googleLogin(callback: (res: any) => void) {
  //@ts-ignore
  if (!window.googlelogincallbacks) {
    //@ts-ignore
    window.googlelogincallbacks = [];
    //@ts-ignore
  }

  if (callback) {
    //@ts-ignore
    window.googlelogincallbacks.push(callback);
  }

  const buttonWrapper = document.createElement("div");
  buttonWrapper.id = "google-login-button-wrapper";

  buttonWrapper.style.display = "none";
  document.body.appendChild(buttonWrapper);
  // @ts-ignore
  google.accounts.id.renderButton(buttonWrapper, {
    theme: "outline",
    size: "large",
    click_listener: (e: Error) => {
      console.log(e);
    },
  });
  // @ts-ignore
  document
    .querySelector('#google-login-button-wrapper div[role="button"]')
    // @ts-ignore
    .click();
  document.body.removeChild(buttonWrapper);
}
