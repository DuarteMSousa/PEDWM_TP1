package websocket

import "encoding/json"

// ClientMessage representa uma mensagem recebida de um cliente WebSocket.
// O campo Type identifica o tipo de comando e o Payload contém os dados
// específicos do comando em formato JSON.
type ClientMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// ServerMessage representa uma mensagem de resposta enviada ao cliente
// após o processamento de um comando.
type ServerMessage struct {
	Type    string `json:"type"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}
