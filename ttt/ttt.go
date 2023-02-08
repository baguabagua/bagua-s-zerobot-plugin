package ttt

import (
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

type tttstate struct {
	open bool
	board [9]int
	OorX bool
}

var state tttstate

func ttt_init() {
	state.open = false
}

func (s tttstate) print_board() string {
	m := ""
	for i := 0; i < 3; i++ {
		m += "\n"
		for j := 0; j < 3; j++ {
			switch s.board[i*3+j] {
			case 0:
				m += "_"
			case 1:
				m += "O"
			case -1:
				m += "X"
			}
		}
	}
	return m
}

func (s tttstate) line(a int, b int, c int) (bool, int) {
	if s.board[a] == 0 {
		return false, 0
	}
	return s.board[a] == s.board[b] && s.board[a] == s.board[c], s.board[a]
}

func (s tttstate) endgame() (bool, int) {
	lines := [][]int {{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, {0, 3, 6}, {1, 4, 7}, {2, 5, 8}, {0, 4, 8}, {2, 4, 6}}
	for _, l := range lines {
		end, winner := s.line(l[0], l[1], l[2])
		if end {
			return end, winner
		}
	}
	for i := 0; i < 9; i++ {
		if s.board[i] == 0 {
			return false, 0
		}
	}
	return true, 0
}

func tic_tac_toc(ctx *zero.Ctx) {
	if state.open {
		ctx.Send(message.Text("游戏正在进行中" + state.print_board()))
		return
	}
	state.open = true
	state.OorX = true
	for i := 0; i < 9; i++ {
		state.board[i] = 0
	}
	ctx.Send(message.Text("游戏开始" + state.print_board()))
}

func ttt_one_step(ctx *zero.Ctx) {
	if !state.open {
		ctx.Send(message.Text("游戏还未开始"))
		return
	}
	m := ctx.ExtractPlainText()
	n := m[2] - '1'
	if n < 0 || n >= 9 {
		ctx.Send(message.Text("没有这样的格子" + state.print_board()))
	} else {
		if state.board[n] != 0 {
			ctx.Send(message.Text("这个格子已经有棋子了" + state.print_board()))
		} else {
			if state.OorX {
				state.board[n] = 1
			} else {
				state.board[n] = -1
			}
			state.OorX = !state.OorX
			end, winner := state.endgame()
			if end {
				info := "游戏结束，"
				switch winner {
				case 0:
					info += "平局"
				case 1:
					info += "O获胜"
				case -1:
					info += "X获胜"
				}
				ctx.Send(message.Text(info + state.print_board()))
				state.open = false
			} else {
				ctx.Send(message.Text(state.print_board()))
			}
		}
	}
}

func ttt_exit(ctx *zero.Ctx) {
	if !state.open {
		ctx.Send(message.Text("游戏还未开始"))
		return
	}
	state.open = false
	ctx.Send(message.Text("已结束游戏"))
}

func init() {
	ttt_init()
	zero.OnCommand("下井字棋").Handle(tic_tac_toc)
	zero.OnCommand("a").Handle(ttt_one_step)
	zero.OnCommand("结束游戏").Handle(ttt_exit)
}