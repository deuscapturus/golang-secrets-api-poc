import { Keys } from "../models/Keys.js";

//component
var keysList = {
	oninit: function(vnode) {
		Keys.getList()
	},
	view: function(vnode) {
		return m("div", [
			m("table[class=table table-striped]", [
				m("thead", [
					m("tr", [
						m("td", "ID"),
						m("td", "Name"),
						m("td", "Created"),
						m("td", "Download")
					])	
				]),
				m("tbody", 
					Keys.available.map( function(key) {
						return (key.Id != "ALL") ? m("tr", [
							m("th[scope=row]", key.Id),
							m("td", key.Name),
							m("td", key.CreationTime),
							m("td", m("a[href=#]", { onclick: function() {
								Keys.getKey(key.Id)
							}}, "Download"))
						]) : null
					})
				)
			])
		])
	}
}

export { keysList }
