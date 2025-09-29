package Fachada

import (
	"fmt"
	"io"
	"os"
)

// TransmitirCancion se encarga de la lógica de negocio: abrir un archivo y leerlo en fragmentos.
// Utiliza una función callback 'enviarFragmento' para devolver cada fragmento a quien lo llamó (el controlador).
// De esta manera, la fachada no sabe nada sobre gRPC, solo sabe cómo enviar bytes.
func TransmitirCancion(titulo string, enviarFragmento func(fragmento []byte) error) error {
	fmt.Printf("Fachada: Procesando solicitud para la canción: %s\n", titulo)

	filePath := fmt.Sprintf("canciones/%s.mp3", titulo)

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error: No se pudo abrir el archivo %s: %v\n", filePath, err)
		return err
	}
	defer file.Close()

	buffer := make([]byte, 65536) // 64 KB
	fragmentoNum := 1

	for {
		bytesLeidos, err := file.Read(buffer)
		if err == io.EOF {
			fmt.Println("Fachada: Fin del archivo. Streaming completado.")
			break
		}
		if err != nil {
			fmt.Printf("Fachada: Error leyendo el archivo: %v\n", err)
			return err
		}

		// Llamamos al callback proporcionado por el controlador para enviar el fragmento.
		if err := enviarFragmento(buffer[:bytesLeidos]); err != nil {
			fmt.Printf("Fachada: Error enviando fragmento a través del callback: %v\n", err)
			return err
		}

		fmt.Printf("Fachada: Fragmento #%d leído (%d bytes) y enviado al controlador...\n", fragmentoNum, bytesLeidos)
		fragmentoNum++
	}

	return nil
}
