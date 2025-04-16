import pdfplumber
from io import BytesIO
import sys
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

def main():
    pdf = sys.stdin.buffer.read()
    print(json.dumps(extraer_texto_y_tablas(BytesIO(pdf))))

if __name__ == "__main__":
    main()
