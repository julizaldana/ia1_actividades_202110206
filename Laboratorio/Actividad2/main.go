package main

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/ichiban/prolog"
)

func main() {
	// Crear una nueva aplicación Fyne
	myApp := app.New()
	myWindow := myApp.NewWindow("Clase 2 - Prolog")

	// Crear elementos de la interfaz
	output := widget.NewMultiLineEntry()
	output.SetPlaceHolder("Resultados de Prolog aparecerán aquí...")
	output.SetMinRowsVisible(10)
	output.Disable()

	queryInput := widget.NewEntry()
	queryInput.SetPlaceHolder("Ingrese una consulta: ej: hermano(maria, jose).")

	// Crear nueva máquina Prolog
	p := prolog.New(os.Stdin, os.Stdout)

	loadButton := widget.NewButton("Cargar archivo .pl", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				output.SetText(fmt.Sprintf("Error al abrir el archivo: %v", err))
				return
			}
			defer reader.Close()

			data := make([]byte, 1024*100) // 100 KB
			n, err := reader.Read(data)
			if err != nil {
				output.SetText(fmt.Sprintf("error al leer el archivo: %v", err))
				return
			}
			code := string(data[:n])
			output.SetText("Archivo cargado:\n" + reader.URI().Name())

			// Cargar el código Prolog desde el archivo
			if err := p.Exec(code); err != nil {
				output.SetText(fmt.Sprintf("error al cargar el código Prolog: %v", err))
			} else {
				output.SetText("Código Prolog cargado correctamente.")
			}
		}, myWindow)
	})

	queryButton := widget.NewButton("Consultar", func() {
		if p == nil {
			output.SetText("Por favor, cargue un archivo Prolog primero.")
			return
		}
		query := queryInput.Text
		solutions, err := p.Query(query)
		if err != nil {
			output.SetText(fmt.Sprintf("Error en la consulta: %v", err))
			return
		}
		defer solutions.Close()

		if solutions.Next() {
			output.SetText("Sí, la consulta es verdadera.")
		} else {
			output.SetText("No, la consulta es falsa.")
		}
	})

	// Organizar los widgets en un layout
	myWindow.SetContent(container.NewVBox(
		loadButton,
		queryInput,
		queryButton,
		output,
	))

	myWindow.Resize(fyne.NewSize(600, 400))
	myWindow.ShowAndRun()
}
