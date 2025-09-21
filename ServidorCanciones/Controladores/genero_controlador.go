package Controladores

import (
	"Streamify/ServidorCanciones/Fachada"
	sc "Streamify/ServidorCanciones/servicios_canciones"
	"context"
	"fmt"
)

// ServidorDeCanciones es el struct que implementará la interfaz del servicio gRPC.
// Debe incluir 'UnimplementedServicioCancionesServer' para compatibilidad futura.
// Definimos el struct aquí, y sus métodos pueden estar en otros archivos del mismo paquete.
type ServidorDeCanciones struct {
	sc.UnimplementedServicioCancionesServer
}

// ListarGeneros implementa el método RPC para obtener la lista de géneros.
// Este método pertenece al struct 'ServidorDeCanciones'.
func (s *ServidorDeCanciones) ListarGeneros(ctx context.Context, req *sc.GetGenerosRequest) (*sc.GetGenerosResponse, error) {
	// Imprime en la consola del servidor para cumplir con el requisito de "eco".
	fmt.Println("Petición remota recibida: ListarGeneros")

	// Llama a la fachada para obtener los datos.
	generos := Fachada.ObtenerTodosLosGeneros()

	// Devuelve la respuesta en el formato definido en el .proto.
	return &sc.GetGenerosResponse{Generos: generos}, nil
}
