package Controladores

import (
	// Importamos nuestra nueva fachada
	"Streamify/ServidorStreaming/Fachada"
	ss "Streamify/ServidorStreaming/servicios_streaming"
	"fmt"
)

// ServidorDeStreaming es el struct que implementa la interfaz del servicio.
type ServidorDeStreaming struct {
	ss.UnimplementedAudioServiceServer
}

func (s *ServidorDeStreaming) StreamAudio(req *ss.PeticionDTO, stream ss.AudioService_StreamAudioServer) error {
	fmt.Printf("Controlador: Petición remota recibida para la canción: %s\n", req.Titulo)

	// Definimos la función que la fachada usará para enviarnos los fragmentos.
	// Esta función SÍ conoce gRPC.
	enviarFragmentoCallback := func(fragmento []byte) error {
		res := &ss.FragmentoCancion{
			Data: fragmento,
		}
		return stream.Send(res)
	}

	// Llamamos a la fachada y le pasamos el título y nuestra función de callback.
	// El controlador ya no sabe cómo se leen los archivos.
	return Fachada.TransmitirCancion(req.Titulo, enviarFragmentoCallback)
}
