package view

import (
	"bufio"
	"os"
	"strings"
	"strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func LoadModel(fileName string) *ModelData {
	f, err := os.Open(fileName)
	check(err)
	model := &ModelData{}
	
	reader := bufio.NewReader(f)
	reader.ReadLine()
	reader.ReadLine()
	line, _, _ := reader.ReadLine()
	strs := strings.Split(string(line), " ")
	nVec, _ := strconv.Atoi(strs[2])

	for ;string(line) != "end_header"; {
		line, _, _ = reader.ReadLine()
	}
	for i := 0; i < nVec; i++ {
		line, _, _ = reader.ReadLine()
		strs = strings.Split(string(line), " ")
		x, _ := strconv.ParseFloat(strs[0], 32)
		y, _ := strconv.ParseFloat(strs[1], 32)
		z, _ := strconv.ParseFloat(strs[2], 32)
		nx, _ := strconv.ParseFloat(strs[3], 32)
		ny, _ := strconv.ParseFloat(strs[4], 32)
		nz, _ := strconv.ParseFloat(strs[5], 32)
		r, _ := strconv.ParseFloat(strs[6], 32)
		g, _ := strconv.ParseFloat(strs[7], 32)
		b, _ := strconv.ParseFloat(strs[8], 32)
		model.vertices = append(model.vertices, float32(x))
		model.vertices = append(model.vertices, float32(y))
		model.vertices = append(model.vertices, float32(z))
		model.normals = append(model.normals, float32(nx))
		model.normals = append(model.normals, float32(ny))
		model.normals = append(model.normals, float32(nz))
		model.colors = append(model.colors, float32(r / 255))
		model.colors = append(model.colors, float32(g / 255))
		model.colors = append(model.colors, float32(b / 255))
	}
	return model
}