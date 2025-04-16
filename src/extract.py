import re

namepattern = re.compile(r"Nombre:\s([A-z][A-záéíóúÁÉÍÓÚñÑ ]*)")
codePattern = re.compile(r"Código:\s([0-9][0-9]*)")
proyectPattern = re.compile(r"Proyecto Curricular:\s([0-9]*)")
carrera = re.compile(r"Plan de Estudios:\s([0-9]*) - ([A-záéíóúÁÉÍÓÚñÑ ]*)")
periodo = re.compile(r"Horario de Clases Período ([0-9]*)-([0-9])")
materia = re.compile(r"([A-zzáéíóúÁÉÍÓÚñÑ ]*)\sDocente")



def unir_horarios(ocupado):
    def combinar_horarios(horarios):
        rangos = [tuple(map(int, h.split('-'))) for h in horarios]
        rangos.sort()
        combinados = []
        for inicio, fin in rangos:
            if combinados and combinados[-1][1] >= inicio:
                combinados[-1] = (combinados[-1][0], max(combinados[-1][1], fin))
            else:
                combinados.append((inicio, fin))
        return [(f"{inicio}-{fin}") for inicio, fin in combinados]

    for dia in ocupado:
        ocupado[dia] = combinar_horarios(ocupado[dia])

    return ocupado


def extraerDatos(data : str, tabla = []):
    datos = {}
    
    datos['nombre'] = namepattern.search(data).group(1)
    datos['codigo'] = codePattern.search(data).group(1)
    datos['proyecto'] = proyectPattern.search(data).group(1)
    datos['carrera'] = carrera.search(data).group(2)
    datos['periodo'] = [int(periodo.search(data).group(1)), int(periodo.search(data).group(2))]

    datos['materias'] = {}
    ocupado = {i:[] for i in ['Lunes', 'Martes', 'Miercoles', 'Jueves', 'Viernes', 'Sabado', 'Domingo']}

    for fila in tabla:
        if 'Cod.' in fila[0]:
            continue
        print(fila)
        _, name, _, _, _, lun, mar, mie, jue, vie, sab, dom = fila
        m_name = materia.search(name.replace('\n',' ')).group(1).strip()
        datos['materias'][m_name] = {i:j.split('\n')[0] for i, j in zip(['Lunes', 'Martes', 'Miercoles', 'Jueves', 'Viernes', 'Sabado', 'Domingo'], [lun, mar, mie, jue, vie, sab, dom]) if j!=''}
        for i in datos['materias'][m_name]:
            ocupado[i].append(datos['materias'][m_name][i])
    datos['Ocupado'] = unir_horarios(ocupado)
     
    return datos