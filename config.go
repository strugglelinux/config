package config

import (
	"io"
	"os"
	"bufio"
	"strings"
	"fmt"
	"errors"
)

const (
	CRLF     = '\n'
	CommentTag  = "#"
	SpliterTag  = "="
	SectionSTag = "[" //标签开始
	SectionETag = "]" //标签结束
	// memory unit
	Byte = 1
	KB   = 1024 * Byte
	MB   = 1024 * KB
	GB   = 1024 * MB

)
type Config struct {
	data  map[string]*Section
}

type Section struct {
	Name string
	element map[string] string
}

func New() *Config {
	return &Config{data: map[string]*Section{}}
}


func (c *Config) Load (file string) error{
	f , err := os.Open(file)
	if  err != nil {
		return  err
	}
	defer f.Close()
	return c.Parse(f)
}

func ( c *Config) Parse(reader io.Reader) error {
	var (
		rd = bufio.NewReader(reader)
		row  int
		key  string
		value string
		section *Section
	)

	for {
		row++
		line , err := rd.ReadString(CRLF)
		line = strings.Replace(line, "\n", "", -1)
		if err == io.EOF && len(line) == 0 {
			break
		}
		if err != nil &&  err != io.EOF {
			return err
		}
		if strings.HasPrefix(line,CommentTag) {
			continue
		}
		if strings.HasPrefix(line,SectionSTag) {
			if !strings.HasSuffix(line,SectionETag) {
				return errors.New(fmt.Sprintf("no end section: %s at :%d " , SectionETag,row) )
			}

			sectionVal  := line[1:len(line)-1]
			if s,ok := c.data[sectionVal]; !ok {
				s = &Section{ sectionVal,map[string]string{}}
				c.data[sectionVal] = s
				section = s
			}else{
				return errors.New(fmt.Sprintf("section: %s already exists at %d ",sectionVal,row))
			}
			continue
		}
		idx := strings.Index(line,SpliterTag)

		if idx > 0 {
			key = strings.TrimSpace(line[:idx])
			if len(line) > idx {
				value = strings.TrimSpace(line[idx+1:])
			}
		}else {
			return errors.New(fmt.Sprintf("no spliter in key: %s at %d ",line,row))
		}
		if section == nil {
			return errors.New(fmt.Sprintf("no section for key: %s at %d", line, row))
		}
		if _,ok := section.element[key] ;ok {
			return errors.New(fmt.Sprintf("section: %s already has key: %s at %d",section.Name , key, line))
		}
		section.element[key] = value
	}
	return nil
}

func (c *Config) Get(section string ) *Section {
	s ,_ := c.data[section]
	return s
}

func (c *Config) GetAll() map[string]*Section {
	return c.data
}

func (s *Section) Get(key string) string {
	k ,_ := s.element[key]
	return k
}