__hijack = function(){
	if (console._log){return}
	var session = parseInt(sessionStorage.getItem('hijack::session')) || parseInt(Math.random() * 1e16);
	sessionStorage.setItem('hijack::session', session.toString());
	var timeout_id = 0;
	var buffer = [];

	var command = function(cmd){
		var output;
		try {
			output = eval(cmd);
		}catch(e){
			output = e.toString();
		}

		return output;
	}

	var run = function(){
		timeout_id = 0;
		var ws = new WebSocket("ws://localhost:8080/hijacker/ws");

		ws.onopen = function(){
			if (!console._log){
				console._log = console.log;
			}

			ws.send(JSON.stringify({event: "info", data:{browser: window.navigator.userAgent, location: window.location.href, session: session}}));

			console.log = function(){
				console._log.apply(this, arguments);

				if (ws.readyState === 1){

					while (buffer.length){
						JSON.stringify({event: "log", data:{session: session, log: JSON.stringify(buffer.shift())}});
					}

					ws.send(JSON.stringify({event: "log", data:{session: session, log: JSON.stringify(arguments)}}));
				} else {
					if (buffer.length < 1024){
						buffer.push(arguments)
					}
				}
			};

			ws.onmessage = function(event){
				var frame = JSON.parse(event.data);

				switch (frame.event) {
					case 'command':
						console.log(command(frame.data.cmd));
						break;
				}
			};
		};

		ws.onclose = ws.onerror = function(){
			if (!timeout_id) {
				timeout_id = setTimeout(run, 5000);
			}
		};
	};

	run()
};

__hijack();