// Code unique to the receiver client

// Hashing to be used on the question
function sha512(question, pin) {
    const hash = CryptoJS.SHA512(question + pin);
    return hash.toString();
}

function decrypt(eVal, key) {
    const dMsg = CryptoJS.AES.decrypt(eVal, key);
    return dMsg.toString(CryptoJS.enc.Utf8);
}

function update_client_2_output(cMsg) {
    document.getElementById("client.2.output").textContent = cMsg;
}

async function post_rx() {
    event.preventDefault();
    const pin = document.getElementById("client.2.box.pin").value.toString();

    let eMsg = await get_from_rx_handler(pin);

    // decrypt eMsg
    const cMsg = decrypt(eMsg, pin);
    update_client_2_output(cMsg);

}

async function get_from_rx_handler(pin) {
    console.log("get_from_rx_handler");
    let host = window.location.hostname;
    let port = window.location.port;
    let path = "/v1/handler/rx";

    let current_params = new URLSearchParams(window.location.search);
    let id = current_params.get("i");
    let verify = current_params.get("v");
    let question = current_params.get("q");
    let answer = sha512(question, pin);

    let url = `http://${host}:${port}${path}?i=${id}&v=${verify}&a=${answer}`;
    console.log("get_from_rx_handler url: " + url);
    const response = await fetch(url);

    const respJson = await response.json();
    if (response.ok) {
      console.log("get_from_rx_handler: success");
    } else {
      console.log("get_from_rx_handler: failure");
      return "";
    }
    console.log("get_from_rx_handler response json: " + respJson);

    // try to get eMsg from response
    let eMsg = respJson.emsg;
    if (eMsg === "") {
      console.log("Error: eMsg not found");
      return "";
    }

    return eMsg;
}