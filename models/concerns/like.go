package models

// These are not meant to go into the database individually
// They are meant to be attached to another model
type Likes map[string]bool
type Dislikes map[string]bool

// special type of like that you get as a reward
type Stars map[string]bool
