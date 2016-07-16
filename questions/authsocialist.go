package questions

var AuthoritarianSocialistQuestions = QuestionSets{
	QuestionSet{
		// Authoritarian Socialist
		Mask: MAuthSocialist,
		Questions: Questions{
			Question{
				// High Authoritarian
				Prompt: "Centralized power is essential to sustaining political organizations.",
				Map:    MAuthoritarian,
			},
			Question{
				// Middle Authoritarian
				Prompt: "Over centralized power leads to corruption.",
				Map:    MMiddleAuthoritarian,
			},
			Question{
				// Far-Left
				Prompt: "An anti-capitalist struggle should strive for nothing less than a classless, stateless, moneyless society",
				Map:    MFarLeft,
			},
			Question{
				// Middle-Left
				Prompt: "\"Freedom of Speech\" should be protected at the state level, regardless of beliefs.",
				Map:    MMiddleLeft,
			},
		},
	},
	// QuestionSet{
	// 	Mask: MAuthSocialist,
	// 	Questions: Questions{
	// 		Question{
	// 			Prompt: "",
	// 			Map:    MAuthoritarian,
	// 		},
	// 		Question{
	// 			Prompt: "",
	// 			Map:    MMiddleAuthoritarian,
	// 		},
	// 		Question{
	// 			Prompt: "",
	// 			Map:    MFarLeft,
	// 		},
	// 		Question{
	// 			Prompt: "",
	// 			Map:    MMiddleLeft,
	// 		},
	// 	},
	// },
}
