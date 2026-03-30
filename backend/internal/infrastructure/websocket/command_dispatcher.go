package websocket

// CommandDispatcher é o componente central responsável por encaminhar
// mensagens WebSocket recebidas dos clientes para os respetivos handlers
// de comando. Cada tipo de mensagem (identificado pelo campo "type") é
// associado a um CommandHandler registado previamente.
//
// Esta abordagem permite desacoplar a camada de transporte WebSocket da
// lógica de processamento de cada comando, facilitando a extensibilidade
// do sistema — para suportar um novo comando basta registar um novo handler.

import (
	"encoding/json"
	"fmt"
	"log"
)

// CommandHandler define a assinatura de uma função capaz de processar
// um comando recebido via WebSocket. Recebe o contexto do comando
// (identificação do jogador, sala e referência ao cliente) e o payload
// JSON específico do comando.
type CommandHandler func(ctx *CommandContext, payload json.RawMessage) error

// CommandContext encapsula a informação contextual associada a um comando
// recebido: o identificador do jogador, a sala a que pertence e a
// referência ao cliente WebSocket que originou o comando.
type CommandContext struct {
	PlayerID string
	RoomID   string
	Client   *Client
}

// CommandDispatcher mantém um registo de handlers indexados pelo tipo
// de mensagem e encaminha cada mensagem recebida para o handler adequado.
type CommandDispatcher struct {
	handlers map[string]CommandHandler
}

// NewCommandDispatcher cria uma nova instância do dispatcher sem handlers
// registados.
func NewCommandDispatcher() *CommandDispatcher {
	return &CommandDispatcher{
		handlers: make(map[string]CommandHandler),
	}
}

// Register associa um tipo de mensagem a um CommandHandler. Se já existir
// um handler para o mesmo tipo, este é substituído.
func (d *CommandDispatcher) Register(messageType string, handler CommandHandler) {
	if d == nil || handler == nil {
		return
	}
	d.handlers[messageType] = handler
}

// Dispatch encaminha uma ClientMessage para o handler registado para o
// respetivo tipo. Caso o comando provoque um panic (situação comum nos
// comandos de domínio existentes), este é recuperado e devolvido como erro.
// Se não existir handler para o tipo de mensagem, é devolvido um erro.
func (d *CommandDispatcher) Dispatch(ctx *CommandContext, msg ClientMessage) (err error) {
	if d == nil {
		return fmt.Errorf("dispatcher not configured")
	}

	handler, ok := d.handlers[msg.Type]
	if !ok {
		return fmt.Errorf("unknown command: %s", msg.Type)
	}

	// Os comandos de domínio existentes utilizam panic para sinalizar erros.
	// Recuperamos aqui para devolver uma resposta de erro ao cliente em vez
	// de terminar a goroutine.
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[dispatcher] panic recovered for command %q: %v", msg.Type, r)
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("%v", r)
			}
		}
	}()

	return handler(ctx, msg.Payload)
}

// HandleMessage é o ponto de entrada principal invocado pelo ReadPump do
// cliente. Faz o parsing da mensagem JSON, constrói o contexto e despacha
// para o handler adequado. A resposta (sucesso ou erro) é enviada de volta
// ao cliente que originou o comando.
func (d *CommandDispatcher) HandleMessage(client *Client, raw []byte) {
	if d == nil || client == nil {
		return
	}

	var msg ClientMessage
	if err := json.Unmarshal(raw, &msg); err != nil {
		d.sendError(client, "parse_error", "invalid message format")
		return
	}

	if msg.Type == "" {
		d.sendError(client, "parse_error", "missing message type")
		return
	}

	ctx := &CommandContext{
		PlayerID: client.id,
		RoomID:   client.roomID,
		Client:   client,
	}

	if err := d.Dispatch(ctx, msg); err != nil {
		d.sendError(client, msg.Type, err.Error())
		return
	}

	d.sendSuccess(client, msg.Type)
}

func (d *CommandDispatcher) sendError(client *Client, msgType string, errMsg string) {
	resp := ServerMessage{
		Type:    msgType,
		Success: false,
		Error:   errMsg,
	}
	data, err := json.Marshal(resp)
	if err != nil {
		return
	}
	client.Enqueue(data)
}

func (d *CommandDispatcher) sendSuccess(client *Client, msgType string) {
	resp := ServerMessage{
		Type:    msgType,
		Success: true,
	}
	data, err := json.Marshal(resp)
	if err != nil {
		return
	}
	client.Enqueue(data)
}
