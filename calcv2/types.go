package calcv2

import (
	. "calcapp/utils"
)

const (
	COLS       = 56 + 1
	ROWS       = 9
	GROUP_SIZE = 27
	STOP_COL   = 0
	STOP_VALUE = 2047

	CHUNK_SIZE = 3

)

type Env struct {
	CurrentCol Value
	Last       Point // 前一列同行
	S1	   	   Point // 前一列上一行，用于计算符号
	S2 		   Point // 同列上一行，用于计算符号
	Stop       bool
}

type BaseData struct {
	//Inst [COLS]Bpoint
	Bp   [COLS]Bpoint
	Zbp  [COLS]Point
	Data [ROWS][COLS]Point
	Ag   [COLS]Point  // A果
	G1   [COLS]Point  // 果1
}

type BaseGroup struct {
	Data  	[GROUP_SIZE]BaseData
	G1		[CHUNK_SIZE+1][COLS]Point	// intermediate value, plus G1
}