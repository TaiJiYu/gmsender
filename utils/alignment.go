package utils

// SingleAlignmentType为单边对齐方式
type SingleAlignmentType uint8

const (
	L SingleAlignmentType = iota // 左对齐
	M                            // 中对齐
	R                            // 右对齐
)

// 获取对齐控制后纠正的位置
func (a SingleAlignmentType) GetAlignmentP(p float64, halfSize float64) float64 {
	switch a {
	case L:
		return p
	case M:
		return p - halfSize
	case R:
		return p - halfSize*2
	}
	return p
}

// 转换为文字
func (a SingleAlignmentType) ToStr() string {
	switch a {
	case L:
		return "L"
	case M:
		return "M"
	case R:
		return "R"
	}
	return "?"
}

// AlignmentType为对齐方式,十位为X轴对齐方式，个位为Y轴对齐方式
type AlignmentType uint8

const (
	LL    AlignmentType = iota // X左对齐，Y顶对齐，以左上角为基准点
	LM                         // X左对齐，Y居中，以左中为基准点
	LR                         // X左对齐，Y底对齐，以左下为基准点
	ML                         // X居中，Y顶对齐，以中上为基准点
	MM                         // X居中，Y居中，以中心为基准点
	MR                         // X居中，Y底对齐，以中下为基准点
	RL                         // X右对齐，Y顶对齐，以右上为基准点
	RM                         // X右对齐，Y居中，以右中为基准点
	RR                         // X右对齐，Y底对齐，以右下为基准点
	maxAT = RR                 // 最大对齐
)

// 获取对齐控制后纠正的位置
func (a AlignmentType) GetAlignmentPos(pos Point, halfSize Point) Point {
	switch a {
	case LL:
		return pos
	case LM:
		return pos.SubY(halfSize.Y)
	case LR:
		return pos.SubY(halfSize.Y * 2)
	case ML:
		return pos.SubX(halfSize.X)
	case MM:
		return pos.Sub(halfSize)
	case MR:
		return pos.SubX(halfSize.X).SubY(halfSize.Y * 2)
	case RL:
		return pos.SubX(halfSize.X * 2)
	case RM:
		return pos.SubX(halfSize.X * 2).SubY(halfSize.Y)
	case RR:
		return pos.Sub(halfSize.MulF1(2))
	}
	return pos
}

// 转换为文字
func (a AlignmentType) ToStr() string {
	switch a {
	case LL:
		return "LL"
	case LM:
		return "LM"
	case LR:
		return "LR"
	case ML:
		return "ML"
	case MM:
		return "MM"
	case MR:
		return "MR"
	case RL:
		return "RL"
	case RM:
		return "RM"
	case RR:
		return "RR"
	}
	return "??"
}

func (a AlignmentType) Next() AlignmentType {
	return (a + 1) % (maxAT + 1)

}
