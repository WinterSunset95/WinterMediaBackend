package movies

import (
	"encoding/json"
	"fmt"

	"github.com/WinterSunset95/WinterMediaBackend/database"
	"github.com/WinterSunset95/WinterMediaBackend/models"
	"github.com/gin-gonic/gin"
)

func ApplyRoutes(r *gin.RouterGroup) {
	db := database.DB
	_ = db
	movies := r.Group("/movies")
	{
		movies.GET("/list", func(ctx *gin.Context) {
			var movieList []models.Movie
			query := `
			select m.id, m.title, m.poster, m.overview, m.duration, m.releaseDate, m.director_id, json_arrayagg(distinct l.name) as languages,
			json_arrayagg(distinct g.name) as genres,
			json_arrayagg(distinct
			json_object(
			'personId', p.person_id,
			'name', p.name,
			'picture', p.picture,
			'bio', p.bio,
			'role', mc.role
			)
			),
			json_arrayagg(distinct
			json_object(
			'personId', p2.person_id,
			'name', p2.name,
			'picture', p2.picture,
			'bio', p2.bio,
			'role', mr.role
			)
			)
			from Movies m 
			join Movie_Languages ml on m.id=ml.movie_id 
			join Languages l on ml.language_id=l.id 
			join Movie_Genres mg on m.id=mg.movie_id
			join Genres g on mg.genre_id=g.id
			join Movie_Cast mc on m.id=mc.movie_id
			join Persons p on mc.person_id=p.person_id
			join Movie_Crew mr on m.id=mr.movie_id
			join Persons p2 on mr.person_id=p2.person_id
			group by m.id
			`
			rows, err := db.Query(query)
			if err != nil {
				fmt.Println(err)
			}
			for rows.Next() {
				var id string
				var title string
				var poster string
				var overview string
				var duration int
				var releaseDate string
				var director string
				var languages string
				var languagesArr []string
				var genres string
				var genresArr []string
				var cast string
				var castObj []models.Person
				var crew string
				var crewObj []models.Person
				err = rows.Scan(&id, &title, &poster, &overview, &duration, &releaseDate, &director, &languages, &genres, &cast, &crew)
				if err != nil {
					fmt.Println(err)
				}
				json.Unmarshal([]byte(languages), &languagesArr)
				json.Unmarshal([]byte(genres), &genresArr)
				json.Unmarshal([]byte(cast), &castObj)
				json.Unmarshal([]byte(crew), &crewObj)
				movie := models.Movie{
					Id: id,
					Title: title,
					Poster: poster,
					Overview: overview,
					Duration: duration,
					ReleaseDate: releaseDate,
					Languages: languagesArr,
					Cast: castObj,
					Crew: crewObj,
					Genres: genresArr,
					Director: director,
				}
				movieList = append(movieList, movie)
			}
			if err != nil {
				fmt.Println(err)
			}
			ctx.JSON(200, &movieList)
		})
	}
}
