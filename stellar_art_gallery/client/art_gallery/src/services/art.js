export function getArt() {
  return fetch(process.env.REACT_APP_API_URL + "/art").then((data) =>
    data.json()
  );
}
