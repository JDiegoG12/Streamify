// El paquete Controladores actúa como la capa de entrada para las peticiones gRPC.
// Su responsabilidad es recibir las solicitudes remotas, deserializar los datos,
// invocar la lógica de negocio a través de la Fachada, y serializar las respuestas
// o errores de vuelta al cliente.
package Controladores

import (
	"Streamify/ServidorCanciones/Fachada"
	sc "Streamify/ServidorCanciones/servicios_canciones"
	"context"
	"fmt"
)

// ServidorDeCanciones es el struct principal que implementa la interfaz del servicio gRPC
// 'ServicioCancionesServer' generada a partir del archivo .proto. Todos los métodos
// RPC del servicio (como ListarGeneros, ConsultarCancion, etc.) estarán asociados a este
// struct.
//
// Se embébe 'sc.UnimplementedServicioCancionesServer' para garantizar la compatibilidad
// hacia adelante. Si en el futuro se añaden nuevos métodos al servicio en el archivo .proto,
// el servidor seguirá compilando, devolviendo por defecto un error 'Unimplemented' para
// los nuevos métodos hasta que sean explícitamente implementados.
type ServidorDeCanciones struct {
	sc.UnimplementedServicioCancionesServer
}

// ListarGeneros es la implementación del método RPC para obtener la lista completa de géneros.
// Se activa cuando un cliente llama a este procedimiento remoto.
//
// Parámetros:
//   - ctx: El contexto de la llamada RPC, gestionado por gRPC.
//   - req: Un puntero al objeto de solicitud 'GetGenerosRequest'. En este caso,
//     está vacío ya que la operación no requiere parámetros de entrada.
//
// Retorna:
//   - Un puntero a 'GetGenerosResponse' que contiene la lista de todos los géneros.
//   - Siempre retorna 'nil' como error, ya que la operación no tiene puntos de fallo
//     previstos en esta capa (la lista de géneros siempre está disponible).
func (s *ServidorDeCanciones) ListarGeneros(ctx context.Context, req *sc.GetGenerosRequest) (*sc.GetGenerosResponse, error) {
	// Requisito del proyecto: Imprimir un "eco" en la consola del servidor para
	// registrar la recepción de la llamada remota.
	fmt.Println("Petición remota recibida: ListarGeneros")

	// Delega la tarea de obtener los datos a la capa de Fachada.
	generos := Fachada.ObtenerTodosLosGeneros()

	// Construye el objeto de respuesta definido en el .proto, lo puebla con los datos
	// obtenidos de la fachada y lo retorna. gRPC se encarga de la serialización.
	return &sc.GetGenerosResponse{Generos: generos}, nil
}
