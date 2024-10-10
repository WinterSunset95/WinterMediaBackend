package movies

import (
	"fmt"

	"github.com/WinterSunset95/WinterMediaBackend/database"
	"github.com/WinterSunset95/WinterMediaBackend/models"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

func checkAuthState(ctx *gin.Context, jwkSet jwk.Set) (jwt.Token, error) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(200, gin.H{
			"error": "No authorization header",
		})
		return nil, fmt.Errorf("No authorization header")
	}
	_, _ = jwt.Parse([]byte(authHeader), jwt.WithKeySet(jwkSet))
	return nil, nil
}

func ApplyRoutes(r *gin.RouterGroup) {
	db := database.DB
	_ = db
	movies := r.Group("/movies")
	{
		movies.GET("/list", func(ctx *gin.Context) {
			var movieList []models.SmallMovieResult

			query := `
			select id, title, poster from Movies
			`
			rows, err := db.Query(query)
			if err != nil {
				fmt.Println(err)
			}
			for rows.Next() {
				var id string
				var title string
				var poster string
				err = rows.Scan(&id, &title, &poster)
				if err != nil {
					fmt.Println(err)
				}
				movie := models.SmallMovieResult{
					Id: id,
					Title: title,
					Poster: poster,
				}
				movieList = append(movieList, movie)
			}
			if err != nil {
				fmt.Println(err)
			}
			ctx.JSON(200, &movieList)
		})

		movies.POST("/getinfo", func(ctx *gin.Context) {
			var movieInfoRequest models.MovieInfoRequest
			err := ctx.ShouldBindJSON(&movieInfoRequest)
			if err != nil {
				ctx.JSON(200, gin.H{
					"error": "Could not get POST data: " + err.Error(),
				})
				return
			}

			id := movieInfoRequest.Id
			var movie models.Movie

			{
				// Query the movie info
				query := "select id, title, poster, overview, duration, releaseDate from Movies where id='" + id + "'"
				rows, err := db.Query(query)
				if err != nil {
					fmt.Println(err)
					ctx.JSON(200, gin.H{
						"error": "Could not get movie info: " + err.Error(),
					})
				}
				for rows.Next() {
					err = rows.Scan(&movie.Id, &movie.Title, &movie.Poster, &movie.Overview, &movie.Duration, &movie.ReleaseDate)
					if err != nil {
						fmt.Println(err)
					}
				}
			}

			{
				// Query the director
				var director models.Person
				query := "select p.person_id, p.picture, p.name, p.bio from Movies m join Persons p on m.director_id = p.person_id where m.id = '" + id + "'"
				rows, err := db.Query(query)
				if err != nil {
					fmt.Println(err)
					ctx.JSON(200, gin.H{
						"error": "Could not get movie info: " + err.Error(),
					})
				}
				for rows.Next() {
					err = rows.Scan(&director.PersonId, &director.Picture, &director.Name, &director.Bio)
					if err != nil {
						fmt.Println(err)
					}
					director.Role = "Director"
					movie.Director = director
				}
			}
			
			{
				// Query the languages
				query := "select l.name from Languages l join Movie_Languages ml on l.id=ml.language_id where ml.movie_id='" + id + "'"
				rows, err := db.Query(query)
				if err != nil {
					fmt.Println(err)
					ctx.JSON(200, gin.H{
						"error": "Could not get movie languages: " + err.Error(),
					})
				}
				for rows.Next() {
					var language string
					err = rows.Scan(&language)
					if err != nil {
						fmt.Println(err)
					}
					movie.Languages = append(movie.Languages, language)
				}
			}

			{
				// Query the genres
				query := "select g.name from Genres g join Movie_Genres mg on g.id=mg.genre_id where mg.movie_id='" + id + "'"
				rows, err := db.Query(query)
				if err != nil {
					fmt.Println(err)
					ctx.JSON(200, gin.H{
						"error": "Could not get movie genres: " + err.Error(),
					})
				}
				for rows.Next() {
					var genre string
					err = rows.Scan(&genre)
					if err != nil {
						fmt.Println(err)
					}
					movie.Genres = append(movie.Genres, genre)
				}
			}

			{
				// For the cast, it is a little different
				var cast []models.Person
				query := "select p.person_id, p.picture, p.name, p.bio, mc.role from Persons p join Movie_Cast mc on p.person_id=mc.person_id where mc.movie_id='" + id + "'"
				rows, err := db.Query(query)
				if err != nil {
					fmt.Println(err)
					ctx.JSON(200, gin.H{
						"error": "Could not get movie cast: " + err.Error(),
					})
				}
				for rows.Next() {
					var person models.Person
					err = rows.Scan(&person.PersonId, &person.Picture, &person.Name, &person.Bio, &person.Role)
					if err != nil {
						fmt.Println(err)
					}
					cast = append(cast, person)
				}
				movie.Cast = cast
			}

			{
				// Now for the crew
				var crew []models.Person
				query := "select p.person_id, p.picture, p.name, p.bio, mc.role from Persons p join Movie_Crew mc on p.person_id=mc.person_id where mc.movie_id='" + id + "'"
				rows, err := db.Query(query)
				if err != nil {
					fmt.Println(err)
					ctx.JSON(200, gin.H{
						"error": "Could not get movie crew: " + err.Error(),
					})
				}
				for rows.Next() {
					var person models.Person
					err = rows.Scan(&person.PersonId, &person.Picture, &person.Name, &person.Bio, &person.Role)
					if err != nil {
						fmt.Println(err)
					}
					crew = append(crew, person)
				}
				movie.Crew = crew
			}

			ctx.JSON(200, &movie)
		})

		movies.GET("/featured", func(ctx *gin.Context) {
			var featuredMovie models.FeaturedMovie
			query := "select id, title, poster from Movies order by rand() limit 1"
			rows, err := db.Query(query)
			if err != nil {
				fmt.Println(err)
				ctx.JSON(200, gin.H{
					"error": "Could not get featured movie: " + err.Error(),
				})
			}
			for rows.Next() {
				err = rows.Scan(&featuredMovie.Id, &featuredMovie.Title, &featuredMovie.Poster)
				if err != nil {
					fmt.Println(err)
				}
			}

			{
				// For the cast, it is a little different
				var cast []models.Person
				query := "select p.person_id, p.picture, p.name, p.bio, mc.role from Persons p join Movie_Cast mc on p.person_id=mc.person_id where mc.movie_id='" + featuredMovie.Id + "'"
				rows, err := db.Query(query)
				if err != nil {
					fmt.Println(err)
					ctx.JSON(200, gin.H{
						"error": "Could not get movie cast: " + err.Error(),
					})
				}
				for rows.Next() {
					var person models.Person
					err = rows.Scan(&person.PersonId, &person.Picture, &person.Name, &person.Bio, &person.Role)
					if err != nil {
						fmt.Println(err)
					}
					cast = append(cast, person)
				}
				featuredMovie.Cast = cast
			}

			ctx.JSON(200, &featuredMovie)
		})
	}
}
