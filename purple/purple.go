package purple

import (
	"fmt"
	"image/color"
	"sync"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dsvg"
)

type RGBA = color.RGBA

type Purple struct {
	State *State
	Year  string

	Stats      map[string]RGBA
	OutputPath string

	Scale       float64
	StrokeWidth float64
	StrokeColor RGBA
}

type ChanCounty struct {
	Name string
	Path *draw2d.Path
}

func (p *Purple) UseDefaults() {
	p.Scale = 10
	p.StrokeWidth = 0.05
	p.StrokeColor = RGBA{0, 0, 0, 255}
}

func (p *Purple) Draw() {
	svg := draw2dsvg.NewSvg()
	gc := draw2dsvg.NewGraphicContext(svg)

	width := p.State.Bbox.Max.X - p.State.Bbox.Min.X
	height := p.State.Bbox.Max.Y - p.State.Bbox.Min.Y

	svg.Width = fmt.Sprintf("%fpx", width*p.Scale)
	svg.Height = fmt.Sprintf("%fpx", height*p.Scale)

	gc.Scale(p.Scale, p.Scale)
	p.drawState(gc)

	draw2dsvg.SaveToSvgFile(p.OutputPath, svg)
}

func (p *Purple) drawState(gc *draw2dsvg.GraphicContext) {
	gc.SetStrokeColor(p.StrokeColor)
	gc.SetLineWidth(p.StrokeWidth)

	counties := p.projectCounties()

	for county := range counties {
		clr := p.getCountyColor(county.Name)
		gc.SetFillColor(clr)
		gc.Fill(county.Path)
		gc.Stroke(county.Path)
	}
}

func (p *Purple) getCountyColor(county string) RGBA {
	if v, ok := p.Stats[county]; ok {
		return v
	}
	return RGBA{0, 0, 0, 0}
}

func (p *Purple) projectCounties() chan *ChanCounty {
	counties := make(chan *ChanCounty, p.State.CountiesN)

	wg := new(sync.WaitGroup)
	wg.Add(p.State.CountiesN)

	for _, county := range p.State.Counties {
		go p.projectCounty(wg, county, counties)
	}

	wg.Wait()
	close(counties)

	return counties
}

func (p *Purple) projectCounty(wg *sync.WaitGroup, county County, counties chan *ChanCounty) {
	path := new(draw2d.Path)
	path.Components = make([]draw2d.PathCmp, county.PointsN+1)
	path.Points = make([]float64, 2*county.PointsN+2)

	for i, point := range county.Points {
		x := point.X - p.State.Bbox.Min.X
		y := p.State.Bbox.Max.Y - point.Y

		path.Components[i] = draw2d.LineToCmp
		path.Points[i*2] = x
		path.Points[i*2+1] = y
	}

	path.Components[0] = draw2d.MoveToCmp
	path.Components[county.PointsN] = draw2d.LineToCmp
	path.Points[2*county.PointsN] = path.Points[0]
	path.Points[2*county.PointsN+1] = path.Points[1]

	newCounty := ChanCounty{}
	newCounty.Name = county.Name
	newCounty.Path = path

	counties <- &newCounty

	wg.Done()
}
