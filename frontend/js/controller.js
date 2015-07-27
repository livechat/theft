'use strict';

var app = angular.module('Hijacker', ['ngWebsocket']);

app.config(['$logProvider', function($logProvider){
    $logProvider.debugEnabled(true);
}])

app.factory('dataAccess', function($websocket){
	var ws = $websocket.$new({url:'ws://localhost:8080/inspector/ws', reconnect:true, protocols: []})

	return {
		websocket: ws
	};
});

app.controller('HijackerController', function HijackerCtrl($scope, dataAccess, $interval) {
	var MAX_BUFFER = 1024;

	$scope.logs = {
		raw: [],
		filtered: [],
		filter: {
			text: "",
			regexp: false,
		}
	};

	$scope.logs.filter.apply = function() {
		if (this.text){
			if (this.regexp){
				try {
					var expression = new RegExp(this.text);
				}catch(e){
					var expression = new RegExp(".");
				}

				$scope.logs.filtered = $scope.logs.raw.filter(function(log){
					return expression.test(log.text);
				});
			}else{
				$scope.logs.filtered = $scope.logs.raw.filter(function(log){
					return log.text.indexOf($scope.logs.filter.text) !== -1;
				});
			}
		}else{
			$scope.logs.filtered = $scope.logs.raw;
		}

		$scope.$apply();
	};

	$scope.logs.insert = function(data){
		var args = JSON.parse(data.log);
		var values = [];

		for (var i in args){
			values.push(args[i]);
		}


		this.raw.push( {text: values.join(" "), type: "log"});
	};

	$scope.logs.event = function(text, type){
		this.raw.push({text: text, type: type})
	}

	$scope.logs.clear = function(){
		this.raw = [];
		this.filtered = [];
	};

	$scope.captured = null;
	$scope.shouldRecapture = false;
	$scope.disconnected = true;

	$scope.command = {
		batch: {
			cmd: ""
		}
	}

	$scope.command.batch.apply = function(){
		dataAccess.websocket.$emit('command', {cmd: $scope.command.batch.cmd, echo: false, batch: true});
	}

	$scope.capture = function(index){
		if ($scope.shouldRecapture){
			for (var i = 0; i < $scope.hijackers.length; i++ ){
				if ($scope.hijackers[i].session == $scope.captured) {
					index = i;
					break;
				}
			}
		}

		if (index == null){ return }
		$scope.shouldRecapture = false;

		if ($scope.captured != $scope.hijackers[index].session){
			$scope.logs.clear()
			$scope.captured = $scope.hijackers[index].session;
		}

		dataAccess.websocket.$emit('inspect', {session: $scope.captured});
	}

	dataAccess.websocket.$on("hijackers", function(data){
		$scope.hijackers = data.hijackers;
		$scope.capture();
		$scope.$apply();
	});

	dataAccess.websocket.$on("hijacker", function(data){
		switch (data.kind) {
			case "register": 
				$scope.hijackers.push(data.hijacker);
				$scope.capture();
				break;

			case "unregister":
				$scope.hijackers = $scope.hijackers.filter (function (hijacker) {
					if (hijacker.session == $scope.captured){
						$scope.shouldRecapture = true;
						$scope.logs.event("::LEFT", "notice");
					}
					return hijacker.session != data.hijacker.session
				});
				break;
			case "delay":
				for (var i = 0; i< $scope.hijackers.length; i++){
					if ($scope.hijackers[i].session == data.hijacker.session){
						$scope.hijackers[i].delay = data.hijacker.delay;
					}
				}
		}

		$scope.$apply();
	});

	dataAccess.websocket.$on("log", function(data){
		$scope.logs.insert(data);
		$scope.logs.filter.apply();
	});

	dataAccess.websocket.$on('$open', function () {
		$scope.disconnected = false;
	});

	dataAccess.websocket.$on('$close', function(){
		$scope.disconnected = true;
		if ($scope.captured){
			$scope.shouldRecapture = true;
		}
	})
});

app.filter('microping', function() {
  return function(input) {
  	if (!input) {
  		return '-'
  	}

  	if (input > 1000){
  		return parseInt(input/1000) + 'ms'
  	}else{
  		return input + 'µs'
  	}
  };
});

app.filter('userAgent', function(){
	return function(input, modifier){
		switch (modifier){
			case 'browser':
				if (/Chrome\//.test(input)){ return "Chrome" }
				if (/Firefox\//.test(input)){ return "Firefox" }
				if (/Chromium\//.test(input)){ return "Chromium" }
				if (/Safari\//.test(input)){ return "Safari" }
				if (/;MSIE/.test(input)){ return "Internet Explorer" }
				if (/OPR\/|Opera\//.test(input)){ return "Opera" }
				return ""
			
			case 'version':
				return ""

			case 'os':
				return ""

			case 'mobile':
				if (/Mobi/.test(input)){
					return "Mobile"
				}

				return ""
			default:
				return input
		}

	}
});