export function createOrder(order) {
  console.log("order: ", order);
  return fetch(process.env.REACT_APP_API_URL + "/order", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(order),
  }).then((data) => data.json());
}

export function getEventSourceSSEReceiveBigImageUrl() {
  const eventSource = new EventSource(
    process.env.REACT_APP_API_URL + "/stream"
  );

  return eventSource;
}

export function getBigFileUrlFromMemoAndDownloadKey(memo, downloadKey) {
  return fetch(
    process.env.REACT_APP_API_URL +
      "/order/" +
      memo +
      "/big_file_url/" +
      downloadKey
  ).then((data) => data.json());
}
