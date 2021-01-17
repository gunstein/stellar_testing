export function getArt() {
    return fetch("https://0v3jk.sse.codesandbox.io/art").then((data) =>
      data.json()
    );
  }
  