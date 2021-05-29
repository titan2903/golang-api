package libraryloadtemplate

import (
	"path/filepath"

	"github.com/gin-contrib/multitemplate"
)

//! This is a custom HTML render to support multi templates, ie. more than one *template.Template.
func LoadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	layouts, err := filepath.Glob(templatesDir + "/layouts/*") //! mengeload semua folder yang ada di dalam forlder template
	if err != nil {
		panic(err.Error())
	}

	includes, err := filepath.Glob(templatesDir + "/**/*") //! mengeload semua folder yang ada di dalam forlder template selain layouts
	if err != nil {
		panic(err.Error())
	}

	//? Generate our templates map from our layouts/ and includes/ directories
	for _, include := range includes {
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append(layoutCopy, include)
		r.AddFromFiles(filepath.Base(include), files...)
	}
	return r
}