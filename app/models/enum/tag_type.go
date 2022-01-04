package enum

type TagType string

const (
	TagBlank           TagType = ""
	TagStarUser        TagType = "star_user"
	TagProblemUser     TagType = "problem_user"
	TagRecommendStatus TagType = "recommend_status"
)
