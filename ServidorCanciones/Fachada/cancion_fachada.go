// El paquete Fachada implementa el patrón de diseño Facade. Actúa como un punto de
// entrada unificado y simplificado hacia la lógica de negocio y el acceso a datos.
// Los Controladores interactúan con esta fachada en lugar de hacerlo directamente
// con los repositorios, lo que desacopla las capas y centraliza la lógica.
package Fachada

import (
	"Streamify/ServidorCanciones/Acceso_Datos/Repositorios"
	sc "Streamify/ServidorCanciones/servicios_canciones"
)

// ObtenerCancionesPorIdGenero actúa como un intermediario para solicitar canciones
// filtradas por género desde la capa de acceso a datos.
//
// Parámetros:
//   - idGenero: El identificador del género por el cual filtrar las canciones.
//
// Retorna:
//   - Un slice de punteros a 'sc.Cancion' correspondientes al género solicitado.
func ObtenerCancionesPorIdGenero(idGenero int32) []*sc.Cancion {
	return Repositorios.ObtenerCancionesPorGenero(idGenero)
}

// ObtenerCancionPorId sirve como intermediario para buscar una canción específica
// por su ID en la capa de acceso a datos.
//
// Al igual que otras funciones de la fachada, centraliza el acceso y permite futuras
// expansiones de la lógica de negocio sin modificar los controladores o repositorios.
//
// Parámetros:
//   - idCancion: El identificador único de la canción a buscar.
//
// Retorna:
//   - Un puntero al objeto 'sc.Cancion' si se encuentra.
//   - 'nil' si la canción no existe en la fuente de datos.
func ObtenerCancionPorId(idCancion int32) *sc.Cancion {
	return Repositorios.ObtenerCancionPorId(idCancion)
}
