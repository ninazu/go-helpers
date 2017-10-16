package i18n

import (
	"bytes"
	"fmt"
	"go/token"
	"os"
	"strings"
	"go/ast"
	"io/ioutil"
	"path"
	"go/parser"
	"go/format"
	"errors"
	"strconv"
	"sort"
)

type generator struct {
	buf bytes.Buffer
}

func (g *generator) Printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}

func Generate(dirPath string, lang []string) (error) {
	if len(instance.path)==0 {
		return errors.New("you MUST call i18n.Init() prior to using it")
	}
	
	var params = make(map[string]bool)
	
	err := scanFolder(dirPath, &params)
	
	if err != nil {
		return err
	}
	
	instance.load()
	newValues := make(map[string]map[string]struct{})
	
	var text string
	
	for k := range params {
		if params[k] {
			text = k
		} else {
			text, err = strconv.Unquote("\"" + k + "\"")
			
			if err != nil {
				return err
			}
		}
		
		for _, l := range lang {
			if _, ok := instance.cache[l][text]; !ok {
				if _, ok := instance.cache[l]; !ok {
					instance.cache[l] = make(map[string]string)
				}
				
				instance.cache[l][text] = text
			}
			
			if _, ok := newValues[l]; !ok {
				newValues[l] = make(map[string]struct{})
			}
			
			newValues[l][text] = struct{}{}
		}
	}
	
	indexedLang := make(map[string]struct{})
	langKeys := []string{}
	
	for _, v := range lang {
		indexedLang[v] = struct{}{}
		langKeys = append(langKeys, v)
	}
	
	indexedKeys := make(map[string]struct{})
	
	for l, v := range instance.cache {
		if _, ok := indexedLang[l]; !ok {
			indexedLang[l] = struct{}{}
			langKeys = append(langKeys, l)
		}
		
		for k := range v {
			indexedKeys[k] = struct{}{}
		}
	}
	
	var uniqueKeys []string
	
	for k := range indexedKeys {
		uniqueKeys = append(uniqueKeys, k)
	}
	
	sort.Strings(uniqueKeys)
	sort.Strings(langKeys)
	
	g := &generator{}
	
	g.Printf("package i18n\n")
	g.Printf("func (i *extends) load() {\n")
	g.Printf("i.cache = make(map[string]map[string]string)\n")
	g.Printf("\n")
	
	for l := range indexedLang {
		comment := ""
		
		if len(newValues[l]) == 0 {
			comment = "//"
		}
		
		g.Printf("%si.cache[\"%s\"] = make(map[string]string)\n", comment, l)
	}
	
	g.Printf("\n")
	
	for _, k := range uniqueKeys {
		for l := range langKeys {
			if v, ok := instance.cache[langKeys[l]][k]; ok {
				comment := ""
				
				if _, ok := newValues[langKeys[l]][k]; !ok {
					comment = "//"
				}
				
				g.Printf("%si.cache[\"%s\"][%s]=%s\n", comment, langKeys[l], strconv.Quote(k), strconv.Quote(v))
			}
		}
		
		if len(lang) > 1 {
			g.Printf("\n")
		}
	}
	
	g.Printf("}\n")
	
	src, err := g.format()
	output := string(src)
	
	if err != nil {
		fmt.Println(output)
		
		return err
	}
	
	outputName := path.Join(instance.path, "source.go")
	err = ioutil.WriteFile(outputName, src, 0644)
	
	if err != nil {
		return err
	}
	
	return nil
}

func (g *generator) format() ([]byte, error) {
	src, err := format.Source(g.buf.Bytes())
	
	if err != nil {
		return g.buf.Bytes(), errors.New(fmt.Sprintf("warning: compile the package to analyze the error, %s", err))
	}
	
	return src, nil
}

func scanFolder(dirPath string, params *map[string]bool) (error) {
	filSet := token.NewFileSet()
	
	d, err := parser.ParseDir(filSet, dirPath, func(fi os.FileInfo) bool {
		return path.Ext(fi.Name()) == ".go"
	}, 0)
	
	if err != nil {
		return err
	}
	
	for _, p := range d {
		for _, f := range p.Files {
			getFuncFirstParam(f, params)
		}
		
	}
	
	dir, err := os.Open(dirPath)
	
	if err != nil {
		return err
	}
	
	info, err := dir.Readdir(0)
	
	if err != nil {
		return err
	}
	
	for _, fi := range info {
		if fi.IsDir() && !strings.HasPrefix(fi.Name(), ".") {
			err := scanFolder(path.Join(dirPath, fi.Name()), params)
			
			if err != nil {
				return err
			}
		}
	}
	
	return nil
}

func getFuncFirstParam(f ast.Node, m *map[string]bool) {
	//TODO HARDCODED
	var nameSpace = "i18n"
	var funcName = "Translate"
	
	var lParen = 0
	var rParen = 0
	var funcDetected bool
	var litDetected bool
	var concatDetected bool
	var namespaceDetected bool
	var result = ""
	var quote = ""
	
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.CallExpr:
			lParen = int(x.Lparen)
			rParen = int(x.Rparen)
		
		case *ast.BasicLit:
			if funcDetected && x.Kind == token.STRING {
				litDetected = true
				quote = x.Value[0:1]
				value := x.Value[1:len(x.Value)-1]
				
				if concatDetected {
					result += value
				} else if result == "" {
					result = value
				}
			}
		
		case *ast.Ident:
			switch x.Name {
			case funcName:
				if namespaceDetected {
					funcDetected = true
				}
			case nameSpace:
				namespaceDetected = true
			}
		
		case *ast.BinaryExpr:
			if n.(*ast.BinaryExpr).Op == token.ADD {
				concatDetected = true
			}
		
		default:
			if funcDetected && litDetected && n != nil {
				if rParen > 0 && int(n.Pos()) > rParen {
					(*m)[result] = quote == "`"
					namespaceDetected = false
					funcDetected = false
					litDetected = false
					concatDetected = false
					result = ""
				}
			}
		}
		
		return true
	})
}
