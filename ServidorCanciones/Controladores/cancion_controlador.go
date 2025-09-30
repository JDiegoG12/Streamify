package Controladores

import (
	"Streamify/ServidorCanciones/Fachada"
	sc "Streamify/ServidorCanciones/servicios_canciones"
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (s *ServidorDeCanciones) ConsultarCancion(ctx context.Context, req *sc.ConsultarCancionRequest) (*sc.Cancion, error) {
	// Requisito: Imprimir un "eco" en la consola del servidor por cada llamada remota.
	fmt.Printf("Petición remota recibida: ConsultarCancion para la canción ID: %d\n", req.IdCancion)

	// Llama a la fachada para obtener la canción.
	cancion := Fachada.ObtenerCancionPorId(req.IdCancion)

	// Manejo de error si la canción no existe.
	if cancion == nil {
		return nil, status.Errorf(codes.NotFound, "La canción con ID %d no fue encontrada", req.IdCancion)
	}

	// Devuelve la canción encontrada.
	return cancion, nil
}
