package engine

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type EngineState uint8

const (
	Initialized EngineState = iota
	WaitingReadyOk
	WaitingUSIOk
	Idling
	Thinking
)

type USIInfo struct {
	ScoreCp   int
	ScoreMate int
	MultiPv   int
	Depth     int
	SelDepth  int
	Nodes     int
	Nps       int
	Time      int
	HashFull  int
	CurrMove  string
	Pv        []string

	Upperbound bool
	Lowerbound bool

	IsCp   bool
	IsMate bool
}

func NewUSIInfo(line string) (USIInfo, error) {
	info := USIInfo{
		MultiPv: 1,
	}

	scanner := bufio.NewScanner(strings.NewReader(line))
	scanner.Split(bufio.ScanWords)

	atoi := func(a string) (int, error) {
		sign := 1
		if a[0] == '+' {
			a = a[1:]
		}
		if a[0] == '-' {
			sign = -1
			a = a[1:]
		}
		i, err := strconv.Atoi(a)
		return sign * i, err
	}

	isMove := func(m string) bool {
		if !(len(m) == 4 || len(m) == 5) {
			return false
		}
		if '0' <= m[0] && m[0] <= '9' {
			return true
		}
		switch m[0] {
		case 'P', 'L', 'N', 'S', 'G', 'B', 'R':
			return true
		}
		return false
	}

	var err error
	for scanner.Scan() {
		switch scanner.Text() {
		case "info":
		case "multipv":
			scanner.Scan()
			info.MultiPv, err = atoi(scanner.Text())
			if err != nil {
				break
			}
		case "cp":
			scanner.Scan()
			info.ScoreCp, err = atoi(scanner.Text())
			if err != nil {
				break
			}
			info.IsCp = true
		case "mate":
			scanner.Scan()
			info.ScoreMate, err = atoi(scanner.Text())
			if err != nil {
				break
			}
			info.IsMate = true
		case "depth":
			scanner.Scan()
			info.Depth, err = atoi(scanner.Text())
			if err != nil {
				break
			}
		case "seldepth":
			scanner.Scan()
			info.SelDepth, err = atoi(scanner.Text())
			if err != nil {
				break
			}
		case "nodes":
			scanner.Scan()
			info.Nodes, err = atoi(scanner.Text())
			if err != nil {
				break
			}
		case "nps":
			scanner.Scan()
			info.Nps, err = atoi(scanner.Text())
			if err != nil {
				break
			}
		case "time":
			scanner.Scan()
			info.Time, err = atoi(scanner.Text())
			if err != nil {
				break
			}
		case "hashfull":
			scanner.Scan()
			info.HashFull, err = atoi(scanner.Text())
			if err != nil {
				break
			}
		case "upperbound":
			info.Upperbound = true
		case "lowerbound":
			info.Lowerbound = true
		case "currmove":
			scanner.Scan()
			info.CurrMove = scanner.Text()
		case "pv":
			// 最後まで読む
			for scanner.Scan() {
				move := scanner.Text()
				if !isMove(move) {
					break
				}
				info.Pv = append(info.Pv, move)
			}
			break
		default:
		}
	}

	if err != nil {
		return USIInfo{}, err
	}
	return info, nil
}

type USIBestMove struct {
	BestMove string
	Ponder   string
}

func NewUSIBestMove(line string) (USIBestMove, error) {
	bestMove := USIBestMove{}

	scanner := bufio.NewScanner(strings.NewReader(line))
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		switch scanner.Text() {
		case "bestmove":
			scanner.Scan()
			bestMove.BestMove = scanner.Text()
		case "ponder":
			scanner.Scan()
			bestMove.Ponder = scanner.Text()
		default:
			return USIBestMove{}, fmt.Errorf("invalid bestmove: %s", line)
		}
	}

	return bestMove, nil
}

type Engine struct {
	cmd    *exec.Cmd
	stdout *bufio.Reader
	stdin  *bufio.Writer

	Name   string
	Author string
	State  EngineState

	InfoC     chan<- USIInfo
	BestMoveC chan<- USIBestMove

	ReadyOk bool
	USIOk   bool
}

func NewEngine(path string) (*Engine, error) {
	cmd := exec.Command(path)
	cmd.Dir = filepath.Dir(path)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	engine := &Engine{
		cmd:    cmd,
		stdin:  bufio.NewWriter(stdin),
		stdout: bufio.NewReader(stdout),
		State:  Initialized,
	}

	go engine.readLines()

	return engine, nil
}

func (e *Engine) Close() error {
	if err := e.cmd.Process.Kill(); err != nil {
		return err
	}
	return nil
}

func (e *Engine) readLines() {
	for {
		line, err := e.stdout.ReadString('\n')
		if err == io.EOF {
			return
		}
		if err != nil {
			panic(err)
		}
		line = strings.Trim(line, "\n")

		switch {
		case line == "usiok":
			e.State = Idling
			e.USIOk = true
		case line == "readyok":
			e.State = Idling
			e.ReadyOk = true
		case strings.HasPrefix(line, "id name"):
			e.Name = line[8:]
		case strings.HasPrefix(line, "id author"):
			e.Author = line[10:]
		case strings.HasPrefix(line, "info string"):
			// 無視
		case strings.HasPrefix(line, "info"):
			if e.InfoC != nil {
				info, err := NewUSIInfo(line)
				if err != nil {
					panic(err)
				}
				e.InfoC <- info
			}
		case strings.HasPrefix(line, "bestmove"):
			e.State = Idling
			if e.BestMoveC != nil {
				bestmove, err := NewUSIBestMove(line)
				if err != nil {
					panic(err)
				}
				e.BestMoveC <- bestmove
			}
		}
	}
}

func (e *Engine) Send(command string) error {
	_, err := e.stdin.WriteString(command + "\n")
	if err != nil {
		return err
	}
	e.stdin.Flush()
	return nil
}

func (e *Engine) SendUSI() error {
	e.State = WaitingUSIOk
	return e.Send("usi")
}

func (e *Engine) SendIsReady() error {
	e.State = WaitingReadyOk
	return e.Send("isready")
}

func (e *Engine) GoInfinite() error {
	e.State = Thinking
	return e.Send("go infinite")
}

func (e *Engine) SendSFEN(sfen string, moves []string) error {
	command := fmt.Sprintf("position sfen %s", sfen)
	if len(moves) != 0 {
		command += " moves"
		for _, move := range moves {
			command += " " + move
		}
	}
	return e.Send(command)
}

func (e *Engine) SendSetOption(name string, value string) error {
	return e.Send(fmt.Sprintf("setoption name %s value %s", name, value))
}

func (e *Engine) SendStop() error {
	return e.Send("stop")
}
