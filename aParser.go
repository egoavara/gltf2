package gltf2

import (
	"encoding/json"
	"fmt"
	"github.com/iamGreedy/glog"

	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type parser struct {
	dst    *GLTF
	src    *SpecGLTF
	parsed bool
	//
	cause error
	err   error
	//
	strictness Strictness
	rd         io.Reader
	dir        http.Dir
	// [progress]
	// + logging task
	logger *glog.Glogger
	//
	pres []PreTask
	//whiles map[string][]WhilePlugin
	posts []PostTask
}

func (s *parser) Reader(reader io.Reader) *parser {
	if reader == nil {
		s.setCauseError(ErrorParserOption, errors.Errorf("nil reader"))
		return s
	}
	s.rd = reader
	if f, ok := s.rd.(*os.File); ok {
		return s.Directory(filepath.Dir(f.Name()))
	}
	return s
}
func (s *parser) Strictness(strictness Strictness) *parser {
	if strictness < LEVEL0 {
		strictness = LEVEL0
	}
	if strictness > LEVEL3 {
		strictness = LEVEL3
	}
	s.strictness = strictness
	return s
}
func (s *parser) Plugin(plugins ...Task) *parser {
	for _, plugin := range plugins {
		if prp, ok := plugin.(PreTask); ok {
			s.pres = append(s.pres, prp)
		}
		if pop, ok := plugin.(PostTask); ok {
			s.posts = append(s.posts, pop)
		}
	}
	return s
}
func (s *parser) Directory(path string) *parser {
	fi, err := os.Stat(path)
	if err != nil {
		s.setCauseError(ErrorParserOption, err)
	}
	if !fi.IsDir() {
		s.setCauseError(ErrorParserOption, errors.Errorf("Not directory '%s'", path))
	} else {
		s.dir = http.Dir(path)
	}
	return s
}
func (s *parser) Logger(dst io.Writer) *parser {
	if dst != nil {
		s.logger = glog.New(log.New(dst, "[ glTF 2.0 ] ", log.LstdFlags), "    ")
	} else {
		s.setCauseError(ErrorParserOption, errors.Errorf("Parser().Logger(<not nillable>)"))
	}
	return s
}
func (s *parser) GetLogger() *glog.Glogger {
	return s.logger
}

//

func (s *parser) Error() error {
	if s.src == nil {
		return errors.WithMessage(ErrorParser, "Closed Parse")
	}
	if s.cause != nil {
		s.logger.Println(s.err)
		return errors.WithMessage(s.cause, s.err.Error())
	}
	if s.err != nil {
		s.logger.Println(s.err)
		return errors.WithMessage(ErrorParser, s.err.Error())
	}
	return nil
}
func (s *parser) setError(err error) {
	s.cause = nil
	s.err = err
}
func (s *parser) setCauseError(cause error, err error) {
	s.cause = cause
	s.err = err
}
func (s *parser) Parsed() bool {
	return s.parsed
}
func (s *parser) Parse() (*GLTF, error) {
	if err := s.Error(); err != nil {
		return nil, err
	}
	if s.parsed {
		return s.dst, nil
	}
	//====================================================//
	// constraint
	if s.rd == nil {
		s.setCauseError(ErrorParser, errors.Errorf("No reader"))
		return nil, s.Error()
	}
	//
	ctx := &parserContext{
		ref: s,
	}
	//====================================================//
	// json parse here
	dec := json.NewDecoder(s.rd)
	s.logger.Println("json decode start...")
	if err := dec.Decode(s.src); err != nil {
		s.setCauseError(ErrorJSON, err)
		return nil, err
	}
	s.logger.Println("json decode complete")
	//====================================================//
	// pre plugin
	if s.pres != nil {
		s.logger.Println("PrePlugins start")
		for i, pre := range s.pres {
			s.logger.Printf("PreTask(%d/%d) '%s'\n", i+1, len(s.pres), pre.TaskName())
			pre.PreLoad(ctx, s.src, s.logger.Indent())
		}
		s.logger.Println("PrePlugins complete")
	}
	//====================================================//
	// Syntax check
	s.logger.Println("syntax check")
	if err := recurSyntax(s.src, s.src, s.strictness); err != nil {
		s.setCauseError(ErrorGLTFSpec, err)
		return nil, s.Error()
	}
	s.logger.Println("syntax valid")
	//====================================================//
	// Convert
	s.dst = s.src.To(nil).(*GLTF)
	s.logger.Println("structure setup")
	recurTo(s.dst, s.src, ctx)
	s.logger.Println("structure setup complete")
	//====================================================//
	// Link
	s.logger.Println("glTFid reference linking start")
	if err := recurLink(s.dst, nil, s.dst, s.src); err != nil {
		s.setCauseError(ErrorGLTFLink, err)
		return nil, s.Error()
	}
	s.logger.Println("glTFid reference linked")
	//====================================================//
	if s.posts != nil {
		s.logger.Println("PostTask start")
		// post plugin
		for i, post := range s.posts {
			s.logger.Printf("PostTask(%d/%d) '%s'\n", i+1, len(s.posts), post.TaskName())
			if err := post.PostLoad(ctx, s.dst, s.logger.Indent()); err != nil {
				s.setCauseError(ErrorPlugin, errors.WithMessage(err, fmt.Sprintf("plugin name : %s", post.TaskName())))
				return nil, s.Error()
			}
		}
		s.logger.Println("PostTask complete")
	}
	s.logger.Println("complete.")
	return s.dst, nil
}
func (s *parser) Close() error {
	s.src = nil
	return nil
}

func Parser() *parser {
	res := &parser{
		src: new(SpecGLTF),

		strictness: LEVEL1,
		dir:        http.Dir("."),
	}
	return res
}

func recurSyntax(root, target ToGLTF, strictness Strictness) error {
	if target == nil{
		return nil
	}
	fmt.Println(target)
	if err := target.Syntax(strictness, root); err != nil {
		return err
	}
	if tc, ok := target.(ChildrunToGLTF); ok {
		for i := 0; i < tc.LenChild(); i++ {
			if err := recurSyntax(root, tc.GetChild(i), strictness); err != nil {
				return err
			}
		}
	}
	return nil
}
func recurTo(data interface{}, target ToGLTF, ctx *parserContext) {
	if tc, ok := target.(ChildrunToGLTF); ok {
		for i := 0; i < tc.LenChild(); i++ {
			child := tc.GetChild(i)
			childData := child.To(ctx)
			recurTo(childData, child, ctx)
			tc.SetChild(i, data, childData)
		}
	}
}
func recurLink(root *GLTF, parent, data interface{}, target ToGLTF) error {
	if tl, ok := target.(LinkToGLTF); ok {
		if err := tl.Link(root, parent, data); err != nil {
			return err
		}
	}
	if tc, ok := target.(ChildrunToGLTF); ok {
		for i := 0; i < tc.LenChild(); i++ {
			if err := recurLink(root, data, tc.ImpleGetChild(i, data), tc.GetChild(i)); err != nil {
				return err
			}
		}
	}
	return nil
}

//func recurWhile(data interface{}, target ToGLTF, plugins map[string][]WhilePlugin) (err error, pluginName string){
//	for _, plugin := range plugins[target.Scheme()] {
//		if err := plugin.WhileLoad(target, data); err != nil{
//			return err, plugin.TaskName()
//		}
//	}
//	if tc, ok := target.(ChildrunToGLTF); ok{
//		for i := 0; i < tc.LenChild(); i++ {
//			if err, pluginName = recurWhile(tc.ImpleGetChild(i, data), tc.GetChild(i), plugins); err != nil{
//				return err, pluginName
//			}
//		}
//	}
//	return nil, ""
//}

type parserContext struct {
	ref *parser
}

func (s *parserContext) Strictness() Strictness {
	return s.ref.strictness
}
func (s *parserContext) Directory() string {
	return string(s.ref.dir)
}
func (s *parserContext) Specification() *SpecGLTF {
	return s.ref.src
}
func (s *parserContext) GLTF() *GLTF {
	return s.ref.dst
}
