export function createOrder(order) {
    console.log("order: ", order)
    return fetch(process.env.REACT_APP_API_URL + "/order", {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(order)
    })
      .then(data => data.json())
   }