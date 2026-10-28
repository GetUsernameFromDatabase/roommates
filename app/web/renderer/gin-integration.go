package renderer

import "github.com/gin-gonic/gin"

// ReplaceGinRenderer replaces [gin.Engine.HTMLRender] with [HTMLTemplRenderer]
//
// Uses [gin.Engine.HTMLRender] as a fallback html renderer, see [example]
//
// [example]: https://github.com/a-h/templ/blob/main/examples/integration-gin/main.go#L16
func ReplaceGinRenderer(r *gin.Engine) {
	ginHTMLRenderer := r.HTMLRender // store previous one
	r.HTMLRender = &HTMLTemplRenderer{FallbackHTMLRenderer: ginHTMLRenderer}
}
