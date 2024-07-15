package scanner

import (
	"bufio"
	"io"
	"unicode/utf8"

	"github.com/arsenzairov/cs153/token"
)

var eof = rune(0)

// represents lexical scanner
type Scanner struct {
	r       *bufio.Reader
	pos     token.Position
	prevPos token.Position // Track previous position
	source  []rune
}

func NewScanner(r io.Reader, filename string, estimatedSize int) *Scanner {
	position := token.Position{Filename: filename}
	return &Scanner{r: bufio.NewReader(r), pos: position, source: make([]rune, 0, estimatedSize)}
}

func (s *Scanner) Source() string {
	return string(s.source)
}

func (s *Scanner) movePos(ch rune) {
	s.prevPos = s.pos // Save the current position before advancing
	size := utf8.RuneLen(ch)
	s.pos.Offset += size
	if ch == '\n' {
		s.pos.Line += 1
		s.pos.Column = 0
	} else {
		s.pos.Column += size
	}
}

// read the next rune from the buffered characters
func (s *Scanner) advance() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	s.source = append(s.source, ch)
	s.movePos(ch)
	return ch
}

// place previous rune back into reader
func (s *Scanner) unread() {
	if len(s.source) > 0 {
		s.source = s.source[:len(s.source)-1]
	}
	_ = s.r.UnreadRune()
	s.pos = s.prevPos
}

func (s *Scanner) Scan() token.Token {
	s.skipWhiteSpace()

	// Read the next rune
	start := s.pos.Offset

	switch ch := s.advance(); {
	case isLetter(ch):
		return s.scanIdent(start)
	case isDigit(ch):
		return s.scanNumber(start)

	default:
		switch ch {
		case '(':
			return s.makeToken(token.LPAREN, token.NewLoc(start, s.pos.Offset))
		case ')':
			return s.makeToken(token.RPAREN, token.NewLoc(start, s.pos.Offset))
		case '{':
			return s.makeToken(token.LBRACE, token.NewLoc(start, s.pos.Offset))
		case '}':
			return s.makeToken(token.RBRACE, token.NewLoc(start, s.pos.Offset))
		case '[':
			return s.makeToken(token.RBRACK, token.NewLoc(start, s.pos.Offset))
		case ']':
			return s.makeToken(token.LBRACK, token.NewLoc(start, s.pos.Offset))
		case ';':
			return s.makeToken(token.SEMICOLON, token.NewLoc(start, s.pos.Offset))
		case '.':
			return s.makeToken(token.PERIOD, token.NewLoc(start, s.pos.Offset))
		case ',':
			return s.makeToken(token.COMMA, token.NewLoc(start, s.pos.Offset))
		case '-':
			if s.match('-') {
				return s.makeToken(token.DEC, token.NewLoc(start, s.pos.Offset))
			}
			if s.match('=') {
				return s.makeToken(token.SUB_ASSIGN, token.NewLoc(start, s.pos.Offset))
			}
			return s.makeToken(token.SUB, token.NewLoc(start, s.pos.Offset))
		case '+':
			if s.match('+') {
				return s.makeToken(token.INC, token.NewLoc(start, s.pos.Offset))
			}
			if s.match('=') {
				return s.makeToken(token.ADD_ASSIGN, token.NewLoc(start, s.pos.Offset))
			}
			return s.makeToken(token.ADD, token.NewLoc(start, s.pos.Offset))
		case '/':
			if s.match('=') {
				return s.makeToken(token.QUO_ASSIGN, token.NewLoc(start, s.pos.Offset))
			}
			if s.match('/') {
				for ch := s.advance(); ch != '\n' && ch != eof; ch = s.advance() {
				}

				if ch != '\n' || ch != eof {
					s.unread()
				}

				// Skip the comment line
				s.skipWhiteSpace()
				return s.Scan()
			}
			return s.makeToken(token.QUO, token.NewLoc(start, s.pos.Offset))
		case '*':
			if s.match('=') {
				return s.makeToken(token.MUL_ASSIGN, token.NewLoc(start, s.pos.Offset))
			}
			return s.makeToken(token.MUL, token.NewLoc(start, s.pos.Offset))
		case '%':
			if s.match('=') {
				return s.makeToken(token.REM_ASSIGN, token.NewLoc(start, s.pos.Offset))
			}
			return s.makeToken(token.REM, token.NewLoc(start, s.pos.Offset))
		case '^':
			if s.match('=') {
				return s.makeToken(token.XOR_ASSIGN, token.NewLoc(start, s.pos.Offset))
			}
			return s.makeToken(token.XOR, token.NewLoc(start, s.pos.Offset))
		case '|':
			if s.match('|') {
				return s.makeToken(token.LOR, token.NewLoc(start, s.pos.Offset))
			}
			if s.match('=') {
				return s.makeToken(token.OR_ASSIGN, token.NewLoc(start, s.pos.Offset))
			}
			return s.makeToken(token.OR, token.NewLoc(start, s.pos.Offset))
		case '&':
			if s.match('&') {
				return s.makeToken(token.LAND, token.NewLoc(start, s.pos.Offset))
			}
			if s.match('=') {
				return s.makeToken(token.AND_ASSIGN, token.NewLoc(start, s.pos.Offset))
			}
			return s.makeToken(token.AND, token.NewLoc(start, s.pos.Offset))
		case '!':
			if s.match('=') {
				return s.makeToken(token.NEQ, token.NewLoc(start, s.pos.Offset))
			}
			return s.makeToken(token.NOT, token.NewLoc(start, s.pos.Offset))
		case '=':
			if s.match('=') {
				return s.makeToken(token.EQL, token.NewLoc(start, s.pos.Offset))
			}
			return s.makeToken(token.DEFINE, token.NewLoc(start, s.pos.Offset))
		case '<':
			if s.match('=') {
				return s.makeToken(token.LEQ, token.NewLoc(start, s.pos.Offset))
			}
			return s.makeToken(token.LSS, token.NewLoc(start, s.pos.Offset))
		case '>':
			if s.match('=') {
				return s.makeToken(token.GEQ, token.NewLoc(start, s.pos.Offset))
			}
			return s.makeToken(token.GTR, token.NewLoc(start, s.pos.Offset))
		case '"':
			return s.scanString(start + 1)
		case '`':
			return s.scanRawString(start + 1)
		case eof:
			return s.makeToken(token.EOF, token.NewLoc(start, s.pos.Offset))
		}

		return s.errorToken(token.NewLoc(start, s.pos.Offset))
	}
}

func (s *Scanner) makeToken(t token.Tag, loc token.Loc) token.Token {
	return token.Token{Tag: t, Pos: s.pos, Loc: loc}
}

func (s *Scanner) errorToken(loc token.Loc) token.Token {
	return token.Token{Tag: token.ILLEGAL, Pos: s.pos, Loc: loc}
}

func (s *Scanner) match(char rune) bool {
	ch := s.advance()
	if ch == eof || ch != char {
		s.unread()
		return false
	}
	return true
}

func (s *Scanner) skipWhiteSpace() {
	for {
		ch := s.advance()
		if !isWhiteSpace(ch) {
			s.unread()
			break
		}
	}
}

func (s *Scanner) scanNumber(start int) token.Token {
	var ch rune

	for {
		ch = s.advance()
		if ch == eof || !isDigit(ch) {
			s.unread()
			break
		}
	}

	if ch != '.' {
		end := s.pos.Offset
		newLoc := token.NewLoc(start, end)
		return token.Token{Tag: token.INT, Pos: s.pos, Loc: newLoc}
	}

	// Handle the decimal point
	ch = s.advance()
	if !isDigit(ch) {
		s.unread() // Unread non-digit character after the decimal point
		s.unread() // Unread the decimal point
		end := s.pos.Offset
		newLoc := token.NewLoc(start, end)
		return token.Token{Tag: token.INT, Pos: s.pos, Loc: newLoc}
	}

	for {
		ch = s.advance()
		if ch == eof || !isDigit(ch) {
			s.unread()
			break
		}
	}

	end := s.pos.Offset
	newLoc := token.NewLoc(start, end)
	return token.Token{Tag: token.FLOAT, Pos: s.pos, Loc: newLoc}
}

func (s *Scanner) scanRawString(start int) token.Token {
	for {
		ch := s.advance()
		if ch == eof || ch == '`' {
			break
		}
	}

	end := s.pos.Offset
	return token.Token{Tag: token.STRING, Loc: token.NewLoc(start, end), Pos: s.pos}
}

func (s *Scanner) scanString(start int) token.Token {
	for {
		ch := s.advance()
		if ch == eof || ch == '"' {
			break
		}
		if ch == '\\' {
			next := s.advance()
			switch next {
			case 'n':
				s.source[len(s.source)-1] = '\n' // Replace '\\' with '\n'
			case 't':
				s.source[len(s.source)-1] = '\t' // Replace '\\' with '\t'
			case '"':
				s.source[len(s.source)-1] = '"' // Replace '\\' with '"'
			case '\\':
				s.source[len(s.source)-1] = '\\' // Keep '\\'
			default:
				s.source = append(s.source, next) // Add the next character as it is
			}
		}
	}

	end := s.pos.Offset
	return token.Token{Tag: token.STRING, Loc: token.NewLoc(start, end), Pos: s.pos}
}

func (s *Scanner) scanIdent(start int) token.Token {
	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		ch := s.advance()
		if ch == eof || (!isLetter(ch) && !isDigit(ch) && ch != '_') {
			s.unread()
			break
		}
	}

	end := s.pos.Offset
	newLoc := token.NewLoc(start, end)
	identifier := string(s.source[start:end])

	tag := token.Lookup(identifier)
	return token.Token{Tag: tag, Pos: s.pos, Loc: newLoc}
}
