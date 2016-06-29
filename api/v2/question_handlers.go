package v2

import (
	"net/http"

	"partisan/auth"
	q "partisan/questions"

	"github.com/gin-gonic/gin"
)

var questionSets q.QuestionSets

func init() {
	questionSets = q.QuestionSets{
		q.QuestionSet{
			// Initial
			Questions: q.Questions{
				q.Question{
					// ProCapital
					Prompt: "Of all economic systems, Capitalism is most compatible with human nature.",
					Map:    []int{2, 3, 6, 7, 10, 11, 14, 15},
				},
				q.Question{
					// Anti-State
					Prompt: "Authority should always be questioned.",
					Map:    []int{8, 9, 10, 11, 12, 13, 14, 15},
				},
				q.Question{
					// Pro-State
					Prompt: "The police, in general, are good.",
					Map:    []int{1, 2, 3, 5, 6, 7, 10, 11},
				},
				q.Question{
					// Anti-Capital
					Prompt: "People over profits.",
					Map:    []int{0, 1, 4, 5, 8, 9, 12, 13},
				},
			},
		},
		q.QuestionSet{
			// Initial
			Questions: q.Questions{
				q.Question{
					// Right-Wing
					Prompt: "Men are more suited for positions of leadership, while women are more suited for positions of nurturing.",
					Map:    []int{2, 3, 6, 7, 10, 11, 14, 15},
				},
				q.Question{
					// Pro-State
					Prompt: "You should always support the troops.",
					Map:    []int{2, 3, 6, 7, 10, 11, 15},
				},
				q.Question{
					// Anti-Capital
					Prompt: "Healthcare should be free of charge.",
					Map:    []int{0, 1, 4, 5, 8, 9, 12, 13},
				},
				q.Question{
					// Anti-State
					Prompt: "Governments care about control, not the good of their people.",
					Map:    []int{8, 9, 10, 11, 12, 13, 14, 15},
				},
			},
		},
		q.QuestionSet{
			// Initial
			Questions: q.Questions{
				q.Question{
					// Pro-Capital
					Prompt: "Hard work leads to upward social mobility",
					Map:    []int{2, 3, 6, 7, 10, 11, 14, 15},
				},
				q.Question{
					// Pro-State
					Prompt: "Good social organization is centralized",
					Map:    []int{0, 1, 2, 3, 4, 5, 6, 7},
				},
				q.Question{
					// Anti-Capital
					Prompt: "Corporations should not self-regulate",
					Map:    []int{0, 1, 4, 5, 8, 9, 12, 13},
				},
				q.Question{
					// Anti-State
					Prompt: "Dissent is a vitrue",
					Map:    []int{8, 9, 10, 11, 12, 13, 14, 15},
				},
			},
		},
		q.QuestionSet{
			// Left-Wing
			// Mask:     []int{0, 1, 4, 5, 8, 9, 12, 13},
			Questions: q.Questions{
				q.Question{
					// Far-Left
					Prompt: "Markets are a bad way to distribute resources.",
					Map:    []int{0, 4, 5, 8, 12},
				},
				q.Question{
					// Middle-Left
					Prompt: "Communism doesn't work, socialism is more practical.",
					Map:    []int{1, 5, 9, 13},
				},
				q.Question{
					// Pro-State
					Prompt: "Government/State is the best way to complete large-scale projects such as building roads.",
					Map:    []int{0, 1, 2, 3, 4, 5, 6, 7, 10, 11},
				},
				q.Question{
					// Anti-State
					Prompt: "The State is always oppressive, no matter who is running it.",
					Map:    []int{8, 9, 10, 11, 12, 13, 14, 15},
				},
			},
		},
		q.QuestionSet{
			// Left-Wing
			Questions: q.Questions{
				q.Question{
					// Far-Left
					Prompt: "Power grows out of the barrel of a gun",
					Map:    q.MFarLeft,
				},
				q.Question{
					// Middle-Left
					Prompt: "Capitalism can be tamed by strong socialist reforms",
					Map:    q.MMiddleLeft,
				},
				q.Question{
					// Pro-State
					Prompt: "Discipline is a virtue",
					Map:    q.MProState,
				},
				q.Question{
					// Anti-State
					Prompt: "Absolute power corrupts, absolutely",
					Map:    q.MAntiState,
				},
			},
		},
		q.QuestionSet{
			// Right-Wing
			Questions: q.Questions{
				q.Question{
					// Far-Right
					Prompt: "Free Markets = Free People",
					Map:    []int{3, 7, 11, 15},
				},
				q.Question{
					// Liberal (Middle-Right)
					Prompt: "Regulations are necessary for a strong economy.",
					Map:    []int{2, 6, 10, 14},
				},
				q.Question{
					// Pro-State
					Prompt: "Prisons make society safer.",
					Map:    []int{1, 2, 3, 5, 6, 7, 10, 11},
				},
				q.Question{
					// Anti-State
					Prompt: "The goverment is always oppressive, no matter who is running it.",
					Map:    []int{8, 9, 10, 11, 12, 13, 14, 15},
				},
			},
		},
		// q.QuestionSet{
		// 	// Authoritarian Socialist
		// 	Mask: q.SocialistMask,
		// 	Questions: q.Questions{
		// 		q.Question{
		// 			// High Authoritarian
		// 			Prompt: "Centralized power is essential to sustaining political organizations",
		// 			Map:    q.MAuthoritarian,
		// 		},
		// 		q.Question{
		// 		// Middle Authoritarian
		// 		},
		// 		q.Question{
		// 		// Far-Left
		// 		},
		// 		q.Question{
		// 		// Middle-Left
		// 		},
		// 	},
		// },
		q.QuestionSet{
			// Libertarian Socialist
			Mask: q.AnarchistMask,
			Questions: q.Questions{
				q.Question{
					// Far-Left
					Prompt: "Work should be abolished.",
					Map:    q.MFarLeft,
				},
				q.Question{
					// Middle-Left
					Prompt: "The workplace is the primary arena of struggle against the excesses of capitalism.",
					Map:    q.MMiddleLeft,
				},
				q.Question{
					// Anti-organization
					Prompt: "Large organizations are inherently authoritarian.",
					Map:    q.MAntiAuthoritarian,
				},
				q.Question{
					// Pro-organization
					Prompt: "It may be necessary to negotiate with capitalists and politicians in furthering political goals.",
					Map:    q.MMiddleAuthoritarian,
				},
			},
		},
		// q.QuestionSet{
		// 		Question{
		// 			// Right-Wing
		// 			Prompt: "The best way to institute change in the system is to work within it.",
		// 			Map:    []int{1, 3, 5, 6, 7, 10, 11, 14, 15},
		// 		},
		// 		Question{
		// 			// Right-Wing
		// 			Prompt: "The borders should be patrolled to prevent illegal immigration.",
		// 			Map:    []int{2, 3, 6, 7, 10, 11, 14, 15},
		// 		},
		// 		Question{
		// 			// Pro-Capital
		// 			Prompt: "When someone is poor, it is mostly their fault.",
		// 			Map:    []int{2, 3, 6, 7, 10, 11, 14, 15},
		// 		},
		// 		Question{
		// 			// Anti-State
		// 			Prompt: "People should be able to lead their own lives, free of government intervention.",
		// 			Map:    []int{8, 9, 10, 11, 12, 13, 14, 15},
		// 		},
		// 		Question{
		// 			// Anti-Capital
		// 			Prompt: "Housing should be considered a human right.",
		// 			Map:    []int{0, 1, 4, 5, 8, 9, 12, 13},
		// 		},
		// 		Question{
		// 			// Anti-Capital
		// 			Prompt: "All education, including higher education, should be free",
		// 			Map:    []int{0, 1, 4, 5, 8, 9, 12, 13},
		// 		},
		// },
		// Question{
		// 	// Traditional (Far-Right)
		// 	Prompt: "The traditional institution of marriage is sacred.",
		// 	Map:    []int{2, 3, 7, 11},
		// },
		// Question{
		// 	// Pro-Liberal (sort of)
		// 	Prompt: "Voting is not just a right; it's a duty.",
		// 	Map:    []int{2, 3, 6, 7, 10, 11},
		// },
		// QuestionSet{
		// Question{
		// 	// Pro-Liberal
		// 	Prompt: "Loyalty to one's country is of utmost importance.",
		// 	Map:    []int{2, 3, 6, 7},
		// },
		// Question{
		// 	// Pro-Capital
		// 	Prompt: "Social welfare programs, such as food stamps, just make people lazy.",
		// 	Map:    []int{2, 3, 6, 7, 10, 11, 14, 15},
		// },
		// Question{
		// 	// Pro-State
		// 	Prompt: "Despite corruption, governments are, at their core, for the good of society.",
		// 	Map:    []int{0, 1, 2, 3, 4, 5, 6, 7},
		// },
		// Question{
		// 	// Pro-Capital
		// 	Prompt: "Management/CEOs deserve a higher salary than their employees.",
		// 	Map:    []int{2, 3, 6, 7, 10, 11, 14, 15},
		// },
		// 	Question{
		// 		// Far-Right
		// 		Prompt: "Some races/ethnicities are superior to others.",
		// 		Map:    []int{2, 3, 7, 11, 14, 15},
		// 	},
		// 	Question{
		// 		// Far-Right
		// 		Prompt: "The best way to end a recession is to remove restrictions on wealth creators.",
		// 		Map:    []int{3, 7, 11, 15},
		// 	},
		// 	Question{
		// 		// Pro-Capital
		// 		Prompt: "Those with higher income should be allowed to pay for better healthcare than those with lower income.",
		// 		Map:    []int{2, 3, 6, 7, 10, 11, 14, 15},
		// 	},
		// 	Question{
		// 		// Pro-Capital
		// 		Prompt: "The primary function of education should be to prepare students for their careers.",
		// 		Map:    []int{2, 3, 6, 7, 10, 11, 14, 15},
		// 	},
		// 	Question{
		// 		// Far-Right
		// 		Prompt: "A country shouldn't concern itself with the problems of refugees.",
		// 		Map:    []int{3, 7, 11, 15},
		// 	},
		// 	Question{
		// 		// Far-Right
		// 		Prompt: "Free Markets = Free People",
		// 		Map:    []int{3, 7, 11, 15},
		// 	},
		// }
	}

}

// QuestionIndex finds a random QuestionSet, shuffles, and shows it
func QuestionIndex(c *gin.Context) {
	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	set, err := questionSets.NextSet(user.CenterX, user.CenterY)
	if err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	c.JSON(http.StatusOK, set)
}
