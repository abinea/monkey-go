package lexer

import (
	"monkey/token"
)

/*
Todo:
	该词法分析器仅支持ASCII字符，不能支持所有的Unicode字符;
	这么做也是为了让事情保持简单，让我们能够专注于解释器的基础部分
	如果要完全支持Unicode和UTF-8，就要将l.ch的类型从byte改为rune；
	同时还要修改读取下一个字符的方式，因为字符此时可能为多字节。
*/

type Lexer struct {
	input        string
	position     int  // 所输入字符串的位置（指向当前字符）
	readPosition int  // 所输入字符串的当前读取位置（指向当前字符之后的一个字符）
	ch           byte //当前正在查看的字符
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// 读取input中的下一个字符，并前移其在input中的位置
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // 表示“尚未读取任何内容”或“文件结尾”
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.Equal, Literal: literal}
		} else {
			tok = newToken(token.Assign, l.ch)
		}
	case '+':
		tok = newToken(token.Plus, l.ch)
	case '-':
		tok = newToken(token.Minus, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NotEqual, Literal: literal}
		} else {
			tok = newToken(token.Bang, l.ch)
		}
	case '*':
		tok = newToken(token.Asterisk, l.ch)
	case '/':
		tok = newToken(token.Slash, l.ch)
	case '<':
		tok = newToken(token.LessThan, l.ch)
	case '>':
		tok = newToken(token.GreaterThan, l.ch)
	case ',':
		tok = newToken(token.Comma, l.ch)
	case ';':
		tok = newToken(token.Semicolon, l.ch)
	case ':':
		tok = newToken(token.Colon, l.ch)
	case '(':
		tok = newToken(token.LeftParen, l.ch)
	case ')':
		tok = newToken(token.RightParen, l.ch)
	case '{':
		tok = newToken(token.LeftBrace, l.ch)
	case '}':
		tok = newToken(token.RightBrace, l.ch)
	case '[':
		tok = newToken(token.LeftBracket, l.ch)
	case ']':
		tok = newToken(token.RightBracket, l.ch)
	case '"':
		tok.Type = token.String
		tok.Literal = l.readString()
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdentifierType(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.Int
			return tok
		} else {
			tok = newToken(token.Illegal, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readString() string {
	pos := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[pos:l.position]
}

func (l *Lexer) skipWhitespace() {
	// 读取到空白字符时，就调用readChar函数读取下一个字符，直到遇到非空白字符为止
	// Todo:有些语言会为换行符创建词法单元，如果它们不在词法单元流中的正确位置，就会抛出解析错误，这里简化处理
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	// 读取到数字字符时，就调用readChar函数读取下一个字符，直到遇到非数字字符为止
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	// Todo: 支持浮点数、十六进制、八进制
	return '0' <= ch && ch <= '9'
}
