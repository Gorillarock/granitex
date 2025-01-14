// emulate DB document
let _db_document = {
    id: "",
    verify: "",
    answer: "",
    eMsg: "",
  };

  // emulation of response payload info to be used to construct the rx GET request
  let _serv_query_resp_payload = {
    id: "",
    question: "",
    verify: "",
  };

  // Generate the verification code, TODO: will occur in server.
  function gen_verify() {
    const array = new Uint8Array(16);
    return self.crypto.getRandomValues(array).join("");
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

  // Hashing to be used on the question
  function sha512(question, pin) {
    const hash = CryptoJS.SHA512(question + pin);
    return hash.toString();
  }

  // Simulate saving to the DB on the server
  function serv_save_to_db_emulation(id, verify, answer, eMsg) {
    _db_document.id = id;
    _db_document.verify = verify;
    _db_document.answer = answer;
    _db_document.eMsg = eMsg;
  }
  
  // simulate the server submission and response
  function serv_tx_post_emulation(question, answer, eMsg) {
    const verify = gen_verify();  // happens in server
    const id = crypto.randomUUID();  // happens in server

    // simulate saving to the DB on the server
    serv_save_to_db_emulation(id, verify, answer, eMsg);
    // double check the DB
    console.log("DB id: " + _db_document.id);
    console.log("DB verify: " + _db_document.verify);
    console.log("DB answer: " + _db_document.answer);

    // Build fake response to store globally for now. TODO: remove.
    this.id = id;
    this.question = question;
    this.verify = verify;
  }

  // Update UI with the URL response
  function log_url(id, question, verify) {
    console.log(`URL: path/v1/query?i=${id}&q=${question}&v=${verify}`);
  }

  // simulate the server response to query GET
  function serv_query_get_emulation(id, question, verify) {
    _serv_query_resp_payload.id = id;
    _serv_query_resp_payload.question = question;
    _serv_query_resp_payload.verify = verify;
  }

  // onclick event for submission from client 1
  function post_tx() {
    const cMsg = document.getElementById("client.1.box.cMsg").value.toString();
    const pin = document.getElementById("client.1.box.pin").value.toString();

    const eMsg = encrypt(cMsg, pin);
    const question = gen_question();  // Not stored in DB
    const answer = sha512(question, pin);  // Stored in DB

    console.log("post_tx");
    console.log("cMsg: " + cMsg);
    console.log("pin: " + pin);
    console.log("submitted eMsg: " + eMsg);
    console.log("submitted question: " + question);
    console.log("submitted answer: " + answer);

    const resp = new serv_tx_post_emulation(question, answer, eMsg);
    log_url(resp.id, resp.question, resp.verify);

    // TODO: remove. here is where we are building the simulated payload from the query GET resp
    serv_query_get_emulation(resp.id, resp.question, resp.verify);

    // Call func to actually post the data to the server
    post_to_tx_handler(question, answer, eMsg);

    // clear the input
    let form = document.getElementById("client.1.form");
    form.reset();

    let host = window.location.hostname;
    let port = window.location.port;
    let path = "/v1/rx";
    let url = `http://${host}:${port}${path}?i=${resp.id}&q=${resp.question}&v=${resp.verify}`;

    window.location.href = url;
    // Test for updating the URL query params dynamically
    // let params = new URLSearchParams(window.location.search);
    // params.set('i', resp.id);
    // params.set('q', resp.question);
    // params.set('v', resp.verify);
    // window.location.search = params.toString();
    // console.log("URL: " + window.location.search);
  }


  // CLIENT 2 below

  // simulate the server response to requesting the encrypted message
  function serv_rx_get_emulation(id, verify, answer) {
    console.log("serv_rx_get_emulation id: " + id);
    console.log("serv_rx_get_emulation verify: " + verify);
    console.log("serv_rx_get_emulation answer: " + answer);
    if (_db_document.id !== id) {
      console.log("Error: id does not match what is in db");
      return "";
    }
    if (_db_document.verify !== verify) {
      console.log("Error: verify does not match what is in db");
      return "";
    }
    if (_db_document.answer === answer) {
      return _db_document.eMsg;
    }
    return "";
  }

  function decrypt(eVal, key) {
    const dMsg = CryptoJS.AES.decrypt(eVal, key);
    return dMsg.toString(CryptoJS.enc.Utf8);
  }

  function update_client_2_output(cMsg) {
    document.getElementById("client.2.output").textContent = cMsg;
  }

  function post_rx() {
    const pin = document.getElementById("client.2.box.pin").value.toString();
    const id = _serv_query_resp_payload.id;
    const question = _serv_query_resp_payload.question;
    const verify = _serv_query_resp_payload.verify;
    const answer = sha512(question, pin);

    console.log("post_rx id: " + id);
    console.log("post_rx question: " + question);
    console.log("post_rx verify: " + verify);
    console.log("post_rx answer: " + answer);
    console.log("post_rx pin: " + pin);

    get_from_rx_handler(answer);
    const eMsg = serv_rx_get_emulation(id, verify, answer);
    if (eMsg === "") {
      console.log("Error: eMsg not found");
      return;
    }
    console.log("rx eMsg: " + eMsg);
    const cMsg = decrypt(eMsg, pin);
    console.log("rx cMsg: " + cMsg);
    update_client_2_output(cMsg);

  }

  function on_load() {
    const urlParams = new URLSearchParams(window.location.search);
    const id = urlParams.get("i");
    console.log("on_load id: " + id);
  }

  function post_to_tx_handler(question, answer, eMsg) {
    console.log("post_to_tx_handler");
    let host = window.location.hostname;
    console.log("host: " + host);
    let port = window.location.port;
    console.log("port: " + port);

    fetch(`http://${host}:${port}/v1/handler/tx`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        question: question,
        answer: answer,
        emsg: eMsg,
      }),
    })
    .then((response) => response.json())
    .then((json) => console.log(json));
  }

  function get_from_rx_handler(answer) {
    console.log("get_from_rx_handler");
    let host = window.location.hostname;
    console.log("host: " + host);
    let port = window.location.port;
    console.log("port: " + port);
    let path = "/v1/handler/rx";

    let current_params = new URLSearchParams(window.location.search);
    let id = current_params.get("i");
    let verify = current_params.get("v");

    let url = `http://${host}:${port}${path}?i=${id}&v=${verify}&a=${answer}`;
    console.log("get_from_rx_handler url: " + url);
    fetch(url)
    .then((response) => response.json())
    .then((json) => console.log(json));
  }