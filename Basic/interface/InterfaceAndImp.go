package main

type Animal interface {
	Eat()
	Run()
}

type Dog struct {
	Name string
}

func (d Dog) Eat() {
	println(d.Name + " eats")
}
func (d Dog) Run() {
	println(d.Name + " runs")
}
func (d Dog) Bark() {
	println(d.Name + " barks")
}

func main() {
	dog := Dog{Name: "dog bob"}
	dog.Eat()
	dog.Run()
	dog.Bark()
	Animal(dog).Eat()
	Animal(dog).Run()
}
