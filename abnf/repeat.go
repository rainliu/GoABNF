package abnf

import (
	"bytes"
	"strconv"
)

// repeat         =  1*DIGIT / (*DIGIT "*" *DIGIT)
type Repeat struct {
	min     int
	max     int
	starred bool
}

func NewRepeat(min, max int, starred bool) *Repeat {
	this := &Repeat{}
	this.min = min
	this.max = max
	this.starred = starred
	return this
}

func (this *Repeat) GetMin() int {
	return this.min
}

func (this *Repeat) GetMax() int {
	return this.max
}

func (this *Repeat) String() string {
	var s bytes.Buffer
	if this.starred {
		if this.min != 0 {
			s.WriteString(strconv.Itoa(this.min))
		}
		s.WriteString("*")
		if this.max != -1 {
			s.WriteString(strconv.Itoa(this.max))
		}
	} else {
		s.WriteString(strconv.Itoa(this.min))
	}
	return s.String()
}
