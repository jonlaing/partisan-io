package questions

var InitialQuestions = QuestionSets{
	QuestionSet{
		// Initial
		Questions: Questions{
			Question{
				// ProCapital
				Prompt: "Of all economic systems, Capitalism is most compatible with human nature.",
				Map:    MProCapital,
			},
			Question{
				// Anti-State
				Prompt: "Authority should always be questioned.",
				Map:    MAntiState,
			},
			Question{
				// Pro-State
				Prompt: "The police, in general, are good.",
				Map:    MProState,
			},
			Question{
				// Anti-Capital
				Prompt: "People over profits.",
				Map:    MAntiCapital,
			},
		},
	},
	QuestionSet{
		// Initial
		Questions: Questions{
			Question{
				// Right-Wing
				Prompt: "Men are more suited for positions of leadership, while women are more suited for positions of nurturing.",
				Map:    MProCapital,
			},
			Question{
				// Pro-State
				Prompt: "You should always support the troops.",
				Map:    MProState,
			},
			Question{
				// Anti-Capital
				Prompt: "Healthcare should be free of charge.",
				Map:    MAntiCapital,
			},
			Question{
				// Anti-State
				Prompt: "Governments care about control, not the good of their people.",
				Map:    MAntiState,
			},
		},
	},
	QuestionSet{
		// Initial
		Questions: Questions{
			Question{
				// Pro-Capital
				Prompt: "Hard work leads to upward social mobility",
				Map:    MProCapital,
			},
			Question{
				// Pro-State
				Prompt: "Good social organization is centralized",
				Map:    MProState,
			},
			Question{
				// Anti-Capital
				Prompt: "Corporations should not self-regulate",
				Map:    MAntiCapital,
			},
			Question{
				// Anti-State
				Prompt: "Dissent is a vitrue",
				Map:    MAntiState,
			},
		},
	},
	QuestionSet{
		Questions: Questions{
			Question{
				// Pro-Capital
				Prompt: "The best way to institute change in the system is to work within it.",
				Map:    MProCapital,
			},
			Question{
				// Pro-State
				Prompt: "Despite corruption, governments are, at their core, for the good of society.",
				Map:    MProState,
			},
			Question{
				// Anti-State
				Prompt: "People should be able to lead their own lives, free of government intervention.",
				Map:    MAntiState,
			},
			Question{
				// Anti-Capital
				Prompt: "All education, including higher education, should be free",
				Map:    MAntiCapital,
			},
		},
	},
	QuestionSet{
		Questions: Questions{
			Question{
				// Right-Wing
				Prompt: "The borders should be patrolled to prevent illegal immigration.",
				Map:    MProCapital,
			},
			Question{
				// Anti-Capital
				Prompt: "Housing should be considered a human right.",
				Map:    MAntiCapital,
			},
			Question{
				// Anti-State
				Prompt: "The government should not be trusted, ever.",
				Map:    MAntiState,
			},
			Question{
				Prompt: "Being an informed voter can institute change.",
				Map:    MProState,
			},
		},
	},
	QuestionSet{
		Questions: Questions{
			Question{
				// Pro-Capital
				Prompt: "When someone is poor, it is mostly their fault.",
				Map:    MProCapital,
			},
			Question{
				//Pro-State
				Prompt: "The military exists to protect our freedoms.",
				Map:    MProState,
			},
			Question{
				//Anti-Capital
				Prompt: "Capitalist industry is destroying the planet.",
				Map:    MAntiCapital,
			},
			Question{
				//AntiState
				Prompt: "\"The War on Terror\" is a ruse to further control the population.",
				Map:    MAntiState,
			},
		},
	},
	QuestionSet{
		Questions: Questions{
			Question{
				// Pro-Capital
				Prompt: "Social welfare programs, such as food stamps, just make people lazy.",
				Map:    MProCapital,
			},
			Question{
				// Pro-State
				Prompt: "If you haven't done anything wrong, you have nothing to hide.",
				Map:    MProState,
			},
			Question{
				// Anti-Capital
				Prompt: "\"From each according to their ablitity, to each according to their need\" is a good way to organize society.",
				Map:    MAntiCapital,
			},
			Question{
				// Anti-State
				Prompt: "No one chooses where they were born, so pride in one's country is foolish.",
				Map:    MAntiState,
			},
		},
	},
	QuestionSet{
		Questions: Questions{
			Question{
				// Pro-Capital
				Prompt: "Management/CEOs deserve a much higher salary than their employees.",
				Map:    MProCapital,
			},
			Question{
				// Pro-Liberal (sort of)
				Prompt: "Voting is not just a right; it's a duty.",
				Map:    MProState,
			},
			Question{
				// Anti-State
				Prompt: "What two consenting adults do in the bedroom is their business and should not be legislated.",
				Map:    MAntiState,
			},
			Question{
				// Anti-Capital
				Prompt: "The system is rigged in favor of the rich.",
				Map:    MAntiCapital,
			},
		},
	},
}

//
// 	Question{
// 		// Pro-Capital
// 		Prompt: "Those with higher income should be allowed to pay for better healthcare than those with lower income.",
// 		Map:    []int{2, 3, 6, 7, 10, 11, 14, 15},
// 	},
