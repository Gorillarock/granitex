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

function update_client_2_output(words) {
    document.getElementById("client.2.output").textContent = words;
}

function goto_home() {
  let host = window.location.hostname;
  let port = window.location.port;

  let endpoint = `http://${host}:${port}`
  let txUrl = `${endpoint}/v1/tx`
  window.location.href = txUrl;
}

function update_client_2_output_msg(cMsg) {
    let e = document.getElementById("client.2.output.msg");
    e.textContent = cMsg;
    e.style.display = "block";
}

async function post_rx() {
    event.preventDefault();

    document.getElementById("client.2.output").textContent = "";
    document.getElementById("client.2.output.msg").textContent = "";
    document.getElementById("client.2.output.msg").style.display = "none";


    const pin = document.getElementById("client.2.box.pin").value.toString();

    let eMsg = await get_from_rx_handler(pin);
    if (eMsg === false) {
        console.log("Error: eMsg not obtained");
        return false
    }
    // decrypt eMsg
    const cMsg = decrypt(eMsg, pin);
    if (cMsg === "") {
        console.log("Error: failed to decrypt message");
        return false;
    }
    update_client_2_output("Succes:");
    update_client_2_output_msg(cMsg);

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

    let response = null;
    try {
      response = await fetch(url).then(resp => {
        if (!resp.ok) {
          let err = new Error("HTTP status code: " + resp.status + " msg: " + resp.statusText);
          err.response = resp;
          err.status = resp.status;
          throw err;
        }
        return resp;
      });
    } catch (err) {
      console.log("get_from_rx_handler: fetch error: " + err);
      update_client_2_output("Error: " + err);
      return false;
    }
    
    let respJson = await response.json();
    if (response.ok) {
      console.log("get_from_rx_handler: success");
    } else {
      console.log("get_from_rx_handler: failure");
      return false;
    }
    console.log("get_from_rx_handler response json: " + respJson);

    // try to get eMsg from response
    let eMsg = respJson.emsg;
    if (eMsg === "") {
      console.log("Error: eMsg not found");
      return false;
    }

    return eMsg;
}