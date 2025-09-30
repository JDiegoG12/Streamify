// El paquete Controladores es la capa de entrada para las peticiones gRPC del servidor
// de streaming. Su función es manejar las solicitudes RPC, interactuar con la capa de
// negocio (Fachada) y gestionar la comunicación de streaming de vuelta al cliente.
package Controladores

import (
	"Streamify/ServidorStreaming/Fachada"
	ss "Streamify/ServidorStreaming/servicios_streaming"
	"fmt"
)

// ServidorDeStreaming es el struct que implementa la interfaz 'AudioServiceServer'
// generada por gRPC a partir del archivo streaming.proto.
type ServidorDeStreaming struct {
	ss.UnimplementedAudioServiceServer
}

// StreamAudio es la implementación del método RPC de tipo "server-streaming".
// Cuando un cliente invoca este método, se establece un canal de comunicación persistente
// a través del cual el servidor puede enviar múltiples mensajes (fragmentos de audio)
// de vuelta al cliente.
//
// Parámetros:
//   - req: La petición inicial del cliente, que contiene el título de la canción deseada.
//   - stream: Un objeto especial proporcionado por gRPC que representa el stream del servidor.
//     Tiene el método 'Send()' para enviar mensajes al cliente.
//
// Retorna:
//   - Un error si la operación de streaming falla de manera irrecuperable. gRPC se
//     encargará de notificar al cliente.
func (s *ServidorDeStreaming) StreamAudio(req *ss.PeticionDTO, stream ss.AudioService_StreamAudioServer) error {
	// Requisito: Imprimir un "eco" para registrar la petición entrante.
	fmt.Printf("Controlador: Petición remota recibida para la canción: %s\n", req.Titulo)

	// Se define una función anónima (closure) que encapsula la lógica específica de gRPC.
	// Esta función sabe cómo tomar un fragmento de bytes y enviarlo a través del stream gRPC.
	enviarFragmentoCallback := func(fragmento []byte) error {
		// Se construye el mensaje de respuesta de gRPC.
		res := &ss.FragmentoCancion{
			Data: fragmento,
		}
		// Se utiliza el método 'Send' del stream para enviar el fragmento al cliente.
		return stream.Send(res)
	}

	// Se delega toda la lógica de negocio (leer el archivo, dividirlo en fragmentos, etc.)
	// a la fachada. El controlador le pasa dos cosas:
	//   1. Los datos que necesita la fachada (el título de la canción).
	//   2. La 'herramienta' (el callback) que la fachada debe usar para devolver los datos.
	// De esta manera, la fachada permanece completamente agnóstica a gRPC.
	return Fachada.TransmitirCancion(req.Titulo, enviarFragmentoCallback)
}
