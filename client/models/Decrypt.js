import { Token } from "./Token.js";

//model
var Decrypt = {
  error: false,
  input: "",
  output: "",
  encoding: "armor",
  decrypt: function() {
    return m
      .request({
        method: "POST",
        url: "decrypt",
        data: {
          token: Token.current,
          encsecret: Decrypt.input,
          encoding: Decrypt.encoding
        },
        deserialize: function(value) {
          return value;
        }
      })
      .then(function(result) {
        Decrypt.output = result;
        Decrypt.error = false;
      })
      .catch(function(e) {
        Decrypt.output = e.message;
        Decrypt.error = true;
      });
  }
};

export { Decrypt };
