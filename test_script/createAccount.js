var http = require('http');

function createAccount() {
	var options = {
		hostname: "localhost",
		port: 8080,
		path: "/test/create-account-directly",
	}
	http.get(options, function(response) {
		//another chunk of data has been recieved, so append it to `str`
		var str = "";
		response.on('data', function (chunk) {
			str += chunk;
		});

		//the whole response has been recieved, so we just print it out here
		response.on('end', function () {
			console.log(str);
		});
	})
}

createAccount();
