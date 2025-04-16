package main

import (
	"regexp"
	"strconv"
	"strings"
)

func extraerDatos(texto string, tabla [][]string) map[string]interface{} {
	datos := make(map[string]interface{})

	// Expresiones regulares
	namePattern := regexp.MustCompile(`Nombre:\s([A-Za-záéíóúÁÉÍÓÚñÑ ]+)`)
	codePattern := regexp.MustCompile(`Código:\s([0-9]+)`)
	proyectPattern := regexp.MustCompile(`Proyecto Curricular:\s([0-9]+)`)
	carreraPattern := regexp.MustCompile(`Plan de Estudios:\s([0-9]+) - ([A-Za-záéíóúÁÉÍÓÚñÑ ]+)`)
	periodoPattern := regexp.MustCompile(`Horario de Clases Período ([0-9]+)-([0-9])`)
	materiaPattern := regexp.MustCompile(`([A-zzáéíóúÁÉÍÓÚñÑ ]*)\sDocente`)

	// Extraer datos del texto
	datos["nombre"] = namePattern.FindStringSubmatch(texto)[1]
	datos["codigo"] = codePattern.FindStringSubmatch(texto)[1]
	datos["proyecto"] = proyectPattern.FindStringSubmatch(texto)[1]
	datos["carrera"] = carreraPattern.FindStringSubmatch(texto)[2]

	periodo := periodoPattern.FindStringSubmatch(texto)
	periodoInt := []int{}
	for _, p := range periodo[1:] {
		val, _ := strconv.Atoi(p)
		periodoInt = append(periodoInt, val)
	}
	datos["periodo"] = periodoInt

	// Procesar materias y horarios
	materias := make(map[string]map[string]string)
	ocupado := map[string][]string{
		"Lunes":     {},
		"Martes":    {},
		"Miercoles": {},
		"Jueves":    {},
		"Viernes":   {},
		"Sabado":    {},
		"Domingo":   {},
	}

	for _, fila := range tabla {
		if strings.Contains(fila[0], "Cod.") {
			continue
		}
		materia := materiaPattern.FindStringSubmatch(strings.ReplaceAll(fila[1], "\n", " "))[1]
		dias := []string{"Lunes", "Martes", "Miercoles", "Jueves", "Viernes", "Sabado", "Domingo"}
		horarios := fila[5:12]

		materiaData := make(map[string]string)
		for i, horario := range horarios {
			if horario != "" {
				materiaData[dias[i]] = strings.Split(horario, "\n")[0]
				ocupado[dias[i]] = append(ocupado[dias[i]], materiaData[dias[i]])
			}
		}
		materias[materia] = materiaData
	}

	datos["materias"] = materias
	datos["ocupado"] = unirHorarios(ocupado)

	return datos
}

func unirHorarios(ocupado map[string][]string) map[string][]string {
	for dia, horarios := range ocupado {
		rangos := [][]int{}
		for _, h := range horarios {
			parts := strings.Split(h, "-")
			inicio, _ := strconv.Atoi(parts[0])
			fin, _ := strconv.Atoi(parts[1])
			rangos = append(rangos, []int{inicio, fin})
		}

		// Ordenar y combinar rangos
		rangos = combinarRangos(rangos)
		horariosCombinados := []string{}
		for _, r := range rangos {
			horariosCombinados = append(horariosCombinados, strconv.Itoa(r[0])+"-"+strconv.Itoa(r[1]))
		}
		ocupado[dia] = horariosCombinados
	}
	return ocupado
}

func combinarRangos(rangos [][]int) [][]int {
	if len(rangos) == 0 {
		return rangos
	}

	// Ordenar rangos por inicio
	for i := 0; i < len(rangos)-1; i++ {
		for j := i + 1; j < len(rangos); j++ {
			if rangos[i][0] > rangos[j][0] {
				rangos[i], rangos[j] = rangos[j], rangos[i]
			}
		}
	}

	// Combinar rangos
	resultado := [][]int{rangos[0]}
	for _, r := range rangos[1:] {
		ultimo := resultado[len(resultado)-1]
		if r[0] <= ultimo[1] {
			ultimo[1] = max(ultimo[1], r[1])
			resultado[len(resultado)-1] = ultimo
		} else {
			resultado = append(resultado, r)
		}
	}
	return resultado
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
