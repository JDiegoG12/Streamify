package Controladores

import (
	"Streamify/ServidorCanciones/Fachada"
	sc "Streamify/ServidorCanciones/servicios_canciones"
	"context"
	"fmt"
)

// ListarCancionesPorGenero implementa el método RPC para obtener las canciones de un género específico.
// Este método se adjunta al struct 'ServidorDeCanciones', que fue definido en 'genero_controlador.go'.
func (s *ServidorDeCanciones) ListarCancionesPorGenero(ctx context.Context, req *sc.GetCancionesPorGeneroRequest) (*sc.GetCancionesPorGeneroResponse, error) {
	// Imprime en la consola del servidor.
	fmt.Printf("Petición remota recibida: ListarCancionesPorGenero para el género ID: %d\n", req.IdGenero)

	// Llama a la fachada para obtener las canciones filtradas.
	canciones := Fachada.ObtenerCancionesPorIdGenero(req.IdGenero)

	// Devuelve la respuesta.
	return &sc.GetCancionesPorGeneroResponse{Canciones: canciones}, nil
}
