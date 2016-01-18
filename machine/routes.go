package machine

type HandlerRoutes struct {
	IncomingCallRoute			string    			// route to the welcome message
	RecordingSuccessRoute		string       		// invoked after message are recorded
	RecordingIncompleteRoute	string   	 		// invoked if a timeout occurs
	RecordingFailedRoute		string    			// invoked if the recording failed due to communication issues between Tropo and the AnsweringMachine
	CheckMessagesRoute			string				// endpoint to browse voice messages
	ConfigurationRoute			string				// endpoint to read runtime effective configuration
}

func GetDefaultRoutes() (routes *HandlerRoutes) {
	return &HandlerRoutes{ "/tropo", "/tropo/recordingSuccess", "/tropo/recordingIncomplete", "/tropo/recordingFailed", "/messages", "/conf" }
}
