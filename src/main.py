import pdfplumber
from extract import extraerDatos
import json

LIMIT = 4

import io

def extraer_texto_y_tablas(pdf_input):
    resultado = []

    if isinstance(pdf_input, (str, bytes, io.IOBase)):
        pdf = pdfplumber.open(pdf_input)
    else:
        raise ValueError("El argumento debe ser una ruta de archivo o un stream de bytes.")

    with pdf as pdf:
        for i, pagina in enumerate(pdf.pages):
            if i >= LIMIT:
                break
            texto = pagina.extract_text()
            tablas = pagina.extract_tables()

            resultado.append({
                'pagina': i + 1,
                'texto': texto.strip() if texto else "",
                'tablas': tablas
            })

    return resultado

if __name__ == '__main__':
    with open("horario.pdf", "rb") as f:
        resultado = extraer_texto_y_tablas(f)
    pagina = resultado[0]
    r = extraerDatos(pagina['texto'], pagina['tablas'][0])

    print(json.dumps(r, indent=4, ensure_ascii=False))