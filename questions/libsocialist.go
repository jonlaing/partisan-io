package questions

var LibertarianSocialistQuestions = QuestionSets{
	QuestionSet{
		// Libertarian Socialist
		Mask: MLibSocialist,
		Questions: Questions{
			Question{
				// Far-Left
				Prompt: "Work should be abolished.",
				Map:    MFarLeft,
			},
			Question{
				// Middle-Left
				Prompt: "The workplace is the primary arena of struggle against the excesses of capitalism.",
				Map:    MMiddleLeft,
			},
			Question{
				// Anti-organization
				Prompt: "Large bureaucratic organizations are inherently authoritarian.",
				Map:    MAntiAuthoritarian,
			},
			Question{
				// Pro-organization
				Prompt: "It may be necessary to negotiate with capitalists and politicians in furthering political goals.",
				Map:    MMiddleAntiAuthoritarian,
			},
		},
	},
}
