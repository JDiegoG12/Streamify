// Archivo: ServidorStreaming/Controladores/streaming_controlador.go
package Controladores

import (
	ss "Streamify/ServidorStreaming/servicios_streaming"
	"fmt"
	"io"
	"os"
)

// ServidorDeStreaming es el struct que implementa la interfaz del servicio.
type ServidorDeStreaming struct {
	ss.UnimplementedAudioServiceServer
}

// StreamAudio implementa el método RPC que lee un archivo y lo envía en fragmentos.
func (s *ServidorDeStreaming) StreamAudio(req *ss.PeticionDTO, stream ss.AudioService_StreamAudioServer) error {
	fmt.Printf("Petición remota recibida: StreamAudio para la canción: %s\n", req.Titulo)

	// Construimos la ruta al archivo de audio.
	// ¡Asegúrate de que los archivos MP3 estén en una carpeta 'canciones'
	// dentro del directorio del ServidorStreaming!
	filePath := fmt.Sprintf("canciones/%s.mp3", req.Titulo)

	// Abrimos el archivo.
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error: No se pudo abrir el archivo %s: %v\n", filePath, err)
		return err
	}
	defer file.Close()

	// Definimos un buffer para leer el archivo en fragmentos.
	// 64 KB es un tamaño de fragmento razonable para streaming de audio.
	buffer := make([]byte, 65536) // 64 * 1024
	fragmentoNum := 1

	// Bucle para leer y enviar el archivo fragmento por fragmento.
	for {
		// Leemos un fragmento del archivo.
		bytesLeidos, err := file.Read(buffer)
		if err == io.EOF {
			// Llegamos al final del archivo.
			fmt.Println("Fin del archivo. Streaming completado.")
			break
		}
		if err != nil {
			fmt.Printf("Error leyendo el archivo: %v\n", err)
			return err
		}

		// Creamos el mensaje gRPC con el fragmento leído.
		fragmento := &ss.FragmentoCancion{
			Data: buffer[:bytesLeidos],
		}

		// Enviamos el fragmento al cliente a través del stream.
		if err := stream.Send(fragmento); err != nil {
			fmt.Printf("Error enviando fragmento al cliente: %v\n", err)
			return err
		}

		fmt.Printf("Fragmento #%d leído (%d bytes) y enviando...\n", fragmentoNum, bytesLeidos)
		fragmentoNum++
	}

	return nil
}
