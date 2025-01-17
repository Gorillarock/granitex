// Code unique to the transmitter client

// Hashing to be used on the question
function sha512(question, pin) {
    const hash = CryptoJS.SHA512(question + pin);
    return hash.toString();
}

// Generate the question code
function gen_question() {
    const array = new Uint8Array(16);  // TODO: Make len variable from 8-20
    return self.crypto.getRandomValues(array).join("");
}

// Only for encrypting the secret message
function encrypt(cMsg, pin) {
    const eMsg = CryptoJS.AES.encrypt(cMsg.toString(), pin.toString());
    return eMsg.toString();
}

function refresh() {
  window.location.reload();
}

// onclick event for submission from client 1
async function post_tx() {
    event.preventDefault();
    const cMsg = document.getElementById("client.1.box.cMsg").value.toString();
    const pin = document.getElementById("client.1.box.pin").value.toString();

    const eMsg = encrypt(cMsg, pin);
    const question = gen_question();  // Not stored in DB
    const answer = sha512(question, pin);  // Stored in DB

    // Call func to actually post the data to the server
    let response_url = await post_to_tx_handler(question, answer, eMsg);
    console.log("response_url: " + response_url);
    window.location.href = response_url;
}


async function post_to_tx_handler(question, answer, eMsg) {
    console.log("post_to_tx_handler");
    let host = window.location.hostname;
    let port = window.location.port;

    let endpoint = `http://${host}:${port}`
    let fetchUrl = `${endpoint}/v1/handler/tx`

    const response = await fetch(fetchUrl, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          question: question,
          answer: answer,
          emsg: eMsg,
        }),
      });

      const respJson = await response.json();
      let path = respJson.path;
      let response_url = `${endpoint}${path}`;
      console.log("response_url: " + response_url);
      return response_url;
}
