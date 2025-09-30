// El paquete Controladores actúa como la capa de entrada para las peticiones gRPC.
// Su responsabilidad es recibir las solicitudes remotas, deserializar los datos
// (lo cual gRPC hace automáticamente), invocar la lógica de negocio a través de la
// Fachada, y serializar las respuestas o errores de vuelta al cliente.
package Controladores

import (
	"Streamify/ServidorCanciones/Fachada"
	sc "Streamify/ServidorCanciones/servicios_canciones"
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ListarCancionesPorGenero es la implementación del método RPC definido en el archivo .proto.
// Funciona como un 'handler' que se ejecuta cuando un cliente invoca este procedimiento remoto.
// Este método está asociado al struct 'ServidorDeCanciones' para cumplir con la interfaz
// generada por gRPC.
//
// Parámetros:
//   - ctx: El contexto de la llamada RPC, proporcionado por gRPC. Contiene metadatos
//     y permite manejar timeouts o cancelaciones.
//   - req: Un puntero al objeto de solicitud 'GetCancionesPorGeneroRequest', que contiene
//     el ID del género a buscar.
//
// Retorna:
//   - Un puntero a 'GetCancionesPorGeneroResponse' con la lista de canciones.
//   - Un error, que en este caso siempre es 'nil' ya que se retorna una lista vacía
//     si no hay resultados, en lugar de un error.
func (s *ServidorDeCanciones) ListarCancionesPorGenero(ctx context.Context, req *sc.GetCancionesPorGeneroRequest) (*sc.GetCancionesPorGeneroResponse, error) {
	// Requisito del proyecto: Imprimir un "eco" en la consola del servidor para
	// denotar que se ha recibido una llamada a un método remoto.
	fmt.Printf("Petición remota recibida: ListarCancionesPorGenero para el género ID: %d\n", req.IdGenero)

	// Delega la responsabilidad de obtener los datos a la capa de Fachada.
	canciones := Fachada.ObtenerCancionesPorIdGenero(req.IdGenero)

	// Construye el objeto de respuesta definido en el .proto y lo retorna.
	// gRPC se encargará de serializarlo y enviarlo al cliente.
	return &sc.GetCancionesPorGeneroResponse{Canciones: canciones}, nil
}

// ConsultarCancion es la implementación del método RPC para buscar una canción por su ID.
//
// Parámetros:
//   - ctx: El contexto de la llamada RPC.
//   - req: Un puntero al objeto de solicitud 'ConsultarCancionRequest' que contiene el ID
//     de la canción a consultar.
//
// Retorna:
//   - Un puntero al objeto 'Cancion' si se encuentra.
//   - Un error gRPC específico si la canción no se encuentra. El uso de 'status.Errorf'
//     con 'codes.NotFound' permite al cliente identificar programáticamente el tipo de error.
func (s *ServidorDeCanciones) ConsultarCancion(ctx context.Context, req *sc.ConsultarCancionRequest) (*sc.Cancion, error) {
	// Requisito: Imprimir un "eco" en la consola del servidor.
	fmt.Printf("Petición remota recibida: ConsultarCancion para la canción ID: %d\n", req.IdCancion)

	// Delega la lógica de búsqueda a la Fachada.
	cancion := Fachada.ObtenerCancionPorId(req.IdCancion)

	// Manejo de un caso de negocio específico: el recurso no fue encontrado.
	if cancion == nil {
		// En lugar de devolver un error genérico, se construye un error gRPC estandarizado.
		// Esto es una buena práctica en APIs, ya que proporciona al cliente un contexto
		// claro sobre por qué falló la solicitud.
		return nil, status.Errorf(codes.NotFound, "La canción con ID %d no fue encontrada", req.IdCancion)
	}

	// Si la canción se encuentra, se retorna directamente. El tipo de retorno de la función
	// coincide con el objeto de respuesta definido en el .proto.
	return cancion, nil
}
