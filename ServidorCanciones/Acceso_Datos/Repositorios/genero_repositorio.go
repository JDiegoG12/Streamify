// El paquete Repositorios implementa el patrón de diseño Repository. Su función
// es actuar como la capa de acceso a datos (Data Access Layer - DAL) de la aplicación,
// abstrayendo el origen de los datos del resto del servidor. En esta implementación,
// se simula una base de datos utilizando slices en memoria.
package Repositorios

import (
	sc "Streamify/ServidorCanciones/servicios_canciones"
)

// generos es un slice en memoria que actúa como una tabla de "géneros" en una base
// de datos simulada. Contiene la lista maestra de todos los géneros musicales
// que el sistema reconoce.
var generos = []*sc.Genero{
	{Id: 1, Nombre: "Salsa"},
	{Id: 2, Nombre: "Cumbia"},
	{Id: 3, Nombre: "Rock"},
}

// ObtenerGeneros proporciona acceso a la lista completa de géneros musicales.
// Esta función abstrae el origen de los datos; las capas superiores simplemente
// solicitan la lista sin necesidad de saber si proviene de la memoria, un archivo
// o una base de datos remota.
//
// Retorna:
//   - Un slice con todos los objetos 'sc.Genero' definidos en el sistema.
func ObtenerGeneros() []*sc.Genero {
	return generos
}
