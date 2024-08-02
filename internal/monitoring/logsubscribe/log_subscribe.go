package logsubscribe

// Previous notes : ### WHEN MOVING TO GO => HANLDE HTTP ERRORS GRACEFULLY + LOGGING

// Log subscribe

// What will it manage:
//1) On open -> ws.send the log subscribe message
//2) On message -> ws.recv the log message (this is where you'll
//                  send the succesful log messages to the get tx function)
// 2.1) If the log message has a error -> skip it
// 3) on error -> handle the error gracefully
// 4) on close -> handle the close gracefully

//send request to get tx if no error -> check if there’s an instruction with 14 accounts -> decrypt instruction data’s first 8bytes  -> if matches create return webhook with name, symbol and ipfs url

//Import which websocket we wanna use aswell as OS and how we wanna load enviroment variables

func On_Message(websocket, message string) {
	// handle message
}

func On_Open(websocket string, event string) {
	// handle open
}

func Main() {
	// connect to websocket
}
