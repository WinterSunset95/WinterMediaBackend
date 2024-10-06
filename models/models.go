package models

// Structs
type Person struct {
	PersonId string `gorm:"primaryKey"`
	Picture string;
	Name string;
	Bio string;
	Role string;
}

type MovieCast struct {
	MovieId string `gorm:"primaryKey"`
	PersonId string `gorm:"primaryKey"`
}

type Movie struct {
	Id string `gorm:"primaryKey"`
	Title string
	Poster string
	Overview string
	Duration int
	ReleaseDate string
	Languages []string
	Genres []string
	Cast []Person
	Crew []Person
	Director string
}

type SmallMovieResult struct {
	Id string
	Title string
	Poster string
}

type FeaturedMovie struct {
	Id string
	Title string
	Poster string
	Cast []Person
}

type MovieInfoRequest struct {
	Id string
}
