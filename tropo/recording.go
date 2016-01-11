// Copyright 2015, Stève Sfartz
// Licensed under the MIT License

package tropo


type recordingResultWrapper struct {
	RecordingResult `json:"result"`
}

type RecordingResult struct {
   SessionID string `json:"sessionId"`
   CallID string `json:"callId"`
   State string `json:"state"`
   SessionDuration int `json:"sessionDuration"`
   Sequence int `json:"sequence"`
   Complete bool `json:"complete"`
   Error interface{} `json:"error"`
   CalledID string `json:"calledid"`
   Actions struct {
		 Name string `json:"name"`
		 Attempts int `json:"attempts"`
		 Disposition string `json:"disposition"`
		 Confidence int `json:"confidence"`
		 Interpretation string `json:"interpretation"`
		 Utterance string `json:"utterance"`
		 Concept string `json:"concept"`
		 Value string `json:"value"`
		 XML string `json:"xml"`
		 Duration int `json:"duration"`
		 URL string `json:"url"`
	 } `json:"actions"`
}