package word

type Repository interface {
	SaveWord(word *Word) (*Word, error)
}
