package questions

var RightWingQuestions = QuestionSets{
	QuestionSet{
		// Right-Wing
		Mask: MProCapital,
		Questions: Questions{
			Question{
				// Far-Right
				Prompt: "Free Markets = Free People",
				Map:    MFarRight,
			},
			Question{
				// Liberal (Middle-Right)
				Prompt: "Regulations are necessary for a strong economy.",
				Map:    MMiddleRight,
			},
			Question{
				// Pro-State
				Prompt: "Prisons make society safer.",
				Map:    MProState,
			},
			Question{
				// Anti-State
				Prompt: "The goverment is always oppressive, no matter who is running it.",
				Map:    MAntiState,
			},
		},
	},
	QuestionSet{
		Mask: MProCapital,
		Questions: Questions{
			Question{
				Prompt: "The traditional institution of marriage is sacred.",
				Map:    MFarRight,
			},
			Question{
				Prompt: "Corporations cannot be trusted to regulate themselves.",
				Map:    MMiddleRight,
			},
			Question{
				Prompt: "If you've done nothing wrong, you have nothing to fear.",
				Map:    MProState,
			},
			Question{
				Prompt: "Governments are the biggest source of oppression.",
				Map:    MAntiState,
			},
		},
	},
	QuestionSet{
		Mask: MProCapital,
		Questions: Questions{
			Question{
				Prompt: "The best way to end a recession is to remove restrictions on wealth creators.",
				Map:    MFarRight,
			},
			Question{
				Prompt: "Laws should be in place to protect ethnic, religious, and other minority groups.",
				Map:    MMiddleRight,
			},
			Question{
				Prompt: "The police are mostly good.",
				Map:    MProState,
			},
			Question{
				Prompt: "End drug prohibition, and end the drug war.",
				Map:    MAntiState,
			},
		},
	},
	QuestionSet{
		Mask: MProCapital,
		Questions: Questions{
			Question{
				Prompt: "The primary function of education should be to prepare students for their careers.",
				Map:    MFarRight,
			},
			Question{
				Prompt: "Environmental protections are fundamentally a good thing.",
				Map:    MMiddleRight,
			},
			Question{
				Prompt: "Voting isn't just a right, it's a duty.",
				Map:    MProState,
			},
			Question{
				Prompt: "The US government is, in one way or another, responsible for 9/11.",
				Map:    MAntiState,
			},
		},
	},
	QuestionSet{
		Mask: MProCapital,
		Questions: Questions{
			Question{
				Prompt: "Affirmative action is reverse-racism",
				Map:    MFarRight,
			},
			Question{
				Prompt: "Taxes should be collected to fund public services.",
				Map:    MMiddleRight,
			},
			Question{
				Prompt: "Refugees and immigrants should not be given assylum.",
				Map:    MProState,
			},
			Question{
				Prompt: "There is secret world order controlling governments and world events.",
				Map:    MAntiState,
			},
		},
	},
	QuestionSet{
		Mask: MProCapital,
		Questions: Questions{
			Question{
				Prompt: "Government would be more efficient if it were run like a business.",
				Map:    MFarRight,
			},
			Question{
				Prompt: "The 1% should be taxed higher.",
				Map:    MMiddleRight,
			},
			Question{
				Prompt: "Loyalty to one's country is of utmost importance.",
				Map:    MProState,
			},
			Question{
				Prompt: "The NSA programs are concerned controlling the population, not catching terrorists.",
				Map:    MAntiState,
			},
		},
	},
}
