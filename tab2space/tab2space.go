/*
 * Copyright 2017 Alvaro Stagg [alvarostagg@protonmail.com]
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stdout, "%s: falta un operando\nPruebe '%s --help' para más información.\n", os.Args[0], os.Args[0])
		os.Exit(1)
	}

	vflag := flag.Bool("version", false, "muestra la versión y finaliza")
	deleteOld := flag.Bool("delete-old", false, "borra el archivo original")
	tabSize := flag.Int("tab-size", 4, "número de espacios")
	flag.Parse()

	const VERSION = "15.06.2017"
	if *vflag {
		fmt.Printf("%s v%s\n", os.Args[0], VERSION)
		os.Exit(0)
	}

	if *tabSize <= 0 {
		fmt.Fprintf(os.Stderr, "%s: el número de espacios debe ser mayor a 0.\n", os.Args[0])
		os.Exit(1)
	}

	for _, fileName := range flag.Args() {
		file, err := os.Open(fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: no se pudo abrir el archivo \"%s\". Error: %s\n", os.Args[0], fileName, err)
			continue
		}

		defer file.Close()

		fileStat, err := file.Stat()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: no se pudo obtener el tamaño de \"%s\". Omitiendo.\n", os.Args[0], fileName)
			continue
		}

		if fileStat.IsDir() {
			fmt.Fprintf(os.Stderr, "%s: omitiendo el directorio \"%s\".\n", os.Args[0], file.Name())
			continue
		}

		data := make([]byte, fileStat.Size())
		if _, err = file.Read(data); err != nil {
			fmt.Fprintf(os.Stderr, "%s: problema al leer \"%s\". Error: %s\n", os.Args[0], fileStat.Name(), err)
			continue
		}

		newContent := make([]byte, 0)
		for c := 0; c < len(data); c++ {
			if data[c] == '\t' {
				for i := 0; i < *tabSize; i++ {
					newContent = append(newContent, 32)
				}
				continue
			}

			newContent = append(newContent, data[c])
		}

		newFileName := getNewName(fileName)
		newFile, err := os.Create(newFileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: no se pudo crear el archivo \"%s\". Error: %s\n", os.Args[0], newFileName, err)
			continue
		}

		defer newFile.Close()

		if _, err = newFile.Write(newContent); err != nil {
			fmt.Fprintf(os.Stderr, "%s: no se pudo escribir el nuevo contenido en el nuevo archivo \"%s\". Error: %s\n", os.Args[0], newFileName, err)
			continue
		}

		fmt.Printf("%s: nuevo contenido escrito en el archivo \"%s\".", os.Args[0], newFileName)

		if *deleteOld {
			if err = os.Remove(fileName); err != nil {
				fmt.Fprintf(os.Stderr, "%s: No se pudo borrar el archivo original \"%s\". Error: %s\n", os.Args[0], fileName, err)
			} else {
				fmt.Printf("%s: archivo \"%s\" eliminado y remplazado con \"%s\".\n", os.Args[0], fileName, newFileName)
			}
		}
	}
}

func getNewName(fileName string) string {
	parts := strings.Split(fileName, ".")
	var newFileName string

	for i := 0; i < len(parts)-1; i++ {
		newFileName += parts[i]
	}

	return newFileName + ".fixed." + parts[len(parts)-1]
}
