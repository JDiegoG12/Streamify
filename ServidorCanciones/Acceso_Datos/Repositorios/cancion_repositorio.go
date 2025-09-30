// El paquete Repositorios implementa el patrón de diseño Repository. Su función
// es actuar como la capa de acceso a datos (Data Access Layer - DAL) de la aplicación,
// abstrayendo el origen de los datos del resto del servidor. En esta implementación,
// se simula una base de datos utilizando slices en memoria, lo que permite un desarrollo
// rápido sin la necesidad de una base de datos real.
package Repositorios

import (
	sc "Streamify/ServidorCanciones/servicios_canciones"
)

// canciones es un slice en memoria que actúa como una base de datos simulada o "mock".
// Contiene un conjunto predefinido de objetos 'Cancion' que la aplicación utilizará.
var canciones = []*sc.Cancion{
	// Género Rock (ID: 3)
	{Id: 1, Titulo: "De musica ligera", Artista: "Soda Stereo", AnioLanzamiento: 1990, Duracion: "3:31", Genero: &sc.Genero{Id: 3, Nombre: "Rock"}},
	{Id: 2, Titulo: "Tren al sur", Artista: "Los Prisioneros", AnioLanzamiento: 1990, Duracion: "5:39", Genero: &sc.Genero{Id: 3, Nombre: "Rock"}},
	{Id: 3, Titulo: "Flaca", Artista: "Andrés Calamaro", AnioLanzamiento: 1997, Duracion: "4:32", Genero: &sc.Genero{Id: 3, Nombre: "Rock"}},
	{Id: 4, Titulo: "Lamento boliviano", Artista: "Enanitos Verdes", AnioLanzamiento: 1994, Duracion: "3:51", Genero: &sc.Genero{Id: 3, Nombre: "Rock"}},
	{Id: 5, Titulo: "Afuera", Artista: "Caifanes", AnioLanzamiento: 1994, Duracion: "4:50", Genero: &sc.Genero{Id: 3, Nombre: "Rock"}},

	// Género Salsa (ID: 1)
	{Id: 6, Titulo: "Pedro Navaja", Artista: "Ruben Blades", AnioLanzamiento: 1978, Duracion: "4:47", Genero: &sc.Genero{Id: 1, Nombre: "Salsa"}},
	{Id: 7, Titulo: "Lloraras", Artista: "Oscar D'León", AnioLanzamiento: 1989, Duracion: "3:42", Genero: &sc.Genero{Id: 1, Nombre: "Salsa"}},

	// Género Cumbia (ID: 2)
	{Id: 8, Titulo: "Como te voy a olvidar", Artista: "Los Ángeles Azules", AnioLanzamiento: 1996, Duracion: "4:28", Genero: &sc.Genero{Id: 2, Nombre: "Cumbia"}},
}

// ObtenerCancionesPorGenero realiza una búsqueda en la fuente de datos simulada para
// encontrar todas las canciones que pertenecen a un género específico.
//
// Parámetros:
//   - idGenero: El identificador único del género a filtrar.
//
// Retorna:
//   - Un slice de punteros a 'sc.Cancion'. Si no se encuentra ninguna canción para
//     ese género, retorna un slice vacío, no nulo.
func ObtenerCancionesPorGenero(idGenero int32) []*sc.Cancion {
	var cancionesFiltradas []*sc.Cancion
	for _, cancion := range canciones {
		if cancion.Genero.Id == idGenero {
			cancionesFiltradas = append(cancionesFiltradas, cancion)
		}
	}
	return cancionesFiltradas
}

// ObtenerCancionPorId busca una única canción en la fuente de datos utilizando su
// identificador único.
//
// Parámetros:
//   - idCancion: El ID de la canción a buscar.
//
// Retorna:
//   - Un puntero al objeto 'sc.Cancion' si se encuentra.
//   - 'nil' si no existe ninguna canción con el ID especificado.
func ObtenerCancionPorId(idCancion int32) *sc.Cancion {
	for _, cancion := range canciones {
		if cancion.Id == idCancion {
			return cancion
		}
	}
	return nil
}
