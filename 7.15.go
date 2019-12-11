package main

import (
	"fmt"
	"text/scanner"
	"strings"
	"strconv"
	"math"
	"os"
	"bufio"
)

func main() {
	fmt.Printf("Expression: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	expr := scanner.Text()
	ast, err := Parse(expr)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Printf("Variables(eg: x=3 y=4): ")
	scanner.Scan()
	env, err := parseEnv(scanner.Text())
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	fmt.Printf("result=%v\n", ast.Eval(env))
}

func parseEnv(s string) (Env, error) {
	var scan scanner.Scanner
	scan.Init(strings.NewReader(s))
	scan.Mode = scanner.ScanIdents|scanner.ScanInts|scanner.ScanFloats
	env := make(Env)
	for {
		token := scan.Scan()
		if token == scanner.EOF {
			break
		}
		if token != scanner.Ident {
			return nil, fmt.Errorf("unexpected %q, want Variables", scan.TokenText())
		}
		id := Var(scan.TokenText())
		token = scan.Scan()
		if token != '=' {
			return nil, fmt.Errorf("unexpected %q, want '='", scan.TokenText())
		}
		token = scan.Scan()
		if token != scanner.Float && token != scanner.Int {
			return nil, fmt.Errorf("unexpected %q, want number", scan.TokenText())
		}
		value, _ := strconv.ParseFloat(scan.TokenText(), 64)
		env[id] = value
	}
	return env, nil

}
//类型定义

type Var string
type literal float64
type unary struct {
	op rune
	x Expr
}
type binary struct {
	op rune
	x Expr
	y Expr
}
//增加一个三元操作符
type ternary struct {
	condition Expr
	x Expr
	y Expr
}
type call struct {
	fn string
	args []Expr
}
type Expr interface {
	Eval(env Env) float64
	String() string
}
type Env map[Var]float64

func (v Var) Eval(env Env) float64 {
	return env[v]
}
func (l literal) Eval(env Env) float64 {
	return float64(l)
}
func (expr unary) Eval(env Env) float64 {
	switch expr.op {
	case '+':
		return +expr.x.Eval(env)
	case '-':
		return -expr.x.Eval(env)
	}
	panic(fmt.Sprintf("unsupported unary operator %s", expr.op))
}
func (expr binary) Eval(env Env) float64 {
	switch expr.op {
	case '+':
		return expr.x.Eval(env) + expr.y.Eval(env)
	case '-':
		return expr.x.Eval(env) - expr.y.Eval(env)
	case '*':
		return expr.x.Eval(env) * expr.y.Eval(env)
	case '/':
		return expr.x.Eval(env) / expr.y.Eval(env)
	}
	panic(fmt.Sprintf("unsupported binary operator %s", expr.op))
}
func (expr call) Eval(env Env) float64 {
	switch expr.fn {
	case "pow":
		return math.Pow(expr.args[0].Eval(env), expr.args[1].Eval(env))
	case "sin":
		return math.Sin(expr.args[0].Eval(env))
	case "sqrt":
		return math.Sqrt(expr.args[0].Eval(env))
	}
	panic(fmt.Sprintf("unsupported function call %s", expr.fn))
}
func (expr ternary) Eval(env Env) float64 {
	if expr.condition.Eval(env) > 0 {
		return expr.x.Eval(env)
	}
	return expr.y.Eval(env)
}

func (v Var) String() string {
	return string(v)
}
func (l literal) String() string {
	return strconv.FormatFloat(float64(l), 'f', -1, 64)
}
func (expr unary) String() string {
	return fmt.Sprintf(`unary{
    op: %c, 
    x: %s
}`, expr.op, indent(expr.x.String()))
}
func (expr binary) String() string {
	return fmt.Sprintf(`binary{
    op: %c,
    x: %s,
    y: %s
}`, expr.op, indent(expr.x.String()), indent(expr.y.String()))
}
func (expr ternary) String() string {
	return fmt.Sprintf(`ternary{
    condition: %s,
    x: %s,
    y: %s
}`, indent(expr.condition.String()), indent(expr.x.String()), indent(expr.y.String()))
}
func (expr call) String() string {
	var args []string
	for _, arg := range expr.args {
		args = append(args, indent(arg.String()))
	}
	return fmt.Sprintf(`call{
    fn: %s,
    args: [%s]
}`, expr.fn, strings.Join(args, ", "))
}
func indent(s string) string {
	pad := fmt.Sprintf("%*s", 4, "")
	return strings.Replace(s, "\n", "\n" + pad, -1)
}
//解析器
type lexer struct {
	scan scanner.Scanner
	token rune
}
//解析的错误类型
type lexPanic string

//scan返回的是token的类型，scanner.(Ident|Int|Float|EOF)，如果不是这几类，就返回对应的rune字符
func (lex *lexer) next() {
	lex.token = lex.scan.Scan()
}
//返回的是具体的token字符
func (lex *lexer) text() string{
	return lex.scan.TokenText()
}
func (lex *lexer) describe() string{
	switch lex.token {
	case scanner.EOF: 
		return "end of file"
	case scanner.Ident:
		return fmt.Sprintf("identifier %s", lex.text())
	case scanner.Int, scanner.Float:
		return fmt.Sprintf("number %s", lex.text())
	}
	return fmt.Sprintf("%q", lex.token)
}

//返回操作符的优先级
func precedence(op rune) int {
	switch op {
	case '*', '/':
		return 2
	case '-', '+':
		return 1
	}
	return 0
}

func Parse(str string) (_ Expr, err error) {
	defer func() {
		switch x := recover().(type) {
		case nil:
		case lexPanic: 
			err = fmt.Errorf("%s", x)
		default:
			panic(x)
		}
	}()
	lex := new(lexer)
	lex.scan.Init(strings.NewReader(str))
	lex.scan.Mode = scanner.ScanIdents|scanner.ScanInts|scanner.ScanFloats
	//初始化token
	lex.next()
	expr := parseExpr(lex)
	//没有解析完全
	if lex.token != scanner.EOF {
		return nil, fmt.Errorf("unexpected %s", lex.describe())
	}
	return expr, nil
}

func parseExpr(lex *lexer) Expr {
	return parseTernary(lex)
}
//解析二元操作符
//binary = unary (+ binary)*
//parseBinay会在遇到优先级小于lastPrec的操作符时停止
func parseBinary(lex *lexer, lastPrec int) Expr {
	lhs := parseUnary(lex)
	for prec := precedence(lex.token); prec >= lastPrec; prec-- {
		for precedence(lex.token) == prec {
			op := lex.token	
			lex.next()
			rhs := parseBinary(lex, prec+1)
			lhs = binary{op, lhs, rhs}
		}
	}
	return lhs
}

//解析三元操作符
//ternary = binary (? ternary : binary)+
func parseTernary(lex *lexer) Expr {
	expr := parseBinary2(lex, 1)
	for {
		if lex.token != '?' {
			break
		}
		//消费?
		lex.next()
		x := parseTernary(lex)
		if lex.token != ':' {
			msg := fmt.Sprintf("unexpected %s, expected ':'", lex.describe())
			panic(lexPanic(msg))
		}
		//消费:
		lex.next()
		y := parseBinary2(lex, 1)
		expr = ternary{expr, x, y}
	}
	return expr
}
/**
 * parseBinary的另一个版本，稍简单一些
 * parseBinary解析二元操作，默认传入的操作符优先级为1
 * 如果op的优先级低于lastPrec, 要么是非法操作符，要么小于上一层次的二元操作，直接返回
 */
func parseBinary2(lex *lexer, lastPrec int) Expr {
	lhs := parseUnary(lex)
	for {
		op := lex.token
		prec := precedence(op)
		if prec < lastPrec {
			break
		}
		lex.next()
		rhs := parseBinary2(lex, prec)
		lhs = binary{op, lhs, rhs}
	}
	return lhs
}

//解析一元操作符 unary = +expr|primary
func parseUnary(lex *lexer) Expr {
	if lex.token == '-' || lex.token == '+' {
		op := lex.token
		lex.next()
		return unary{op, parseUnary(lex)}
	}
	return parsePrimary(lex)
}
//primary = id
//        | id(expr,...expr)
//        | (expr)
func parsePrimary(lex *lexer) Expr {
	switch lex.token {
	case scanner.Ident:
		id := lex.text()
		lex.next()
		if lex.token != '(' {
			return Var(id)
		}
		//消费掉'('
		lex.next()
		var args []Expr
		for lex.token != ')' {
			args = append(args, parseExpr(lex))
			//消费掉分隔符
			if lex.token != ',' {
				break
			}
			lex.next()
		}
		if lex.token != ')' {
			msg := fmt.Sprintf("got %q, want ')'", lex.describe())
			panic(lexPanic(msg))
		}
		//消费 ')'
		lex.next()
		return call{id, args}
	case scanner.Int, scanner.Float:
		f, err := strconv.ParseFloat(lex.text(), 64)
		if err != nil {
			panic(lexPanic(err.Error()))
		}
		lex.next()
		return literal(f)
	case '(':
		lex.next()
		expr := parseExpr(lex)
		if lex.token != ')' {
			msg := fmt.Sprintf("got %q, want ')'", lex.describe())
			panic(lexPanic(msg))
		}
		//消费 ')'
		lex.next()
		return expr
	}
	msg := fmt.Sprintf("unexpected %s", lex.describe())
	panic(lexPanic(msg))
}