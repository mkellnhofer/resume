package main

import (
	"bufio"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/goccy/go-yaml"
)

const (
	templateName    = "resume.tmpl"
	stringsDirName  = "strings"
	stringsFileName = "^strings-[a-z]{2}\\.yaml$"
	inputDirName    = "data"
	inputFileName   = "^resume-[a-z]{2}\\.yaml$"
	outputDirName   = "web"
	outputFileName  = "index-[code].html"
)

// Resume holds resume data.
type Resume struct {
	IsDefault bool               `yaml:"isDefault"`
	Basics    *ResumeBasics      `yaml:"basics"`
	Intro     string             `yaml:"intro"`
	Skills    *ResumeSkills      `yaml:"skills"`
	Work      []*ResumeWork      `yaml:"work"`
	Education []*ResumeEducation `yaml:"education"`
}

// ResumeBasics holds basic resume data.
type ResumeBasics struct {
	Name  string        `yaml:"name"`
	Title string        `yaml:"title"`
	Links []*ResumeLink `yaml:"links"`
}

// ResumeLinkType defines available link types.
type ResumeLinkType string

// Available link types.
const (
	Email   ResumeLinkType = "email"
	Link    ResumeLinkType = "link"
	Address ResumeLinkType = "address"
)

// ResumeLink holds link data.
type ResumeLink struct {
	Type        ResumeLinkType `yaml:"type"`
	Icon        string         `yaml:"ico"`
	Text        string         `yaml:"text"`
	URL         string         `yaml:"url"`
	Description string         `yaml:"desc"`
}

// ResumeSkills holds skills data.
type ResumeSkills struct {
	Major []*ResumeMajorSkill `yaml:"major"`
	Minor []*ResumeMinorSkill `yaml:"minor"`
}

// ResumeMajorSkill holds data for a major skill.
type ResumeMajorSkill struct {
	Name     string   `yaml:"name"`
	Refs     []string `yaml:"refs"`
	Progress int      `yaml:"progress"`
}

// ResumeMinorSkill holds data for a minor skill.
type ResumeMinorSkill struct {
	Name string   `yaml:"name"`
	Refs []string `yaml:"refs"`
}

// ResumeWork holds work data.
type ResumeWork struct {
	Company    string   `yaml:"company"`
	StartDate  string   `yaml:"startDate"`
	EndDate    string   `yaml:"endDate"`
	Location   string   `yaml:"location"`
	Title      string   `yaml:"title"`
	SkillRef   string   `yaml:"skillRef"`
	Highlights []string `yaml:"highlights"`
}

// ResumeEducation holds education data.
type ResumeEducation struct {
	Institution string                 `yaml:"institution"`
	StartDate   string                 `yaml:"startDate"`
	EndDate     string                 `yaml:"endDate"`
	Location    string                 `yaml:"location"`
	Degree      string                 `yaml:"degree"`
	SkillRef    string                 `yaml:"skillRef"`
	Work        []*ResumeEducationWork `yaml:"work"`
}

// ResumeEducationWork holds education work data.
type ResumeEducationWork struct {
	Type       string   `yaml:"type"`
	StartDate  string   `yaml:"startDate"`
	EndDate    string   `yaml:"endDate"`
	Location   string   `yaml:"location"`
	SkillRef   string   `yaml:"skillRef"`
	Highlights []string `yaml:"highlights"`
}

// Strings holds general strings for a language.
type Strings struct {
	Language string            `yaml:"language"`
	Strings  map[string]string `yaml:"strings"`
}

// Language defines the language of input/output data.
type Language struct {
	Code      string
	Name      string
	IsDefault bool
}

// InputData holds input data.
type InputData struct {
	Lng  *Language
	Res  *Resume
	Strs map[string]string
}

// OutputData holds output data.
type OutputData struct {
	Lng     *Language
	AltLngs []*Language
	Res     *Resume
	Strs    map[string]string
}

func main() {
	// Read input data
	inputData := readInputData()
	// Create output data
	outputData := createOutputData(inputData)

	// Read template
	templ := readTemplate()
	// Generate HTML
	generateHTMLs(templ, outputData)
}

func readInputData() []*InputData {
	// Read files
	resumes := readResumeFiles()
	strings := readStringsFiles()

	// Create input data structures
	var data []*InputData
	defCnt := 0
	for lc, r := range resumes {
		if r.IsDefault {
			defCnt++
		}

		s, _ := strings[lc]
		if s == nil {
			log.Fatalf("Strings for '%s' is missing!", lc)
		}

		lng := Language{lc, s.Language, r.IsDefault}
		d := InputData{&lng, r, s.Strings}

		data = append(data, &d)
	}
	if defCnt > 1 {
		log.Fatal("Only one resume can be the default!")
	}

	return data
}

func readResumeFiles() map[string]*Resume {
	// Get files
	fs := getFiles(inputDirName)

	// Read/parse files
	m := make(map[string]*Resume)
	readFiles(inputDirName, inputFileName, fs, func(dn string, fn string) {
		lc := getFileLangCode(fn)
		bytes := readFile(dn, fn)
		var obj Resume
		parseFile(dn, fn, bytes, &obj)
		m[lc] = &obj
	})

	return m
}

func readStringsFiles() map[string]*Strings {
	// Get files
	fs := getFiles(stringsDirName)

	// Read/parse files
	m := make(map[string]*Strings)
	readFiles(stringsDirName, stringsFileName, fs, func(dn string, fn string) {
		lc := getFileLangCode(fn)
		bytes := readFile(dn, fn)
		var obj Strings
		parseFile(dn, fn, bytes, &obj)
		m[lc] = &obj
	})

	return m
}

func getFiles(dirName string) []os.FileInfo {
	fs, err := ioutil.ReadDir(dirName)
	if err != nil {
		log.Fatalf("Could read directry '%s': %s", dirName, err)
	}
	return fs
}

func readFiles(dirName string, fileNameTmpl string, files []os.FileInfo, readFunc func(string,
	string)) {
	r := regexp.MustCompile(fileNameTmpl)
	for _, file := range files {
		if !file.IsDir() && r.MatchString(file.Name()) {
			readFunc(dirName, file.Name())
		}
	}
}

func getFileLangCode(fileName string) string {
	return strings.Split(strings.Split(fileName, "-")[1], ".")[0]
}

func readFile(dirName string, fileName string) []byte {
	file, ofErr := os.Open(dirName + "/" + fileName)
	if ofErr != nil {
		log.Fatalf("Could not open file '%s'/'%s': %s", dirName, fileName, ofErr)
	}
	defer file.Close()

	bytes, rfErr := ioutil.ReadAll(file)
	if rfErr != nil {
		log.Fatalf("Could not read file '%s'/'%s': %s", dirName, fileName, rfErr)
	}
	return bytes
}

func parseFile(dirName string, fileName string, bytes []byte, obj interface{}) {
	pfErr := yaml.Unmarshal(bytes, obj)
	if pfErr != nil {
		log.Fatalf("Could not parse file '%s'/'%s': %s", dirName, fileName, pfErr)
	}
}

func createOutputData(inData []*InputData) []*OutputData {
	outData := make([]*OutputData, 0)
	for _, in := range inData {
		out := OutputData{}
		out.Lng = in.Lng
		als := make([]*Language, 0)
		for _, in2 := range inData {
			if in2.Lng != in.Lng {
				als = append(als, in2.Lng)
			}
		}
		out.AltLngs = als
		out.Res = in.Res
		out.Strs = in.Strs
		outData = append(outData, &out)
	}
	return outData
}

func readTemplate() *template.Template {
	// Read template
	b, rErr := ioutil.ReadFile(templateName)
	if rErr != nil {
		log.Fatalf("Could load template '%s': %s", templateName, rErr)
	}

	// Register template
	tmpl := template.New(templateName)

	// Parse template
	tmpl, pErr := tmpl.Parse(string(b))
	if pErr != nil {
		log.Fatalf("Could parse template '%s': %s", templateName, pErr)
	}

	return tmpl
}

func generateHTMLs(template *template.Template, data []*OutputData) {
	for _, d := range data {
		generateHTML(template, d)
	}
}

func generateHTML(template *template.Template, data *OutputData) {
	// Creat filename
	var name string
	if !data.Res.IsDefault {
		name = outputDirName + "/" + strings.Replace(outputFileName, "[code]", data.Lng.Code, -1)
	} else {
		name = outputDirName + "/" + strings.Replace(outputFileName, "-[code]", "", -1)
	}

	// Create file
	f, cErr := os.Create(name)
	if cErr != nil {
		log.Fatalf("Could create output file '%s': %s", name, cErr)
	}
	defer f.Close()

	// Generate HTML
	w := bufio.NewWriter(f)
	rErr := template.Execute(w, data)
	if rErr != nil {
		log.Fatalf("Could not generate HTML! %s", rErr)
	}
	w.Flush()
}
