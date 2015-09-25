package abnf

import (
	"io"
	//"fmt"
)

const (
	PEEKER_EOF = -1
)

type Peeker struct {
	pos  int
	line int
	/**
	 * The underlying stream.
	 */
	stream io.Reader

	/**
	 * Bytes that have been peeked at.
	 */
	peekBytes []byte

	/**
	 * How many bytes have been peeked at.
	 */
	peekLength int
}

/**
 * The constructor accepts an InputStream to setup the
 * object.
 *
 * @param is
 *          The InputStream to parse.
 */
func NewPeeker(reader io.Reader) *Peeker {
	this := &Peeker{}
	this.stream = reader
	this.peekBytes = make([]byte, 10)
	this.peekLength = 0
	this.pos = 1
	this.line = 1
	return this
}

func (this *Peeker) GetPos() int {
	return this.pos
}

func (this *Peeker) GetLine() int {
	return this.line
}

//when read LF, set pos = 1
//when read CR, line++
func (this *Peeker) UpdatePosition(value byte) {
	if value == 0x0D {
		this.pos = 1
	} else if value == 0x0A {
		this.line++
	} else {
		this.pos++
	}
}

/**
 * Peek at the next character from the stream.
 *
 * @return The next character.
 * @throws IOException
 *           If an I/O exception occurs.
 */
/*func (this *Peeker) Peek() (IOException error){
    return peek(0);
}*/

/**
 * Peek at a specified depth.
 *
 * @param depth
 *          The depth to check.
 * @return The character peeked at.
 * @throws IOException
 *           If an I/O exception occurs.
 */
func (this *Peeker) Peek(depth int) int {
	// does the size of the peek buffer need to be extended?
	if len(this.peekBytes) <= depth {
		temp := make([]byte, depth+10)
		for i := 0; i < len(this.peekBytes); i++ {
			temp[i] = this.peekBytes[i]
		}
		this.peekBytes = temp
	}

	// does more data need to be read?
	if depth >= this.peekLength {
		offset := this.peekLength
		length := (depth - this.peekLength) + 1
		_, IOException := this.stream.Read(this.peekBytes[offset:offset+length])
		
		if IOException != nil{
			//fmt.Printf("Peek(%d): peekLength=%d, length=%d, readLength=%d, err=%s\n", depth, this.peekLength, length, readLength, IOException);
			return PEEKER_EOF;//panic(IOException.Error())
		}

		this.peekLength = depth + 1
	}

	return int(this.peekBytes[depth])
}

/*
 * Read a single byte from the stream. @throws IOException
 * If an I/O exception occurs. @return The character that
 * was read from the stream.
 */

func (this *Peeker) Read() int {
	if this.peekLength == 0 {
		var value [1]byte
		_, IOException := this.stream.Read(value[:])
		if IOException != nil{
			//fmt.Printf("Read(): peekLength=%d, length=%d, readLength=%d, err=%s\n", this.peekLength, 1, readLength, IOException);
			return PEEKER_EOF;//panic(IOException.Error())
		}

		this.UpdatePosition(value[0])
		return int(value[0])
	}

	result := this.peekBytes[0]
	this.UpdatePosition(result)
	this.peekLength--
	for i := 0; i < this.peekLength; i++ {
		this.peekBytes[i] = this.peekBytes[i+1]
	}

	return int(result)
}
