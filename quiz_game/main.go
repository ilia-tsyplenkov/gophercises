package main

type QuizGame struct {
}

func (g *QuizGame) CheckAnswers(got, want string) bool {
	return true
}

func (g *QuizGame) Launch() error {
	return nil
}

func main() {

}
