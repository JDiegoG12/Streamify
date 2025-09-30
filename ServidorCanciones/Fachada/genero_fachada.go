// El paquete Fachada implementa el patrón de diseño Facade. Actúa como un punto de
// entrada unificado y simplificado hacia la lógica de negocio y el acceso a datos.
// Los Controladores interactúan con esta fachada en lugar de hacerlo directamente
// con los repositorios, lo que desacopla las capas y centraliza la lógica.
package Fachada

import (
	"Streamify/ServidorCanciones/Acceso_Datos/Repositorios"
	sc "Streamify/ServidorCanciones/servicios_canciones"
)

// ObtenerTodosLosGeneros sirve como el punto de acceso principal para solicitar la lista
// completa de géneros musicales del sistema.
//
// Retorna:
//   - Un slice con todos los objetos 'sc.Genero' disponibles.
func ObtenerTodosLosGeneros() []*sc.Genero {
	return Repositorios.ObtenerGeneros()
}
