package common

type Person struct {
	personId string
	picture string;
	name string;
	bio string;
}

type ParsedMovie struct {
	id string;
	title string;
	poster string;
	overview string;
	duration int;
	releaseDate string;
	languages []string;
	genres []string;
	cast []Person;
	crew []Person;
	director Person;
}

type Movie struct {
	id uint
	title uint
	poster string
	overview string
	duration int
	releaseDate string
	director string
}
