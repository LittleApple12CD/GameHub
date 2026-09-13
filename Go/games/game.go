package games

type Game interface {
	Run()
	HandleEvents() bool
	Update()
	Draw()
	Reset()
}