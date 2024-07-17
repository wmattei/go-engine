package engine

type Scene struct {
	models map[string]*Model
	camera Camera
}

func NewScene() Scene {
	return Scene{
		models: make(map[string]*Model),
	}
}

func (s *Scene) AddModel(model *Model) {
	s.models[model.ID] = model
}

func (s *Scene) SetCamera(camera Camera) {
	s.camera = camera
}
