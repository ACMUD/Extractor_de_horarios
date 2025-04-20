from fastapi import FastAPI, UploadFile, File
from fastapi.middleware.cors import CORSMiddleware
from pdfutil import extraer_texto_y_tablas
from extract import extraerDatos
from pydantic import BaseModel

class Student(BaseModel):
    nombre: str
    codigo: int
    proyecto: int
    carrera: str
    periodo: tuple[int, int]
    materias: dict[str, dict[str, str]]
    ocupado: dict[str, list[str]]

app = FastAPI()

app.add_middleware(
    CORSMiddleware,
    allow_origins=['*'],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.post("/", response_model=Student)
async def materias(file: UploadFile = File(...)):
    resultado = extraer_texto_y_tablas(file.file)

    pagina = resultado[0]
    r:dict = extraerDatos(pagina['texto'], pagina['tablas'][0])

    return r