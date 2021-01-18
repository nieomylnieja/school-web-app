package student

type Student struct {
	Name      string `json:"name" bson:"name"`
	Surname   string `json:"surname" bson:"surname"`
	Age       int    `json:"age" bson:"age"`
	TeacherID string `json:"teacherId,omitempty" bson:"teacherId,omitempty"`
}
