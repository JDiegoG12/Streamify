package Repositorios

import (
	sc "Streamify/ServidorCanciones/servicios_canciones"
)

// Simulamos una tabla de canciones en una base de datos.
var canciones = []*sc.Cancion{
	// Rock (Genero ID: 3)
	{Id: 1, Titulo: "De musica ligera", Artista: "Soda Stereo", AnioLanzamiento: 1990, Duracion: "3:31", Genero: &sc.Genero{Id: 3, Nombre: "Rock"}},
	{Id: 2, Titulo: "Tren al sur", Artista: "Los Prisioneros", AnioLanzamiento: 1990, Duracion: "5:39", Genero: &sc.Genero{Id: 3, Nombre: "Rock"}},
	{Id: 3, Titulo: "Flaca", Artista: "Andrés Calamaro", AnioLanzamiento: 1997, Duracion: "4:32", Genero: &sc.Genero{Id: 3, Nombre: "Rock"}},
	{Id: 4, Titulo: "Lamento boliviano", Artista: "Enanitos Verdes", AnioLanzamiento: 1994, Duracion: "3:51", Genero: &sc.Genero{Id: 3, Nombre: "Rock"}},
	{Id: 5, Titulo: "Afuera", Artista: "Caifanes", AnioLanzamiento: 1994, Duracion: "4:50", Genero: &sc.Genero{Id: 3, Nombre: "Rock"}},

	// Salsa (Genero ID: 1)
	{Id: 6, Titulo: "Pedro Navaja", Artista: "Ruben Blades", AnioLanzamiento: 1978, Duracion: "4:47", Genero: &sc.Genero{Id: 1, Nombre: "Salsa"}},
	{Id: 7, Titulo: "Lloraras", Artista: "Oscar D'León", AnioLanzamiento: 1989, Duracion: "3:42", Genero: &sc.Genero{Id: 1, Nombre: "Salsa"}},

	// Cumbia (Genero ID: 2)
	{Id: 8, Titulo: "Como te voy a olvidar", Artista: "Los Ángeles Azules", AnioLanzamiento: 1996, Duracion: "4:28", Genero: &sc.Genero{Id: 2, Nombre: "Cumbia"}},
}

// ObtenerCancionesPorGenero busca en la lista de canciones y devuelve solo las que coinciden con el ID del género.
func ObtenerCancionesPorGenero(idGenero int32) []*sc.Cancion {
	var cancionesFiltradas []*sc.Cancion
	for _, cancion := range canciones {
		if cancion.Genero.Id == idGenero {
			cancionesFiltradas = append(cancionesFiltradas, cancion)
		}
	}
	return cancionesFiltradas
}

// Devuelve nil si no se encuentra.
func ObtenerCancionPorId(idCancion int32) *sc.Cancion {
	for _, cancion := range canciones {
		if cancion.Id == idCancion {
			return cancion
		}
	}
	return nil
}
