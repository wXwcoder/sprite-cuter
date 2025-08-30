package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Point 表示一个二维坐标点
type Point struct {
	X, Y int
}

// Rect 表示一个矩形区域
type Rect struct {
	LT, LB, RT, RB Point // 左上、左下、右上、右下
}

func main() {
	// 解析命令行参数
	pngFile := flag.String("input", "", "PNG文件路径")
	flag.Parse()

	if *pngFile == "" {
		log.Fatal("请提供PNG文件路径作为参数，使用 -input 参数")
	}

	if !fileExists(*pngFile) {
		log.Fatalf("文件不存在: %s", *pngFile)
	}

	// 读取PNG文件
	img, err := readPNG(*pngFile)
	if err != nil {
		log.Fatalf("读取图片时发生错误: %v", err)
	}

	// 获取输出目录名
	outDir := strings.TrimSuffix(filepath.Base(*pngFile), ".png")

	// 创建输出目录
	if err := createDir("export"); err != nil {
		log.Fatal(err)
	}
	if err := createDir("export/" + outDir); err != nil {
		log.Fatal(err)
	}

	// 检测并提取精灵
	spritesArray := getSprites(img)

	// 生成CSS文件
	css := getCSS(spritesArray, *pngFile)
	if err := os.WriteFile("export/"+outDir+"/"+outDir+".css", []byte(css), 0644); err != nil {
		log.Fatal(err)
	}
	fmt.Println("CSS文件已保存!")

	// 生成JSON文件
	json := getJson(spritesArray, *pngFile)
	if err := os.WriteFile("export/"+outDir+"/"+outDir+".json", []byte(json), 0644); err != nil {
		log.Fatal(err)
	}
	fmt.Println("JSON文件已保存!")

	// 切割并保存精灵图
	for i, rect := range spritesArray {
		fmt.Printf("精灵 %d: %+v\n", i, rect)
		if err := saveSprite(img, rect, outDir, i); err != nil {
			log.Printf("保存精灵 %d 时出错: %v", i, err)
		}
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func readPNG(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return png.Decode(file)
}

func createDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.Mkdir(path, 0755)
	}
	return nil
}

// getSprites 检测图像中的所有精灵
func getSprites(img image.Image) []Rect {
	bounds := img.Bounds()
	imgWidth, imgHeight := bounds.Dx(), bounds.Dy()

	// 创建像素数据副本
	data := make([]uint8, imgWidth*imgHeight*4)
	for y := 0; y < imgHeight; y++ {
		for x := 0; x < imgWidth; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			idx := (y*imgWidth + x) * 4
			data[idx] = uint8(r >> 8)
			data[idx+1] = uint8(g >> 8)
			data[idx+2] = uint8(b >> 8)
			data[idx+3] = uint8(a >> 8)
		}
	}

	var spritesArray []Rect
	contourVector := marchingSquares(data, imgHeight, imgWidth)

	for len(contourVector) > 3 {
		rect := getRect(contourVector)
		width := max(1, rect.RT.X-rect.LT.X)
		height := max(1, rect.RB.Y-rect.RT.Y)

		if width > 3 && height > 3 {
			spritesArray = append(spritesArray, rect)
		}

		// 清除已处理的区域
		for y := rect.RT.Y; y < rect.RB.Y; y++ {
			for x := rect.LB.X; x < rect.RB.X; x++ {
				if x >= 0 && x < imgWidth && y >= 0 && y < imgHeight {
					idx := (y*imgWidth + x) * 4
					data[idx] = 0
					data[idx+1] = 0
					data[idx+2] = 0
					data[idx+3] = 0
				}
			}
		}

		contourVector = marchingSquares(data, imgHeight, imgWidth)
	}

	return spritesArray
}

// getCSS 生成CSS样式
func getCSS(spritesArray []Rect, pngName string) string {
	css := ".sprite {display:inline-block; overflow:hidden; background-repeat: no-repeat;background-image:url(" + pngName + ");}"
	for i, rect := range spritesArray {
		css += getSpriteCSS(fmt.Sprintf("sprite%d", i), rect)
	}
	return css
}

// getJson 生成JSON样式
func getJson(spritesArray []Rect, pngName string) string {
	imgWidth := 0
	imgHeight := 0
	json := fmt.Sprintf(`{"sprite":{
		"width":%d,
		"height":%d,
		"image":"%s",
		"frames":[`, imgWidth, imgHeight, pngName)
	for i, rect := range spritesArray {
		json += getSpriteJson(fmt.Sprintf("sprite%d", i), rect)
	}
	json = json[:len(json)-1] + "]}}"
	return json
}

// getSpriteCSS 生成单个精灵的CSS样式
func getSpriteCSS(spriteName string, rect Rect) string {
	width := rect.RT.X - rect.LT.X
	height := rect.RB.Y - rect.RT.Y
	return fmt.Sprintf(".%s {width:%dpx; height:%dpx; background-position: %dpx %dpx}",
		spriteName, width, height, -rect.LT.X, -rect.LT.Y)
}

// getSpriteJson 生成单个精灵的JSON样式
func getSpriteJson(spriteName string, rect Rect) string {
	width := rect.RT.X - rect.LT.X
	height := rect.RB.Y - rect.RT.Y
	return fmt.Sprintf(`{"name":"%s","x":%d,"y":%d,"width":%d,"height":%d},`,
		spriteName, -rect.LT.X, -rect.LT.Y, width, height)
}

// marchingSquares 实现marching squares算法检测轮廓
func marchingSquares(data []uint8, height, width int) []Point {
	var contourVector []Point

	// 获取起始像素
	startPoint := getStartingPixel(data, height, width)
	if startPoint == nil || startPoint.X < 0 || startPoint.Y < 0 {
		return contourVector
	}

	pX, pY := startPoint.X, startPoint.Y
	var stepX, stepY int
	var prevX, prevY int
	closedLoop := false
	iteration := 0

	for !closedLoop && iteration < 200000 {
		squareValue := getSquareValue(data, pX, pY, width, height)

		switch squareValue {
		case 1, 5, 13:
			stepX, stepY = 0, -1
		case 8, 10, 11:
			stepX, stepY = 0, 1
		case 4, 12, 14:
			stepX, stepY = -1, 0
		case 2, 3, 7:
			stepX, stepY = 1, 0
		case 6:
			if prevX == 0 && prevY == -1 {
				stepX, stepY = -1, 0
			} else {
				stepX, stepY = 1, 0
			}
		case 9:
			if prevX == 1 && prevY == 0 {
				stepX, stepY = 0, -1
			} else {
				stepX, stepY = 0, 1
			}
		default:
			stepX, stepY = 0, 0
		}

		pX += stepX
		pY += stepY
		contourVector = append(contourVector, Point{X: pX, Y: pY})
		prevX, prevY = stepX, stepY
		iteration++

		// 如果回到起点，循环结束
		if pX == startPoint.X && pY == startPoint.Y {
			closedLoop = true
		}
	}

	return contourVector
}

// getSquareValue 获取2x2网格的方值
func getSquareValue(data []uint8, pX, pY, width, height int) int {
	squareValue := 0

	if pX-1 >= 0 && pY-1 >= 0 && !isAlpha(data, pX-1, pY-1, width, height) {
		squareValue += 1
	}
	if pY-1 >= 0 && !isAlpha(data, pX, pY-1, width, height) {
		squareValue += 2
	}
	if pX-1 >= 0 && !isAlpha(data, pX-1, pY, width, height) {
		squareValue += 4
	}
	if !isAlpha(data, pX, pY, width, height) {
		squareValue += 8
	}

	return squareValue
}

// getStartingPixel 扫描找到第一个非透明像素作为起始点
func getStartingPixel(data []uint8, height, width int) *Point {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			idx := (y*width + x) * 4
			if idx+3 < len(data) && data[idx+3] > 0 {
				fmt.Printf("起始点: %d %d\n", x, y)
				return &Point{X: x, Y: y}
			}
		}
	}
	return nil
}

// isAlpha 检查像素是否透明
func isAlpha(data []uint8, x, y, width, height int) bool {
	if x < 0 || y < 0 || x >= width || y >= height {
		return true
	}
	idx := (y*width + x) * 4
	if idx < 0 || idx >= len(data) {
		return true
	}
	return data[idx+3] == 0
}

// getRect 从轮廓点中提取矩形边界
func getRect(points []Point) Rect {
	if len(points) == 0 {
		return Rect{}
	}

	maxX, minX := points[0].X, points[0].X
	maxY, minY := points[0].Y, points[0].Y

	for _, p := range points {
		if p.X > maxX {
			maxX = p.X
		}
		if p.X < minX {
			minX = p.X
		}
		if p.Y > maxY {
			maxY = p.Y
		}
		if p.Y < minY {
			minY = p.Y
		}
	}

	return Rect{
		LT: Point{X: minX, Y: minY},
		LB: Point{X: minX, Y: maxY},
		RT: Point{X: maxX, Y: minY},
		RB: Point{X: maxX, Y: maxY},
	}
}

// saveSprite 保存切割后的精灵图
func saveSprite(img image.Image, rect Rect, outDir string, index int) error {
	width := max(1, rect.RT.X-rect.LT.X)
	height := max(1, rect.RB.Y-rect.RT.Y)

	// 创建新的图像
	newImg := image.NewRGBA(image.Rect(0, 0, width, height))

	// 复制像素数据
	srcX := max(0, rect.LT.X)
	srcY := max(0, rect.LT.Y)
	copyWidth := min(width, img.Bounds().Dx()-srcX)
	copyHeight := min(height, img.Bounds().Dy()-srcY)

	for y := 0; y < copyHeight; y++ {
		for x := 0; x < copyWidth; x++ {
			newImg.Set(x, y, img.At(srcX+x, srcY+y))
		}
	}

	// 保存文件
	filename := fmt.Sprintf("export/%s/%s%d.png", outDir, outDir, index)
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, newImg)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
