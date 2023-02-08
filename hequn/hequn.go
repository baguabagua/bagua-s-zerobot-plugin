package hequn

import (
	"strconv"

	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/imgfactory"
	"github.com/fogleman/gg"
)

type hequnBoard struct {
	piece [][]int // 1-white 0-empty -1-black
	color [][]int
	width int
	height int
	end bool
	result int
	player int
	step int
}

func (board *hequnBoard) init(width int, height int) {
	board.width = width
	board.height = height
	board.end = false
	board.player = 1
	board.piece = make([][]int, width)
	board.color = make([][]int, width)
	for i := 0; i < width; i++ {
		board.piece[i] = make([]int, height)
		board.color[i] = make([]int, height)
	}
	board.step = 0
}

func (board *hequnBoard) outOfField(x int, y int) bool {
	return x <= 0 || x > board.width || y <= 0 || y > board.height
}

func (board *hequnBoard) paint(x int, y int) {
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if !board.outOfField(x + i, y + j) {
				board.color[x + i][y + j] = board.player
			}
		}
	}
}

func (board *hequnBoard) line(x1 int, y1 int, x2 int, y2 int) bool {
	return board.piece[x1][y1] == board.player && board.piece[x2][y2] == board.player
}

func (board *hequnBoard) play(x int, y int) bool {
	var directions = [][]int {{1, 0}, {1, 1}, {1, -1}, {0, -1}, {-1, -1}, {-1, 0}, {-1, 1}, {0, 1}}
	if board.outOfField(x, y) || board.piece[x][y] != 0 {
		return false
	}
	board.piece[x][y] = board.player
	for i := 0; i < 8; i++ {
		dx, dy := directions[i][0], directions[i][1]
		if !board.outOfField(x + dx + dx, y + dy + dy) && board.line(x + dx, y + dy, x + dx + dx, y + dy + dy) {
			board.paint(x + dx, y + dy)
		}
	}
	for i := 0; i < 4; i++ {
		dx, dy := directions[i][0], directions[i][1]
		if !board.outOfField(x + dx, y + dy) && !board.outOfField(x - dx, y - dy) && board.line(x + dx, y + dy, x - dx, y - dy) {
			board.paint(x, y)
		}
	}
	board.player = -board.player
	board.step++
	if (board.step == board.width * board.height) {
		board.end = true
		white, black := 0, 0
		for _, list := range board.color {
			for _, c := range list {
				if c == 1 {
					white++
				} else if c == -1 {
					black++
				}
			}
		}
		if white == black {
			board.result = 0
		} else if white > black {
			board.result = 1
		} else {
			board.result = -1
		}
	}
	return true
}

func (board *hequnBoard) print_board() []byte {
	canvas := gg.NewContext(1300, 1000)
	canvas.SetRGB(1.0, 1.0, 1.0)
	canvas.Clear()

	grid_width := 800.0 / float64(board.width)
	grid_height := 800.0 / float64(board.height)
	var min_wh float64
	if grid_width < grid_height {
		min_wh = grid_width
	} else {
		min_wh = grid_height
	}
	piece_radians := min_wh * 0.8
	for i := 0; i < board.width; i++ {
		for j := 0; j < board.height; j++ {
			canvas.DrawRectangle(100 + float64(i)*grid_width, 100 + float64(j)*grid_height, grid_width, grid_height)
			switch board.color[i][j] {
			case 1:
				canvas.SetRGB(0.5, 0.5, 1.0)
			case 0:
				canvas.SetRGB(0.5, 0.5, 0.5)
			case -1:
				canvas.SetRGB(1.0, 0.5, 0.5)
			}
			canvas.Fill()
			switch board.piece[i][j] {
			case 1:
				canvas.DrawCircle(100 + float64(i)*grid_width + grid_width/2, 100 + float64(j)*grid_height + grid_height/2, piece_radians)
				canvas.SetRGB(1.0, 1.0, 1.0)
				canvas.Fill()
			case -1:
				canvas.DrawCircle(100 + float64(i)*grid_width + grid_width/2, 100 + float64(j)*grid_height + grid_height/2, piece_radians)
				canvas.SetRGB(0.0, 0.0, 0.0)
				canvas.Fill()
			}
		}
	}
	canvas.DrawRectangle(1000.0, 100.0, 100.0, 100.0)
	canvas.SetRGB(0.5, 0.5, 1.0)
	canvas.Fill()
	canvas.DrawRectangle(1100.0, 100.0, 100.0, 100.0)
	canvas.SetRGB(1.0, 0.5, 0.5)
	canvas.Fill()
	canvas.DrawRectangle(1000.0, 200.0, 100.0, 100.0)
	canvas.SetRGB(1.0, 1.0, 1.0)
	canvas.Fill()
	canvas.DrawRectangle(1100.0, 200.0, 100.0, 100.0)
	canvas.SetRGB(1.0, 1.0, 1.0)
	canvas.Fill()
	canvas.DrawCircle(1050.0, 150.0, 40.0)
	canvas.SetRGB(1.0, 1.0, 1.0)
	canvas.Fill()
	canvas.DrawCircle(1150.0, 150.0, 40.0)
	canvas.SetRGB(0.0, 0.0, 0.0)
	canvas.Fill()

	canvas.DrawRectangle(1000, 800, 100, 100)
	if board.player == 1 {
		canvas.SetRGB(0.5, 0.5, 1.0)
	} else {
		canvas.SetRGB(1.0, 0.5, 0.5)
	}
	canvas.Fill()
	canvas.DrawCircle(1050, 850, 40)
	if board.player == 1 {
		canvas.SetRGB(1.0, 1.0, 1.0)
	} else {
		canvas.SetRGB(0.0, 0.0, 0.0)
	}
	canvas.Fill()
	canvas.DrawString("行棋", 1100, 800)
	base64Bytes, err := imgfactory.ToBase64(canvas.Image())
	if err != nil {
		log.Println("[hequnboard2img]", err)
		return nil
	}
	return base64Bytes
}

func init() {
	var board hequnBoard
	start := false
	zero.OnCommand("hequn").Handle(func (ctx *zero.Ctx) {
		if start {
			ctx.Send(message.Text("游戏已经开始"))
			return
		}
		start = true
		board.init(10, 10)
		ctx.SendChain(message.Image("base64://" + binary.BytesToString(board.print_board())))
	})
	zero.OnCommand("结束游戏").Handle(func (ctx *zero.Ctx) {
		if !start {
			ctx.Send(message.Text("游戏尚未开始"))
			return
		}
		start = false
		ctx.Send("游戏已结束")
	})
	zero.OnRegex(`^([abcdefghij])([0123456789]+)$`).Handle(func (ctx *zero.Ctx) {
		if !start {
			ctx.Send(message.Text("游戏尚未开始"))
			return
		}
		xstr := ctx.State["regex_matched"].([]string)[1]
		ystr := ctx.State["regex_matched"].([]string)[2]
		x := xstr[0] - 'a'
		y, _ := strconv.Atoi(ystr)
		if board.play(int(x), y) {
			ctx.SendChain(message.Image("base64://" + binary.BytesToString(board.print_board())))
			if board.end {
				start = false
				ctx.Send(message.Text("游戏已结束"))
			}
		} else {
			ctx.Send(message.Text("不合法的行动"))
		}
	})
}