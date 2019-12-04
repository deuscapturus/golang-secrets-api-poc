import { Keys } from "./Keys.js";

//model
var Token = {
  current: sessionStorage.getItem("token"),
  error: sessionStorage.getItem("tokenError"),
  info: JSON.parse(sessionStorage.getItem("tokenInfo")) || {},
  adminString: sessionStorage.getItem("adminString") || "",
  numKeys: sessionStorage.getItem("numKeys") || 0,
  getInfo: function(token) {
    return m
      .request({
        method: "POST",
        url: "token/info",
        data: { token: token }
      })
      .then(function(result) {
        Token.info = result;
        sessionStorage.setItem("tokenInfo", JSON.stringify(result));
        Token.adminString = result.admin == 1 ? "Admin" : "Non-Admin";
        sessionStorage.setItem("adminString", Token.adminString);
        Token.numKeys = Token.info.keys[0] == "ALL"
          ? "ALL"
          : Token.info.keys.length;
        sessionStorage.setItem("numKeys", Token.numKeys);
        Keys.getList();
      })
      .catch(function(e) {
        Token.error = e.message;
        sessionStorage.setItem("tokenError", e.message);
        Token.info = {};
        sessionStorage.setItem("tokenInfo", JSON.stringify({}));
      });
  }
};

export { Token };
