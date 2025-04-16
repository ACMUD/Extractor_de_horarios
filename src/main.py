import pdfplumber
from extract import extraerDatos
import json

LIMIT:int = 4

def extraer_texto_y_tablas(pdf_input: bytes):
    resultado = []

    pdf = pdfplumber.open(pdf_input)

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


#Zona test
if __name__ == '__main__':
    with open("Horario 2025 - 1.pdf", "rb") as f:
        resultado:dict = extraer_texto_y_tablas(f)
    pagina = resultado[0]
    r:dict = extraerDatos(pagina['texto'], pagina['tablas'][0])

    print(json.dumps(r, indent=4, ensure_ascii=False))