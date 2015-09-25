package abnf

import (
	"fmt"
	"strconv"
)

type MatchException struct { //extends Exception {
	actual   int
	pos      int
	line     int
	expected string
}

func NewMatchException(expected string, actual, pos, line int) *MatchException {
	this := &MatchException{}
	this.expected = expected
	this.actual = actual
	this.pos = pos
	this.line = line
	return this
}

func (this *MatchException) String() string {
	return "Mismatch with '" + string(byte(this.actual)) +
		"' [" + fmt.Sprintf("%02X", this.actual) + "] at position " +
		strconv.Itoa(this.pos) + ": line " + strconv.Itoa(this.line) +
		". Expected value is " + this.expected
}

type CollisionException struct { //extends Exception {
	collision string
	pos       int
	line      int
}

func NewCollisionException(collision string, pos, line int) *CollisionException {
	this := &CollisionException{}
	this.collision = collision
	this.pos = pos
	this.line = line
	return this
}

func (this *CollisionException) String() string {
	return "Collision at position " + strconv.Itoa(this.pos) + ": line " + strconv.Itoa(this.line) + ". Description: " + this.collision
}
