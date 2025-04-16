package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
)

// Debes importar la función si está en otro paquete, por ahora asumimos que está en el mismo
// Si está en otro paquete, ej: "github.com/tuusuario/tuapp/analizador", usa:
// import "github.com/tuusuario/tuapp/analizador"

type PageData struct {
	Pagina int          `json:"pagina"`
	Texto  string       `json:"texto"`
	Tablas [][][]string `json:"tablas"`
}

func main() {
	pdfPath := "../Estudiantes (3).pdf"

	// Abrir PDF
	pdfFile, err := os.Open(pdfPath)
	if err != nil {
		log.Fatalf("Error abriendo el PDF: %v", err)
	}
	defer pdfFile.Close()

	// Ejecutar script Python
	cmd := exec.Command("python3", "extract_pdf.py")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalf("Error creando stdin: %v", err)
	}

	go func() {
		defer stdin.Close()
		_, err := pdfFile.WriteTo(stdin)
		if err != nil {
			log.Fatalf("Error escribiendo al stdin: %v", err)
		}
	}()

	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error ejecutando el script Python: %v", err)
	}

	// Parsear JSON
	var paginas []PageData
	err = json.Unmarshal(output, &paginas)
	if err != nil {
		log.Fatalf("Error parseando JSON: %v", err)
	}

	if len(paginas) == 0 {
		log.Fatalf("No se encontró ninguna página")
	}

	texto := paginas[0].Texto
	tabla := [][]string{}
	if len(paginas[0].Tablas) > 0 {
		tabla = paginas[0].Tablas[0]
	}

	// Llamar a extraerDatos
	datos := extraerDatos(texto, tabla)

	// Imprimir resultado
	//fmt.Printf("\n📄 Datos extraídos:\n%v\n", datos)

	jsonBytes, err := json.MarshalIndent(datos, "", "  ")
	if err != nil {
		log.Fatalf("Error convirtiendo a JSON: %v", err)
	}

	// Imprimir salida como JSON
	fmt.Println(string(jsonBytes))
}
