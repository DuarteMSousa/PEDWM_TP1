package events_infrastructure

// CommandDispatcher é o componente central responsável por encaminhar
// mensagens WebSocket recebidas dos clientes para os respetivos handlers
// de comando. Cada tipo de mensagem (identificado pelo campo "type") é
// associado a um CommandHandler registado previamente.
//
// Esta abordagem permite desacoplar a camada de transporte WebSocket da
// lógica de processamento de cada comando, facilitando a extensibilidade
// do sistema — para suportar um novo comando basta registar um novo handler.

import (
	"backend/internal/domain/events"
	"encoding/json"
	"fmt"
	"sync"
)

// CommandHandler define a assinatura de uma função capaz de processar
// um comando recebido via WebSocket. Recebe o contexto do comando
// (identificação do jogador, sala e referência ao cliente) e o payload
// JSON específico do comando.
type EventHandler func(payload json.RawMessage) error

// EventDispatcher mantém um registo de handlers indexados pelo tipo
// de mensagem e encaminha cada mensagem recebida para o handler adequado.
type EventDispatcher struct {
	mu       sync.RWMutex
	handlers map[string]EventHandler
}

var (
	eventDispatcherInstance *EventDispatcher
	onceEventDispatcher     sync.Once
)

func GetEventDispatcherInstance() *EventDispatcher {
	onceEventDispatcher.Do(func() {
		eventDispatcherInstance = &EventDispatcher{
			handlers: make(map[string]EventHandler),
		}
	})
	return eventDispatcherInstance
}

// Register associa um tipo de mensagem a um EventHandler. Se já existir
// um handler para o mesmo tipo, este é substituído.
func (d *EventDispatcher) Register(messageType string, handler EventHandler) {
	if d == nil || handler == nil {
		return
	}
	d.handlers[messageType] = handler
}

func (d *EventDispatcher) Dispatch(event events.Event) error {
	if d == nil {
		return fmt.Errorf("dispatcher not configured")
	}

	handler, ok := d.handlers[string(event.Type)]
	if !ok {
		return fmt.Errorf("unknown event type: %s", event.Type)
	}

	payload, ok := event.Payload.(json.RawMessage)
	if !ok {
		return fmt.Errorf("invalid payload type")
	}
	return handler(payload)
}

// HandleMessage é o ponto de entrada principal invocado pelo ReadPump do
// cliente. Faz o parsing da mensagem JSON, constrói o contexto e despacha
// para o handler adequado. A resposta (sucesso ou erro) é enviada de volta
// ao cliente que originou o comando.
func (d *EventDispatcher) HandleMessage(event events.Event) {
	if d == nil {
		return
	}

	d.Dispatch(event)
}
