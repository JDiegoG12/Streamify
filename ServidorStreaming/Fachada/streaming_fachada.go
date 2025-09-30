// El paquete Fachada implementa la lógica de negocio principal para el servicio de streaming.
// Su responsabilidad es manejar la lectura de archivos de audio del sistema de archivos,
// dividirlos en fragmentos manejables (chunks), y pasarlos a la capa superior a través
// de un mecanismo de callback, manteniendo así un total desconocimiento del protocolo de red (gRPC).
package Fachada

import (
	"fmt"
	"io"
	"os"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TransmitirCancion orquesta la lectura de un archivo de audio y su envío fragmentado.
// Esta función implementa el patrón de Inversión de Control: en lugar de devolver
// los datos directamente, acepta una función ('enviarFragmento') como parámetro y la
// invoca por cada fragmento de datos que procesa.
//
// Parámetros:
//   - titulo: El nombre del archivo de la canción (sin extensión) a transmitir.
//   - enviarFragmento: Una función de callback proporcionada por el llamador (el Controlador).
//     Esta fachada ejecutará esta función por cada fragmento de audio leído del disco.
//
// Retorna:
//   - Un error si ocurre un problema durante la lectura del archivo.
//   - Retorna 'nil' si el streaming se completa con éxito o si es cancelado
//     ordenadamente por el cliente.
func TransmitirCancion(titulo string, enviarFragmento func(fragmento []byte) error) error {
	fmt.Printf("Fachada: Procesando solicitud para la canción: %s\n", titulo)

	// Construye la ruta al archivo de audio en la carpeta 'canciones'.
	filePath := fmt.Sprintf("canciones/%s.mp3", titulo)

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error: No se pudo abrir el archivo %s: %v\n", filePath, err)
		return err
	}
	defer file.Close() // Asegura que el archivo se cierre al finalizar la función.

	// Define un búfer para leer el archivo en fragmentos de 64 KB.
	// Este tamaño es un buen compromiso entre eficiencia de E/S y uso de memoria.
	buffer := make([]byte, 65536)
	fragmentoNum := 1

	for {
		// Lee el siguiente fragmento del archivo en el búfer.
		bytesLeidos, err := file.Read(buffer)

		// Comprueba si se llegó al final del archivo.
		if err == io.EOF {
			fmt.Println("Fachada: Fin del archivo. Streaming completado.")
			break // Termina el bucle de lectura de forma normal.
		}
		// Comprueba si hubo un error de lectura diferente al final del archivo.
		if err != nil {
			fmt.Printf("Fachada: Error leyendo el archivo: %v\n", err)
			return err
		}

		// Invoca el callback con el fragmento de datos leído. El control pasa
		// temporalmente al controlador para que envíe los datos por la red.
		if err := enviarFragmento(buffer[:bytesLeidos]); err != nil {
			// El callback puede devolver un error (por ejemplo, si el cliente se desconectó).
			// Es crucial inspeccionar este error.
			if st, ok := status.FromError(err); ok && st.Code() == codes.Canceled {
				// Si el error es una cancelación explícita del cliente, no es un fallo del
				// servidor. Se considera una terminación normal de la operación.
				fmt.Println("Fachada: El cliente canceló la conexión. Deteniendo el envío.")
				return nil // Se retorna 'nil' para indicar un final exitoso/esperado.
			}

			// Si el error es de otro tipo (la conexión de red se cayó, etc.),
			// sí se considera un fallo del servidor y se propaga el error.
			fmt.Printf("Fachada: Error inesperado al enviar fragmento: %v\n", err)
			return err
		}

		fmt.Printf("Fachada: Fragmento #%d leído (%d bytes) y enviado al controlador...\n", fragmentoNum, bytesLeidos)
		fragmentoNum++
	}

	// Si el bucle 'for' termina sin errores, el streaming fue exitoso.
	return nil
}
