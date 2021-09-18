package server

//Chunk is a smaller piece of storage that contains the messages that are written to it.
//Is incomplete if the Chunk is being currently written to.
type Chunk struct {
	Name     string `json:"name"`
	Complete bool   `json:"complete"`
}
