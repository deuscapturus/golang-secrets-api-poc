//component
var tISM = {};


//model
tISM = { 
	Key: function(data) {
		this.Id = m.prop(data.Id);
		this.Name = m.prop(data.Name);
		this.CreationTime = m.prop(data.CreationTime);
	},

	keys:  function(token) {
		return m.request({
			method: "POST",
		       	url: "key/list",
			data: { "token": token },
			type: tISM.Key
		})
	},
};

//controller
tISM.controller = function() {
	// when this controller is updated perform a new POST for data by calling message?
	var token = m.prop("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6MSwiZXhwIjo5OTk5OTk5OTk5OSwianRpIjoiNGI2dmk2aTV0NWphZCIsImtleXMiOlsiQUxMIl19.zi5Ei-ZJcl-8dUVW27vGZLIM6qOFfwk94r_aauqr-tQ")

	var keys = function() {
		tISM.keys(this.token)
	}
	return { token, keys }
};

//view
tISM.view = function(ctrl) {
	return m("div"), [
		m("ol", ctrl.keys().map( function(key, index) {
			return m("li", key.Id, key.Name, key.CreationTime)
		})),
		m("input", {
			onchange: m.withAttr("value", ctrl.token)
		})
	]
};



//initialize
m.mount(document.getElementById("app"), tISM);

